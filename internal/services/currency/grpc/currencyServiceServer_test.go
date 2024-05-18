package grpc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CurrencyServiceServerUnitTestSuite struct {
	suite.Suite
	StartingNumber int
}

// this function executes before the test suite begins execution
func (suite *CurrencyServiceServerUnitTestSuite) SetupSuite() {
	// set StartingNumber to one
	fmt.Println(">>> From SetupSuite")
	suite.StartingNumber = 1
}

// this function executes after all tests executed
func (suite *CurrencyServiceServerUnitTestSuite) TearDownSuite() {
	fmt.Println(">>> From TearDownSuite")
}

// this function executes before each test case
func (suite *CurrencyServiceServerUnitTestSuite) SetupTest() {
	// reset StartingNumber to one
	fmt.Println("-- From SetupTest")
	suite.StartingNumber = 1
}

// this function executes after each test case
func (suite *CurrencyServiceServerUnitTestSuite) TearDownTest() {
	fmt.Println("-- From TearDownTest")
}

func (suite *CurrencyServiceServerUnitTestSuite) TestAddOne() {
	fmt.Println("From TestAddOne")
	suite.StartingNumber += 1
	suite.Equal(2, suite.StartingNumber)
}

func (suite *CurrencyServiceServerUnitTestSuite) TestSubtractOne() {
	fmt.Println("From TestSubtractOne")
	suite.StartingNumber -= 1
	suite.Equal(0, suite.StartingNumber)
}

func TestCurrencyServiceServerTestSuite(t *testing.T) {
	suite.Run(t, new(CurrencyServiceServerUnitTestSuite))
}
