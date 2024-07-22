package components

import (
	"context"
	"orderin-server/internal/services"
	"orderin-server/pkg/common/utils"
	"strings"
)

const (
	CONFIG_NAME = "example_config_name"
)

type MyConfig struct {
	ConfigService services.SysConfig
	ConfigMap     map[string]string
}

func (e MyConfig) ReloadAll() {
	configs, err := e.ConfigService.GetAll()
	if err != nil {
		panic(err)
	}
	configMap := make(map[string]string)
	if configs != nil {
		for _, config := range *configs {
			configMap[config.Name] = config.Value
		}
	}
	e.ConfigMap = configMap
}

func (e MyConfig) GetConfigByName(ctx context.Context, name string) (*string, error) {
	value := e.ConfigMap[name]
	return &value, nil
}

func (e MyConfig) GetConfigByNameWithDefaultValue(ctx context.Context, name string, defaultValue string) *string {
	value, err := e.GetConfigByName(ctx, name)
	if value == nil || err != nil || *value == "" {
		return &defaultValue
	}
	return value
}

func (e MyConfig) GetIntConfigByName(ctx context.Context, name string) (*int, error) {
	value, err := e.GetConfigByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	intValue := utils.StringToInt(*value)
	return &intValue, err
}

func (e MyConfig) GetBoolConfigByName(ctx context.Context, name string) (*bool, error) {
	value, err := e.GetConfigByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	var boolValue *bool
	if strings.ToLower(*value) == "true" {
		*boolValue = true
	} else if strings.ToLower(*value) == "false" {
		*boolValue = false
	}
	return boolValue, nil
}
