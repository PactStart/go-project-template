package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

const (
	FileName          = "config.yaml"
	DefaultFolderPath = "../config/"
	Version           = "1.0.0"
)

// return absolude path join ../config/, this is k8s container config path
func GetDefaultConfigPath() string {
	b, err := filepath.Abs(os.Args[0])
	if err != nil {
		fmt.Println("filepath.Abs error,err=", err)
		return ""
	}
	return filepath.Join(filepath.Dir(b), DefaultFolderPath)
}

// getProjectRoot returns the absolute path of the project root directory
func GetProjectRoot() string {
	// 获取当前可执行文件所在的路径
	rootDir, _ := os.Getwd()
	return rootDir
}

func initConfig(config interface{}, configName, configFolderPath string) error {
	configFolderPath = filepath.Join(configFolderPath, configName)
	_, err := os.Stat(configFolderPath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("stat config path error:", err.Error())
			return fmt.Errorf("stat config path error: %w", err)
		}
		configFolderPath = filepath.Join(GetProjectRoot(), "config", configName)
		fmt.Println("flag's path,environment's path,default path all is not exist,using project path:", configFolderPath)
	}
	data, err := os.ReadFile(configFolderPath)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}
	if err = yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("unmarshal yaml error: %w", err)
	}
	fmt.Println("use config", configFolderPath)

	return nil
}

func InitConfig(configFolderPath string) error {
	if configFolderPath == "" {
		envConfigPath := os.Getenv("XXXJZ_CONFIG")
		if envConfigPath != "" {
			configFolderPath = envConfigPath
		} else {
			configFolderPath = GetDefaultConfigPath()
		}
	}
	return initConfig(&Config, FileName, configFolderPath)
}
