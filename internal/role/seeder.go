package role

import (
	"auth/config"
	"fmt"
)

func SeedRole() {
	roles := []RoleModel{
		{Name: "Admin"},
		{Name: "Staff"},
	}

	for _, role := range roles {
		if err := config.DB.Create(&role).Error; err != nil {
			fmt.Println("Failed to seed role:", err)
			return
		}
	}

	fmt.Println("Role seeded successfully!")
}
