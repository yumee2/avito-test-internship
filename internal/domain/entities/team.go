package entities

type Team struct {
	TeamName    string `gorm:"size:16;primaryKey;not null;unique"`
	TeamMembers []User `gorm:"foreignKey:TeamName;references:TeamName"`
}

type TeamMember struct {
	UserId   string `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
}
