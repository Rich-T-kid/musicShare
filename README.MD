
# Music application for sharing playlist as well as getting songs of the day
---
##activate env -> source env/bin/activate

# 🎵 Music Recommendation App - Core Requirements
> **A smart music recommendation system that helps users discover, rate, and reflect on their music choices!**  

## 🚀 Key Features
### 🔗 Connect to Music Providers
Users can **connect their accounts** to the following platforms:
- 🎧 **Spotify**  
- 🍏 **Apple Music**  
- 🔴 **SoundCloud**  
- 🎦 **YouTube Music**  

---

### 🎶 Song/Album of the Day
- A **song or album of the day** is recommended to users.
- After reading a **brief review**, users decide **whether or not to listen**.
- The app **tracks user engagement** with the song/album.

---

### 📊 Collecting Music Stats
The system gathers key **listening analytics**:
- **⏳ Time Played:** How long the user listens to the song.  
- **⏯ Pauses & Skips:** Number of pauses and skips.  
- **🔁 Replays:** Number of times a song was repeated.  
- **🎧 Total Listening Time:** Overall time spent listening.

```
| Stat Type        | Description                          |
|-----------------|----------------------------------|
| Time Played     | Tracks how long the song was played |
| Pauses         | Counts number of times paused |
| Skips         | Tracks if the user skipped the song |
| Replays         | Number of times the song was replayed |
| Total Time Listened | Measures overall listening time |
```

---

### 📝 User Reviews & Ratings
- Users can **write a review** and provide **ratings**.
- System **remembers the most recent** **positive & negative feedback** for future recommendations.

```
Example Review Format:
- Song Name: "Bohemian Rhapsody"
- Artist: Queen
- Rating: ⭐⭐⭐⭐ (4/5)
- Review: "This track is a masterpiece. The vocals, the instrumentation, and the overall storytelling are just unmatched. It’s a song that never gets old, and every listen feels like an experience. Whether you’re into rock or not, you gotta respect the genius behind it."
```

---

### 🎨 UI Customization Based on Music Provider
Each music provider will have a **unique UI theme**:
```
| Music Provider | UI Theme |
|--------------|-------------------|
| Spotify | 🟢 Green & Black |
| Apple Music | 🎀 Pink & White |
| SoundCloud | 🔴 Red & Black |
| YouTube Music | 🟠 Orange & Black |
```

> **Example:** If a user is connected to **Spotify**, the UI will be **green & black**.

---

### 📈 Feedback-Based Music Discovery
- **Most Positive Feedback Songs:** Users can discover songs that received the highest ratings.
- **Most Negative Feedback Songs:** Users can review songs that weren’t well received.
- System **remembers past feedback** to improve future recommendations.

---

### Music Reflection Feature
- When a song/album is recommended, a **random reflection prompt** appears.
- Users reflect on whether they **enjoyed the music** or not.
- Recommendations are also **genre-based** to match user preferences.

```
> 🎵 *"Did this song make you feel happy, nostalgic, or inspired?"*
>
> 🎸 *"Would you listen to more music from this artist?"*
```

---

### 💈 Data Storage & Recommendation System
- **All user data is securely stored** in a database.
- User feedback & listening habits are **analyzed to improve recommendations**.
- **Genre-based matching** enhances the personalized experience.

---

## features TBD
- streaks
- add friends/connections
- view friends songs, reiews, ect

---
# Routes
## [x] indicates completed routes
### Non Auth routes
    - [] /
        -landing page, sign a user into their music providers 
    - [] /home
        - home page after user logs in
    - [] /song/{userID}
        -  sends song of the day
        - need to pass in user info somehow (tracking purposes)
    - [] /review/{songId}
        - user reviews
        - Get {get a review of the song}
        - Post {Upload a review of the song}
        - CR(UD) 
### Auth routes
    - [] /



---





