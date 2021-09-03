package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type User_Email struct {
	Email             string    `gorm:"size:100;not null;unique" json:"email"`
	ID                uuid.UUID `gorm:"primary_key;type:uuid" json:"id"`
	Confirmed         bool      `gorm:"default:false" json:"enabled"`
	ConfirmationToken uuid.UUID `gorm:"type:uuid;not null;unique" json:"confirmation_token"`
}

type Confirm_email struct {
	ConfirmationToken uuid.UUID `gorm:"size:100;not null" json:"confirmation_token"`
}

func (u *User) UpdateAnEmail(db *gorm.DB, uid uuid.UUID, tuid uuid.UUID) (*User, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"user_name":  u.UserName,
			"first_name": u.FirstName,
			"last_name":  u.LastName,
			"updated_at": time.Now(),
			"updated_by": tuid,
			"email":      u.Email,
			"password":   u.Password,
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err := db.Debug().Model(&User{}).Where("id = ?", uid).Preload("Roles").Preload("Roles.Permissions").Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}
