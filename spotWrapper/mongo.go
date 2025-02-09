package spotwrapper

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// cache interface defining methods for caching
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
		fmt.Printf("Could not convert key %s of value %s to type %v Value of reddis implentation cache", key, str)
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
type DocumentStore interface {
	// Connectivity Check
	Connected() bool // Pings the database and returns whether it's connected.
	UserStore
	SongStore
	CommentStore
	// CRUD Operations
}
type UserStore interface {
	GetUserByID(userID string) (*UserMongoDocument, error)
	SaveUser(user *UserMongoDocument) error
	GetUserSongs(userID string) (*SongTypes, error)
	GetUserComments(userID string) ([]UserComments, error)
}

// TODO: Change this after merge
type SongStore interface {
	GetSongByID(songID string) (*SongTypes, error)
	InsertSong(song *SongTypes) error
	UpdateSong(song *SongTypes) error
	DeleteSong(songID string) error
}
type CommentStore interface {
	SubmitComment(songID spotifyURI, comment UserComments) error
	GetComments(songID spotifyURI, limit, offset int) ([]UserComments, error)
	UpdateComment(oldComment string, newComment UserComments) (bool, error)
	DeleteComment(commentID string) error
	GetComment(commentID string) (*UserComments, error)
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
	return &MongoDBStore{client: client}
}

type MongoDBStore struct {
	client *mongo.Client
}

func (m *MongoDBStore) Connected() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := m.client.Ping(ctx, nil)
	if err != nil {
		return false
	}
	return true
}

// Implement UserStore methods
func (m *MongoDBStore) GetUserByID(userID string) (*UserMongoDocument, error) { return nil, nil }
func (m *MongoDBStore) SaveUser(user *UserMongoDocument) error                { return nil /* implementation */ }
func (m *MongoDBStore) GetUserSongs(userID string) (*SongTypes, error) {
	return nil, nil /* implementation */
}
func (m *MongoDBStore) GetUserComments(userID string) ([]UserComments, error) {
	return nil, nil /* implementation */
}

// Implement SongStore methods
func (m *MongoDBStore) GetSongByID(songID string) (*SongTypes, error) {
	return nil, nil /* implementation */
}
func (m *MongoDBStore) InsertSong(song *SongTypes) error { return nil /* implementation */ }
func (m *MongoDBStore) UpdateSong(song *SongTypes) error { return nil /* implementation */ }
func (m *MongoDBStore) DeleteSong(songID string) error   { return nil /* implementation */ }

// Implement CommentStore methods
func (m *MongoDBStore) SubmitComment(songID spotifyURI, comment UserComments) error {
	return nil /* implementation */
}
func (m *MongoDBStore) GetComments(songID spotifyURI, limit, offset int) ([]UserComments, error) {
	return nil, nil /* implementation */
}
func (m *MongoDBStore) UpdateComment(oldComment string, newComment UserComments) (bool, error) {
	return false, nil /* implementation */
}
func (m *MongoDBStore) DeleteComment(commentID string) error { return nil /* implementation */ }
func (m *MongoDBStore) GetComment(commentID string) (*UserComments, error) {
	return nil, nil /* implementation */
}

/*
End of mongoDB implementation
*/
func SaveUser() {
	fmt.Println("Prentending to save user info ->  implement later")
}

func SubmitComment(songid spotifyURI, comment UserComments) error { return nil }
func GetComments(songid spotifyURI, limit, offset int) []UserComments {
	var empty []UserComments
	return empty
}

// find old and update it with new comment. if old cant be found return error
func UpdateComment(oldComment string, new UserComments) (bool, error) {
	return true, nil
}

func GetComment(commentID string) *UserComments {
	return nil
}

func DeleteComment(commentID string) error { return nil }

type SongTypes struct {
	AllSongs      []spotifyURI `json:"songs"`          // all songs
	LikedSongs    []spotifyURI `json:"liked_songs"`    /// positively rated songs 3/5 <= out of 5
	DislikedSongs []spotifyURI `json:"disliked_songs"` /// less than 2/5 out of  5
}

// if any of these return nil it means it wasnt found
func GetUserDocument(userid string) (*UserMongoDocument, error) { return nil, nil }
func GetUserSongs(userid string) (*SongTypes, error)            { return nil, nil }
func GetUserComments(userid string) ([]UserComments, error)     { return nil, nil }
