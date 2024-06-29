package architecture

import (
	"fmt"
	"strings"
	"testing"
)

func NewArchtestErr(data interface{}) *archtestErr {
	str, ok := data.(string)
	if !ok {
		return nil
	}

	lines := strings.Split(str, "\n")
	stack := lines[1 : len(lines)-1]
	if len(stack) < 2 {
		return nil
	}

	return &archtestErr{
		basePackage: strings.TrimSpace(stack[0]),
		dependency:  strings.TrimSpace(stack[1]),
	}
}

type archtestErr struct {
	basePackage string
	dependency  string
}

func (e *archtestErr) String() string {
	return fmt.Sprintf("\tpackage\n\t\t\"%s\"\n\tshould not depend on\n\t\t\"%s\"", e.basePackage, e.dependency)
}

type testingT struct {
	errors []*archtestErr
}

func (tt *testingT) Error(args ...interface{}) {
	tt.errors = append(tt.errors, NewArchtestErr(args[0]))
}

func (tt *testingT) errored() bool {
	return len(tt.errors) > 0
}

func (tt *testingT) message() *archtestErr {
	return tt.errors[0]
}

func (tt *testingT) AssertNoError(t *testing.T, mockT *testingT) {
	t.Helper()

	if mockT.errored() {
		t.Fatalf("\n\tunexpected error occured:\n\n%s", mockT.message())
	}
}
