package models

type User_Password_Reset struct {
	Email   string `gorm:"primary_key; size:100;not null;unique" json:"email"`
	PIN     string `gorm:"not null;" json:"pin,omitempty"`
	Enabled bool   `gorm:"default:false" json:"enabled"`
}

type Reset_PIN_Confirm struct {
	PIN string `json:"pin"`
}
