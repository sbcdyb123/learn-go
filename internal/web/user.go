package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/sbcdyb123/learn-go/internal/domain"
	"github.com/sbcdyb123/learn-go/internal/service"
	"net/http"
)

type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/user")
	ug.POST("/signup", u.Signup)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
}

func (u *UserHandler) Signup(c *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}
	var req SignUpReq
	//bind方法会根据content-type来解析你的数据到req里面
	//解析错误，直接就会返回4xx错误
	if err := c.Bind(&req); err != nil {
		return
	}

	ok, err := u.emailExp.MatchString(req.Email)
	//ok, err := regexp.Match(emailRegexPattern, []byte(req.Email))
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		c.String(http.StatusOK, "邮箱格式错误")
		return
	}
	if req.Password != req.ConfirmPassword {
		c.String(http.StatusOK, "两次输入的密码不一致")
		return
	}
	ok, err = u.passwordExp.MatchString(req.Password)
	//ok, err = regexp.Match(passwordRegexPattern, []byte(req.Password))
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		c.String(http.StatusOK, "密码必须大于8位，且包含大小写字母、数字和特殊符号")
	}
	fmt.Printf("req:%+v\n", req)
	err = u.svc.SignUp(c.Request.Context(), domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	c.String(http.StatusOK, "注册成功")
}
func (u *UserHandler) Login(c *gin.Context) {

}
func (u *UserHandler) Edit(c *gin.Context) {

}
func (u *UserHandler) Profile(c *gin.Context) {

}
