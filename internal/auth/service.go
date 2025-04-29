package auth

import (
	"auth/config"
	"auth/internal/user"
	"fmt"
)

func GetUserFromDB(email string) (*user.UserModel, error) {
	var userModel user.UserModel
	if err := config.DB.Where("email = ?", email).First(&userModel).Error; err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	return &userModel, nil
}
