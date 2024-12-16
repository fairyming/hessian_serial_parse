package utils

import (
	"encoding/json"
)

func EncodeJson(obj interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(obj, "", "  ") // 第二个参数是前缀，第三个参数是缩进
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
