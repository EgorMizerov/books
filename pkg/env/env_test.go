package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EnvTest struct {
	suite.Suite
	envKey string
}

func TestEnv(t *testing.T) {
	suite.Run(t, new(EnvTest))
}

func (self *EnvTest) SetupTest() {
	self.envKey = "ENV_KEY"
}

func (self *EnvTest) TearDownTest() {
	os.Clearenv()
}

func (self *EnvTest) TestGetString() {
	envValue := "value"
	err := os.Setenv(self.envKey, envValue)
	self.NoError(err)

	result := GetString(self.envKey, "")

	self.Equal(envValue, result)
}

func (self *EnvTest) TestGetStringDefault() {
	defaultValue := "default"

	result := GetString(self.envKey, defaultValue)

	self.Equal(defaultValue, result)
}

func (self *EnvTest) TestGetInt() {
	envValue := "10"
	envValueInt := 10
	err := os.Setenv(self.envKey, envValue)
	self.NoError(err)

	result := GetInt(self.envKey, 0)

	self.Equal(envValueInt, result)
}

func (self *EnvTest) TestGetIntDefault() {
	defaultValue := 20

	result := GetInt(self.envKey, defaultValue)

	self.Equal(defaultValue, result)
}

func (self *EnvTest) TestGetBool() {
	envValue := "true"
	envValueBool := true
	err := os.Setenv(self.envKey, envValue)
	self.NoError(err)

	result := GetBool(self.envKey, false)

	self.Equal(envValueBool, result)
}

func (self *EnvTest) TestGetBoolDefault() {
	defaultValue := true

	result := GetBool(self.envKey, defaultValue)

	self.Equal(defaultValue, result)
}
