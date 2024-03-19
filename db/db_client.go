package db

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"prathameshj.dev/passhash-gateway/models"
)

type DatabaseClient interface {
	Ready() bool
	FindByEmail(ctx context.Context, emailAddress string) (*models.User, error)
	AddUser(ctx context.Context, user *models.User) (*models.User, error) 
}

type Client struct {
	DB *gorm.DB
}

func NewDataBaseClient() (DatabaseClient, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
	 "localhost", 
	 "postgres",
	 "postgres",
	 "postgres",
	 5432,
	 "disable")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "passhash.",
		},
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		QueryFields: true,
	})

	if err != nil {
		return nil, err
	}

	client := Client{
		DB: db,
	}

	return client, nil
}

func (c Client) Ready() bool {
	var ready string
	tx := c.DB.Raw("SELECT 1 as ready").Scan(&ready)

	if tx.Error != nil{
		return false
	}

	if ready == "1"{
		return true
	}

	return false
}