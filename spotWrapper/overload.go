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

	"loveShare/logs"
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

var (
	logger = logs.NewLogger() //TODO: Implementing logging in here everywhere
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

func (o *Overload) RetryRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	delay := o.defaultDelay
	userid, ok := ctx.Value(UsernameKey{}).(string)
	if !ok {
		logger.Critical("UserName was not properly set in the context.") 
		return nil, fmt.Errorf("UserName was not properly set in the context.context \n")
	}

	// returns true if access Token is valid. if false generate new one. if error occures as in cant generate a valid access token then return error
	// this is so that we dont have to deal with handling a 401 response code, if the access token is valid here it will remain valid for the duration of the retrys below
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		logger.Critical(fmt.Sprintf("No Authorization header found for user: %s", userid)) // 
		return nil, errors.New("forgot to attach Authorization header before making Spotify request")
	}

	// 🔥 Extract Access Token
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")
	if exist := validateAccessToken(accessToken); !exist {
		// Get refresh token from cache
		_, refreshToken, err := tokenStore.GetTokens(userid)
		if err != nil || refreshToken == "" {
			logger.Critical(fmt.Sprintf("Refresh token not found or is empty for user: %s, err: %v", userid, err)) // <-- added logging
			return nil, errors.New("refresh token not found, user must re-authenticate")
		}

		validAccessToken := generateAccessToken(refreshToken)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validAccessToken))

		err = tokenStore.StoreTokens(userid, validAccessToken, refreshToken)
		if err != nil {
			logger.Critical(fmt.Sprintf("Failed to store new tokens in cache for user %s, err: %v", userid, err)) // <-- added logging
			return nil, fmt.Errorf("failed to store new tokens in cache: %w", err)
		}
	}

	for attempt := 1; attempt <= o.maxTries; attempt++ {
		// Clone the request to ensure all retries send the exact same data
		clonedReq, err := cloneRequest(req)
		if err != nil {
			logger.Critical(fmt.Sprintf("Attempt %d for user %s: error cloning request - %v", attempt, userid, err)) // <-- added logging
			return nil, fmt.Errorf("error cloning request: %w", err)
		}

		// Send HTTP request
		resp, err := o.client.Do(clonedReq)
		if err != nil {
			logger.Warning(fmt.Sprintf("Attempt %d for user %s: Request failed - %v", attempt, userid, err)) // <-- added logging
		} else {
			// Handle rate limits (429) and server errors (500, 503, 408)
			if resp.StatusCode == 429 || resp.StatusCode == 500 || resp.StatusCode == 503 || resp.StatusCode == 408 {
				logger.Warning(fmt.Sprintf(
					"Attempt %d for user %s: Received %d response - applying backoff",
					attempt, userid, resp.StatusCode,
				)) 

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
				logger.Info(fmt.Sprintf("User %s did not have a valid access token. Attempting to obtain a new one.", userid))
				fmt.Println("Write to log file as well as getting new refresh token")
			} else {
				// 400, 403, 404 -> retrying won't help. Break out early & return the error
				bodyBytes, readErr := io.ReadAll(resp.Body)
				if readErr != nil {
					logger.Critical(fmt.Sprintf("Failed to read response body on attempt %d for user %s: %v", attempt, userid, readErr)) // <-- added logging
					fmt.Println(readErr)
				}
				bodyString := string(bodyBytes)
				logger.Route(fmt.Sprintf(
					"Request to %s for user %s returned a %d StatusCode with a body of: %s",
					req.URL, userid, resp.StatusCode, bodyString,
				)) // <-- added logging

				resp.Body.Close()
				return resp, unexpectedStatus
			}
		}

		// Apply exponential backoff with jitter (randomized delay)
		jitter := rand.Intn(delay) // Add randomness to prevent bursts
		sleepDuration := time.Duration(delay+jitter) * time.Second

		logger.Info(fmt.Sprintf("Attempt %d for user %s failed, retrying in %v...", attempt, userid, sleepDuration)) // <-- added logging

		select {
		case <-time.After(sleepDuration): // Wait before retrying
		case <-ctx.Done():
			// Clean up resources
			logger.Warning(fmt.Sprintf("Context canceled for user %s on attempt %d, stopping retries.", userid, attempt)) // <-- added logging
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
			return nil, ctx.Err()
		}

		// Double the delay for exponential backoff (max cap at 20s)
		delay = min(delay*2, 20)
	}

	logger.Critical(fmt.Sprintf("Max retry attempts (%d) reached for user %s, request failed.", o.maxTries, userid)) // <-- added logging
	return nil, failedRetry
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

/* retrys request up to {maxTries} times and returns a copy of the body in bytes*/