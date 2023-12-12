package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        string `gorm:"primary_key;column:id"`
	Account   string `gorm:"type:varchar(20); unique_index; column:account; default:'' "`
	Username  string `gorm:"type:varchar(20); unique_index; column:username; default:'默认名' " `
	Password  string `gorm:"type:varchar(255); column:password"`
	Avatar    string `gorm:"default:''; column:avatar"`
	Gender    uint8  `gorm:"type:tinyint(3);default:0; column:gender"`
	Mobile    string `gorm:"type:varchar(20);default:'';unique_index;column:mobile"`
	Email     string `gorm:"type:varchar(30);default:'';unique_index;column:email"`
	Introduce string `gorm:"default:'';column:introduce"`
	Status    uint8  `gorm:"type:tinyint(3);default:1;column:status"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate 在插入记录之前调用，设置创建时间和更新时间
func (user *User) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

// BeforeUpdate 在更新记录之前调用，设置更新时间
func (user *User) BeforeUpdate(tx *gorm.DB) error {
	user.UpdatedAt = time.Now()
	return nil
}
