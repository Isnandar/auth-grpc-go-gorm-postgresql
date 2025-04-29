package auth

import (
	"auth/config"
	"auth/internal/user"
	pbauth "auth/pb/auth"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
)

type AuthHandler struct {
	pbauth.UnimplementedAuthServiceServer
}

func ExtractTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md["authorization"]) == 0 {
		return "", fmt.Errorf("missing authorization metadata")
	}

	authHeader := md["authorization"][0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid authorization format")
	}

	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

func checkSession(ctx context.Context) (map[string]interface{}, error) {
	token, err := ExtractTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	redisKey := fmt.Sprintf("user_session:%s", token)
	sessionData, err := config.GetRedisValue(redisKey)
	if err != nil {
		return nil, fmt.Errorf("session not found or expired")
	}

	var sessionMap map[string]interface{}
	if err := json.Unmarshal([]byte(sessionData), &sessionMap); err != nil {
		return nil, fmt.Errorf("failed to parse session data")
	}

	return sessionMap, nil
}

var jwtSecret = []byte("!#2025@%")

func (h *AuthHandler) Login(ctx context.Context, req *pbauth.LoginRequest) (*pbauth.LoginResponse, error) {
	currentSessionData, currentErr := checkSession(ctx)
	if currentErr == nil && currentSessionData != nil {
		return &pbauth.LoginResponse{
			Status:  false,
			Message: "Unauthorized: User already logged in",
		}, nil
	}

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
	sessionData, err := checkSession(ctx)
	if err != nil {
		return &pbauth.LogoutResponse{Status: false, Message: "Unauthorized"}, nil
	}

	redisKey := fmt.Sprintf("user_session:%s", sessionData["access_token"])
	err = config.DeleteRedisValue(redisKey)
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
