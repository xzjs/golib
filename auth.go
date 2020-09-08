package lib

import "github.com/gin-gonic/gin"

// GetUserID 获取用户id
func GetUserID(c *gin.Context) uint {
	temp, exist := c.Get("userID")
	if exist {
		return temp.(uint)
	}
	return 0

}
