package user

import (
	"auth/internal/role"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	ID       uint           `gorm:"primarykey"`
	Email    string         `gorm:"type:varchar(255);unique;not null"`
	Name     string         `gorm:"type:varchar(255);not null"`
	Password string         `gorm:"type:varchar(255);not null"`
	RoleID   uint           `gorm:"not null"`
	Role     role.RoleModel `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
