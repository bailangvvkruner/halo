package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/halo-dev/halo-go/internal/middleware"
	"github.com/halo-dev/halo-go/internal/model"
	"github.com/halo-dev/halo-go/internal/service"
)

type AuthHandler struct {
	authSvc service.AuthService
}

func NewAuthHandler(authSvc service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Email       string `json:"email" binding:"required"`
	DisplayName string `json:"displayName"`
}

type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expiresIn"`
}

type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error(), "data": nil})
		return
	}
	token, err := h.authSvc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data":    TokenResponse{Token: token, ExpiresIn: 86400},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error(), "data": nil})
		return
	}
	user := &model.User{
		Spec: model.UserSpec{
			UserName:    req.Username,
			RawPassword: req.Password,
			Email:       req.Email,
			DisplayName: req.DisplayName,
		},
	}
	token, err := h.authSvc.Register(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "注册失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "注册成功",
		"data":    TokenResponse{Token: token, ExpiresIn: 86400},
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error(), "data": nil})
		return
	}
	token, err := h.authSvc.RefreshToken(c.Request.Context(), req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "刷新令牌失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "令牌刷新成功",
		"data":    TokenResponse{Token: token, ExpiresIn: 86400},
	})
}

func (h *AuthHandler) PasswordReset(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "密码重置请求已提交，请查看邮箱",
		"data":    nil,
	})
}

func (h *AuthHandler) CurrentUser(c *gin.Context) {
	username := c.GetString(middleware.ContextKeyUsername)
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未认证", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    gin.H{"username": username},
	})
}
