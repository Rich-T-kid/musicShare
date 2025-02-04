# Spotify API Specification

## Overview
This API provides endpoints for user authentication via Spotify, template rendering for web pages, and user data management. It includes endpoints for checking username availability, handling Spotify login redirects, and retrieving the song of the day.

## Base URL
```
http://localhost:8080
```

## Endpoints

### **1. `/test`**
- **Method:** `GET`
- **Description:** A simple test endpoint to check if the server is running.
- **Response:**
  ```json
  "Hello world"
  ```

---
### **2. `/login`**
- **Method:** `GET`
- **Description:** Serves the `login.html` template for user login.
- **Response:** HTML page.

---
### **3. `/signIn`**
- **Method:** `POST`
- **Description:** Handles user sign-in and adds them to the database.
- **Flow:**
  - Stores user details in the database.
  - Redirects to `/auth` after successful sign-in.
- **Response:** Redirect to `/auth`.

---
### **4. `/link`**
- **Method:** `GET`
- **Description:** Generates a Spotify login URL for user authentication.
- **Headers:**
  - `X-username: {Username}` (Required)
- **Response:** JSON containing the login link.
  ```json
  {
    "link": "https://accounts.spotify.com/authorize?..."
  }
  ```

---
### **5. `/callback`**
- **Method:** `GET`
- **Description:** Handles Spotify login callback and retrieves the access token.
- **Query Parameters:**
  - `code`: Authorization code from Spotify.
  - `state`: Username passed during authorization.
  - `error`: Error message (if any).
- **Flow:**
  - If an error occurs, returns `400 Bad Request`.
  - If successful, retrieves and caches the Spotify token.
  - Fetches the user's top tracks from Spotify.
- **Response:** JSON with Spotify top tracks.

---
### **6. `/loveShare`**
- **Method:** `GET`
- **Description:** Renders `SongofDay.html` template.
- **Response:** HTML page.

---
### **7. `/auth`**
- **Method:** `GET`
- **Description:** Renders `index.html` template for authenticated users.
- **Response:** HTML page.

---
### **8. `/Songs`**
- **Method:** `GET`
- **Description:** Redirects to `index.html`, serves as an API placeholder.
- **Response:** HTML page.

---
### **9. `/exist/{name}`**
- **Method:** `GET`
- **Description:** Checks whether a username is available.
- **Path Parameters:**
  - `{name}`: The username to check.
- **Flow:**
  - If username exists, returns `409 Conflict`.
  - If username is available, returns `200 OK`.
- **Response:**
  ```json
  {
    "message": "Username {name} is available"
  }
  ```
  or
  ```json
  {
    "message": "Username {name} already exists. Choose another one"
  }
  ```

---
### **10. `/refresh`**
- **Method:** `POST`
- **Description:** Refreshes the Spotify access token using the refresh token.
- **Request Body:**
  ```json
  {
    "refresh_token": "your_refresh_token"
  }
  ```
- **Response:**
  ```json
  {
    "access_token": "new_access_token",
    "expires_in": 3600
  }
  ```

## Code Flow Summary

1. **User Authentication Flow:**
   - User visits `/login`.
   - User clicks login with Spotify → Redirected to `/link`.
   - Spotify redirects back to `/callback` with `code`.
   - Access token is retrieved and cached.
   - User is added to the database if new.
   - Redirects to `/auth` for the home page.

2. **Username Availability Flow:**
   - Client requests `/exist/{name}`.
   - Server checks if the username exists in cache.
   - If username exists → `409 Conflict`.
   - If username is available → `200 OK`.

3. **Token Refresh Flow:**
   - Client sends a `POST` request to `/refresh` with a refresh token.
   - Server returns a new access token with expiration time.

This API is structured for seamless Spotify authentication and user management, ensuring a smooth login experience while maintaining efficient username availability checking.