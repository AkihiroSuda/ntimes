package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testSpec struct {
	args        []string
	errExpected bool
}

func testCLI(t *testing.T, spec *testSpec) {
	t.Logf("testing: %+v", spec)
	stdin := strings.NewReader("")
	var stdout, stderr bytes.Buffer
	args := append([]string{"ntests"}, spec.args...)
	err := xmain(args, stdin, &stdout, &stderr)
	t.Logf("stdout:\n%s\n", stdout.String())
	t.Logf("stderr:\n%s\n", stderr.String())
	t.Logf("err: %v", err)
	if spec.errExpected {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
}

func TestCLIVersion(t *testing.T) {
	testCLI(t, &testSpec{args: []string{"--version"}})
}

func TestCLIBadArgs(t *testing.T) {
	testCLI(t, &testSpec{args: []string{}, errExpected: true})
	testCLI(t, &testSpec{args: []string{"--bad-flag"}, errExpected: true})
}
