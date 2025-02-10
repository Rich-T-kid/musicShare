from pymongo import MongoClient

# Use the new credentials
uri = "mongodb://admin:secretpassword@localhost:27017/?authSource=admin"
client = MongoClient(uri)

try:
    
    client.admin.command('ping')
    print("Pinged your deployment. You successfully connected to MongoDB!")

    
    db = client["musicLibrary"]
    songs = db["songs"]

    
    songs.insert_one({"name": "tyler"})

    print("Current documents in the 'songs' collection:")
    for song in songs.find():
        print(song)

except Exception as e:
    print(f"An error occurred: {e}")