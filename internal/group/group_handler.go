package group

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type groupHandler struct {
	service *groupService
}

func NewGroupHandler(service *groupService) *groupHandler {
	return &groupHandler{service: service}
}

func (h *groupHandler) GetGroupProgress(c *gin.Context) {
	// 1. Get userID from JWT context
	userID := c.GetInt("userID")

	// 2. Get groupID from URL parameter
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	// 3. Get group progress
	progress, err := h.service.GetGroupProgress(groupID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"progress": progress,
	})
}

func (h *groupHandler) GetMemberContribution(c *gin.Context) {
	// 1. Get userID from JWT context
	userID := c.GetInt("userID")

	// 2. Get groupID from URL parameter
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	// 3. Get contribution
	contribution, err := h.service.GetMemberContribution(groupID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"group_id":          groupID,
		"user_id":           userID,
		"user_contribution": contribution,
	})
}
