package middleware

import (
	"net/http"
	"strings"

	"halo/internal/service"

	"github.com/gin-gonic/gin"
)

func PermissionMiddleware(roleService *service.RoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userId")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		uid, ok := userID.(uint)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid user id"})
			return
		}

		path := c.FullPath()
		method := c.Request.Method

		resource, action := parseResourceAction(path, method)
		if resource == "" {
			c.Next()
			return
		}

		if !roleService.HasPermission(uid, resource, action) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "permission denied"})
			return
		}

		c.Next()
	}
}

func parseResourceAction(path, method string) (resource, action string) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return "", ""
	}

	resource = parts[1]
	action = method
	return
}
