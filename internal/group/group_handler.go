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

type CreateGroupRequest struct {
	Name         string  `json:"name"          binding:"required"`
	TargetAmount float64 `json:"target_amount" binding:"required,gt=0"`
}

func (h *groupHandler) CreateGroup(c *gin.Context) {

	userID := c.GetInt("userID")

	//Parse request body
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Create group
	group, err := h.service.CreateGroup(req.Name, req.TargetAmount, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "group created successfully",
		"group":   group,
	})
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

	userID := c.GetInt("userID")

	// Get groupID from URL parameter
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	//Get contribution
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