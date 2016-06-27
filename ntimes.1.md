% NTIMES(1) General Commands Manual
% ntimes
% JUNE 2016
# NAME
ntimes \- **time(1)** with average time, flaky rate, ..

# SYNOPSIS
**ntimes** [OPTIONS] COMMAND [arg...]

# DESCRIPTION

ntimes is an improved version of **time(1)**.

You can execute an command`N times, and measure the average/max/min time taken for the execution.
You can also measure the "flaky" rate (i.e. failure rate).

# OPTIONS
**--help**
  Print usage statement

**-f**, **--format**=*""*
  Format string (in golang text/template, e.g. "{{json .}}")

**-n**, **--repeat-n-timesg**=*1*
  Number of times

**--storage**
  Path to stdout,stderr storage	

**--version**
  Print version to stdout and exit

**--warm-up**
  Skip first n iterations for stat

# Examples

Example usage:

	$ ntimes -n 10 bash -c 'sleep=$((RANDOM%5)); fail=$((RANDOM%2)); echo "id=$NTIMES_ID, sleep=$sleep, fail=$fail"; sleep $sleep; exit $fail'
    id=0, sleep=4, fail=0
    id=1, sleep=0, fail=0
    id=2, sleep=3, fail=1
    id=3, sleep=4, fail=0
    id=4, sleep=3, fail=1
    id=5, sleep=0, fail=1
    id=6, sleep=3, fail=1
    id=7, sleep=0, fail=0
    id=8, sleep=3, fail=1
    id=9, sleep=3, fail=1
    
    real average: 2.303595053s, max: 4.004607594s, min: 3.455164ms, std dev: 1.636182569s
    real 99 percentile: 4.004607594s, 95 percentile: 4.004607594s, 50 percentile: 3.002190904s
    user average: 0, max: 0, min: 0, std dev: 0
    sys  average: 0, max: 0, min: 0, std dev: 0
    flaky: 60%
    
You can specify the report format using Go's **text/template** syntax.
Additionally to the standard functions provided by **text/template**, the **json** function is available.
Note that a **time.Duration** value is expressed in nanoseconds.

    $ ntimes --format "{{json .}}" -n 10 dd if=/dev/urandom of=/dev/null bs=512 count=1000
    1000+0 records in
    1000+0 records out
    512000 bytes (512 kB, 500 KiB) copied, 0.0448978 s, 11.4 MB/s
	...
    {"real":{"average":44155207,"max":68222928,"min":38143407,"stddev":9421337,"percentiles":{"50":39855284,"95":68222928,"99":68222928}},"user":{"average":0,"max":0,"min":0,"stddev":0},"system":{"average":36000000,"max":36000000,"min":36000000,"stddev":0},"flaky":0}

Practical example for debugging flaky tests with Namazu ( **nmz(1)**, https://github.com/osrg/namazu):

    $ cd some_maven_project
    $ sudo ntimes -n 10 --storage /tmp/logs nmz inspectors -cmd "mvn test"
    ...
    Flaky: 10%
    
    $ find /tmp/logs -name result.json | xargs jq .successful
    true
    true
    true
    false
    true
    true
    true
    true
    true
    true


# SEE ALSO
**time(1)**, **nmz(1)**

# AUTHOR
June 2016, writen by Akihiro Suda
https://github.com/AkihiroSuda/ntimes
