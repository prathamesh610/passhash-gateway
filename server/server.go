package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"prathameshj.dev/passhash-gateway/db"
	"prathameshj.dev/passhash-gateway/models"
)

type Server interface {
	Start()
	Readiness(ctx *gin.Context)
	Liveness(ctx *gin.Context)

	SignUp(ctx *gin.Context)
	SignIn(ctx *gin.Context)

	Authenticate(ctx *gin.Context)
}

type GinServer struct {
	gin *gin.Engine
	DB  db.DatabaseClient
}

func StartServer(db db.DatabaseClient) Server {
	server := &GinServer{
		gin: gin.Default(),
		DB:  db,
	}
	return server
}

func (s *GinServer) Start() {
	s.registerRoutes()
	s.gin.Run(":8080")
}

func (s *GinServer) registerRoutes() {
	s.gin.GET("/readiness", s.Readiness)
	s.gin.GET("/liveness", s.Liveness)

	s.gin.POST("/signup", s.SignUp)
	s.gin.POST("/signin", s.SignIn)

	s.gin.GET("/authenticate", s.Authenticate)

}

func (s *GinServer) Readiness(ctx *gin.Context) {
	ready := s.DB.Ready()
	if ready {
		ctx.JSON(http.StatusOK, models.Health{Status: "OK"})
		return 
	}

	ctx.JSON(http.StatusInternalServerError, models.Health{Status: "Failure"})
}

func (s *GinServer) Liveness(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, models.Health{Status: "OK"})
}
