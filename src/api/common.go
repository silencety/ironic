package api

import (
	"mime"
	"src/github.com/Sirupsen/logrus"
)

const defaultApiVersion = "1.0"



// MatchesContentType validates the content type against the expected one
func MatchesContentType(contentType, expectedType string) bool {
	mimetype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		logrus.Errorf("Error parsing media type: %s error: %v", contentType, err)
	}
	return err == nil && mimetype == expectedType
}


func Version() string{
	return defaultApiVersion
}

