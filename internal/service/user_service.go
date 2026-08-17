package service

import (
	"errors"
	"gostart/internal/dao/model"
	"gostart/internal/dao/query"
	"gostart/internal/pkg"
	"strconv"
)

// handler 层：接收请求、解析参数、调用 service、返回响应（薄层）
// service 层：业务逻辑编排、校验、事务管理、调用 DAO/模型
// model/dao 层：数据持久化（GORM 查询）

func Login(username string, password string) (*model.User, string, error) {
	if username == "" || password == "" {
		return nil, "", errors.New("用户名或密码不能为空")
	}
	// errors.Is(err, gorm.ErrRecordNotFound)
	user, err := query.User.Where(query.User.Username.Eq(username)).First()
	if err != nil {
		return nil, "", errors.New("用户名不存在")
	}

	if !pkg.Compare(user.Password, password) {
		return nil, "", errors.New("密码错误")
	}

	token, err := pkg.GenerateToken(strconv.FormatInt(user.ID, 10))
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}
func CreateUser(userParam *model.User) (*model.User, string, error) {

	if userParam.Username == "" || userParam.Password == "" {
		return nil, "", errors.New("用户名或密码不能为空")
	}
	// tx := pkg.DB.Where("username = ?", userParam.Username).First(&result)

	user, userErr := query.User.Where(query.User.Username.Eq(userParam.Username)).First()

	if userErr != nil {
		return nil, "", userErr
	}
	if user.Username != "" {
		return nil, "", errors.New("用户名已存在")
	}
	encryptPwd, err := pkg.Encrypt(userParam.Password)
	if err != nil {
		return nil, "", err
	}
	userParam.Password = encryptPwd

	createErr := query.User.Create(userParam)
	if createErr != nil {
		return nil, "", createErr
	}

	token, err := pkg.GenerateToken(strconv.FormatInt(userParam.ID, 10))
	if err != nil {
		return nil, "", err
	}

	return userParam, token, nil
}

func GetUserById(id int) (*model.User, error) {
	user, err := query.User.Where(query.User.ID.Eq(int64(id))).First()
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, errors.New("用户名不存在")
	}
	return user, nil
}

func GetAllUsers() ([]*model.User, error) {
	users, err := query.User.Find()
	if err != nil {
		return nil, err
	}
	return users, nil
}
