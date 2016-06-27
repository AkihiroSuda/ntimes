package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMath(t *testing.T) {
	data := []Result{
		{ID: 0, Successful: true, Real: 42 * time.Second},
		{ID: 1, Successful: true, Real: 4 * time.Second},
		{ID: 2, Successful: false, Real: 4242 * time.Second},
		{ID: 3, Successful: true, Real: 42 * time.Second},
		{ID: 4, Successful: false, Real: 24 * time.Second},
	}
	sels := []selector{selReal}
	stat := Stat{}
	stat.Real = &TimeStat{}
	assertAfterPhase0 := func() {
		assert.InEpsilon(t, 40, stat.Flaky, 0.001)
		assert.Equal(t, 4242*time.Second, stat.Real.Max)
		assert.Equal(t, 4*time.Second, stat.Real.Min)
	}
	updateStatPhase0(&stat, data, sels)
	t.Logf("stat after phase0: %+v, %+v", stat, stat.Real)
	assertAfterPhase0()

	updateStatPhase1(&stat, data, sels)
	t.Logf("stat after phase1: %+v, %+v", stat, stat.Real)
	// check: `python3 -c "import statistics; print(statistics.stdev((42,4,4242,42,24)))"`
	assert.InEpsilon(t,
		float64(1884*time.Second), float64(stat.Real.StdDev),
		float64(1*time.Second))
	// again
	assertAfterPhase0()
}

func TestMathStrange(t *testing.T) {
	data := []Result{
		{ID: 0, Successful: false, Real: 0 * time.Second},
	}
	sels := []selector{selReal}
	stat := Stat{}
	stat.Real = &TimeStat{}
	assertAfterPhase0 := func() {
		assert.InEpsilon(t, 100, stat.Flaky, 0.001)
		assert.Equal(t, 0*time.Second, stat.Real.Max)
		assert.Equal(t, 0*time.Second, stat.Real.Min)
	}
	updateStatPhase0(&stat, data, sels)
	t.Logf("stat after phase0: %+v, %+v", stat, stat.Real)
	assertAfterPhase0()

	updateStatPhase1(&stat, data, sels)
	t.Logf("stat after phase1: %+v, %+v", stat, stat.Real)
	// stddev is time.Duration(i.e. int64)
	assert.Zero(t, stat.Real.StdDev)
	// again
	assertAfterPhase0()
}

func TestMathPercentile(t *testing.T) {
	// from the wikipedia "Worked Example of the First Variant"
	// https://en.wikipedia.org/w/index.php?title=Percentile&oldid=724036224#Worked_Example_of_the_First_Variant
	data := []Result{
		{Real: 15 * time.Second},
		{Real: 20 * time.Second},
		{Real: 35 * time.Second},
		{Real: 40 * time.Second},
		{Real: 50 * time.Second},
	}
	t.Logf("data=%v", data)
	expected := map[int]time.Duration{
		5:  15 * time.Second,
		30: 20 * time.Second,
		40: 27*time.Second + 500*time.Millisecond,
		95: 50 * time.Second,
	}
	for i, d := range expected {
		p, err := percentile(data, i, selReal)
		assert.NoError(t, err)
		t.Logf("%d percentile=%v", i, p)
		assert.Equal(t, d, p)
	}
	_, err := percentile(data, 0, selReal)
	assert.Error(t, err)
	_, err = percentile(data, 100, selReal)
	assert.Error(t, err)
}
