package handler

import (
	"fmt"
	"gostart/internal/model"
	"gostart/internal/pkg/response"
	"gostart/internal/service"
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

// @Summary      登录
// @Description  登录
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        username  formData string true "用户名"
// @Param        password  formData string true "用户名"
// @Success      200   {object}   response.LoginResp  "登录成功"
// @Failure      400   {string}  string  "请求参数错误"
// @Router       /api/login [post]
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

// @Summary      创建用户
// @Description  创建用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        username  formData string true "用户名"
// @Param        password  formData string true "用户密码"
// @Success      200   {object}   response.LoginResp  "注册成功"
// @Failure      400   {string}  string  "请求参数错误"
// @Router       /api/regist [post]
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

// @Summary      获取用户列表
// @Description  获取用户列表
// @Tags         用户管理
// @Produce      json
// @Success      200   {object}   []model.UserDO
// @Failure      400   {string}  string
// @Router       /api/user [get]
func GetUsers(c *gin.Context) {
	ret, err := service.GetAllUsers()
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	OK(c, ret)
}

// @Summary      获取单个用户
// @Description  获取单个用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id path int  true "用户id"
// @Success      200   {object}   []model.UserDO
// @Failure      400   {string}  string
// @Router       /api/user [get]
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
