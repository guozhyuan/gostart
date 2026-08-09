package handler

import (
	"fmt"
	"myserver/internal/model"
	"myserver/internal/pkg/response"
	"myserver/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, model.Assemble(http.StatusOK, "success", data))
}

func Fail(c *gin.Context, code int, message string) {
	c.JSON(code, model.Assemble(code, message, nil))
}

// 登录
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	userDO, tokenDO, err := service.Login(username, password)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	loginResp := &response.LoginResp{
		ID:           userDO.ID,
		Username:     userDO.Username,
		Email:        userDO.Email,
		Age:          userDO.Age,
		AccessToken:  tokenDO.AccessToken,
		RefreshToken: tokenDO.RefreshToken,
	}
	OK(c, loginResp)
}

func Logout(c *gin.Context) {

}

// 创建用户
func Registe(c *gin.Context) {
	var userParam model.UserDO
	// if err := c.ShouldBindJSON(&userParam); err != nil {
	// 	Fail(c, http.StatusBadRequest, err.Error())
	// 	return
	// }
	userParam.Username = c.PostForm("username")
	userParam.Password = c.PostForm("password")
	ret, tokenDO, err := service.CreateUser(&userParam)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	loginResp := &response.LoginResp{
		ID:           ret.ID,
		Username:     ret.Username,
		Email:        ret.Email,
		Age:          ret.Age,
		AccessToken:  tokenDO.AccessToken,
		RefreshToken: tokenDO.RefreshToken,
	}
	OK(c, loginResp)
}

// 获取用户列表
func GetUsers(c *gin.Context) {
	ret, err := service.GetAllUsers()
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	OK(c, ret)
}

// 获取单个用户
func GetUser(c *gin.Context) {
	intId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Fail(c, http.StatusBadRequest, "id can not be empty")
	}

	userId, _ := c.Get("userId")
	fmt.Println("token中获取的userId : ", userId)

	ret, err := service.GetUserById(intId)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	OK(c, ret)
}

// 更新用户
func UpdateUser(c *gin.Context) {
	/* id := c.Param("id")
	var user model.UserDO
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.Error("user not found"))
		return
	}
	var req model.UserDO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	if err := config.DB.Model(&user).Updates(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.OK(user)) */
}

// 删除用户
func DeleteUser(c *gin.Context) {
	/* id := c.Param("id")
	if err := config.DB.Delete(&model.UserDO{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.OK(nil)) */
}
