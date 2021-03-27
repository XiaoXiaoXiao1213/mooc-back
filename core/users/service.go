package users

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"management/infra/base"
	"sync"
	"time"
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

func (u *userService) Edit(user User) error {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := UserDao{runner: runner}
		oldUser := dao.GetOne(user.Phone, user.UserType)
		if oldUser == nil {
			err := errors.New("用户不存在")
			log.Error(err)
			return err
		}
		createEditUser(oldUser, &user)
		oldUser.UpdatedAt = time.Now()
		update, err := dao.Update(oldUser)
		if update < 1 || err != nil {
			log.Error(err, fmt.Sprintf("update num %d", update))
			err := errors.New("更新失败")
			return err
		}
		return nil
	})
	return err
}

func (u *userService) Login(phone, password string, userType int) error {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := UserDao{runner: runner}

		user := dao.GetOne(phone, userType)
		//创建用户
		if user == nil {
			err := errors.New("用户不存在")
			log.Error(err)
			return err
		}
		if user.Password != password {
			err := errors.New("密码错误")
			log.Error(err)
			return err
		}
		return nil

	})
	if err != nil {
		log.Error(err)

	}
	return err
}

// 创建用户
func (u *userService) Create(user User) error {
	user.Password = user.Id_code[len(user.Id_code)-6:]
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := UserDao{runner: runner}
		_, err := dao.runner.Insert(user)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
func (u *userService) ResetPassword(user User) error {
	var newUser *User
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := UserDao{runner: runner}
		newUser = dao.GetOne(user.Phone, user.UserType)
		if newUser == nil {
			err := errors.New("用户不存在")
			log.Error(err)
			return err
		}
		newUser.Password = newUser.Id_code[len(newUser.Id_code)-6:]
		newUser.UpdatedAt = time.Now()
		updateCount, err := dao.Update(newUser)
		if err != nil {
			log.Error(err)
			return err
		}
		if updateCount < 1 {
			err := errors.New("重置密码失败")
			log.Error(err)
			return err
		}
		return nil
	})
	return err
}
func (u *userService) GetUserByPhone(phone string, userType int) (user *User, err error) {

	err = base.Tx(func(runner *dbx.TxRunner) error {
		dao := UserDao{runner: runner}
		user = dao.GetOne(phone, userType)
		if user == nil {
			err := errors.New("用户不存在")
			log.Error(err)
			return err
		}
		return nil

	})
	return user, err
}

func createEditUser(editUser *User, user *User) {
	if user.Password != "" {
		editUser.Password = user.Password
	}
	if user.Wechat != "" {
		editUser.Wechat = user.Wechat
	}
	if user.Email != "" {
		editUser.Email = user.Email
	}
	if user.Skills != "" {
		editUser.Skills = user.Skills
	}
	if user.State != 0 {
		editUser.State = user.State
	}
	if user.Score != 0 {
		editUser.Score = user.Score
	}
}
