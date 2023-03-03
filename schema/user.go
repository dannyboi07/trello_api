package schema

import (
	"fmt"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRegister struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (u *UserRegister) Validate() (int, error) {
	if u.Email == nil {
		return http.StatusBadRequest, fmt.Errorf("Field 'email' is empty")
	} else if u.Password == nil {
		return http.StatusBadRequest, fmt.Errorf("Field 'password' is empty")
	} else if trimmedPw := strings.Trim(*u.Password, " "); len(trimmedPw) < 4 || len(trimmedPw) > 10 {
		return http.StatusBadRequest, fmt.Errorf("Password must be between 4 & 10 characters")
	}

	return 0, nil
}

type UserLogin struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (u *UserLogin) Validate() (int, error) {
	if u.Email == nil {
		return http.StatusBadRequest, fmt.Errorf("Field 'email' is empty")
	} else if u.Password == nil {
		return http.StatusBadRequest, fmt.Errorf("Field 'password' is empty")
	}

	return 0, nil
}

type User struct {
	Id    primitive.ObjectID `json:"id"`
	Email string             `json:"email"`
}

type UserForToken struct {
	Id    primitive.ObjectID `json:"id"`
	Email string             `json:"email"`
}
