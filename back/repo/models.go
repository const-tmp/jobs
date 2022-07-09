package repo

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type City struct {
	BaseModel
	Name string `gorm:"unique"`
}

const (
	RoleUser      = "user"
	RoleHirer     = "hirer"
	RoleWorker    = "worker"
	RoleModerator = "moderator"
	RoleAdmin     = "admin"
)

type User struct {
	BaseModel
	Name  string
	Phone string
	Role  string `gorm:"default:'user'"`
}

type CV struct {
	BaseModel
	UserID      uint64
	User        User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Raw         string
	Desc        string
	CityID      uint64
	City        City         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Professions []Profession `gorm:"many2many:cv_professions;"`
	Skills      []Skill      `gorm:"many2many:cv_skills;"`
	Moderated   bool
	ModeratedBy uint64
}

type Ad struct {
	BaseModel
	UserID      uint64
	User        User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Raw         string
	Desc        string
	CityID      uint64
	City        City         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Professions []Profession `gorm:"many2many:ad_professions;"`
	Skills      []Skill      `gorm:"many2many:ad_skills;"`
}

type Profession struct {
	BaseModel
	Name string `gorm:"unique"`
	Ads  []Ad   `gorm:"many2many:ad_professions;"`
	CVs  []CV   `gorm:"many2many:cv_professions;"`
}

type Skill struct {
	BaseModel
	Name string `gorm:"unique"`
	Ads  []Ad   `gorm:"many2many:ad_skills;"`
	CVs  []CV   `gorm:"many2many:cv_skills;"`
}
