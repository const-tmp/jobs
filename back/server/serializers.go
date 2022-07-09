package server

import (
	"back/proto"
	"back/repo"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"log"
)

func BaseDeserialize(base *proto.Base) *repo.BaseModel {
	b := repo.BaseModel{}
	if base != nil {
		b.ID = base.Id
		b.CreatedAt = base.CreatedAt.AsTime()
		b.UpdatedAt = base.UpdatedAt.AsTime()
		if base.DeletedAt != nil {
			b.DeletedAt = gorm.DeletedAt{Time: base.DeletedAt.AsTime(), Valid: true}
		}
	}
	return &b
}

func BaseSerialize(base *repo.BaseModel) *proto.Base {
	return &proto.Base{
		Id:        base.ID,
		CreatedAt: timestamppb.New(base.CreatedAt),
		UpdatedAt: timestamppb.New(base.UpdatedAt),
		DeletedAt: timestamppb.New(base.DeletedAt.Time),
	}
}

func UserDeserialize(user *proto.User) *repo.User {
	u := repo.User{
		Name:      user.Name,
		Phone:     user.Phone,
		Role:      user.Role,
		BaseModel: *BaseDeserialize(user.Base),
	}
	return &u
}

func UserSerialize(user *repo.User) *proto.User {
	return &proto.User{
		Name:  user.Name,
		Phone: user.Phone,
		Role:  user.Role,
		Base:  BaseSerialize(&user.BaseModel),
	}
}

func CityDeserialize(city *proto.City) *repo.City {
	u := repo.City{
		Name:      city.Name,
		BaseModel: *BaseDeserialize(city.Base),
	}
	return &u
}

func CitySerialize(city *repo.City) *proto.City {
	return &proto.City{
		Name: city.Name,
		Base: BaseSerialize(&city.BaseModel),
	}
}

func SkillDeserialize(skill *proto.Skill) *repo.Skill {
	u := repo.Skill{
		Name:      skill.Name,
		BaseModel: *BaseDeserialize(skill.Base),
	}
	return &u
}

func SkillSerialize(skill *repo.Skill) *proto.Skill {
	return &proto.Skill{
		Name: skill.Name,
		Base: BaseSerialize(&skill.BaseModel),
	}
}

func ProfessionDeserialize(prof *proto.Profession) *repo.Profession {
	u := repo.Profession{
		Name:      prof.Name,
		BaseModel: *BaseDeserialize(prof.Base),
	}
	return &u
}

func ProfessionSerialize(prof *repo.Profession) *proto.Profession {
	return &proto.Profession{
		Name: prof.Name,
		Base: BaseSerialize(&prof.BaseModel),
	}
}

func CVDeserialize(cv *proto.CV) *repo.CV {
	log.Println("CVDeserialize: proto:", cv)
	log.Println("CVDeserialize: proto: professions:", cv.Professions)
	log.Println("CVDeserialize: proto: skills:", cv.Skills)
	var profs []repo.Profession
	if cv.Professions != nil {
		profs = make([]repo.Profession, len(cv.Professions))
		for i, profession := range cv.Professions {
			profs[i] = *ProfessionDeserialize(profession)
		}
	}

	var skills []repo.Skill
	if cv.Skills != nil {
		skills = make([]repo.Skill, len(cv.Skills))
		for i, skill := range cv.Skills {
			skills[i] = *SkillDeserialize(skill)
		}
	}

	newCV := repo.CV{
		BaseModel:   *BaseDeserialize(cv.Base),
		UserID:      cv.UserId,
		Raw:         cv.Raw,
		Desc:        cv.Desc,
		CityID:      cv.CityId,
		Professions: profs,
		Skills:      skills,
		Moderated:   cv.Moderated,
		ModeratedBy: cv.ModeratedBy,
	}
	if cv.User != nil {
		newCV.User = *UserDeserialize(cv.User)
	}
	if cv.City != nil {
		newCV.City = *CityDeserialize(cv.City)
	}
	log.Println("CVDeserialize: db:", newCV)
	log.Println("CVDeserialize: db: professions:", newCV.Professions)
	log.Println("CVDeserialize: db: skills:", newCV.Skills)

	return &newCV
}

func CVSerialize(cv *repo.CV) *proto.CV {
	log.Println("CVSerialize: db:", cv)
	var profs []*proto.Profession
	if cv.Professions != nil {
		profs = make([]*proto.Profession, len(cv.Professions))
		for i, profession := range cv.Professions {
			profs[i] = ProfessionSerialize(&profession)
		}
	}

	var skills []*proto.Skill
	if cv.Professions != nil {
		skills = make([]*proto.Skill, len(cv.Skills))
		for i, skill := range cv.Skills {
			skills[i] = SkillSerialize(&skill)
		}
	}

	tmp := proto.CV{
		Base:        BaseSerialize(&cv.BaseModel),
		UserId:      cv.UserID,
		User:        UserSerialize(&cv.User),
		Raw:         cv.Raw,
		Desc:        cv.Desc,
		CityId:      cv.CityID,
		City:        CitySerialize(&cv.City),
		Professions: profs,
		Skills:      skills,
	}
	log.Println("CVSerialize: proto:", cv)

	return &tmp
}
