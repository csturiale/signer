package signer

import (
	"fmt"
	goalone "github.com/bwmarrin/go-alone"
	"strings"
	"time"
)

// Signature is the type for the package. Secret is the signer secret, a lengthy
// and hard to guess string we use to sign things.
type Signature struct {
	Secret string
}

// GenerateTokenFromString generates a signed token and returns it
func (s *Signature) GenerateTokenFromString(data string) string {
	var urlToSign string

	pen := goalone.New([]byte(s.Secret), goalone.Timestamp)

	if strings.Contains(data, "?") {
		// handle case where URL contains query parameters
		urlToSign = fmt.Sprintf("%s&hash=", data)
	} else {
		// no query parameters
		urlToSign = fmt.Sprintf("%s?hash=", data)
	}

	tokenBytes := pen.Sign([]byte(urlToSign))
	token := string(tokenBytes)

	return token
}

// VerifyToken verifies a signed token and returns true if it is valid,
// false if it is not.
func (s *Signature) VerifyToken(token string) bool {
	pen := goalone.New([]byte(s.Secret), goalone.Timestamp)
	_, err := pen.Unsign([]byte(token))

	if err != nil {
		// signature is not valid. Token was tampered with, forged, or maybe it's
		// not even a token at all! Either way, it's not safe to use it.
		return false
	}

	// valid hash
	return true

}

// Expired checks to see if a token has expired. It returns true if
// the token was created within minutesUntilExpire, and false otherwise.
func (s *Signature) Expired(token string, minutesUntilExpire int) bool {
	pen := goalone.New([]byte(s.Secret), goalone.Timestamp)
	ts := pen.Parse([]byte(token))

	return time.Since(ts.Timestamp) > time.Duration(minutesUntilExpire)*time.Minute
}