package server

import (
	"encoding/json"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"prathameshj.dev/passhash-gateway/constants"
	"prathameshj.dev/passhash-gateway/models"
	"prathameshj.dev/passhash-gateway/service"
)

func (s *GinServer) SignUp(ctx *gin.Context) {
	var user models.NewUser

	err := json.NewDecoder(ctx.Request.Body).Decode(&user)
	ctx.Header("Content-Type", "application/json")
	if err != nil {
		fmt.Printf("Error reading body, %v", err)
		ctx.JSON(http.StatusInternalServerError, "Error reading body")

		return
	}

	err = service.SignUp(s.DB, ctx, &user)

	if err != nil {
		fmt.Printf("Error creating user, %v", err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "Successfully creted user. Please login to continue")
}

func (s *GinServer) SignIn(ctx *gin.Context) {
	var authdetails models.Authentication
	err := json.NewDecoder(ctx.Request.Body).Decode(&authdetails)

	ctx.Header("Content-Type", "application/json")
	if err != nil {
		fmt.Printf("Error reading body, %v", err)
		ctx.JSON(http.StatusInternalServerError, "Error reading body")

		return
	}

	token, err := service.SignIn(s.DB, ctx, &authdetails)

	if err != nil {
		fmt.Printf("Error logging user")
		ctx.JSON(http.StatusForbidden, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, *token)
}

func (s *GinServer) Authenticate(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	token := ctx.Request.Header.Get("Authorization")
	// fmt.Printf("Token is %s", token)
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, constants.NO_AUTH_TOKEN)
		return
	}

	role, err := service.Authenticate(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx.JSON(http.StatusOK, role)
}
