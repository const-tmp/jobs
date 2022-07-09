package repo

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

var allTables = []interface{}{&User{}, &City{}, &CV{}, &Ad{}, &Profession{}, &Skill{}}

type RepoTestSuite struct {
	suite.Suite
	Repo *Repo
}

func TestRepo(t *testing.T) {
	suite.Run(t, new(RepoTestSuite))
}

func (s *RepoTestSuite) SetupSuite() {
	r, err := New()
	s.Require().NoError(err)
	s.Repo = r
}

func (s *RepoTestSuite) TestRepo() {
	s.Run("migrations", func() {
		err := s.Repo.DB.Debug().AutoMigrate(allTables...)
		s.Require().NoError(err)
	})

	cities := []City{
		{Name: "Одесса"},
		{Name: "Киев"},
		{Name: "Львов"},
	}
	s.Run("create cities", func() {
		res := s.Repo.DB.Debug().Create(&cities)
		s.Require().NoError(res.Error)
	})
	s.Run("get all cities", func() {
		var c []City
		res := s.Repo.DB.Debug().Find(&c)
		s.Require().NoError(res.Error)
		s.T().Logf("cities: %#v", cities)
	})

	professions := []Profession{
		{Name: "Автомеханик"},
		{Name: "Продавец"},
		{Name: "Няня"},
	}
	s.Run("create professions", func() {
		res := s.Repo.DB.Debug().Create(&professions)
		s.Require().NoError(res.Error)
	})
	s.Run("get all professions", func() {
		var p []Profession
		res := s.Repo.DB.Debug().Find(&p)
		s.Require().NoError(res.Error)
		s.T().Logf("professions: %#v", professions)
	})

	skills := []Skill{
		{Name: "Работа с детьми"},
		{Name: "Решение проблем"},
		{Name: "Мастер на все руки"},
	}
	s.Run("create skills", func() {
		res := s.Repo.DB.Debug().Create(&skills)
		s.Require().NoError(res.Error)
	})
	s.Run("get all skills", func() {
		var p []Profession
		res := s.Repo.DB.Debug().Find(&p)
		s.Require().NoError(res.Error)
		s.T().Logf("skills: %#v", skills)
	})

	testUser := User{
		BaseModel: BaseModel{ID: 12345},
		Name:      "Test User",
		Role:      "+380576235476",
	}

	s.Run("create user", func() {
		u, err := s.Repo.CreateUser(testUser.ID, testUser.Name, testUser.Role)
		s.Require().NoError(err)
		testUser = *u
	})

	s.Run("get all users", func() {
		u, err := s.Repo.GetAllUsers()
		s.Require().NoError(err)
		s.T().Logf("users: %#v", u)
	})

	s.Run("create cv", func() {
		cv := CV{
			User:        testUser,
			Raw:         "нууу я эммм это, тыры пыры там",
			Desc:        "нормальное описание",
			City:        cities[0],
			Professions: professions,
			Skills:      skills,
		}
		res := s.Repo.DB.Debug().Create(&cv)
		s.Require().NoError(res.Error)
		s.T().Logf("new CV: %#v", cv)
		s.T().Logf("CV user: %#v", cv.User)
		s.T().Logf("CV city: %#v", cv.City)
		s.T().Logf("CV profs: %#v", cv.Professions)
		s.T().Logf("CV skills: %#v", cv.Skills)
	})

	//s.Run("query CV", func() {
	//	var cvs []CV
	//	err := s.Repo.DB.Debug().Model(&skills).Preload(clause.Associations).Association("CVs").Find(&cvs)
	//	s.Require().NoError(err)
	//	s.T().Logf("cvs: %#v", cvs)
	//	for i, cv := range cvs {
	//		s.T().Logf("%d %#v", i, cv)
	//	}
	//})
}

//func (s *RepoTestSuite) TearDownSuite() {
//	err := s.Repo.DB.Debug().Migrator().DropTable(allTables...)
//	s.Require().NoError(err)
//}
