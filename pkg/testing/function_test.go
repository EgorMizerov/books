package testing

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type FuncsTests struct {
	suite.Suite
	testSuite       *suite.Suite
	exampleFunction func()
}

func TestFuncs(t *testing.T) {
	suite.Run(t, new(FuncsTests))
}

func (self *FuncsTests) SetupTest() {
	self.testSuite = &suite.Suite{}
	self.testSuite.SetT(&testing.T{})
}

func (self *FuncsTests) TestFuncsEqual() {
	FuncsEqual(self.T(), self.exampleFunction, self.exampleFunction)

	self.False(self.testSuite.T().Failed())
}

func (self *FuncsTests) TestFuncsUnequal() {
	FuncsEqual(self.testSuite.T(), self.exampleFunction, func() {})

	self.True(self.testSuite.T().Failed())
}
