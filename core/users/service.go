package users

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"management/core/dao"
	"management/core/domain"
	"management/infra/base"
	"sync"
)

var _ UserService = new(userService)
var once sync.Once

func init() {
	once.Do(func() {
		IUserService = new(userService)
	})
}

type userService struct {
}

func (u *userService) Login(phone, password string) (*domain.User, error) {
	db := base.MgoDatabase()
	userDao := dao.UserDao{db}
	user, err := userDao.GetUserByPhone(phone)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if user == nil {
		return nil, errors.New("该手机号未注册")
	}
	if user.Password != password {
		err := errors.New("密码错误")
		return nil, err
	}
	return user, nil

}

// 创建用户
func (u *userService) Create(user domain.User) (*domain.User, error) {
	if user.Phone == "" {
		return nil, errors.New("手机号不能为空")
	}
	db := base.MgoDatabase()
	userDao := dao.UserDao{db}
	userRes, err := userDao.GetUserByPhone(user.Phone)
	if userRes != nil {
		log.Error(err)
		return nil, errors.New("该手机号被注册")
	}

	err = userDao.Insert(&user)
	if err != nil {
		log.Error(err)
		return nil, errors.New("注册失败")
	}
	userRes, err = userDao.GetUserByPhone(user.Phone)
	if err != nil {
		log.Error(err)
		return nil, errors.New("注册失败")
	}
	return userRes, err
}
func (u *userService) ResetPassword(user domain.User, newPassword string) error {

	db := base.MgoDatabase()
	userDao := dao.UserDao{db}
	userRes, err := userDao.GetUserByUserId(string(user.Id))
	if err != nil {
		log.Error(" userService ResetPassword:", err)
		return errors.New("修改密码失败")
	}
	if userRes == nil {
		log.Error(" userService userRes == nil", err)
		return errors.New("修改密码失败")
	}

	if userRes.Password != user.Password {
		return errors.New("旧密码错误")
	}
	userRes.Password = newPassword
	err = userDao.Update(userRes)
	if err != nil {
		log.Error(err)
		return errors.New("修改密码失败")
	}
	return nil
}

func (u *userService) Update(user domain.User) error {
	db := base.MgoDatabase()
	userDao := dao.UserDao{db}
	err := userDao.Update(&user)
	if err != nil {
		log.Error(err)
	}
	return err
}
func (u *userService) GetUserById(userId string) (*domain.User, error) {
	db := base.MgoDatabase()
	userDao := dao.UserDao{db}
	userRes, err := userDao.GetUserByUserId(userId)
	if err != nil {
		log.Error(err)
	}
	return userRes, err
}
