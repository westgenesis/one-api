package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/westgenesis/one-api/common/logger"
	billingratio "github.com/westgenesis/one-api/relay/billing/ratio"
)

type GroupRequest struct {
	Name  string  `json:"name" binding:"required"`
	Ratio float64 `json:"ratio" binding:"required"`
}

type DeleteGroupRequest struct {
	Name string `json:"name" binding:"required"`
}

func GetGroups(c *gin.Context) {
	groupNames := make([]string, 0)
	for groupName := range billingratio.GroupRatio {
		groupNames = append(groupNames, groupName)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    groupNames,
	})
}

func CreateGroupHandler(c *gin.Context) {
	var req GroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request: " + err.Error(),
		})
		return
	}

	billingratio.GroupRatio[req.Name] = req.Ratio
	if err := billingratio.SaveGroupRatioToFile(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to save group ratio: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "group created/updated",
		"data":    req,
	})
}

func DeleteGroupHandler(c *gin.Context) {
	var req DeleteGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request: " + err.Error(),
		})
		return
	}

	if _, exists := billingratio.GroupRatio[req.Name]; exists {
		delete(billingratio.GroupRatio, req.Name)
		logger.SysLog("group '" + req.Name + "' has been deleted")

		if err := billingratio.SaveGroupRatioToFile(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "failed to save group ratio: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "group deleted",
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "group not found",
		})
	}
}

func GetAllGroupNames() []string {
	billingratio.GroupRatioLock.RLock()
	defer billingratio.GroupRatioLock.RUnlock()

	names := make([]string, 0, len(billingratio.GroupRatio))
	for name := range billingratio.GroupRatio {
		names = append(names, name)
	}
	return names
}
