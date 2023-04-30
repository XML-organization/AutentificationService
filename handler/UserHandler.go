package handler

import (
	"autentification_service/model"
	"autentification_service/service"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserService *service.UserService
}

const SecretKeyForJWT = "v123v1iy2v321sdasada8"

func (loginHandler *UserHandler) Login(writer http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	var user model.UserDTO
	err := decoder.Decode(&user)
	if err != nil {
		panic(err)
	}

	loggedUser, err := loginHandler.UserService.FindByEmail(user.Email)

	if err != nil {
		//writer.WriteHeader(http.StatusNotFound)
		message := model.RequestMessage{
			Message: "User not found!",
		}
		json.NewEncoder(writer).Encode(message)
		return
	}

	//password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	if err := bcrypt.CompareHashAndPassword(loggedUser.Password, []byte(user.Password)); err != nil {
		//writer.WriteHeader(http.StatusBadRequest)
		message := model.RequestMessage{
			Message: "Incorrect password!",
		}
		json.NewEncoder(writer).Encode(message)
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&model.JwtClaims{
			Id:   loggedUser.ID,
			Role: int(loggedUser.Role),
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			}})

	token, err := claims.SignedString([]byte(SecretKeyForJWT))

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Could not login, JWT token can not be created!"))
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}

	http.SetCookie(writer, &cookie)

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(loggedUser)
	return
}

func (loginHandler *UserHandler) User(writer http.ResponseWriter, req *http.Request) {

	cookies := req.Cookies()

	if len(cookies) == 0 {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte("Unautenticated!"))
		return
	}

	claims := &model.JwtClaims{}

	token, err := jwt.ParseWithClaims(cookies[0].Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKeyForJWT), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte("JWT invalid!"))
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Bad request!"))
		return
	}

	if !token.Valid {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	//ObjectID(\"64199dfe142552ce4d25b90b\")
	id := strings.Split(claims.Id.String(), "\"")[1]

	user, err := loginHandler.UserService.FindUser(id)

	json.NewEncoder(writer).Encode(user)
	writer.WriteHeader(http.StatusOK)
}

func (loginHandler *UserHandler) Logout(writer http.ResponseWriter, req *http.Request) {
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(writer, &cookie)

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Logout success!"))
}

func (handler *UserHandler) Create(writer http.ResponseWriter, req *http.Request) {
	var userDTO model.User
	err := json.NewDecoder(req.Body).Decode(&userDTO)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var user model.UserCredentials
	//hesovanje passworda
	password, _ := bcrypt.GenerateFromPassword([]byte(userDTO.Password), 14)

	user.Password = password
	user.Email = userDTO.Email
	user.Role = model.Role(userDTO.Role)

	err = handler.UserService.Create(&user)
	if err != nil {
		println("Error while creating a new user")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
}
