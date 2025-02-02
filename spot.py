#this will be for helper functions bewtween
# MongoDB storage (Local to machine) and api request handlers

from pymongo import MongoClient
import requests
import base64
MONGO_URI = "mongodb://admin:secretpassword@localhost:27017"


baseEncoded = "OGIyNzdmYjE2NzIxNDQwMWJiOTQ4NmU1M2QxODM5NjM6ZGI5NzY3MTc5MWVjNDYxZjkyMmM1MjM1OWQ4OWNkZGY="
clientID = "8b277fb167214401bb9486e53d183963"
clientSecret = "db97671791ec461f922c52359d89cddf"
accessToken = "BQBd31LsGuVq69jeKh279rc1JsQNb7WlvUHNnoQQ5fBNrRB3WrjGf_phCA993sZtR-ba3QB9aJovlKMhSlDE42Q8zUCZOpRHBkH2rPcYmNbbJmZcc-_NiMoQquy9w98mh-sgXlfhge1fBSrE5ZyPCKtYQ71RvewYLhfaBd6B6s3cQRkLk59YSnLLyWZQC3-P4pYnRCveglw2g7IG1FsqiTpRZS8d_sG20Eg2oXCLzTaZFravU_BVI-aIJIjXgLfcZWaWaLbGiGgs40Bk4g6f0Jp0LudEq8Mlj9PGTA3Lje6b1uHtsGMFuSpHBhytIpqq2SID3PHnClOYqEU6QA"
refreshToken = "AQDhizdD4iolQSHFRjsdWXbCJmVZlG2TWSx_SoVjuVqUQQyl85wh52Adbm0nc4ihdupD2ioickQnl2F8OTpJt_h_D552a4GSoAD_53JxO-tfPv2EmKaOMIunTI4Exww-Wko"
class Overloader:
    """
    A Singleton class to handle '429 Too Many Requests' rate-limiting from the Spotify API.
    You can expand this class to keep track of rate-limits, 
    sleep for a certain time, or queue requests as needed.
    """
    _instance = None

    def __new__(cls):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
            # Initialize any additional properties here
        return cls._instance

    def check_overload(self):
        """
        Check if we are currently overloaded (e.g., recently got a 429).
        Return True if we need to wait, False otherwise.
        """
        # TODO: implement your logic
        return False
# Select collection
# Insert a test document
class Spotifyhelper:
    def __init__(self,access_Token:str,refresh_Token:str,encoded_creds:str):
        self.base = "https://api.spotify.com/v1/" 
        self.token = access_Token
        self.refresh = refresh_Token
        self.header = {
            "Authorization" : f"Bearer {self.token}"
        }
        self.Overload = Overloader() # all instances share this overloader class. If we ever get a 429 its managed here
        self.encoded_creds = encoded_creds
    def Refresh(self,user_id):
        spotify_token_url = "https://accounts.spotify.com/api/token"
        headers = {
            "Content-Type" : "application/x-www-form-urlencoded",
            "Authorization": f"Basic {self.encoded_creds}"
        }        
        body = {
            "grant_type" : "refresh_token",
            "refresh_token": self.refresh
        }
        tokenurl = "https://accounts.spotify.com/api/token"
        response = requests.post(tokenurl,data=body,headers=headers)
        if response.status_code == 200:
            token_data = response.json()
            access_token = token_data.get("access_token")

            # Store the updated token in MongoDB
            mongo = MongoTune()
            mongo.update_access_token(user_id, access_token)

            return access_token
        else:
            return None
class MongoTune:
    """
    A simple MongoDB wrapper to handle user data.
    """
    def __init__(self):
        self.client = MongoClient(MONGO_URI)
        self.db = self.client["my_database"]
        self.collection = self.db["users"]

    def clear_users(self):
        """
        Deletes all documents in the 'users' collection, effectively clearing the users database.
        """
        result = self.collection.delete_many({})  # Deletes all documents
        print(f"Deleted {result.deleted_count} users from the database.")
    def insert_user(self, user_data: dict) -> bool:
        """
        Insert a user record into the 'users' collection.
        Returns True on success, False otherwise.
        """
        try:
            self.collection.insert_one(user_data)
            return True
        except Exception as e:
            print(f"Error inserting user: {e}")
            return False

    def all_users(self) -> list:
        """
        Retrieve all user documents from the 'users' collection.
        """
        return list(self.collection.find())

    def find_user(self, query: dict) -> dict:
        """
        Find a single user matching the provided query.
        """
        return self.collection.find_one(query)
    def update_access_token(self, user_id: str, access_token: str):
        """
        Updates the access token (and optionally the refresh token) for a given user.
        """
        update_data = {
            "$set": {
                "spotify.access_token": access_token,
            }
        }

        result = self.collection.update_one({"user_id": user_id}, update_data)

        if result.modified_count > 0:
            return True
        else:
            return False
# manages 429 responses 
spot = Spotifyhelper(accessToken,refreshToken,baseEncoded)
mongo = MongoTune()
#print(mongo.all_users())
print(spot.Refresh("richard"))
print(mongo.all_users())

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