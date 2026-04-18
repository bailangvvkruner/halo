package httpserver

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func parseUintParam(c *gin.Context, name string) (uint, error) {
	var id uint
	if _, err := fmt.Sscan(c.Param(name), &id); err != nil {
		return 0, fmt.Errorf("invalid %s", name)
	}
	return id, nil
}
