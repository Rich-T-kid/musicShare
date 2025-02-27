import grpc
from protobuff import song_pb2
from protobuff import song_pb2_grpc
from concurrent import futures
import random
from pymongo import MongoClient
import time
# MongoDB connection details <- old
#MONGO_HOST = "localhost"
#MONGO_PORT = 27017
#MONGO_USER = "admin"
#MONGO_PASS = "secretpassword"
#MONGO_DB = "test_db"

# MongoDB URI with authentication
#MONGO_URI = f"mongodb://{MONGO_USER}:{MONGO_PASS}@{MONGO_HOST}:{MONGO_PORT}/"

MONGO_HOST = "cluster0.avlxk.mongodb.net"
MONGO_USER = "rbb98"
MONGO_PASS = "cfxARjWMSnojKSjj"
MONGO_DB = "test_db"

# MongoDB URI with authentication (Atlas Cloud Connection)
MONGO_URI = f"mongodb+srv://{MONGO_USER}:{MONGO_PASS}@{MONGO_HOST}/{MONGO_DB}?retryWrites=true&w=majority"

def userCollection():
    client = MongoClient(MONGO_URI)
    db = client[MONGO_DB]
    collection = db["users"]
    return collection 

def get_random_songs(uuid):
    """Fetch user data and return 5 random songs from TopSinglesTracks."""
    try:
        collection = userCollection()
        user = collection.find_one({"uuid": uuid})
        if not user:
            return []

        user_music_info = user.get("user_music_info", {})
        top_tracks = user_music_info.get("toptracks", {})
        top_singles = top_tracks.get("topsingles", [])

        if not isinstance(top_singles, list) or not top_singles:
            return []

        return random.sample(top_singles, min(10, len(top_singles)))

    except Exception as e:
        print(f"MongoDB query failed: {str(e)}")
        return []
test_song = {
    "name": "Blinding Lights//THis is dummy Data implement fr later",
    "artist": "The Weeknd",
    "song_uri": "spotify:track:0VjIjW4GlUZAMYd2vXMi3b",
    "rank": 1  # Will be randomized in your server logic
}

class SongService(song_pb2_grpc.SongServiceServicer):
    def GetSong(self, request, context):
        uuid = request.user_id
        recommended_songs = get_random_songs(uuid)
        print(recommended_songs)
        n = len(recommended_songs)

        song_list = [
            song_pb2.songBody(
                name=song["trackname"],
                artist=song["artist"],
                song_uri=song["tracklink"],
                rank=random.randint(1, n)  # Assign a random rank between 1 and n
            ) for song in recommended_songs
        ]
        print(song_list)
        return song_pb2.SongResponse(songs=song_list)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    song_pb2_grpc.add_SongServiceServicer_to_server(SongService(), server)

    port = 9000  # Running on port 9000
    server.add_insecure_port(f"[::]:{port}")
    print(f"gRPC server is running on port {port}...")

    server.start()
    server.wait_for_termination()

if __name__ == "__main__":
    serve()
    """ sv = get_random_songs("8eeb1587-0d15-4602-9287-d3dd4cb48631")
    for song in sv:
        print(song["trackname"],song["artist"],song["tracklink"])
        
    print("done")"""