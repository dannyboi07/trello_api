package util

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"trelloBE/common"
	"trelloBE/schema"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type jwtCustomClaims struct {
	*jwt.RegisteredClaims
	schema.UserForToken
}

func CreateAccessToken(userDetails schema.UserForToken) (string, int, error) {
	var token *jwt.Token = jwt.New(jwt.GetSigningMethod("RS256"))

	var createdTime time.Time = time.Now()
	var expireAtTime time.Time = createdTime.Add(time.Minute * 15)

	token.Claims = jwtCustomClaims{
		&jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAtTime),
		},
		userDetails,
	}

	signedToken, err := token.SignedString(common.PrivateKey)
	if err != nil {
		return "", 0, err
	}

	return signedToken, int(expireAtTime.Sub(createdTime).Seconds()), nil
}

func CreateRefreshToken(userDetails schema.UserForToken) (string, int, error) {
	var token *jwt.Token = jwt.New(jwt.GetSigningMethod("RS256"))

	var createdTime time.Time = time.Now()
	var expireAtTime time.Time = createdTime.AddDate(0, 0, 7)

	token.Claims = jwtCustomClaims{
		&jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAtTime),
		},
		userDetails,
	}

	signedToken, err := token.SignedString(common.PrivateKey)
	if err != nil {
		return "", 0, err
	}

	return signedToken, int(expireAtTime.Sub(createdTime).Seconds()), nil
}

func VerifyJwtToken(token string) (jwt.MapClaims, int, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if tokenAlg := t.Method.Alg(); tokenAlg != "RS256" {
			return nil, fmt.Errorf("Unexpected signing method: %s", tokenAlg)
		}

		return common.PublicKey, nil
	})

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if parsedToken.Valid {
		return parsedToken.Claims.(jwt.MapClaims), http.StatusOK, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, http.StatusBadRequest, errors.New("Malformed token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, http.StatusUnauthorized, errors.New("Token expired")
		}
		// else {
		// 	Log.Println("Couldn't handle this token:", err)
		// 	return nil, http.StatusInternalServerError, errors.New("Sorry, something went wrong")
		// }
	}

	Log.Println("Couldn't handle this token:", err)
	return nil, http.StatusInternalServerError, errors.New("Sorry, something went wrong")
}

func ParseJwtClaims(jwtClaims jwt.MapClaims) (map[string]interface{}, int, error) {
	var embeddedDetails map[string]interface{} = make(map[string]interface{})

	if userIdString, ok := jwtClaims["id"].(string); ok {
		userId, err := primitive.ObjectIDFromHex(userIdString)
		if err != nil {
			return nil, http.StatusUnauthorized, fmt.Errorf("Malformed token")
		}
		embeddedDetails["id"] = userId
	} else {
		return nil, http.StatusUnauthorized, fmt.Errorf("Malformed token")
	}

	if userEmail, ok := jwtClaims["email"].(string); ok {
		embeddedDetails["email"] = userEmail
	} else {
		return nil, http.StatusUnauthorized, fmt.Errorf("Malformed token")
	}

	return embeddedDetails, 0, nil
}
