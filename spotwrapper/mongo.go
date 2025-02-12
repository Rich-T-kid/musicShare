package spotwrapper

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Rich-T-kid/musicShare/pkg"
)

const DatabaseName = "test_db"

type cache[T comparable, V any] interface {
	Get(ctx context.Context, key string) V
	Set(ctx context.Context, key string, data T, expire int) // Handles inserts & updates
	Delete(ctx context.Context, key string)
	Exist(ctx context.Context, key string) bool
	StoreTokens(id, access, refresh string) error
	GetTokens(id string) (string, string, error) // if more flexibility is needed later return a slice of strings
}

// Factory function to create a new cache instance
func NewCache[T comparable, V ~string]() cache[T, V] {
	return newSpotCache[T, V]()
}

// spotish struct implementing the cache interface
type spotish[T comparable, V any] struct {
	client *redis.Client
}

// Function to initialize a new Redis client
func newSpotCache[T comparable, V any]() *spotish[T, V] {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
	})
	return &spotish[T, V]{client: client}
}

// Get method retrieves a value from Redis cache
func (s *spotish[T, V]) Get(ctx context.Context, key string) V {
	str, err := s.client.Get(key).Result()
	if err != nil {
		var zero V
		fmt.Printf("Get key: %s , Error From cache %e \n ", key, err)
		return zero
	}

	// Convert string back to type V
	value, ok := any(str).(V)
	if !ok {
		var zero V
		fmt.Printf("Could not convert key %s of value %s to type (Cache[key,**Value]) Value of reddis implentation cache", key, str)
		return zero
	}

	return value
}

// Set method stores a value in Redis cache
func (s *spotish[T, V]) Set(ctx context.Context, key string, data T, expire int) {
	err := s.client.Set(key, fmt.Sprintf("%v", data), time.Duration(expire)*time.Hour).Err()
	if err != nil {
		fmt.Println("Error setting cache:", err)
	}
}

// Delete method removes a key from Redis cache
func (s *spotish[T, V]) Delete(ctx context.Context, key string) {
	err := s.client.Del(key).Err()
	if err != nil {
		fmt.Println("Error deleting cache key:", err)
	}
}

// Exist method checks if a key exists in Redis cache
func (s *spotish[T, V]) Exist(ctx context.Context, key string) bool {
	exists, err := s.client.Exists(key).Result()
	if err != nil {
		return false
	}
	return exists > 0
}
func (s *spotish[T, V]) StoreTokens(userID, accessToken, refreshToken string) error {
	key := fmt.Sprintf("user:%s", userID) // Store by user ID

	err := s.client.HSet(key, "access_token", accessToken).Err()
	if err != nil {
		return fmt.Errorf("failed to store access token in Redis: %w", err)
	}

	err = s.client.HSet(key, "refresh_token", refreshToken).Err()
	if err != nil {
		return fmt.Errorf("failed to store refresh token in Redis: %w", err)
	}

	return nil

}

