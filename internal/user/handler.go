package user

import (
	"auth/config"
	"auth/internal/roleRight"
	pbuser "auth/pb/user"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
)

const UserServiceRoute = "UserService"

type UserHandler struct {
	pbuser.UnimplementedUserServiceServer
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
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

func checkPermission(roleId uint, route string, permissionType string) (bool, error) {
	var roleRight roleRight.RoleRightModel
	if err := config.DB.Where("role_id = ? AND route = ?", roleId, route).First(&roleRight).Error; err != nil {
		return false, fmt.Errorf("role rights not found for route: %s", route)
	}

	switch permissionType {
	case "r_create":
		return roleRight.RCreate, nil
	case "r_read":
		return roleRight.RRead, nil
	case "r_update":
		return roleRight.RUpdate, nil
	case "r_delete":
		return roleRight.RDelete, nil
	default:
		return false, fmt.Errorf("invalid permission type")
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pbuser.CreateUserRequest) (*pbuser.CreateUserResponse, error) {
	sessionData, err := checkSession(ctx)
	if err != nil {
		return &pbuser.CreateUserResponse{Status: false, Message: "Unauthorized"}, nil
	}
	roleId := uint(sessionData["role_id"].(float64))
	canCreate, err := checkPermission(roleId, UserServiceRoute, "r_create")
	if err != nil || !canCreate {
		return &pbuser.CreateUserResponse{Status: false, Message: "Unauthorized"}, nil
	}

	user := UserModel{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
		RoleID:   roleId,
	}
	if err := config.DB.Create(&user).Error; err != nil {
		return &pbuser.CreateUserResponse{Status: false, Message: "Fail"}, err
	}
	return &pbuser.CreateUserResponse{Status: true, Message: "Successfully"}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, _ *pbuser.GetUserRequest) (*pbuser.GetUserResponse, error) {
	sessionData, err := checkSession(ctx)
	if err != nil {
		return &pbuser.GetUserResponse{Status: false, Message: "Unauthorized"}, nil
	}
	roleId := uint(sessionData["role_id"].(float64))
	canRead, err := checkPermission(roleId, UserServiceRoute, "r_read")
	if err != nil || !canRead {
		return &pbuser.GetUserResponse{Status: false, Message: "Unauthorized"}, nil
	}

	var user UserModel
	if err := config.DB.First(&user).Error; err != nil {
		return &pbuser.GetUserResponse{Status: false, Message: "Not Found"}, err
	}
	return &pbuser.GetUserResponse{
		Status:  true,
		Message: "Successfully",
		Data:    &pbuser.Data{User: toProto(&user)},
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pbuser.UpdateUserRequest) (*pbuser.UpdateUserResponse, error) {
	sessionData, err := checkSession(ctx)
	if err != nil {
		return &pbuser.UpdateUserResponse{Status: false, Message: "Unauthorized"}, nil
	}
	roleId := uint(sessionData["role_id"].(float64))
	canUpdate, err := checkPermission(roleId, UserServiceRoute, "r_update")
	if err != nil || !canUpdate {
		return &pbuser.UpdateUserResponse{Status: false, Message: "Unauthorized"}, nil
	}

	var user UserModel
	if err := config.DB.First(&user).Error; err != nil {
		return &pbuser.UpdateUserResponse{Status: false, Message: "Not Found"}, err
	}
	user.Name = req.Name
	if err := config.DB.Save(&user).Error; err != nil {
		return &pbuser.UpdateUserResponse{Status: false, Message: "Fail"}, err
	}
	return &pbuser.UpdateUserResponse{Status: true, Message: "Successfully"}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pbuser.DeleteUserRequest) (*pbuser.DeleteUserResponse, error) {
	sessionData, err := checkSession(ctx)
	if err != nil {
		return &pbuser.DeleteUserResponse{Status: false, Message: "Unauthorized"}, nil
	}
	roleId := uint(sessionData["role_id"].(float64))
	canDelete, err := checkPermission(roleId, UserServiceRoute, "r_delete")
	if err != nil || !canDelete {
		return &pbuser.DeleteUserResponse{Status: false, Message: "Unauthorized"}, nil
	}

	var user UserModel
	if err := config.DB.Where("id = ?", req.UserId).First(&user).Error; err != nil {
		return &pbuser.DeleteUserResponse{Status: false, Message: "Not Found"}, err
	}
	if err := config.DB.Delete(&user).Error; err != nil {
		return &pbuser.DeleteUserResponse{Status: false, Message: "Fail"}, err
	}
	return &pbuser.DeleteUserResponse{Status: true, Message: "Successfully"}, nil
}
