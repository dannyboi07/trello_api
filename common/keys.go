package common

import (
	"crypto/rsa"
	"os"
)

var PrivateKey *rsa.PrivateKey
var PublicKey *rsa.PublicKey

var AccessTokenPath string
var RefreshTokenPath string

func Init() {
	AccessTokenPath = "/trelloapi/" + os.Getenv("API_VERSION")
	RefreshTokenPath = "/trelloapi/" + os.Getenv("API_VERSION") + "/auth"
}
