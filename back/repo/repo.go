package repo

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

const (
	Host     = "localhost"
	Port     = 5432
	Username = "postgres"
	Password = "jobs"
	DBName   = "postgres"
	TZ       = "Europe/Kiev"
)

type Repo struct {
	DB *gorm.DB
}

func (r *Repo) Migrate() error {
	return r.DB.AutoMigrate(&User{}, &City{}, &Profession{}, &Skill{}, &Ad{}, &CV{})
}

func New() (*Repo, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		TZ,
	)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}
	return &Repo{DB: db}, nil
}

func (r *Repo) CreateUser(id uint64, name, role string) (*User, error) {
	var user = User{
		BaseModel: BaseModel{ID: id},
		Name:      name,
		Role:      role,
	}
	res := r.DB.Debug().Create(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

func (r *Repo) GetAllUsers() ([]User, error) {
	var users []User
	res := r.DB.Debug().Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}
	return users, nil
}

func (r Repo) GetUserById(id uint64) (*User, error) {

}
