package configuration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Clearenv()
	os.Exit(m.Run())
}

func TestConfigurationManager_GetRequiredValue_ReturnsError_WhenNoValueGetterIsConfigured(t *testing.T) {
	configManager := New()

	val, err := configManager.GetRequiredValue("SOME_KEY")

	assert.NotNil(t, err)
	assert.Equal(t, val, "")
}

func TestConfigurationManager_GetRequiredValue_LoadsConfigurationValue_WhenValueExists(t *testing.T) {
	const CONFIG_KEY string = "CONFIG_KEY"
	const CONFIG_VALUE string = "config value"

	configManager := New(
		WithInMemoryGetter(
			map[string]string{
				CONFIG_KEY: CONFIG_VALUE,
			},
		),
	)

	val, err := configManager.GetRequiredValue(CONFIG_KEY)

	assert.Nil(t, err)
	assert.Equal(t, val, CONFIG_VALUE)
}

func TestConfigurationManager_GetRequiredValue_ReturnsError_WhenValueDoesNotExists(t *testing.T) {
	const CONFIG_KEY string = "CONFIG_KEY"
	const CONFIG_VALUE string = "config value"

	os.Setenv("SOME_OTHER_KEY", "some other value")

	configManager := New(
		WithInMemoryGetter(
			map[string]string{
				CONFIG_KEY: CONFIG_VALUE,
			},
		),
		WithEnvironmentVariables(),
	)

	val, err := configManager.GetRequiredValue("YET_SOME_ANOTHER_CONFIG")

	assert.NotNil(t, err)
	assert.Equal(t, "", val)
}

func TestConfigurationManager_GetRequiredValue_LoadsConfigurationValue_RespectingConfigurationValueGetterOrderPrecedence(t *testing.T) {
	type GetterType string
	const (
		InMemoryGetterType GetterType = "InMemory"
		EnvVarsGetterType  GetterType = "EnvVarsGetteerType"
	)

	for _, testCase := range []struct {
		Name           string
		InMemoryConfig map[string]string
		EnvVarsConfig  map[string]string
		ComesFirst     GetterType
		FetchKey       string
		ExpectedValue  string
	}{
		{
			Name: "gets first configuration value found from multiple configuration getters when the same key is present in more than one configuration geteters",
			InMemoryConfig: map[string]string{
				"SOME_CONFIG": "SOME_VALUE_1",
			},
			EnvVarsConfig: map[string]string{
				"SOME_CONFIG": "SOME_VALUE_2",
			},
			ComesFirst:    EnvVarsGetterType,
			FetchKey:      "SOME_CONFIG",
			ExpectedValue: "SOME_VALUE_2",
		},
		{
			Name: "gets the value in the second configuration getter when value is not present in the first",
			InMemoryConfig: map[string]string{
				"SOME_CONFIG": "SOME_VALUE",
			},
			EnvVarsConfig: map[string]string{
				"SOME_WANTED_CONFIG": "SOME_WANTED_VALUE",
			},
			ComesFirst:    InMemoryGetterType,
			FetchKey:      "SOME_WANTED_CONFIG",
			ExpectedValue: "SOME_WANTED_VALUE",
		},
	} {

		t.Run(testCase.Name, func(t *testing.T) {
			// Setting the environment variables
			for k, v := range testCase.EnvVarsConfig {
				os.Setenv(k, v)
			}

			var (
				firstConfig  configurationValueGetter
				secondConfig configurationValueGetter
			)

			if testCase.ComesFirst == EnvVarsGetterType {
				firstConfig = WithEnvironmentVariables()
				secondConfig = WithInMemoryGetter(testCase.InMemoryConfig)
			} else {
				firstConfig = WithInMemoryGetter(testCase.InMemoryConfig)
				secondConfig = WithEnvironmentVariables()
			}

			configManager := New(
				firstConfig,
				secondConfig,
			)

			val, err := configManager.GetRequiredValue(testCase.FetchKey)
			assert.Nil(t, err)
			assert.Equal(t, testCase.ExpectedValue, val)

			os.Clearenv()
		})
	}
}

func TestConfigurationManager_GetOptionalValue_ReturnsEmptyString_WhenNoValueGetterIsConfigured(t *testing.T) {
	configManager := New()

	val := configManager.GetOptionalValue("SOME_KEY")

	assert.Equal(t, val, "")
}

func TestConfigurationManager_GetOptionalValue_LoadsConfigurationValue_WhenValueExists(t *testing.T) {
	const CONFIG_KEY string = "CONFIG_KEY"
	const CONFIG_VALUE string = "config value"

	configManager := New(
		WithInMemoryGetter(
			map[string]string{
				CONFIG_KEY: CONFIG_VALUE,
			},
		),
	)

	val := configManager.GetOptionalValue(CONFIG_KEY)

	assert.Equal(t, val, CONFIG_VALUE)
}
