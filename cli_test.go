package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ensureCmd(t *testing.T, name string) {
	_, err := exec.LookPath(name)
	if err != nil {
		t.Skipf("%s unavailable: %v", name, err)
	}
}

type testSpec struct {
	args        []string
	errExpected bool
}

func testCLI(t *testing.T, spec *testSpec) ([]byte, []byte, interface{}) {
	t.Logf("testing: %+v", spec)
	stdin := strings.NewReader("")
	var stdout, stderr bytes.Buffer
	args := append([]string{"ntests"}, spec.args...)
	intf, err := xmain(args, stdin, &stdout, &stderr)
	t.Logf("stdout:\n%s\n", stdout.String())
	t.Logf("stderr:\n%s\n", stderr.String())
	t.Logf("err: %v", err)
	if spec.errExpected {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
	return stdout.Bytes(), stderr.Bytes(), intf
}

func TestCLIVersion(t *testing.T) {
	_, _, _ = testCLI(t, &testSpec{args: []string{"--version"}})
}

func TestCLIBadArgs1(t *testing.T) {
	_, _, _ = testCLI(t, &testSpec{args: []string{}, errExpected: true})
	_, _, _ = testCLI(t, &testSpec{args: []string{"--bad-flag"}, errExpected: true})
	_, _, _ = testCLI(t, &testSpec{args: []string{"-n", "3", "badcommandxxxxx"}, errExpected: true})
}

func TestCLIBadArgs2(t *testing.T) {
	ensureCmd(t, "true")
	_, _, _ = testCLI(t, &testSpec{args: []string{"-n", "0", "true"}, errExpected: true})
	_, _, _ = testCLI(t, &testSpec{args: []string{"-n", "3", "--warm-up", "10", "true"}, errExpected: true})
}

// TestCLI1 tests the CLI with `bash` and `sleep`.
// Tested flags are "-n" and "--warm-up".
// It takes about 10 seconds.
func TestCLI1(t *testing.T) {
	ensureCmd(t, "bash")
	ensureCmd(t, "sleep")
	cmd := `
sleep=$NTIMES_ID
fail=$((NTIMES_ID%2))
echo "id=$NTIMES_ID, sleep=$sleep, fail=$fail"; sleep $sleep; exit $fail
`
	_, _, intf := testCLI(t, &testSpec{args: []string{"-n", "5",
		"--warm-up", "2",
		"bash", "-c", cmd}})
	// stat for [2,3,4] (excluding [0,1])
	stat := intf.(*Stat)
	epsilon := float64(500 * time.Millisecond)
	assert.InEpsilon(t, float64(3*time.Second), float64(stat.Real.Average), epsilon)
	assert.InEpsilon(t, float64(4*time.Second), float64(stat.Real.Max), epsilon)
	assert.InEpsilon(t, float64(2*time.Second), float64(stat.Real.Min), epsilon)
	assert.InEpsilon(t, float64(1*time.Second), float64(stat.Real.StdDev), epsilon)
	assert.InEpsilon(t, float64(4*time.Second), float64(stat.Real.Percentiles["99"]), epsilon)
	assert.InEpsilon(t, float64(4*time.Second), float64(stat.Real.Percentiles["95"]), epsilon)
	assert.InEpsilon(t, float64(3*time.Second), float64(stat.Real.Percentiles["50"]), epsilon)
	assert.InEpsilon(t, 100.0/3.0, stat.Flaky, 0.01)
}

// TestCLI2 tests the CLI with `bash` and `echo`.
// Tested flags are "-n", "--debug", "--format" and "--storage".
// It should not take so long.
func TestCLI2(t *testing.T) {
	ensureCmd(t, "bash")
	ensureCmd(t, "echo")
	storage, err := ioutil.TempDir("", "ntimes-TestCLI2")
	assert.NoError(t, err)
	t.Logf("storage: %s", storage)
	defer func() {
		t.Logf("removing %s", storage)
		assert.NoError(t, os.RemoveAll(storage))
	}()
	stdout, stderr, _ := testCLI(t, &testSpec{args: []string{"-n", "5",
		"--debug",
		"--format", "{{json .}}",
		"--storage", storage,
		"bash", "-c", "echo hello"}})
	assert.Empty(t, stdout)
	var stat Stat
	err = json.Unmarshal(stderr, &stat)
	assert.NoError(t, err)
	assert.NotZero(t, stat.Real.Max)
	for i := 0; i < 5; i++ {
		dir := filepath.Join(storage, strconv.Itoa(i))
		xstdout, err := ioutil.ReadFile(filepath.Join(dir, "stdout"))
		assert.NoError(t, err)
		assert.Regexp(t, "hello\\s*", string(xstdout))
		xstderr, err := ioutil.ReadFile(filepath.Join(dir, "stderr"))
		assert.NoError(t, err)
		assert.Empty(t, xstderr)
		xresult, err := ioutil.ReadFile(filepath.Join(dir, "result.json"))
		var result Result
		err = json.Unmarshal(xresult, &result)
		assert.NoError(t, err)
		assert.Equal(t, i, result.ID)
		assert.True(t, result.Successful)
	}
}