func (s *spotish[T, V]) GetTokens(userID string) (string, string, error) {

	key := fmt.Sprintf("user:%s", userID)

	// Fetch both access and refresh tokens
	accessToken, err := s.client.HGet(key, "access_token").Result()
	if err != nil {
		return "", "", fmt.Errorf("error fetching access token: %w", err)
	}

	refreshToken, err := s.client.HGet(key, "refresh_token").Result()
	if err != nil {
		return "", "", fmt.Errorf("error fetching refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// TODO: Implement all of this. In mongo DB fr this time
// MongoDb  Implementation
// functions below should rely on this interface for now
// seperate ticket to imlement all of this

/*
Mongo DB implementation below
*/

type MongoDBStore struct {
	client      *mongo.Client
	databases   []string
	collections []string
}
type DocumentStore interface {
	// Connectivity Check
	Connected(ctx context.Context) error // Pings the database and returns whether it's connected.
	UserStore
	SongStore
	CommentStore
	// CRUD Operations
}
type UserStore interface {
	GetUserByID(userID string) (*pkg.UserMongoDocument, error)
	SaveUser(user *pkg.UserMongoDocument) error
	GetUserSongs(userID string) ([]pkg.SongTypes, error)
	GetUserComments(userID string) ([]pkg.UserComments, error)
}

// TODO: Change this after merge
type SongStore interface {
	GetSongByID(songID string) (*pkg.SongTypes, error)
	InsertSong(song *pkg.SongTypes) error
	UpdateSong(song *pkg.SongTypes) error
	DeleteSong(songID string) error
}
type CommentStore interface {
	SubmitComment(songID string, comment pkg.UserComments) error
	GetComments(songID string) ([]pkg.UserComments, error)
	UpdateComment(oldComment string, newComment pkg.UserComments) (bool, error)
	DeleteComment(commentID string) error
	GetComment(commentID string) (*pkg.UserComments, error)
}

/*
Mongo DB implementation below
*/
func newDocumentStore() DocumentStore {
	// Define MongoDB connection URI (matches Docker container settings)
	mongoURI := "mongodb://admin:secretpassword@localhost:27017/"

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Establish connection
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the database to check if it's reachable
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB connection test failed: %v", err)
	}

	fmt.Println("✅ Successfully connected to MongoDB")
	db := []string{"test_db", "prod_db"}
	collection := []string{"users", "comments", "songs"}
	return &MongoDBStore{client: client,
		databases:   db,
		collections: collection,
	}
}

func (m *MongoDBStore) Connected(ctx context.Context) error {
	// Use the provided context, or create one with a timeout if needed.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Ping the MongoDB server.
	if err := m.client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("mongo ping failed: %w", err)
	}

	// List available databases.
	availableDBs, err := m.client.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		return fmt.Errorf("could not list databases: %w", err)
	}
	//fmt.Println("Available databases:", availableDBs)

	// Verify each expected database and its collections.
	for _, expectedDB := range m.databases {
		if !contains(availableDBs, expectedDB) {
			return fmt.Errorf("expected database %q not found", expectedDB)
		}

		// List collections for the expected database.
		availableCols, err := m.client.Database(expectedDB).ListCollectionNames(ctx, bson.D{})
		if err != nil {
			return fmt.Errorf("failed to list collections for database %q: %w", expectedDB, err)
		}
		//fmt.Printf("Database %q collections: %v\n", expectedDB, availableCols)

		// Check that each expected collection exists.
		for _, expectedCol := range m.collections {
			if !contains(availableCols, expectedCol) {
				return fmt.Errorf("expected collection %q not found in database %q", expectedCol, expectedDB)
			}
		}
	}

	return nil
}

// Implement UserStore methods

func (m *MongoDBStore) GetUserByID(userID string) (*pkg.UserMongoDocument, error) {
	db := m.client.Database(DatabaseName)
	collection := db.Collection("users")

	// Filter by the application-provided user ID stored in the "id" field.
	filter := bson.M{"uuid": userID}

	var userDoc pkg.UserMongoDocument
	err := collection.FindOne(context.TODO(), filter).Decode(&userDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no user found with ID %s", userID)
		}
		return nil, fmt.Errorf("error retrieving user document: %w", err)
	}

	return &userDoc, nil
}
func (m *MongoDBStore) SaveUser(user *pkg.UserMongoDocument) error {
	db := m.client.Database(DatabaseName)
	collection := db.Collection("users")
	if user.UUID == "" {

		user.UUID = newID()
	}
	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)
	fmt.Println("Inserted document ID:", insertResult.InsertedID)

	return nil /* implementation */
}

// GetUserComments retrieves all comments stored in the user document that has the application-provided ID.
func (m *MongoDBStore) GetUserComments(userID string) ([]pkg.UserComments, error) {
	db := m.client.Database(DatabaseName)
	collection := db.Collection("users")

	// Use the application-provided userID stored in the "ID" field.
	filter := bson.M{"uuid": userID}
	var userDoc pkg.UserMongoDocument
	if err := collection.FindOne(context.TODO(), filter).Decode(&userDoc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no user found with ID %s", userID)
		}
		return nil, fmt.Errorf("error retrieving user document: %w", err)
	}

	return userDoc.Comments, nil
}

// GetSongByID retrieves a song document from the "songs" collection using its song UUID or SongURI.
func (m *MongoDBStore) GetSongByID(songID string) (*pkg.SongTypes, error) {
	db := m.client.Database(DatabaseName)
	collection := db.Collection("songs")
	fmt.Println("Looking for song with ID ", songID)
	filter := bson.M{"songURI": songID}
	var song pkg.SongTypes
	if err := collection.FindOne(context.TODO(), filter).Decode(&song); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no song found with id %s", songID)
		}
		return nil, fmt.Errorf("error retrieving song: %w", err)
	}

	return &song, nil
}

