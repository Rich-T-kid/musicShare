import requests
import base64

client_id = "8b277fb167214401bb9486e53d183963"
client_secret = "db97671791ec461f922c52359d89cddf"
redirect = "http://localhost:8080/callback"
SPOTIFY_TOKEN = "https://accounts.spotify.com/api/token"
temp_code = "AQDqVYYbj4ctgsjYC9xz_xxMteHQfmdwAtmrkElm-5m07vJ6Dk0aTRVbl1UIV59TU7GMCdAKzxAJMgxUB9dYmYJ54z0-qJ6y9QwFEwh1J680UhDdxEm45kh28UrWjZQi8o4egcj5AYx2gaMIncFXt08CjtWDWsW90Wn9Go6ufPvD94R-u_-zKwASu5kTZi-Mldsn-gpuaPNzbF89OrbsyVwrK1a20NV0MR-GZ6iMSY25Fsmyz48NytRNPvF7nm30IoX74XZ5qmhsqizDqjozAYGjc2b1RXMCYd3CG2oO8ZIzEkL8S_5Msn-W3wq8jU8A-QR4TK6mVOrb_dU3cHcKWiuDo3PBITQAZgVLWsrAc68y1p9jJucB-Z1Q"
# Encode client_id and client_secret in Base64
credentials = f"{client_id}:{client_secret}".encode("utf-8")
encoded_credentials = base64.b64encode(credentials).decode("utf-8")

print(f"Basic {encoded_credentials}")
# Set request headers
headers = {
    "Authorization": f"Basic {encoded_credentials}",
    "Content-Type": "application/x-www-form-urlencoded"
}

payload = {
        "grant_type": "client_credentials",
    }
print(payload,headers) 
response = requests.post(SPOTIFY_TOKEN, data=payload, headers=headers)
# Send POST request
# Print response
print(response.json())