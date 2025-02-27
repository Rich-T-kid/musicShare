package spotwrapper

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Rich-T-kid/musicShare/pkg"
	"github.com/Rich-T-kid/musicShare/pkg/models"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: "default",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err = client.Ping(context.TODO()).Result()
	if err != nil {
		log.Fatal("Reddis Instance is not returnng correctly (Response) ->  %e", err)
		return nil
	}
	fmt.Println("Reddis Instance is good to go")
	return &spotish[T, V]{client: client}
}

// Get method retrieves a value from Redis cache
func (s *spotish[T, V]) Get(ctx context.Context, key string) V {
	str, err := s.client.Get(ctx, key).Result()
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
	err := s.client.Set(ctx, key, fmt.Sprintf("%v", data), time.Duration(expire)*time.Hour).Err()
	if err != nil {
		fmt.Println("Error setting cache:", err)
	}
}

// Delete method removes a key from Redis cache
func (s *spotish[T, V]) Delete(ctx context.Context, key string) {
	err := s.client.Del(ctx, key).Err()
	if err != nil {
		fmt.Println("Error deleting cache key:", err)
	}
}

// Exist method checks if a key exists in Redis cache
func (s *spotish[T, V]) Exist(ctx context.Context, key string) bool {
	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false
	}
	return exists > 0
}
func (s *spotish[T, V]) StoreTokens(userID, accessToken, refreshToken string) error {
	key := fmt.Sprintf("user-token:%s", userID) // Store by user ID

	err := s.client.HSet(context.TODO(), key, "access_token", accessToken).Err()
	if err != nil {
		return fmt.Errorf("failed to store access token in Redis: %w", err)
	}

	err = s.client.HSet(context.TODO(), "refresh_token", refreshToken).Err()
	if err != nil {
		return fmt.Errorf("failed to store refresh token in Redis: %w", err)
	}

	return nil

}

func (s *spotish[T, V]) GetTokens(userID string) (string, string, error) {

	key := fmt.Sprintf("user-token:%s", userID)

	// Fetch both access and refresh tokens
	accessToken, err := s.client.HGet(context.TODO(), key, "access_token").Result()
	if err != nil {
		return "", "", fmt.Errorf("error fetching access token: %w", err)
	}

	refreshToken, err := s.client.HGet(context.TODO(), key, "refresh_token").Result()
	if err != nil {
		return "", "", fmt.Errorf("error fetching refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

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
	GetUserByID(userID string) (*models.UserMongoDocument, error)
	SaveUser(user *models.UserMongoDocument) error
	GetUserSongs(userID string) ([]models.SongTypes, error)
	GetUserComments(userID string) ([]models.UserComments, error)
}

// TODO: Change this after merge
type SongStore interface {
	GetSongByID(songID string) (*models.SongTypes, error)
	InsertSong(song *models.SongTypes) error
	AddSongtoDB(songURI string) error
	UpdateSong(song *models.SongTypes) error
	DeleteSong(songID string) error
}
type CommentStore interface {
	SubmitComment(songID string, comment models.UserComments) error
	GetComments(songID string) ([]models.UserComments, error)
	UpdateComment(oldComment string, newComment models.UserComments) (bool, error)
	DeleteComment(commentID string) error
	GetComment(commentID string) (*models.UserComments, error)
}

/*
Mongo DB implementation below
*/
func NewDocumentStore() DocumentStore {
	// Define MongoDB connection URI (matches Docker container settings)
	mongoURI := os.Getenv("MONGO_URI")

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Establish connection
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the database to check if it's reachable
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB connection test failed: %v", err)
	}

	fmt.Println("✅ Successfully connected to MongoDB")
	//db := []string{"test_db", "prod_db"}
	//collection := []string{"users", "comments", "songs"}
	return &MongoDBStore{client: client} //databases:   db,
	//collections: collection,

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
	fmt.Println("Available databases:", availableDBs)

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

func (m *MongoDBStore) GetUserByID(userID string) (*models.UserMongoDocument, error) {
	db := m.client.Database(DatabaseName)
	collection := db.Collection("users")

	// Filter by the application-provided user ID stored in the "id" field.
	filter := bson.M{"uuid": userID}

	var userDoc models.UserMongoDocument
	err := collection.FindOne(context.TODO(), filter).Decode(&userDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no user found with ID %s", userID)
		}
		return nil, fmt.Errorf("error retrieving user document: %w", err)
	}

	return &userDoc, nil
}

func (m *MongoDBStore) SaveUser(user *models.UserMongoDocument) error {
	db := m.client.Database(DatabaseName)
	collection := db.Collection("users")
	if user.UUID == "" {
		log.Fatal("A users UUID should always already be generated by this point")
		user.UUID = newID()
	}
	fmt.Printf("Adding user: %s with id: %s to MongoDB database %+v\n", user.UserProfileResponse.DisplayName, user.UUID, user)
	_, err := collection.InsertOne(context.TODO(), user)
	return err
}

// GetUserComments retrieves all comments stored in the user document that has the application-provided ID.
func (m *MongoDBStore) GetUserComments(userID string) ([]models.UserComments, error) {
	db := m.client.Database(DatabaseName)
	collection := db.Collection("users")

	// Use the application-provided userID stored in the "ID" field.
	filter := bson.M{"uuid": userID}
	var userDoc models.UserMongoDocument
	if err := collection.FindOne(context.TODO(), filter).Decode(&userDoc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no user found with ID %s", userID)
		}
		return nil, fmt.Errorf("error retrieving user document: %w", err)
	}

	return userDoc.Comments, nil
}

