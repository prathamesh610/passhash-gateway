package service

import (
	"context"
	"log"
	"net/mail"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"prathameshj.dev/passhash-gateway/constants"
	"prathameshj.dev/passhash-gateway/db"
	"prathameshj.dev/passhash-gateway/dberrors"
	"prathameshj.dev/passhash-gateway/models"
)

func SignUp(db db.DatabaseClient, ctx context.Context, newUser *models.NewUser) error {
	emailBool := checkEmail(newUser.Email)

	if !emailBool {
		log.Fatal()
		return &dberrors.CustomError{
			Err: constants.INVALID_EMAIL_FORMAT,
		}
	}

	_, err := db.FindByEmail(ctx, newUser.Email)

	if err == nil {
		return &dberrors.CustomError{
			Err: constants.EMAIL_ALREADY_EXISTS,
		}
	}

	// happy flow - user doesnot exist.
	hashedPassword, err := generatehashPassword(newUser.Password)
	if err != nil {
		log.Fatalln("error in password hash")
		return &dberrors.CustomError{
			Err: constants.PASSWORD_HASHING,
		}
	}

	// add user
	var user models.User
	user.Name = newUser.Name
	user.Email = newUser.Email
	user.Password = hashedPassword
	user.Role = "User"
	db.AddUser(ctx, &user)
	return nil

}

func SignIn(db db.DatabaseClient, ctx context.Context, authdetails *models.Authentication) (*models.Token, error) {

	emailBool := checkEmail(authdetails.Email)

	if !emailBool {
		log.Fatal()
		return nil, &dberrors.CustomError{
			Err: constants.INVALID_EMAIL_FORMAT,
		}
	}

	user, err := db.FindByEmail(ctx, authdetails.Email)

	if err != nil {
		return nil, err
	}

	check := checkPasswordHash(authdetails.Password, user.Password)

	if !check {
		//raise error
		return nil, &dberrors.CustomError{
			Err: constants.INVALID_PASSWORD,
		}
	}
	validToken, err := generateJWT(user.Email, user.Role)

	if err != nil {
		return nil, err
	}

	var token models.Token
	token.Email = user.Email
	token.Role = user.Role
	token.TokenString = validToken

	return &token, nil
}

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func checkEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func generatehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func generateJWT(email string, role string) (string, error) {
	secretKey := os.Getenv("jwt-secret")
	mySigningKey := []byte(secretKey)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		log.Fatalf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}
