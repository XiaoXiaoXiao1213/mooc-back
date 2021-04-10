package users

import (
	"management/core/domain"
	"management/infra/base"
)

var IUserService UserService

//用于对外暴露账户应用服务，唯一的暴露点
func GetUserService() UserService {
	base.Check(IUserService)
	return IUserService
}

type UserService interface {
	Create(user domain.User) (*domain.User, error)
	Login(phone, password string) (*domain.User, error)
	ResetPassword(user domain.User,newPassword string) error
	GetUserById(userId string) (*domain.User, error)
	Update(user domain.User) error
}
