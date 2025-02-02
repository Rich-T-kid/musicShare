from flask import Flask, render_template, request ,jsonify
import urllib.parse
import requests
import json
import random 
import string
import spot

def generate_random_string(length):
    """Generates a random string of letters and digits."""
    letters_and_digits = string.ascii_letters + string.digits
    return ''.join(random.choice(letters_and_digits) for i in range(length))

clientID = "8b277fb167214401bb9486e53d183963"
clientSecret = "db97671791ec461f922c52359d89cddf"
authURL = "https://accounts.spotify.com/authorize"
redirect = "http://localhost:8080/callback"
token_url = "https://accounts.spotify.com/api/token"
scopes = "user-library-read user-modify-playback-state playlist-modify-public playlist-modify-private playlist-read-private user-top-read user-follow-read"
csrf = generate_random_string(16)
# example url = https://accounts.spotify.com/authorize?response_type=code&client_id=8b277fb167214401bb9486e53d183963&scope=user-read-private%20user-read-email&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback&state=2d62cbdee6510f1f

"""
So no code for working with auth is up and working now just write the helper functions for storing
this data to a persistant database and making requests
"""


app = Flask(__name__)

@app.route("/signIn")
def signIn():
    mg = spot.MongoTune()
    

@app.route("/")
def helloworld():
    
    mg = spot.MongoTune()
    return jsonify(mg.all_users())
    return "Hello World!"

@app.route("/home")
def index():
    # Define query parameters
    params = {
        "response_type": "code",
        "client_id": clientID,
        "scope": scopes,
        "redirect_uri": redirect,
        "state": csrf  # Optional for security
    }

    # Encode parameters into URL format
    auth_url = f"{authURL}?{urllib.parse.urlencode(params)}"

    return render_template("index.html", auth_url=auth_url)

@app.route("/callback")
def callback():
    # Get the authorization code from the URL
    auth_code = request.args.get("code") 
    state = request.args.get("state")
    # Verify state to prevent CSRF (should match what was sent in /home)
    if state != csrf:
        return "State verification failed", 400

    # Exchange the authorization code for an access token
    payload = {
        "grant_type": "authorization_code",
        "code": auth_code,
        "redirect_uri": redirect,
        "client_id": clientID,
        "client_secret": clientSecret,
    }

    headers = {"Content-Type": "application/x-www-form-urlencoded"}
    print(payload,headers) 
    response = requests.post(token_url, data=payload, headers=headers)
    
    # Convert response to JSON
    token_data = response.json()
    
    data = {"username": "richard"}
    # database
    # insert into databae with the name of richard for now. 
    """"
    {name:richard, code:code,access:accessToken,refresh:RefreshToken,expires:expires}
    """
    accessToken = token_data["access_token"]
    refreshToken = token_data["refresh_token"]
    expire = token_data["expires_in"] # int
    data = {"ip": request.remote_addr,"accessToken":accessToken,"refreshToken":refreshToken,"expire":expire}
    mg = spot.MongoTune()
    mg.insert_user(data)
    print("all users ->", mg.all_users())
    #storeToken(accessToken,refreshToken,expire)
    return jsonify(token_data)  # Return the access token JSON

#show users info
@app.route("/user/<userID>")
def userinfo(userID):
    pass

# presents song of the day
@app.route("/song/<userID>")
def song(userID):
    pass


# CRUD operations on reivews/comments
@app.route("/review/<songID>")
def review(songID):
    pass


if __name__ == "__main__":
    app.run(port=8080)
    
