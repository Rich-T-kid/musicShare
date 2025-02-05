package spotwrapper

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
	implent overloader struct

context will contain both the tokens structs as well as
a overloader. overloaders will manaege cool downs for request and 429 responses as well as 503 server down

Caller should not call this for status codes that are not : 500 , 503, 429 mabey 408
Note: retry logic should target server-side errors like 500 (Internal Server Error) or 503 (Service Unavailable) check the retry after header , as they often indicate temporary issues. Client-side errors like 404 (Not Found) or 400 (Bad Request) usually signify issues that retries won't resolve.
*/
// Error messages
var (
	failedRetry      = errors.New("http request retrys were unsuccessful")
	serverDown       = errors.New("target servers are down, cannot complete http request")
	unexpectedStatus = errors.New("received an unexpected HTTP status code")
)

/*
Implementation of exponetial jitter retrys
*/

type Overload struct {
	client *http.Client // performs request
	//mu sync.Mutex // assuming this will run concurrently for now
	defaultDelay int // starts off a 1 second for delays and each time it doubles
	maxTries     int // max # of attempts to retry request
}

func newOverloader() *Overload {
	c := &http.Client{}
	return &Overload{
		client:       c,
		defaultDelay: 1,
		maxTries:     5,
	}
}

// cloneRequest creates a deep copy of the original request
func cloneRequest(req *http.Request) (*http.Request, error) {
	var bodyCopy []byte
	if req.Body != nil {
		var err error
		bodyCopy, err = io.ReadAll(req.Body) // Read and copy request body
		if err != nil {
			return nil, fmt.Errorf("error reading request body: %w", err)
		}
		req.Body = io.NopCloser(bytes.NewReader(bodyCopy)) // Restore original body
	}

	// Create new request with same method, URL, and body
	clonedReq, err := http.NewRequest(req.Method, req.URL.String(), bytes.NewReader(bodyCopy))
	if err != nil {
		return nil, fmt.Errorf("error cloning request: %w", err)
	}

	// Copy headers
	for key, values := range req.Header {
		for _, value := range values {
			clonedReq.Header.Add(key, value)
		}
	}

	return clonedReq, nil
}

/* retrys request up to {maxTries} times and returns a copy of the body in bytes*/
func (o *Overload) RetryRequest(ctx context.Context, req *http.Request, userid string) (*http.Response, error) {
	delay := o.defaultDelay
	// returns true if access Token is valid. if false generate new one. if error occures as in cant generate a valid access token then return error
	// this is so that we dont have to deal with handling a 401 response code, if the access token is valid here it will remain valid for the duration of the retrys below
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("forgot to attach Authorization header before making Spotify request")
	}

	// 🔥 Extract Access Token
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")
	if exist := validateAccessToken(accessToken); !exist {
		// Get refresh token from cache
		_, refreshToken, err := tokenStore.GetTokens(userid)
		if err != nil || refreshToken == "" {
			return nil, errors.New("refresh token not found, user must re-authenticate")
		}
		validAccessToken := generateAccessToken(refreshToken)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validAccessToken))
		err = tokenStore.StoreTokens(userid, validAccessToken, refreshToken)
		if err != nil {
			return nil, fmt.Errorf("failed to store new tokens in Redis: %w", err)
		}
		// make request

		// if accessToken doesnt exist check if refresh token exist and is associated with it (old access token) if there is one
		// get a new access token using that refresh token and return it here. save to cache before returning here
		// if a refresh token is not associated with this access token return nil, an appropate error
	}
	// 🔥 Validate Authorization Header

	for attempt := 1; attempt <= o.maxTries; attempt++ {
		// Clone the request to ensure all retries send the exact same data
		clonedReq, err := cloneRequest(req)
		if err != nil {
			return nil, fmt.Errorf("error cloning request: %w", err)
		}

		// Send HTTP request
		resp, err := o.client.Do(clonedReq)
		if err != nil {
			fmt.Printf("Attempt %d: Request failed - %v\n", attempt, err)
		} else {
			// Handle rate limits (429) and server errors (500, 503, 408)
			if resp.StatusCode == 429 || resp.StatusCode == 500 || resp.StatusCode == 503 || resp.StatusCode == 408 {
				fmt.Printf("Attempt %d: Received %d response\n", attempt, resp.StatusCode)

				// If a `Retry-After` header is present, respect it
				if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
					parsedDelay, err := strconv.Atoi(retryAfter)
					if err == nil {
						delay = parsedDelay // Override delay with `Retry-After`
					}
				}

				// Close response body before retrying
				resp.Body.Close()
			} else if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
				return resp, nil // ✅ Success, return response
			} else if resp.StatusCode == 401 {

			} else { // 400 , 403 , 404 , retrying wont affect these resposne codes so just break out early and return the error
				// clean up resources
				resp.Body.Close()
				return resp, unexpectedStatus
			}
		}

		// Apply exponential backoff with jitter (randomized delay)
		jitter := rand.Intn(delay) // Add randomness to prevent bursts
		sleepDuration := time.Duration(delay+jitter) * time.Second
		fmt.Printf("Retrying in %v seconds...\n", sleepDuration)

		select {
		case <-time.After(sleepDuration): // Wait before retrying
		case <-ctx.Done():
			//clean up resources
			resp.Body.Close()
			fmt.Println("Context canceled, stopping retries")
			return nil, ctx.Err() // If context is canceled, return immediately
		}

		// Double the delay for exponential backoff (max cap at 20s)
		delay = min(delay*2, 20)
	}

	fmt.Println("Max retry attempts reached, request failed")
	return nil, failedRetry
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
