package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/westgenesis/one-api/common/logger"
	billingratio "github.com/westgenesis/one-api/relay/billing/ratio"
)

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

// CreateGroup 创建或更新分组，并保存到文件。
func CreateGroup(name string, ratio float64) error {
	billingratio.groupRatioLock.Lock()
	defer billingratio.groupRatioLock.Unlock()

	billingratio.GroupRatio[name] = ratio

	return billingratio.saveGroupRatioToFile()
}

// DeleteGroup 删除分组，并保存到文件。
func DeleteGroup(name string) {
	billingratio.groupRatioLock.Lock()
	defer billingratio.groupRatioLock.Unlock()

	if _, exists := billingratio.GroupRatio[name]; exists {
		delete(billingratio.GroupRatio, name)
		logger.SysLog("group '" + name + "' has been deleted")

		if err := billingratio.saveGroupRatioToFile(); err != nil {
			logger.SysError("failed to save group ratio after deletion: " + err.Error())
		}
	}
}

func GetAllGroupNames() []string {
	billingratio.groupRatioLock.RLock()
	defer billingratio.groupRatioLock.RUnlock()

	names := make([]string, 0, len(billingratio.GroupRatio))
	for name := range billingratio.GroupRatio {
		names = append(names, name)
	}
	return names
}
