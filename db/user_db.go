package db

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"prathameshj.dev/passhash-gateway/dberrors"
	"prathameshj.dev/passhash-gateway/models"
)

func (c Client) FindByEmail(ctx context.Context, emailAddress string) (*models.User, error) {
	user := &models.User{}

	result := c.DB.WithContext(ctx).Where(&models.User{Email: emailAddress}).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &dberrors.NotFoundError{Entity: "user", ID: emailAddress}
		}
		return nil, result.Error
	}

    return user, nil

}


func (c Client) AddUser(ctx context.Context, user *models.User) (*models.User, error) {

    result := c.DB.WithContext(ctx).Create(&user)
    
    if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey){
			return nil, &dberrors.ConflictError{}
		}
		return nil, result.Error
	}

	return user, nil
}