//+build !windows

package main

import (
	"testing"
)

func TestCLIWithBashSleep(t *testing.T) {
	testCLI(t, &testSpec{args: []string{"-n", "3", "bash", "-c", `sleep=$((RANDOM%2)); fail=$((RANDOM%2)); echo "id=$NTIMES_ID, sleep=$sleep, fail=$fail"; sleep $sleep; exit $fail`}})
}

func TestCLIWithDD_JSON(t *testing.T) {
	testCLI(t, &testSpec{args: []string{"--format", "{{json .}}", "-n", "5", "dd", "if=/dev/urandom", "of=/dev/null", "bs=512", "count=1000"}})
}
