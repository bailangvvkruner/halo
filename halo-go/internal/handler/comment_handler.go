package handler

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/halo-dev/halo-go/internal/model"
	"github.com/halo-dev/halo-go/internal/service"
)

type CommentHandler struct {
	svc service.CommentService
}

func NewCommentHandler(svc service.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

func (h *CommentHandler) Create(c *gin.Context) {
	var comment model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error(), "data": nil})
		return
	}

	if h.isSpam(&comment) {
		comment.Spec.Approved = false
		comment.Spec.Reason = "疑似垃圾评论"
	} else if h.shouldAutoApprove() {
		now := time.Now()
		comment.Spec.Approved = true
		comment.Spec.ApproveTime = &now
	}

	created, err := h.svc.Create(c.Request.Context(), &comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建评论失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": created})
}

func (h *CommentHandler) Get(c *gin.Context) {
	name := c.Param("name")
	comment, err := h.svc.Get(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "评论不存在: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": comment})
}

func (h *CommentHandler) Update(c *gin.Context) {
	name := c.Param("name")
	var comment model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error(), "data": nil})
		return
	}
	comment.Metadata.Name = name
	updated, err := h.svc.Update(c.Request.Context(), &comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新评论失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": updated})
}

func (h *CommentHandler) Delete(c *gin.Context) {
	name := c.Param("name")
	if err := h.svc.Delete(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除评论失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": nil})
}

func (h *CommentHandler) List(c *gin.Context) {
	opts := parseListOptions(c)
	result, err := h.svc.List(c.Request.Context(), &opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询评论列表失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": result})
}

func (h *CommentHandler) Approve(c *gin.Context) {
	name := c.Param("name")
	comment, err := h.svc.Get(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "评论不存在", "data": nil})
		return
	}
	comment.Spec.Approved = true
	now := time.Now()
	comment.Spec.ApproveTime = &now
	updated, err := h.svc.Update(c.Request.Context(), comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "审批失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "审批成功", "data": updated})
}

func (h *CommentHandler) Reply(c *gin.Context) {
	name := c.Param("name")
	var reply model.Reply
	if err := c.ShouldBindJSON(&reply); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error(), "data": nil})
		return
	}
	reply.Spec.CommentName = name
	reply.Spec.AllowNotification = true

	if h.isSpamReply(&reply) {
		reply.Spec.Approved = false
		reply.Spec.Reason = "疑似垃圾回复"
	} else if h.shouldAutoApprove() {
		reply.Spec.Approved = true
		now := time.Now()
		reply.Spec.ApproveTime = &now
	}

	created, err := h.svc.CreateReply(c.Request.Context(), &reply)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建回复失败: " + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": created})
}

var spamPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(viagra|casino|lottery|porn|xxx|sex|gambling)`),
	regexp.MustCompile(`(?i)(free money|click here|best price|act now|limited time)`),
	regexp.MustCompile(`https?://\S+`),
	regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
}

func (h *CommentHandler) isSpam(comment *model.Comment) bool {
	content := strings.ToLower(comment.Spec.RawContent + " " + comment.Spec.IP)
	for _, pattern := range spamPatterns {
		if pattern.MatchString(content) {
			return true
		}
	}
	return len(comment.Spec.RawContent) < 5 || len(comment.Spec.RawContent) > 2000
}

func (h *CommentHandler) isSpamReply(reply *model.Reply) bool {
	content := strings.ToLower(reply.Spec.RawContent)
	for _, pattern := range spamPatterns {
		if pattern.MatchString(content) {
			return true
		}
	}
	return len(reply.Spec.RawContent) < 3 || len(reply.Spec.RawContent) > 1000
}

func (h *CommentHandler) shouldAutoApprove() bool {
	return true
}