// GetUserSongs retrieves all songs associated with the user. It does so by:
// 1. Fetching the user document by the provided ID.
// 2. Combining the song URIs from the Listened, LikedSongs, and DislikedSongs arrays.
// 3. Querying the "songs" collection for documents with a matching songURI.
func (m *MongoDBStore) GetUserSongs(userID string) ([]pkg.SongTypes, error) {
	db := m.client.Database(DatabaseName)
	usersCollection := db.Collection("users")

	// Retrieve the user document using the application-provided ID.
	filter := bson.M{"uuid": userID}
	var userDoc pkg.UserMongoDocument
	if err := usersCollection.FindOne(context.TODO(), filter).Decode(&userDoc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no user found with ID %s", userID)
		}
		return nil, fmt.Errorf("error retrieving user document: %w", err)
	}

	// Combine song URIs from Listened, LikedSongs, and DislikedSongs into a unique set.
	songIDSet := make(map[string]struct{})
	for _, uri := range userDoc.Listened {
		songIDSet[uri.Song] = struct{}{}
	}

	// Convert the set into a slice.
	var songIDs []string
	for id := range songIDSet {
		songIDs = append(songIDs, id)
	}

	// **Fix: Ensure songIDs is not empty before querying**
	if len(songIDs) == 0 {
		var empty []pkg.SongTypes
		return empty, nil // Return an empty list instead of querying MongoDB
	}

	// Query the "songs" collection for documents whose songURI is in songIDs.
	songsCollection := db.Collection("songs")
	query := bson.M{"songURI": bson.M{"$in": songIDs}}

	cursor, err := songsCollection.Find(context.TODO(), query)
	if err != nil {
		return nil, fmt.Errorf("error querying songs: %w", err)
	}
	defer cursor.Close(context.TODO())

	var songs []pkg.SongTypes
	for cursor.Next(context.TODO()) {
		var song pkg.SongTypes
		if err := cursor.Decode(&song); err != nil {
			return nil, fmt.Errorf("error decoding song document: %w", err)
		}
		songs = append(songs, song)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return songs, nil
}

// handles duplicates as well
func (m *MongoDBStore) InsertSong(song *pkg.SongTypes) error {
	db := m.client.Database(DatabaseName)
	collection := db.Collection("songs")

	// Ensure the song has a UUID; if not, generate one.
	if song.UUID == "" {
		song.UUID = newID() // newID() should generate a new random UUID.
	}

	// Use the songURI as the unique identifier for the document.
	filter := bson.M{"songURI": song.SongURI}

	// Use $setOnInsert to set fields only if a new document is inserted.
	// This ensures that if a song with the same SongURI already exists,
	// its uuid and AlternateName remain unchanged.
	update := bson.M{
		"$setOnInsert": bson.M{
			"songURI":       song.SongURI,
			"uuid":          song.UUID,
			"AlternateName": song.AlternateName,
		},
		// $push new comments to the existing array.
		"$push": bson.M{
			"comments": bson.M{
				"$each": song.Comments,
			},
		},
	}

	// Upsert: update if the document exists, otherwise insert a new document.
	updateOptions := options.Update().SetUpsert(true)

	result, err := collection.UpdateOne(context.TODO(), filter, update, updateOptions)
	if err != nil {
		return fmt.Errorf("failed to upsert song: %v", err)
	}

	// The behavior:
	// - If result.MatchedCount > 0, a song with that SongURI existed and only new comments were appended.
	// - If result.UpsertedCount > 0, a new song document was inserted.
	if result.MatchedCount > 0 {
		fmt.Printf("Updated existing document for song %s with new comments.\n", song.SongURI)
	} else if result.UpsertedCount > 0 {
		fmt.Printf("Inserted new document for song %s.\n", song.SongURI)
	} else {
		fmt.Println("No changes made to the document.")
	}

	return nil
}

func (m *MongoDBStore) UpdateSong(song *pkg.SongTypes) error { return nil /* implementation */ }

// DeleteSong deletes a song document from the "songs" collection using its SongURI.
func (m *MongoDBStore) DeleteSong(songID string) error {
	db := m.client.Database(DatabaseName)
	collection := db.Collection("songs")

	filter := bson.M{"songURI": songID}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("failed to delete song: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("no song found with songURI: %s", songID)
	}

	fmt.Printf("Deleted song with songURI: %s\n", songID)
	return nil
}

// Implement CommentStore methods
func (m *MongoDBStore) SubmitComment(songID string, comment pkg.UserComments) error {
	db := m.client.Database(DatabaseName)
	collection := db.Collection("songs")

	// Ensure the comment has a UUID
	if comment.UUID == "" {
		comment.UUID = newID()
	}
	comment.SongID = songID

	// Filter to check if the song exists
	filter := bson.M{"songURI": songID}

	// Update operation to add a comment
	update := bson.M{
		"$push": bson.M{
			"comments": comment,
		},
	}

	// Try updating the song document
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update song: %w", err)
	}

	// If the song does not exist, create a new document
	if updateResult.MatchedCount == 0 {
		newSong := pkg.SongTypes{
			SongURI:       songID,
			Comments:      []pkg.UserComments{comment},
			AlternateName: []string{},
			UUID:          newID(),
		}
		fmt.Println("New comment being added", comment)

		_, err := collection.InsertOne(context.TODO(), newSong)
		if err != nil {
			return fmt.Errorf("failed to insert new song: %w", err)
		}
		fmt.Printf("Created new song document with comment for %s\n", songID)
	} else {
		fmt.Printf("Successfully added comment to song %s\n", songID)
	}

	return nil
}

