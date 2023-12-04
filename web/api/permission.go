package api

import (
	"hiper-backend/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func grantCreationPermission(c *gin.Context, author_ids string) {
	userID := c.MustGet("userID").(int)
	author_id, err := strconv.Atoi(author_ids)
	if err != nil {
		c.JSON(401, gin.H{})
		return
	}
	if userID != 1 || author_id == 1 {
		c.JSON(401, gin.H{})
		return
	}
	_, err = model.GetUserById((uint)(author_id))
	if err != nil {
		c.JSON(404, gin.H{})
	} else {
		err = model.UpdateUserById((uint)(author_id), map[string]interface{}{"authorization": "Secondary Admin"})
		if err != nil {
			c.JSON(500, gin.H{})
			return
		}
		c.JSON(200, gin.H{})
	}
}

func revokeCreationPermission(c *gin.Context, author_ids string) {
	userID := c.MustGet("userID").(int)
	author_id, err := strconv.Atoi(author_ids)
	if err != nil {
		c.JSON(401, gin.H{})
		return
	}
	if userID != 1 || author_id == 1 {
		c.JSON(401, gin.H{})
		return
	}
	_, err = model.GetUserById((uint)(author_id))
	if err != nil {
		c.JSON(404, gin.H{})
	} else {
		err = model.UpdateUserById((uint)(author_id), map[string]interface{}{"authorization": "Regular User"})
		if err != nil {
			c.JSON(500, gin.H{})
			return
		}
		c.JSON(200, gin.H{})
	}
}
