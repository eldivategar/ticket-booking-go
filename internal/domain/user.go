package domain

type User struct {
	BaseModel
	FullName string `gorm:"type:varchar(100);not null" json:"full_name"`
	Username string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Avatar   string `gorm:"type:text" json:"avatar,omitempty"`
}