// GetSongByID retrieves a song document from the "songs" collection using its song UUID or SongURI.
func (m *MongoDBStore) GetSongByID(songID string) (*models.SongTypes, error) {
	db := m.client.Database(DatabaseName)
	collection := db.Collection("songs")
	filter := bson.M{"songURI": songID}
	var song models.SongTypes
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
func (m *MongoDBStore) GetUserSongs(userID string) ([]models.SongTypes, error) {
	db := m.client.Database(DatabaseName)
	usersCollection := db.Collection("users")

	// Retrieve the user document using the application-provided ID.
	filter := bson.M{"uuid": userID}
	var userDoc models.UserMongoDocument
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
		var empty []models.SongTypes
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

	var songs []models.SongTypes
	for cursor.Next(context.TODO()) {
		var song models.SongTypes
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
func (m *MongoDBStore) InsertSong(song *models.SongTypes) error {
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

func (m *MongoDBStore) AddSongtoDB(songURI string) error {
	db := m.client.Database(DatabaseName)
	collection := db.Collection("songs")

	filter := bson.M{"songURI": songURI}

	update := bson.M{
		"$setOnInsert": bson.M{
			"songURI":  songURI,
			"uuid":     newID(), // Generate only if new
			"comments": []models.UserComments{},
		},
	}

	updateOptions := options.Update().SetUpsert(true)

	result, err := collection.UpdateOne(context.TODO(), filter, update, updateOptions)
	if err != nil {
		return fmt.Errorf("failed to upsert song: %v", err)
	}

	if result.MatchedCount > 0 {
		fmt.Printf("Song already exists in DB: %s\n", songURI)
	} else if result.UpsertedCount > 0 {
		fmt.Printf("Inserted new song: %s\n", songURI)
	}
	return nil
}

func (m *MongoDBStore) UpdateSong(song *models.SongTypes) error { return nil /* implementation */ }

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
func (m *MongoDBStore) SubmitComment(songID string, comment models.UserComments) error {
	db := m.client.Database(DatabaseName)
	collection := db.Collection("songs")

	// Ensure the comment has a UUID
	// more imporantly every comment needs to have a new and unique uuid so we cant rely on what gets passed in
	comment.UUID = newID()
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
		newSong := models.SongTypes{
			SongURI:       songID,
			Comments:      []models.UserComments{comment},
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
func (m *MongoDBStore) GetComments(songID string) ([]models.UserComments, error) {
	db := m.client.Database(DatabaseName)
	collection := db.Collection("songs")

	filter := bson.M{"songURI": songID}
	var song models.SongTypes
	err := collection.FindOne(context.TODO(), filter).Decode(&song)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			var empty []models.UserComments
			return empty, nil
			//return nil, fmt.Errorf("no song found with songURI: %s", songID)
		}
		return nil, fmt.Errorf("failed to find song: %w", err)
	}

	return song.Comments, nil
}

// UpdateComment searches for a comment with a matching comment ID
// and updates it to the newComment. Returns true if the comment was updated.
func (m *MongoDBStore) UpdateComment(commentID string, newComment models.UserComments) (bool, error) {
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
func (m *MongoDBStore) GetComment(commentID string) (*models.UserComments, error) {
	db := m.client.Database(DatabaseName)
	collection := db.Collection("songs")

	// Filter to find any document that has a comment with the matching comment ID.
	filter := bson.M{"comments.uuid": commentID}
	// Projection returns only the matching element from the comments array.
	projection := bson.M{"comments.$": 1}
	opts := options.FindOne().SetProjection(projection)

	// Use an inline struct to decode only the comments array.
	var result struct {
		Comments []models.UserComments `bson:"comments"`
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
	return pkg.NewUUID()
}

/*
End of mongoDB implementation
*/
var (
	database DocumentStore
	dbLock   sync.Mutex
)

func CreateNewMongoInstance() DocumentStore {
	// Thread-safe singleton pattern
	dbLock.Lock()
	defer dbLock.Unlock()

	if database == nil {
		fmt.Println("🔄 Initializing MongoDB connection...")
		database = NewDocumentStore()
		fmt.Println("✅ MongoDB connection established")
	}
	return database
}

func SaveUser(user *models.UserMongoDocument) error {
	fmt.Println("Attempting to save user info with name ", user.UserProfileResponse.DisplayName)
	return CreateNewMongoInstance().SaveUser(user)
}

// Assume that anything that has a suffix of id is referring to a uuid, unless otherwise specified
// for now just get working but later on there should be a slight level of
// misdirection so that we can handle errors better
func SubmitComment(songid string, comment models.UserComments) error {
	return CreateNewMongoInstance().SubmitComment(songid, comment)
}
func GetComments(songid string, limit, offset int) ([]models.UserComments, error) {
	return CreateNewMongoInstance().GetComments(songid)
}

// find old and update it with new comment. if old cant be found return error
func UpdateComment(oldCommentID string, new models.UserComments) (bool, error) {
	return CreateNewMongoInstance().UpdateComment(oldCommentID, new)
}

func GetComment(commentID string) (*models.UserComments, error) {
	return CreateNewMongoInstance().GetComment(commentID)
}

func DeleteComment(commentID string) error { return CreateNewMongoInstance().DeleteComment(commentID) }

// if any of these return nil it means it wasnt found
func GetUserDocument(userid string) (*models.UserMongoDocument, error) {
	return CreateNewMongoInstance().GetUserByID(userid)
}
func GetUserSongs(userid string) ([]models.SongTypes, error) {
	return CreateNewMongoInstance().GetUserSongs(userid)
}
func GetUserComments(userid string) ([]models.UserComments, error) {
	return CreateNewMongoInstance().GetUserComments(userid)
}

func AddSongtoDB(songURI string) error {
	return CreateNewMongoInstance().AddSongtoDB(songURI)

}
func ReturnSongbyID(songURI string) (*models.SongTypes, error) {
	return CreateNewMongoInstance().GetSongByID(songURI)
}

func GetUserByID(userUUID string) (*models.UserMongoDocument, error) {
	return CreateNewMongoInstance().GetUserByID(userUUID)
}
