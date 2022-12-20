package configuration

import (
	"fmt"
	"os"
)

type ConfigurationManager struct {
	configurationValueGetter []configurationValueGetter
}

func New(opts ...configurationValueGetter) *ConfigurationManager {
	configManager := &ConfigurationManager{
		configurationValueGetter: []configurationValueGetter{},
	}
	for _, getter := range opts {
		configManager.configurationValueGetter = append(configManager.configurationValueGetter, getter)
	}
	return configManager
}

func (m *ConfigurationManager) GetOptionalValue(key string, defaultValue string) string {
	val, err := m.GetRequiredValue(key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (m *ConfigurationManager) GetRequiredValue(key string) (string, error) {
	if len(m.configurationValueGetter) == 0 {
		return "", fmt.Errorf("there is no configuration getter configured")
	}
	for _, valueGetter := range m.configurationValueGetter {
		val, err := valueGetter(key)
		if err == nil {
			return val, nil
		}
	}
	return "", fmt.Errorf("key %s not found", key)
}

type configurationValueGetter func(key string) (string, error)

// ----------------------------------------------------------------
// ConfigurationManager Options
// ----------------------------------------------------------------

func WithCustomGetter(fnGetter configurationValueGetter) configurationValueGetter {
	return fnGetter
}

func WithInMemoryGetter(mp map[string]string) configurationValueGetter {
	return func(key string) (string, error) {
		val, exists := mp[key]
		if !exists {
			return "", fmt.Errorf("value: %s does not exists in the map", key)
		}
		return val, nil
	}
}

func WithEnvironmentVariables() configurationValueGetter {
	return func(key string) (string, error) {
		val, ok := os.LookupEnv(key)
		if !ok {
			return "", fmt.Errorf("value: %s does not exsits in environment variables", key)
		}
		return val, nil
	}
}
