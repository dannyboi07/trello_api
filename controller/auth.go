package controller

import (
	"encoding/json"
	"net/http"
	"trelloBE/common"
	"trelloBE/db"
	"trelloBE/schema"
	"trelloBE/util"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// /auth/register
func Register(w http.ResponseWriter, r *http.Request) {
	var jDec *json.Decoder = json.NewDecoder(r.Body)
	jDec.DisallowUnknownFields()

	var (
		userRegister schema.UserRegister
		statusCode   int
		err          error
	)
	statusCode, err = util.JsonParseErr(jDec.Decode(&userRegister))

	if err != nil {
		util.WriteApiMessage(w, err.Error(), statusCode, false)
		// http.Error(w, err.Error(), statusCode)
		util.Log.Println("Failed to decode request, err: ", err)
		return
	}
	statusCode, err = userRegister.Validate()
	if err != nil {
		util.WriteApiMessage(w, err.Error(), statusCode, false)
		// http.Error(w, err.Error(), statusCode)
		util.Log.Println("User register request validation err: ", err)
		return
	}

	var userHashedPw string
	userHashedPw, err = util.HashPassword(*userRegister.Password)
	if err != nil {
		util.WriteApiMessage(w, "", 0, false)
		util.Log.Println("Failed to hash user password, err: ", err)
		return
	}

	err = db.InsertUser(userRegister, userHashedPw)
	if err != nil {
		statusCode, message := 0, ""
		if db.IsDuplicateKeyError(err) {
			statusCode, message = http.StatusBadRequest, "Account already exists"
			util.Log.Println("Registration attempt with existing email:", *userRegister.Email, r.RemoteAddr)
		} else {
			util.Log.Println("Failed to register/insert user details into db, err:", err)
		}
		util.WriteApiMessage(w, message, statusCode, false)
		return
	}

	util.WriteApiMessage(w, "Your account has been created!", 0, true)
}

// /auth/login
func Login(w http.ResponseWriter, r *http.Request) {
	var jDec *json.Decoder = json.NewDecoder(r.Body)
	jDec.DisallowUnknownFields()

	var (
		userLogin  schema.UserLogin
		statusCode int
		err        error
	)
	statusCode, err = util.JsonParseErr(jDec.Decode(&userLogin))
	if err != nil {
		util.WriteApiMessage(w, err.Error(), statusCode, false)
		util.Log.Println("Failed to decode request, err:", err)
		return
	}

	statusCode, err = userLogin.Validate()
	if err != nil {
		util.WriteApiMessage(w, err.Error(), statusCode, false)
		util.Log.Println("User login validation err:", err)
		return
	}

	user, found, err := db.SelectUserByEmail(*userLogin.Email)
	if err != nil {
		util.WriteApiMessage(w, "", 0, false)
		util.Log.Println("Failed to get user details for login, err:", err)
		return
	} else if !found {
		util.WriteApiMessage(w, "Check email/password & try again", http.StatusBadRequest, false)
		// util.Log.Println("Login attempt, user's email not found", r.RemoteAddr)
		return
	}

	var isCorrectPw bool = util.VerifyPassword(user.PasswordHash, *userLogin.Password)
	if !isCorrectPw {
		util.WriteApiMessage(w, "Check email/password & try again", http.StatusUnauthorized, false)
		// util.Log.Println("Login attempt, wrong password for email:", userLogin.Email, r.RemoteAddr)
		return
	}

	var userTokenDetails schema.UserForToken = schema.UserForToken{
		Id:    user.Id,
		Email: user.Email,
	}
	accessToken, accessTokenExpiresIn, err := util.CreateAccessToken(userTokenDetails)
	if err != nil {
		util.WriteApiMessage(w, "Sorry, something went wrong", 0, false)
		util.Log.Println("Failed to create access token")
		return
	}
	refreshToken, refreshTokenExpiresIn, err := util.CreateRefreshToken(userTokenDetails)
	if err != nil {
		util.WriteApiMessage(w, "Sorry, something went wrong", 0, false)
		util.Log.Println("Failed to create refresh token")
		return
	}

	var (
		accessTokenCookie  *http.Cookie
		refreshTokenCookie *http.Cookie
	)
	accessTokenCookie = &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		MaxAge:   int(accessTokenExpiresIn),
		Path:     common.AccessTokenPath,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}
	refreshTokenCookie = &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		MaxAge:   int(refreshTokenExpiresIn),
		Path:     common.RefreshTokenPath,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}

	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokenCookie)

	util.WriteApiMessage(w, "Logged in", 0, true)
}

// /auth/refresh
func RefreshAccessToken(w http.ResponseWriter, r *http.Request) {

	var (
		refreshTokenCookie *http.Cookie
		err                error
	)
	refreshTokenCookie, err = r.Cookie("refreshToken")
	if err != nil {
		util.WriteApiMessage(w, "Missing refresh token", http.StatusForbidden, false)
		util.Log.Println("Missing refresh token", r.RemoteAddr)
		return
	}

	var (
		jwtClaims  jwt.MapClaims
		statusCode int
	)
	jwtClaims, statusCode, err = util.VerifyJwtToken(refreshTokenCookie.Value)
	if err != nil {
		util.WriteApiMessage(w, err.Error(), statusCode, false)
		util.Log.Println("Failed to verify refresh token, err:", err, r.RemoteAddr)
		return
	}

	userDetails, statusCode, err := util.ParseJwtClaims(jwtClaims)
	if err != nil {
		util.WriteApiMessage(w, err.Error(), statusCode, false)
		util.Log.Println("Failed to parse jwt, err:", err)
		return
	}

	accessToken, accessTokenExpiresIn, err := util.CreateAccessToken(schema.UserForToken{
		Id:    userDetails["id"].(primitive.ObjectID),
		Email: userDetails["email"].(string),
	})
	if err != nil {
		util.WriteApiMessage(w, "Sorry, something went wrong", 0, false)
		util.Log.Println("Failed to create refresh access token, err:", err)
		return
	}

	var accessTokenCookie *http.Cookie = &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		MaxAge:   int(accessTokenExpiresIn),
		Path:     common.AccessTokenPath,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}

	http.SetCookie(w, accessTokenCookie)

	w.WriteHeader(http.StatusOK)
}
