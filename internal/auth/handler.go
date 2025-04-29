package auth

import (
	"auth/config"
	"auth/internal/user"
	pbauth "auth/pb/auth"
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	pbauth.UnimplementedAuthServiceServer
}

var jwtSecret = []byte("!#2025@%")

func (h *AuthHandler) Login(ctx context.Context, req *pbauth.LoginRequest) (*pbauth.LoginResponse, error) {
	var userModel user.UserModel
	if err := config.DB.Where("email = ?", req.Email).First(&userModel).Error; err != nil {
		return &pbauth.LoginResponse{
			Status:  false,
			Message: "Invalid Credentials",
		}, fmt.Errorf("invalid credentials: %v", err)
	}

	err := bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(req.Password))
	if err != nil {
		return &pbauth.LoginResponse{
			Status:  false,
			Message: "Invalid Credentials",
		}, fmt.Errorf("invalid credentials: %v", err)
	}

	token, err := generateJWT(userModel.ID)
	if err != nil {
		return &pbauth.LoginResponse{
			Status:  false,
			Message: "Fail Token",
		}, fmt.Errorf("fail token: %v", err)
	}

	redisKey := fmt.Sprintf("user_session:%s", token)
	sessionData := map[string]interface{}{
		"role_id": userModel.RoleID,
		"email":   userModel.Email,
		"token":   token,
	}

	err = config.SetRedisValue(redisKey, sessionData, time.Hour*24)
	if err != nil {
		return &pbauth.LoginResponse{
			Status:  false,
			Message: "Failed to store session in Redis",
		}, fmt.Errorf("failed to store session in Redis: %v", err)
	}

	return &pbauth.LoginResponse{
		Status:  true,
		Message: "Successfully logged in",
		Data: &pbauth.Data{
			AccessToken: token,
		},
	}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, req *pbauth.LogoutRequest) (*pbauth.LogoutResponse, error) {

	redisKey := fmt.Sprintf("auth_user_session:%s", req.AccessToken)
	err := config.DeleteRedisValue(redisKey)
	if err != nil {
		return &pbauth.LogoutResponse{
			Status:  false,
			Message: "Failed to log out",
		}, fmt.Errorf("failed to delete session from Redis: %v", err)
	}

	return &pbauth.LogoutResponse{
		Status:  true,
		Message: "Successfully logged out",
	}, nil
}

func generateJWT(userID uint) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Issuer:    "Auth",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
