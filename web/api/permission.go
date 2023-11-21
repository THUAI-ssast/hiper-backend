package api

import (
	"hiper-backend/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func grant_creation_permission(c *gin.Context, author_ids string) {
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
	_, valid := model.SelectMySql("user", map[string]interface{}{"user_id": author_id})
	if !valid {
		c.JSON(404, gin.H{})
	} else {
		model.UpdateMySQL("user", map[string]interface{}{"authorization": "Secondary Admin"}, map[string]interface{}{"user_id": author_id})
		c.JSON(200, gin.H{})
	}
}

func revoke_creation_permission(c *gin.Context, author_ids string) {
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
	_, valid := model.SelectMySql("user", map[string]interface{}{"user_id": author_id})
	if !valid {
		c.JSON(404, gin.H{})
	} else {
		model.UpdateMySQL("user", map[string]interface{}{"authorization": "Regular user"}, map[string]interface{}{"user_id": author_id})
	}
}
