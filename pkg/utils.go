package numeral

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// GetSetting checks for environment variables in system
func GetSetting(setting string) (string, error) {
	value, ok := os.LookupEnv(setting)
	if !ok {
		err := fmt.Errorf("init error, %s ENV var not found", setting)
		return "", err
	}
	return value, nil
}

// AuthMiddleware is the authentication middleware
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the username and password from the request Authorization header.
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte("usertest")) // TODO env var
			expectedPasswordHash := sha256.Sum256([]byte("userpass")) // TODO env var

			// Use the subtle.ConstantTimeCompare() function to check if
			// the provided username and password hashes equal the
			// expected username and password hashes.
			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			// If the username and password are correct, call next handler
			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		// If the Authentication header is not present or invalid, set a WWW-Authenticate
		// header to inform the client and send a 401 Unauthorized response.
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		responseUnauthorized(w)
	})
}

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// Decode decodes body request to an structure represented by the input type
func Decode(r *http.Request, input interface{}) (interface{}, error) {

	switch input.(type) {
	case Payment:
		output := Payment{}
		if err := json.NewDecoder(r.Body).Decode(&output); err != nil {
			return nil, err
		}
		defer r.Body.Close()
		return output, nil
	}
	return nil, nil
}

// GetPaymentFromRequest maps the output from decode method to an User entity
func GetPaymentFromRequest(req *http.Request) (*Payment, error) {
	payment := Payment{}
	output, err := Decode(req, payment)
	if err != nil {
		return nil, fmt.Errorf("malformed input data")
	}
	payment = output.(Payment)
	log.Println(payment)
	return &payment, nil
}
