package main

import (
	"math"
	"time"
)

type selector func(*Stat, *Result) (*TimeStat, time.Duration)

func selReal(stat *Stat, result *Result) (ts *TimeStat, t time.Duration) {
	if stat != nil {
		ts = stat.Real
	}
	if result != nil {
		t = result.Real
	}
	return
}
func selUser(stat *Stat, result *Result) (ts *TimeStat, t time.Duration) {
	if stat != nil {
		ts = stat.User
	}
	if result != nil {
		t = result.User
	}
	return
}
func selSystem(stat *Stat, result *Result) (ts *TimeStat, t time.Duration) {
	if stat != nil {
		ts = stat.System
	}
	if result != nil {
		t = result.System
	}
	return
}

func updateMinMaxSum(stat *Stat, result *Result, sels []selector) {
	for _, sel := range sels {
		tStat, t := sel(stat, result)
		if tStat.Min == 0 ||
			t < tStat.Min {
			tStat.Min = t
		}
		if t > tStat.Max {
			tStat.Max = t
		}
		tStat.sum += t
	}
}

// updateVarianceNumerator updates the numerator of the variance.
// Average values in stat must be already computed.
func updateVarianceNumerator(stat *Stat, result *Result, sels []selector) {
	for _, sel := range sels {
		tStat, t := sel(stat, result)
		a := float64(tStat.Average) - float64(t)
		tStat.varianceNumerator += a * a
	}
}

func updateNFlakies(stat *Stat, result *Result) {
	if !result.Successful {
		stat.nflakies++
	}
}

// updateStatPhase0 fills up
// stat.{Flaky, nflakies, TimeStat.{Average, Min,Max,sum}}
func updateStatPhase0(stat *Stat, results []Result, sels []selector) {
	for _, result := range results {
		updateMinMaxSum(stat, &result, sels)
		updateNFlakies(stat, &result)
	}
	for _, sel := range sels {
		tStat, _ := sel(stat, nil)
		tStat.Average = tStat.sum / time.Duration(len(results))
	}
	stat.Flaky = 100 * float64(stat.nflakies) / float64(len(results))
}

// updateStatPhase1 fills up
// stat.{TimeStat.{StdDev, variance}}
func updateStatPhase1(stat *Stat, results []Result, sels []selector) {
	for _, result := range results {
		updateVarianceNumerator(stat, &result, sels)
	}
	for _, sel := range sels {
		tStat, _ := sel(stat, nil)
		if len(results) > 1 {
			// we use population variance here
			variance := tStat.varianceNumerator / float64(len(results)-1)
			tStat.StdDev = time.Duration(math.Sqrt(variance))
		} else {
			tStat.StdDev = 0
		}
	}
}
