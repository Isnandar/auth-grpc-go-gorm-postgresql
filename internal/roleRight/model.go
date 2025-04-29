package roleRight

import (
	"auth/internal/role"

	"gorm.io/gorm"
)

type RoleRightModel struct {
	gorm.Model
	ID      uint           `gorm:"primarykey"`
	Route   string         `gorm:"type:varchar(255);not null"`
	RCreate bool           `gorm:"type:boolean;not null"`
	RRead   bool           `gorm:"type:boolean;not null"`
	RUpdate bool           `gorm:"type:boolean;not null"`
	RDelete bool           `gorm:"type:boolean;not null"`
	RoleID  uint           `gorm:"not null"`
	Role    role.RoleModel `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
