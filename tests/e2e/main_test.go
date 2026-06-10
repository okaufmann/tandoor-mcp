package e2e_test

import (
	"os"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

var fixture *infra.Fixture

func TestMain(m *testing.M) {
	var err error
	fixture, err = infra.SetupFixture()
	if err != nil {
		panic(err)
	}

	code := m.Run()

	fixture.Teardown()
	os.Exit(code)
}
