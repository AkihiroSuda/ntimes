package main

import (
	"fmt"
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

type resultsByReal []Result

func (a resultsByReal) Len() int           { return len(a) }
func (a resultsByReal) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a resultsByReal) Less(i, j int) bool { return a[i].Real < a[j].Real }

// percentile computes percentile for p (0..99) using "linear interpolation between
// closest ranks method, C=1/2"
// https://en.wikipedia.org/w/index.php?title=Percentile&oldid=724036224#First_Variant.2C_.7F.27.22.60UNIQ--postMath-0000002C-QINU.60.22.27.7F
func percentile(sorted []Result, P int, sel selector) (time.Duration, error) {
	if P <= 0 || P >= 100 {
		return 0, fmt.Errorf("strange P %d", P)
	}
	n := float64(len(sorted))
	// fp computes $f(p)$ described in the "First Variant, C=1/2"
	// section of the wikipedia article.
	// Note that p is up to 1.0 here.
	fp := func(p float64) float64 {
		p1 := 1 / (2 * n)
		pn := (2*n - 1) / (2 * n)
		if p < p1 {
			return 1
		} else if p < pn {
			return n*p + 0.5
		} else {
			return n
		}
	}
	x := fp(float64(P) / 100.0)
	// vl and vh are $v_{\floor(x)}$, $v_{\floor(x)+1$}
	// in the wikipedia article.
	// Note that the wikipedia article uses [1..N]
	// indices for the slice.
	//
	//  v is $v(x) = vl + (x%1)(vh - vl)$.
	idxl := int(math.Floor(x) - 1)
	idxh := idxl + 1
	_, vl := sel(nil, &sorted[idxl])
	v := float64(vl)
	if idxh <= len(sorted)-1 {
		_, vh := sel(nil, &sorted[idxh])
		v += (math.Remainder(x, 1)) * (float64(vh) - float64(vl))
	}
	return time.Duration(v), nil
}
