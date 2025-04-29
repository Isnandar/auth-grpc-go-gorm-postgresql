package main

import (
	"auth/config"
	"auth/internal/role"
	"auth/internal/roleRight"
	"auth/internal/user"
	"fmt"
)

func main() {
	config.InitDB()

	if err := config.DB.AutoMigrate(&role.RoleModel{}); err != nil {
		panic("Auto Migrate Role failed: " + err.Error())
	}
	role.SeedRole()

	if err := config.DB.AutoMigrate(&roleRight.RoleRightModel{}); err != nil {
		panic("Auto Migrate Role Right failed: " + err.Error())
	}
	roleRight.SeedRoleRight()

	if err := config.DB.AutoMigrate(&user.UserModel{}); err != nil {
		panic("Auto Migrate User failed: " + err.Error())
	}
	user.SeedUser()

	fmt.Println("Migration completed successfully!")
}
