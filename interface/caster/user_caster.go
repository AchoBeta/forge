package caster

import (
	"forge/biz/entity"
	"forge/interface/def"
	"github.com/bytedance/gg/gslice"
)

// note: 视图与实体并不是1:1的关系，比如实体有password这个字段，但是这个字段我是不可能给前端的
// 实体对象不关心视图是啥样的，所有东西都是为实体服务的

// CastUserDTO2DO
//
//	@Description: 用户视图视图与实体转化
//	@param dto
//	@return *entity.User
func CastUserDTO2DO(dto *def.User) *entity.User {
	return &entity.User{
		UserID: dto.UserID,
		Name:   dto.Username,
		Dogs:   CastDogDTOs2DOs(dto.Dogs), // 演示领域间的关联关系
	}
}

func CastUserDO2DTO(dto *entity.User) *def.User {
	return &def.User{
		UserID:   dto.UserID,
		Username: dto.Name,
		Dogs:     CastDogDOs2DTOs(dto.Dogs),
	}
}

// CastDogDTOs2DOs
//
//	@Description: 狗视图与实体转换 这里看使用量，如果只有一个地方使用可以直接用gslice.Map就可以了，可以灵活一点
//	@param dtos
//	@return []*entity.Dog
func CastDogDTOs2DOs(dtos []*def.Dog) []*entity.Dog {
	return gslice.Map(dtos, CastDogDTO2DO)
}

// CastDogDTO2DO
//
//	@Description: 狗视图与实体转换
//	@param dto
//	@return *entity.Dog
func CastDogDTO2DO(dto *def.Dog) *entity.Dog {
	return &entity.Dog{
		DogID:   dto.DogID,
		DogName: dto.DogName,
	}
}

// CastDogDO2DTO
//
//	@Description:
//	@param do
//	@return *def.Dog
func CastDogDO2DTO(do *entity.Dog) *def.Dog {
	return &def.Dog{
		DogID:   do.DogID,
		DogName: do.DogName,
	}
}

// CastDogDOs2DTOs
//
//	@Description:
//	@param dos
//	@return []*def.Dog
func CastDogDOs2DTOs(dos []*entity.Dog) []*def.Dog {
	return gslice.Map(dos, CastDogDO2DTO)
}
