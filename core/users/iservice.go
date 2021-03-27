package users

import (
	"management/infra/base"
)



var IUserService UserService

//用于对外暴露账户应用服务，唯一的暴露点
func GetUserService() UserService {
	base.Check(IUserService)
	return IUserService
}

type UserService interface {
	Create(user User) error
	Edit(user User) error
	Login(phone, password string, userType int) error
	ResetPassword(user User) error
	GetUserByPhone(phone string, userType int) (*User, error)
	//GetEnvelopeAccountByUserId(userId string) *AccountDTO
	//GetAccount(accountNo string) *AccountDTO
}