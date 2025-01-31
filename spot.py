#this will be for helper functions bewtween
# MongoDB storage (Local to machine) and api request handlers

from pymongo import MongoClient
MONGO_URI = "mongodb://admin:secretpassword@localhost:27017"

# Connect to MongoDB
client = MongoClient(MONGO_URI)

# Select database
db = client["my_database"]

# Select collection
collection = db["users"]

# Insert a test document
collection.insert_one({"name": "Richard", "age": 25})

# Fetch and print all documents
for user in collection.find():
    print(user)
print("done")

"""
user = {
    spotify UserjsonBlob = {
        directly from spotify
    },
    token information = {
        accessToken,
        refreshToken,
        exp,
    },
    // for reccomending songs, what ever song is sent to user is stored here first, can use date to order what was sent when
    // starts empty but first song could be their top played song on spotify
    songInfo = {
        song1 = {songID, date}
        song2 = {songID,date}
    }
    musicInfo = {
        TopItems = {
            // directly from spotify
        },
        followedArtist = {
          // directly from spotify  
        }
    }
}


Functions to check if a user follows a playlist and/or artist
"""