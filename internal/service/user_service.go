package service

import (
	"errors"
	"gostart/internal/model"
	"gostart/internal/pkg"
	"strconv"

	"gorm.io/gorm"
)

// handler 层：接收请求、解析参数、调用 service、返回响应（薄层）
// service 层：业务逻辑编排、校验、事务管理、调用 DAO/模型
// model/dao 层：数据持久化（GORM 查询）

func Login(username string, password string) (*model.UserDO, *model.TokenDO, error) {
	if username == "" || password == "" {
		return nil, nil, errors.New("用户名或密码不能为空")
	}
	user, err := GetUserByName(username)
	if err != nil {
		return nil, nil, err
	}
	// 没查到
	if user.Username == "" {
		return nil, nil, errors.New("用户名不存在")
	}

	if !pkg.Compare(user.Password, password) {
		return nil, nil, errors.New("密码错误")
	}

	token, err := pkg.GenerateToken(strconv.FormatUint(user.ID, 10))
	if err != nil {
		return nil, nil, err
	}

	tokenDO := &model.TokenDO{
		AccessToken:  token,
		RefreshToken: "",
	}

	return user, tokenDO, nil
}
func CreateUser(userParam *model.UserDO) (*model.UserDO, *model.TokenDO, error) {

	if userParam.Username == "" || userParam.Password == "" {
		return nil, nil, errors.New("用户名或密码不能为空")
	}

	result, err := GetUserByName(userParam.Username)
	if err != nil {
		return nil, nil, err
	}

	if result.Username != "" {
		return nil, nil, errors.New("用户名已存在")
	}
	encryptPwd, err := pkg.Encrypt(userParam.Password)
	if err != nil {
		return nil, nil, err
	}
	userParam.Password = encryptPwd

	if err := DB.Create(userParam).Error; err != nil {
		return nil, nil, err
	}

	token, err := pkg.GenerateToken(strconv.FormatUint(userParam.ID, 10))
	if err != nil {
		return nil, nil, err
	}

	tokenDO := &model.TokenDO{
		AccessToken:  token,
		RefreshToken: "",
	}
	return userParam, tokenDO, nil
}

func GetUserById(id int) (*model.UserDO, error) {
	var user model.UserDO
	if err := DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByName(username string) (*model.UserDO, error) {
	var user model.UserDO
	err := DB.Where("username = ?", username).First(&user).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &user, nil
}

func GetAllUsers() ([]model.UserDO, error) {
	var users []model.UserDO
	if err := DB.Find(&users).Error; err != nil {
		return users, err
	}
	return users, nil
}
