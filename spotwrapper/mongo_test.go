package spotwrapper

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Rich-T-kid/musicShare/pkg/models"
)

// TestCache checks the basic functionality of the cache implementation
func TestCache(t *testing.T) {
	ctx := context.Background()
	cache := NewCache[string, string]() // Initialize cache with string keys & values

	// Test Set()
	cache.Set(ctx, "testKey", "testValue", 5) // Expires in 5 sec

	// Test Get()
	val := cache.Get(ctx, "testKey")
	assert.Equal(t, "testValue", val, "Expected value should match stored value")

	// Test Exist()
	exists := cache.Exist(ctx, "testKey")
	assert.True(t, exists, "Key should exist after being set")

	// Test Delete()
	cache.Delete(ctx, "testKey")
	existsAfterDelete := cache.Exist(ctx, "testKey")
	assert.False(t, existsAfterDelete, "Key should not exist after deletion")

	// Test Get() on missing key
	missingVal := cache.Get(ctx, "nonExistingKey")
	assert.Equal(t, "", missingVal, "Get() should return empty string for missing keys")

	// Wait for key expiration and test existence
	time.Sleep(6 * time.Second)
	expiredVal := cache.Get(ctx, "testKey")
	assert.Equal(t, "", expiredVal, "Key should have expired after timeout")
}

// Global variable to hold the DocumentStore for tests.
var store DocumentStore

// TestMain sets up the testing environment (ensuring that the expected databases and collections exist)
// and tears it down after all tests have run.
func TestMain(m *testing.M) {
	// Create our document store.
	store = NewDocumentStore()
	mongoStore, ok := store.(*MongoDBStore)
	if !ok {
		log.Fatalf("Expected store to be *MongoDBStore")
	}
	client := mongoStore.client
	ctx := context.Background()

	// Ensure that both expected databases exist with the expected collections.
	// (Note: In MongoDB a database exists if it contains at least one collection.)
	expectedDBs := []string{"test_db", "prod_db"}
	expectedCollections := []string{"users", "comments", "songs"}
	for _, dbName := range expectedDBs {
		db := client.Database(dbName)
		for _, colName := range expectedCollections {
			// Try to create the collection. If it exists, we ignore the error.
			_ = db.CreateCollection(ctx, colName)
		}
	}

	// Run the tests.
	code := m.Run()

	// Cleanup: drop the databases used in testing.
	//_ = client.Database("test_db").Drop(ctx)
	//_ = client.Database("prod_db").Drop(ctx)

	os.Exit(code)
}

// TestConnected checks that the Connected method works.
func TestConnected(t *testing.T) {
	ctx := context.Background()
	if err := store.Connected(ctx); err != nil {
		t.Fatalf("Connected() failed: %v", err)
	}
}

// TestUserStore exercises the user-related methods.
func TestUserStore(t *testing.T) {
	// Generate a unique user ID.
	userID := "richard"
	// Create a sample user document.
	user := &models.UserMongoDocument{
		UUID: userID,
		Comments: []models.UserComments{
			{
				Username: userID,
				Rating:   5,
				Review:   "Great song!",
				SongID:   "song1",
			},
		},
	}

	// Test SaveUser.
	if err := store.SaveUser(user); err != nil {
		t.Fatalf("SaveUser failed: %v", err)
	}

	// Test GetUserByID.
	gotUser, err := store.GetUserByID(userID)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if gotUser.UUID != userID {
		t.Errorf("GetUserByID returned wrong user; expected ID %s, got %s", userID, gotUser.UUID)
	}

	// Test GetUserComments.
	comments, err := store.GetUserComments(userID)
	if err != nil {
		t.Fatalf("GetUserComments failed: %v", err)
	}
	if len(comments) != len(user.Comments) {
		t.Errorf("GetUserComments returned %d comments; expected %d", len(comments), len(user.Comments))
	}

	// For GetUserSongs, we need to ensure that songs exist.
	// Create three sample songs.
	s1 := &models.SongTypes{SongURI: "song1", Comments: []models.UserComments{}}
	s2 := &models.SongTypes{SongURI: "song2", Comments: []models.UserComments{}}
	s3 := &models.SongTypes{SongURI: "song3", Comments: []models.UserComments{}}

	// Insert songs (using the InsertSong method, which handles duplicates via upsert).
	for _, s := range []*models.SongTypes{s1, s2, s3} {
		if err := store.InsertSong(s); err != nil {
			t.Fatalf("InsertSong failed for song %s: %v", s.SongURI, err)
		}
	}

	// Test GetUserSongs.
	userSongs, err := store.GetUserSongs(userID)
	if err != nil {
		t.Fatalf("GetUserSongs failed: %v", err)
	}
	// We expect songs "song1", "song2", and "song3".
	expectedSongs := map[string]bool{"song1": true, "song2": true, "song3": true}
	for _, s := range userSongs {
		if !expectedSongs[s.SongURI] {
			t.Errorf("GetUserSongs returned unexpected song %s", s.SongURI)
		}
		delete(expectedSongs, s.SongURI)
	}
	if len(expectedSongs) > 0 {
		t.Errorf("GetUserSongs did not return expected songs; missing: %v", expectedSongs)
	}
}

