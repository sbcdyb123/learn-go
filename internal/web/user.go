package web

import (
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sbcdyb123/learn-go/internal/domain"
	"github.com/sbcdyb123/learn-go/internal/service"
	"net/http"
	"time"
)

const biz = "login"

var _ handler = (*UserHandler)(nil)

type UserHandler struct {
	svc         *service.UserService
	codeSvc     *service.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService, codeSvc *service.CodeService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		codeSvc:     codeSvc,
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
	ug.POST("/login_sms/code/send", u.SendLoginSmsCode)
	ug.POST("/login_sms", u.LoginSms)
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
	if errors.Is(err, service.ErrUserDuplicate) {
		c.String(http.StatusOK, "邮箱已存在")
		return
	}

	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	c.String(http.StatusOK, "注册成功")
}
func (u *UserHandler) Login(c *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := c.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(c.Request.Context(), req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		c.String(http.StatusOK, "用户名或密码错误")
		return
	}
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	err, tokenStr := u.setJwtToken(c, user.Id)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	fmt.Println(tokenStr)
	c.JSON(http.StatusOK, gin.H{
		"token": tokenStr,
	})

}

func (u *UserHandler) setJwtToken(c *gin.Context, uid int64) (error, string) {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		UserId: uid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return err, ""
	}

	c.Header("x-jwt-token", tokenStr)
	return nil, tokenStr
}
func (u *UserHandler) Edit(c *gin.Context) {
	type EditReq struct {
		Id       int64  `json:"id"`
		Username string `json:"username"`
		//	生日
		BirthDay string `json:"birthDay"`
		// 个人简介
		Intro string `json:"intro"`
	}
	var req EditReq
	if err := c.Bind(&req); err != nil {
		return
	}
	birthDay, err := time.Parse("2006-01-02", req.BirthDay)
	if err != nil {
		c.String(http.StatusOK, "生日格式错误")
		return
	}
	err = u.svc.Edit(c.Request.Context(), domain.User{
		Id:       req.Id,
		Username: req.Username,
		BirthDay: birthDay.UnixMilli(),
		Intro:    req.Intro,
	})
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	c.String(http.StatusOK, "修改成功")
}
func (u *UserHandler) Profile(c *gin.Context) {
	cl, ok := c.Get("claims")
	if !ok {
		c.String(http.StatusOK, "系统错误")
		return
	}
	claims, ok := cl.(*UserClaims)
	if !ok {
		c.String(http.StatusOK, "系统错误")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"userId": claims.UserId,
	})
}
func (u *UserHandler) SendLoginSmsCode(c *gin.Context) {
	type SmsCodeRequest struct {
		Phone string `json:"phone"`
	}
	var req SmsCodeRequest
	//bind方法会根据content-type来解析你的数据到req里面
	//解析错误，直接就会返回4xx错误
	if err := c.Bind(&req); err != nil {
		return
	}
	if req.Phone == "" {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "手机号不能为空",
		})
		return
	}
	err := u.codeSvc.Send(c, biz, req.Phone)
	switch {
	case err == nil:
		c.JSON(http.StatusOK, Result{
			//Code: 0,
			Msg: "发送验证码成功",
		})
		return
	case errors.Is(err, service.ErrCacheTooFrequently):
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "发送过于频繁，请稍后再试",
		})
		return
	default:
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}

}
func (u *UserHandler) LoginSms(c *gin.Context) {
	type Request struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Request
	//bind方法会根据content-type来解析你的数据到req里面
	//解析错误，直接就会返回4xx错误
	if err := c.Bind(&req); err != nil {
		return
	}
	ok, err := u.codeSvc.Verify(c, biz, req.Phone, req.Code)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "验证码错误",
		})
		return
	}
	user, err := u.svc.FindOrCreate(c, req.Phone)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if err, _ = u.setJwtToken(c, user.Id); err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	c.JSON(http.StatusOK, Result{
		Msg: "验证码校验成功",
	})
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserId int64
}
