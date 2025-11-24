package ratio

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/westgenesis/one-api/common/logger"
)

const ratioFileName = "group_ratio.json"

var GroupRatioLock sync.RWMutex
var GroupRatio = map[string]float64{
	"default": 1,
	"vip":     1,
	"svip":    1,
}

// init 函数在包被导入时自动执行，用于从文件加载初始数据。
func init() {
	if err := loadGroupRatioFromFile(); err != nil {
		// 如果加载失败，则使用上面定义的默认值，并尝试将默认值保存到文件
		if err := SaveGroupRatioToFile(); err != nil {
			logger.SysError("failed to save default group ratio: " + err.Error())
		}
	} else {
		logger.SysLog("group ratio loaded from " + ratioFileName)
	}
}

// --- 文件操作核心函数 ---

// loadGroupRatioFromFile 从本地文件加载 GroupRatio。
func loadGroupRatioFromFile() error {
	// 不加锁，内部 SaveGroupRatioToFile 已加锁
	if _, err := os.Stat(ratioFileName); os.IsNotExist(err) {
		GroupRatio = map[string]float64{
			"default": 1,
			"vip":     1,
			"svip":    1,
		}
		if err := SaveGroupRatioToFile(); err != nil {
			logger.SysError("failed to create default group ratio file: " + err.Error())
			return err
		}
		logger.SysLog("created default group ratio file: " + ratioFileName)
		return nil
	}

	data, err := os.ReadFile(ratioFileName)
	if err != nil {
		return err
	}

	GroupRatioLock.Lock()
	defer GroupRatioLock.Unlock()
	GroupRatio = make(map[string]float64)
	return json.Unmarshal(data, &GroupRatio)
}

// saveGroupRatioToFile 将 GroupRatio 保存到本地文件。
// 这个函数应该在 groupRatioLock.Lock() 保护下调用。
func SaveGroupRatioToFile() error {
	GroupRatioLock.RLock()
	defer GroupRatioLock.RUnlock()

	jsonBytes, err := json.MarshalIndent(GroupRatio, "", "  ")
	if err != nil {
		return err
	}

	// 使用 0644 权限写入文件
	err = os.WriteFile(ratioFileName, jsonBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GroupRatio2JSONString() string {
	jsonBytes, err := json.Marshal(GroupRatio)
	if err != nil {
		logger.SysError("error marshalling model ratio: " + err.Error())
	}
	return string(jsonBytes)
}

func UpdateGroupRatioByJSONString(jsonStr string) error {
	GroupRatioLock.Lock()
	defer GroupRatioLock.Unlock()
	GroupRatio = make(map[string]float64)
	return json.Unmarshal([]byte(jsonStr), &GroupRatio)
}

func GetGroupRatio(name string) float64 {
	GroupRatioLock.RLock()
	defer GroupRatioLock.RUnlock()
	ratio, ok := GroupRatio[name]
	if !ok {
		logger.SysError("group ratio not found: " + name)
		return 1
	}
	return ratio
}
