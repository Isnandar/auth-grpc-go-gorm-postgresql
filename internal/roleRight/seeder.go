package roleRight

import (
	"auth/config"
	"fmt"
)

func SeedRoleRight() {
	roleRights := []RoleRightModel{
		{
			Route:   "UserService",
			RoleID:  1,
			RCreate: true,
			RRead:   true,
			RUpdate: true,
			RDelete: true,
		},
		{
			Route:   "UserService",
			RoleID:  2,
			RCreate: false,
			RRead:   true,
			RUpdate: false,
			RDelete: false,
		},
	}

	for _, roleRight := range roleRights {
		if err := config.DB.Create(&roleRight).Error; err != nil {
			fmt.Println("Failed to seed role right:", err)
			return
		}
	}

	fmt.Println("Role Right seeded successfully!")
}
