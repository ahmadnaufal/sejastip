package util

import (
	"encoding/base64"
	"log"
	"mime"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// DecodeUploadedBase64File returns one possibility of a file extension
//  and its decoded contents, given a base64-encoded file in string
func DecodeUploadedBase64File(encodedString string) ([]byte, string, error) {
	// first, check if the encodedString is a valid representation of a
	// base64-encoded file
	match, _ := regexp.MatchString("^data:([a-zA-Z0-9]+\\/[a-zA-Z0-9-.+]+).*,.*$", encodedString)
	if !match {
		return nil, "", errors.New("Unknown file uploaded")
	}

	// then, get the mimetype
	mimeType := encodedString[5:strings.Index(encodedString, ";")]
	if len(mimeType) < 1 {
		return nil, "", errors.New("Unknown mime type detected")
	}

	// get the possible extension of the mimeType
	extensions, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return nil, "", err
	}
	if len(extensions) == 0 {
		return nil, "", errors.New("Extension for the mime type not found")
	}

	// decode the base64 to get the file contents
	encodedContent := encodedString[strings.Index(encodedString, ",")+1:]
	log.Println(extensions, mimeType, encodedContent)
	decoded, err := base64.StdEncoding.DecodeString(encodedContent)
	if err != nil {
		return nil, "", err
	}

	return decoded, extensions[0], nil
}
