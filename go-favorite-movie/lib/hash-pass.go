package myPac

import (
	"crypto/sha256"
	"encoding/base64"
)

func HashPass(pass string) string {
	return base64.StdEncoding.EncodeToString(sha256.New().Sum([]byte(pass)))
}
