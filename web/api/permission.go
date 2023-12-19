package api

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"hiper-backend/model"
)

func grantCreationPermission(c *gin.Context) {
	author_ids := c.Param("user_id")
	userID := c.MustGet("userID").(int)
	author_id, err := strconv.Atoi(author_ids)
	if err != nil {
		c.JSON(401, gin.H{})
	}
	if userID != 1 || author_id == 1 {
		c.JSON(401, gin.H{})
	}
	_, err = model.GetUserById((uint)(author_id))
	if err != nil {
		c.JSON(404, gin.H{})
	}
	err = model.UpdateUserById((uint)(author_id), map[string]interface{}{"can_create_game_or_contest": true})
	if err != nil {
		c.JSON(500, gin.H{})
	}
	c.JSON(200, gin.H{})
}

func revokeCreationPermission(c *gin.Context) {
	author_ids := c.Param("user_id")
	userID := c.MustGet("userID").(int)
	author_id, err := strconv.Atoi(author_ids)
	if err != nil {
		c.JSON(401, gin.H{})
	}
	if userID != 1 || author_id == 1 {
		c.JSON(401, gin.H{})
	}
	_, err = model.GetUserById((uint)(author_id))
	if err != nil {
		c.JSON(404, gin.H{})
	}
	err = model.UpdateUserById((uint)(author_id), map[string]interface{}{"can_create_game_or_contest": false})
	if err != nil {
		c.JSON(500, gin.H{})
	}
	c.JSON(200, gin.H{})
}