// GetComments returns all the comments for a given song identified by songID (SongURI) or UUID.
func (m *MongoDBStore) GetComments(songID string) ([]pkg.UserComments, error) {
	db := m.client.Database(DatabaseName)
	collection := db.Collection("songs")

	filter := bson.M{"songURI": songID}
	var song pkg.SongTypes
	err := collection.FindOne(context.TODO(), filter).Decode(&song)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			var empty []pkg.UserComments
			return empty, nil
			//return nil, fmt.Errorf("no song found with songURI: %s", songID)
		}
		return nil, fmt.Errorf("failed to find song: %w", err)
	}

	return song.Comments, nil
}

// UpdateComment searches for a comment with a matching comment ID
// and updates it to the newComment. Returns true if the comment was updated.
func (m *MongoDBStore) UpdateComment(commentID string, newComment pkg.UserComments) (bool, error) {
	db := m.client.Database(DatabaseName)
	collection := db.Collection("songs")

	// Find a song document containing a comment with the matching comment ID.
	filter := bson.M{"comments.uuid": commentID}

	// Update that comment using the positional "$" operator.
	update := bson.M{
		"$set": bson.M{
			"comments.$": newComment,
		},
	}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return false, fmt.Errorf("failed to update comment: %w", err)
	}
	if result.MatchedCount == 0 {
		return false, fmt.Errorf("no comment found with id: %s", commentID)
	}
	return result.ModifiedCount > 0, nil
}

// DeleteComment removes a comment from any song document by matching on its comment ID.
func (m *MongoDBStore) DeleteComment(commentID string) error {
	db := m.client.Database(DatabaseName)
	collection := db.Collection("songs")
	fmt.Println("Attempting to delete comment with ID :", commentID)
	// Filter: Find any song document that contains a comment with the matching uuid.
	filter := bson.M{"comments.uuid": commentID}
	// Update: Remove (pull) from the comments array any comment with that uuid.
	update := bson.M{
		"$pull": bson.M{"comments": bson.M{"uuid": commentID}},
	}

	// Perform the update.
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	fmt.Printf("updated %d documents\n", res.ModifiedCount)
	return nil
}

// GetComment retrieves a single comment by its comment ID.
func (m *MongoDBStore) GetComment(commentID string) (*pkg.UserComments, error) {
	db := m.client.Database(DatabaseName)
	collection := db.Collection("songs")

	// Filter to find any document that has a comment with the matching comment ID.
	filter := bson.M{"comments.uuid": commentID}
	// Projection returns only the matching element from the comments array.
	projection := bson.M{"comments.$": 1}
	opts := options.FindOne().SetProjection(projection)

	// Use an inline struct to decode only the comments array.
	var result struct {
		Comments []pkg.UserComments `bson:"comments"`
	}

	err := collection.FindOne(context.TODO(), filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no comment found with id: %s", commentID)
		}
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	if len(result.Comments) == 0 {
		return nil, fmt.Errorf("no comment found with id: %s", commentID)
	}

	return &result.Comments[0], nil
}
func newID() string {
	return uuid.New().String()
}

/*
End of mongoDB implementation
*/
var (
	database = newDocumentStore()
)

// TODO: Beofore testing on other code base. just implement the below methods using the interfaces
func SaveUser(user *pkg.UserMongoDocument) error {
	fmt.Println("Attempting to save user info with name ", user.UserProfileResponse.DisplayName)
	return database.SaveUser(user)
}

// for now just get working but later on there should be a slight level of
// misdirection so that we can handle errors better
func SubmitComment(songid string, comment pkg.UserComments) error {
	return database.SubmitComment(songid, comment)
}
func GetComments(songid string, limit, offset int) ([]pkg.UserComments, error) {
	return database.GetComments(songid)
}

// find old and update it with new comment. if old cant be found return error
func UpdateComment(oldCommentID string, new pkg.UserComments) (bool, error) {
	return database.UpdateComment(oldCommentID, new)
}

func GetComment(commentID string) (*pkg.UserComments, error) {
	return database.GetComment(commentID)
}

func DeleteComment(commentID string) error { return database.DeleteComment(commentID) }

// if any of these return nil it means it wasnt found
func GetUserDocument(userid string) (*pkg.UserMongoDocument, error) {
	return database.GetUserByID(userid)
}
func GetUserSongs(userid string) ([]pkg.SongTypes, error) {
	return database.GetUserSongs(userid)
}
func GetUserComments(userid string) ([]pkg.UserComments, error) {
	return database.GetUserComments(userid)
}
