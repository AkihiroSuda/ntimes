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
	p5 := percentile(data, 5, selReal)
	t.Logf("5 percentile=%v", p5)
	assert.Equal(t, 15*time.Second, p5)

	p30 := percentile(data, 30, selReal)
	t.Logf("30 percentile=%v", p30)
	assert.Equal(t, 20*time.Second, p30)

	p40 := percentile(data, 40, selReal)
	t.Logf("40 percentile=%v", p40)
	assert.Equal(t, 27*time.Second+500*time.Millisecond, p40)

	p95 := percentile(data, 95, selReal)
	t.Logf("95 percentile=%v", p95)
	assert.Equal(t, 50*time.Second, p95)
}
