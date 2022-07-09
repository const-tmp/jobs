package server

import (
	"back/proto"
	repo2 "back/repo"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"strings"
)

type server struct {
	proto.UnimplementedAPIServer
	repo *repo2.Repo
}

func NewServer(repo *repo2.Repo) *server {
	s := server{repo: repo}
	return &s
}

func (s *server) CreateUser(_ context.Context, req *proto.User) (*proto.User, error) {
	user := UserDeserialize(req)
	tx := s.repo.DB.Debug().Create(&user)
	if tx.Error != nil {
		return nil, status.Errorf(codes.Internal, "create user error: %s", tx.Error)
	}

	return UserSerialize(user), nil
}

func (s *server) GetUserByID(_ context.Context, req *proto.IDRequest) (*proto.User, error) {
	user := repo2.User{}
	tx := s.repo.DB.Debug().Take(&user, req.Id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user %d not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "db error: %s", tx.Error)
	}

	return UserSerialize(&user), nil
}

func (s *server) GetAllUsers(_ context.Context, req *proto.User) (*proto.Users, error) {
	var users []repo2.User
	tx := s.repo.DB.Debug().Where(UserDeserialize(req)).Find(&users)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "users not found")
		}
		return nil, status.Errorf(codes.Internal, "db error: %s", tx.Error)
	}
	log.Println("len users:", len(users))
	log.Println("users:", users)

	u := make([]*proto.User, len(users))
	for i, user := range users {
		u[i] = UserSerialize(&user)
	}
	log.Println("len users:", len(u))
	log.Println("users:", u)
	return &proto.Users{Users: u}, nil
}

func (s *server) SetPhone(_ context.Context, req *proto.SetPhoneRequest) (*proto.Result, error) {
	var phone, prefix string
	if !strings.HasPrefix(req.Phone, "+") {
		prefix = "+"
	}
	phone = prefix + req.Phone
	tx := s.repo.DB.Debug().Model(repo2.User{
		BaseModel: repo2.BaseModel{ID: req.Id},
	}).Update("phone", phone)
	if tx.Error != tx.Error {
		return nil, status.Errorf(codes.Internal, "db error: %s", tx.Error)
	}
	return &proto.Result{Ok: true}, nil
}

func (s *server) GetOrCreateUser(_ context.Context, req *proto.User) (*proto.User, error) {
	user := UserDeserialize(req)
	tx := s.repo.DB.Debug().Where(user).FirstOrCreate(user)
	if tx.Error != nil {
		return nil, status.Errorf(codes.Internal, "db error: %s", tx.Error)
	}
	return UserSerialize(user), nil
}

func (s *server) CreateCity(_ context.Context, req *proto.City) (*proto.City, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCity not implemented")
}

func (s *server) GetOrCreateCity(_ context.Context, req *proto.City) (*proto.City, error) {
	city := CityDeserialize(req)
	tx := s.repo.DB.Debug().Where(city).FirstOrCreate(city)
	if tx.Error != nil {
		return nil, status.Errorf(codes.Internal, "db error: %s", tx.Error)
	}
	return CitySerialize(city), nil
}

func (s *server) GetAllCities(_ context.Context, req *proto.EmptyRequest) (*proto.Cities, error) {
	var cities []repo2.City
	tx := s.repo.DB.Debug().Find(&cities)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "cities not found")
		}
		return nil, status.Errorf(codes.Internal, "db error: %s", tx.Error)
	}

	u := make([]*proto.City, len(cities))
	for _, city := range cities {
		u = append(u, CitySerialize(&city))
	}
	return &proto.Cities{Cities: u}, nil
}

func (s *server) CreateSkill(_ context.Context, req *proto.Skill) (*proto.Skill, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSkill not implemented")
}
func (s *server) GetAllSkills(_ context.Context, req *proto.EmptyRequest) (*proto.Skill, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllSkills not implemented")
}
func (s *server) CreateProfession(_ context.Context, req *proto.Profession) (*proto.Profession, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProfession not implemented")
}
func (s *server) GetAllProfessions(_ context.Context, req *proto.EmptyRequest) (*proto.Professions, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllProfessions not implemented")
}

func (s *server) CreateCV(_ context.Context, req *proto.CV) (*proto.CV, error) {
	log.Println("CreateCV: req:", req)
	cv := CVDeserialize(req)
	log.Println("CreateCV: cv:", cv)

	err := s.repo.DB.Debug().Omit("ModeratedBy").Create(cv).Error
	if err != nil {
		return nil, err
	}

	err = s.repo.DB.Debug().Preload(clause.Associations).Take(cv).Error
	if err != nil {
		return nil, err
	}
	log.Println("CreateCV: cv:", cv)

	return CVSerialize(cv), nil
}

func (s *server) UpdateCV(_ context.Context, req *proto.CV) (*proto.CV, error) {
	cv := CVDeserialize(req)

	err := s.repo.DB.Debug().Omit("User", "City", "DeletedAt").Save(cv).Error
	if err != nil {
		return nil, err
	}
	log.Println("UpdateCV: cv after save:", cv)

	err = s.repo.DB.Debug().Preload(clause.Associations).Take(cv).Error
	if err != nil {
		return nil, err
	}
	log.Println("UpdateCV: cv after take:", cv)

	return CVSerialize(cv), nil
}

func (s *server) GetAllCVs(_ context.Context, req *proto.EmptyRequest) (*proto.CVs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllCVs not implemented")
}
func (s *server) CreateAd(_ context.Context, req *proto.Ad) (*proto.Ad, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAd not implemented")
}
func (s *server) GetAllAds(_ context.Context, req *proto.EmptyRequest) (*proto.Ads, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllAds not implemented")
}

func (s *server) GetCVByID(_ context.Context, req *proto.IDRequest) (*proto.CV, error) {
	cv := repo2.CV{BaseModel: repo2.BaseModel{ID: req.Id}}
	err := s.repo.DB.Debug().Preload(clause.Associations).Take(&cv).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error: %s", err)
	}
	return CVSerialize(&cv), nil
}