// TestSongStore exercises the song-related methods.
func TestSongStore(t *testing.T) {
	// Generate a unique song URI.
	songURI := "SongTest" //"testSong-" + uuid.New().String()
	testUUID := "CustomUUID"
	song := &models.SongTypes{
		SongURI:  songURI,
		Comments: []models.UserComments{},
		UUID:     testUUID,
	}

	// Test InsertSong.
	if err := store.InsertSong(song); err != nil {
		t.Fatalf("InsertSong failed: %v", err)
	}

	// Test GetSongByID.
	gotSong, err := store.GetSongByID(songURI)
	if err != nil {
		t.Fatalf("GetSongByID failed: %v", err)
	}
	if gotSong.SongURI != songURI {
		t.Errorf("GetSongByID returned wrong songURI; expected %s, got %s", songURI, gotSong.SongURI)
	}

	// Test DeleteSong.
	if err := store.DeleteSong(songURI); err != nil {
		t.Fatalf("DeleteSong failed: %v", err)
	}
	// Verify deletion.
	_, err = store.GetSongByID(songURI)
	if err == nil {
		t.Errorf("GetSongByID should have failed after deletion, but it did not")
	}
}

// TestCommentStore exercises the comment-related methods.
func TestCommentStore(t *testing.T) {
	// Generate a unique song URI for comment testing.
	songURI := uuid.New().String()

	song := &models.SongTypes{
		SongURI:  songURI,
		Comments: []models.UserComments{},
	}
	// Insert the song.
	if err := store.InsertSong(song); err != nil {
		t.Fatalf("InsertSong (for comment test) failed: %v", err)
	}

	// Submit a comment.
	commentID := uuid.New().String()
	comment := models.UserComments{
		UUID:     commentID,
		Username: "commenter",
		Rating:   4,
		Review:   "Nice song!",
		SongID:   songURI,
	}
	_, err := store.SubmitComment(songURI, comment)
	if err != nil {
		t.Fatalf("SubmitComment failed: %v", err)
	}

	// GetComments for the song.
	comments, err := store.GetComments(songURI)
	if err != nil {
		t.Fatalf("GetComments failed: %v", err)
	}
	found := false
	for _, c := range comments {
		if c.UUID == commentID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("GetComments did not return the submitted comment with ID %s", commentID)
	}

	// GetComment by commentID.
	gotComment, err := store.GetComment(commentID)
	if err != nil {
		t.Fatalf("GetComment failed: %v", err)
	}
	if gotComment.UUID != commentID {
		t.Errorf("GetComment returned wrong comment; expected ID %s, got %s", commentID, gotComment.UUID)
	}

	// UpdateComment: change the review text.
	newComment := comment
	newComment.Review = "Awesome song!"
	updated, err := store.UpdateComment(commentID, newComment)
	if err != nil {
		t.Fatalf("UpdateComment failed: %v", err)
	}
	if !updated {
		t.Errorf("UpdateComment did not update any document")
	}
	// Verify the update.
	updatedComment, err := store.GetComment(commentID)
	if err != nil {
		t.Fatalf("GetComment after update failed: %v", err)
	}
	if updatedComment.Review != "Awesome song!" {
		t.Errorf("UpdateComment did not update review correctly; expected 'Awesome song!', got '%s'", updatedComment.Review)
	}

	// DeleteComment.
	if err := store.DeleteComment(commentID); err != nil {
		t.Fatalf("DeleteComment failed: %v", err)
	}
	// Verify deletion.
	_, err = store.GetComment(commentID)
	if err == nil {
		t.Errorf("GetComment should have failed after deletion of comment ID %s, but it did not", commentID)
	}
}
