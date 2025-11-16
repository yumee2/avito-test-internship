package entities

type User struct {
	UserID   string `gorm:"primaryKey;default:gen_random_uuid()"`
	Username string `gorm:"size:16;not null;unique"`
	TeamName string `gorm:"size:16;not null"`
	IsActive bool   `gorm:"not null;default:true"`
}
