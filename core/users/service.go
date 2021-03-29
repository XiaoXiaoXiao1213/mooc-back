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
		dao := UserDao{runner}
		oldUser := dao.GetOne(user.Phone, user.UserType)
		if oldUser == nil {
			err := errors.New("用户不存在")
			log.Error(err)
			return err
		}

		createEditUser(oldUser, user)
		log.Error("ea", oldUser)
		oldUser.UpdatedAt = time.Now()
		update, err := dao.Update(oldUser)
		log.Error("ea", update)

		if update < 1 || err != nil {
			log.Error(err, fmt.Sprintf("update num %d", update))
			err := errors.New("更新失败")
			return err
		}
		return nil
	})
	return err
}

func (u *userService) Login(phone, password string, userType int) (*User, error) {
	var user *User
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := UserDao{runner}
		user = dao.GetOne(phone, userType)
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

	return user, err
}

// 创建用户
func (u *userService) Create(user User) error {
	user.Password = user.Id_code[len(user.Id_code)-6:]
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := UserDao{runner}
		_, err := dao.Runner.Insert(user)
		if err != nil {
			log.Error(err)
			return err
		}
		// 创建员工关联表
		score := EmployeeScore{
			EmployeeId: user.Id,
		}
		scoreDao := EmployeeScoreDao{runner}
		_, err = scoreDao.Insert(&score)
		if err != nil {
			log.Error(err)
			return err
		}
		return nil
	})
	return err
}
func (u *userService) ResetPassword(user User) error {
	var newUser *User
	dao := UserDao{}
	newUser = dao.GetOne(user.Phone, user.UserType)
	if newUser == nil {
		err := errors.New("用户不存在")
		log.Error(err)
		return err
	}
	newUser.Password = newUser.Id_code[len(newUser.Id_code)-6:]
	newUser.UpdatedAt = time.Now()
	updateCount, err := dao.Update(newUser)
	if err != nil || updateCount < 1 {
		log.Error(err)
		return errors.New("重置密码失败")
	}
	return nil
}
func (u *userService) GetUserByCond(cond User) (*[]User, int, error) {
	var users *[]User
	var total int
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		dao := UserDao{runner}
		if cond.Page == 0 {
			cond.Page = 1
		}
		if cond.PageSize < 1 || cond.PageSize > 10 {
			cond.PageSize = 10
		}
		users, total = dao.GetByCond(cond)

		return nil
	})
	return users, total, nil
}
func (u *userService) DeleteUserById(userId int64) (err error) {
	err = base.Tx(func(runner *dbx.TxRunner) error {
		dao := UserDao{runner}
		_, err := dao.DeleteByUserId(userId)
		if err != nil {
			log.Error(err)
			return err
		}
		return nil
	})
	return err
}

func (u *userService) GetUserByPhone(phone string, userType int) (user *User, err error) {
	err = base.Tx(func(runner *dbx.TxRunner) error {
		dao := UserDao{runner}
		log.Error(phone, userType)
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
func (u *userService) GetUserById(userId int64) (user *User, err error) {
	err = base.Tx(func(runner *dbx.TxRunner) error {
		dao := UserDao{runner}
		user = dao.GetOneById(userId)
		if user == nil {
			err := errors.New("用户不存在")
			log.Error(err)
			return err
		}
		return nil

	})
	return user, err
}

func createEditUser(editUser *User, user User) {
	if user.Password != "" {
		editUser.Password = user.Password
	}
	if user.Name != "" {
		editUser.Name = user.Name
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
