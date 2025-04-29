package main

import (
	"auth/config"
	"auth/internal/auth"
	"auth/internal/user"
	pbauth "auth/pb/auth"
	pbuser "auth/pb/user"
	"fmt"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		return
	}

	config.InitRedis()
	config.InitDB()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pbuser.RegisterUserServiceServer(grpcServer, &user.UserHandler{})
	pbauth.RegisterAuthServiceServer(grpcServer, &auth.AuthHandler{})

	fmt.Println("gRPC server running on port", port)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
