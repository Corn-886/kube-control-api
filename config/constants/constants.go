package constants

import (
	"os"
	"path/filepath"
)

/**
获取环境变量值，如果空则返回默认值
*/
func GetEnvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

/**
当前目录下 + /config + /resource
*/
func GET_CONFIG_RESOURCE_DIR() string {
	return filepath.Join(GET_CONFIG_DIR(), "resource")
}

/**
当前目录下 + /config
*/
func GET_CONFIG_DIR() string {
	dir, _ := os.Getwd()
	configDir := filepath.Join(dir, "config")
	return configDir
}

func GET_KUBE_VERSION() string {
	return "spray-v2.21.0c_k8s-v1.26.4_v4.4-amd64"
}

func GET_KUBE_DOCKER_IMAGE_ADDRESS() string {
	return "registry.cn-shanghai.aliyuncs.com/kuboard-spray/kuboard-spray-resource"
}

func GET_DATA_DIR() string {
	dir, _ := os.Getwd()
	dataDir := filepath.Join(dir, "data")
	return dataDir
}
