package role

import (
	"gorm.io/gorm"
)

type RoleModel struct {
	gorm.Model
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"type:varchar(255);not null"`
}
