package handler

import (
	"fmt"
	"gostart/internal/model"
	"gostart/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary      登录
// @Description  登录
// @Tags         用户管理
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        username  formData string true "用户名"
// @Param        password  formData string true "用户名"
// @Success      200   {object}   model.LoginResp  "登录成功"
// @Failure      400   {object}  model.Base  "请求参数错误"
// @Router       /api/login [post]
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	userDO, tokenDO, err := service.Login(username, password)
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	loginResp := &model.LoginResp{
		ID:           userDO.ID,
		Username:     userDO.Username,
		Email:        userDO.Email,
		Age:          userDO.Age,
		AccessToken:  tokenDO.AccessToken,
		RefreshToken: tokenDO.RefreshToken,
	}
	OK(c, loginResp)
}

// @Summary      登出
// @Description  登出
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Success      200   {string}   string  "登出成功"
// @Failure      400   {object}  model.Base  "请求参数错误"
// @Router       /api/logout [post]
func Logout(c *gin.Context) {

}

// @Summary      创建用户
// @Description  创建用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        username  formData string true "用户名"
// @Param        password  formData string true "用户密码"
// @Success      200   {object}   model.LoginResp  "注册成功"
// @Failure      400   {object}  model.Base  "请求参数错误"
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

	loginResp := &model.LoginResp{
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
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @Success      200   {object}   []model.UserResp
// @Failure      400   {object}  model.Base
// @Router       /api/user [get]
func GetUsers(c *gin.Context) {
	ret, err := service.GetAllUsers()
	var userResps = []*model.UserResp{}
	for _, v := range ret {
		userResps = append(userResps, &model.UserResp{
			ID:       v.ID,
			Username: v.Username,
			Email:    v.Email,
			Age:      v.Age,
		})
	}
	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	OK(c, userResps)
}

// @Summary      获取单个用户
// @Description  获取单个用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        id path int  true "用户id"
// @Security     ApiKeyAuth
// @Success      200   {object}  model.UserResp
// @Failure      400   {object}  model.Base
// @Router       /api/user/{id} [get]
func GetUser(c *gin.Context) {
	intId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Fail(c, http.StatusBadRequest, "id can not be empty")
	}

	userId, _ := c.Get("userId")
	fmt.Println("token中获取的userId : ", userId)

	ret, err := service.GetUserById(intId)
	var userResp = model.UserResp{
		ID:       ret.ID,
		Username: ret.Username,
		Email:    ret.Email,
		Age:      ret.Age,
	}

	if err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	OK(c, userResp)
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
