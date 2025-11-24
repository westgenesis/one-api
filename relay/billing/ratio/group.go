package ratio

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"

	"github.com/westgenesis/one-api/common/logger"
)

const ratioFileName = "group_ratio.json"

var groupRatioLock sync.RWMutex
var GroupRatio = map[string]float64{
	"default": 1,
	"vip":     1,
	"svip":    1,
}

// init 函数在包被导入时自动执行，用于从文件加载初始数据。
func init() {
	if err := loadGroupRatioFromFile(); err != nil {
		// 如果加载失败，则使用上面定义的默认值，并尝试将默认值保存到文件
		if err := saveGroupRatioToFile(); err != nil {
			logger.SysError("failed to save default group ratio: " + err.Error())
		}
	} else {
		logger.SysLog("group ratio loaded from " + ratioFileName)
	}
}

// --- 文件操作核心函数 ---

// loadGroupRatioFromFile 从本地文件加载 GroupRatio。
func loadGroupRatioFromFile() error {
	groupRatioLock.Lock()
	defer groupRatioLock.Unlock()

	// 检查文件是否存在
	if _, err := os.Stat(ratioFileName); os.IsNotExist(err) {
		// 文件不存在，不认为是错误，使用内存中的默认值
		return nil
	}

	data, err := ioutil.ReadFile(ratioFileName)
	if err != nil {
		return err
	}

	// 重新初始化 GroupRatio 以确保所有旧数据被替换
	GroupRatio = make(map[string]float64)
	return json.Unmarshal(data, &GroupRatio)
}

// saveGroupRatioToFile 将 GroupRatio 保存到本地文件。
// 这个函数应该在 groupRatioLock.Lock() 保护下调用。
func saveGroupRatioToFile() error {
	// 注意：这里我们假设调用方（如 CreateGroup/DeleteGroup）已经持有了 Lock。
	// 但为了健壮性，我们在这里重新获取，以防在其他地方调用。
	groupRatioLock.RLock()
	defer groupRatioLock.RUnlock()

	jsonBytes, err := json.MarshalIndent(GroupRatio, "", "  ")
	if err != nil {
		return err
	}

	// 使用 0644 权限写入文件
	err = ioutil.WriteFile(ratioFileName, jsonBytes, 0644)
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
	groupRatioLock.Lock()
	defer groupRatioLock.Unlock()
	GroupRatio = make(map[string]float64)
	return json.Unmarshal([]byte(jsonStr), &GroupRatio)
}

func GetGroupRatio(name string) float64 {
	groupRatioLock.RLock()
	defer groupRatioLock.RUnlock()
	ratio, ok := GroupRatio[name]
	if !ok {
		logger.SysError("group ratio not found: " + name)
		return 1
	}
	return ratio
}
