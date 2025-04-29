package user

import (
	"auth/config"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func SeedUser() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("adminadmin"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Failed to hash password:", err)
		return
	}

	user := UserModel{
		Email:    "admin@gmail.com",
		Password: string(hashedPassword),
		Name:     "Administrator",
		RoleID:   1,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		fmt.Println("Failed to seed user:", err)
		return
	}

	fmt.Println("User seeded successfully!")
}
