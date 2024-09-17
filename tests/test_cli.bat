@echo off
REM
REM Test command-line interface
REM

REM Nothing, should display usage
REM Usage: jx [options] <user script> file...
REM   -c, --compact-output      Minimize output
REM   -e, --exit-status         Exit status based on user script return value
REM   -f, --from-file string    Read script from scriptfile
REM   -h, --help                Print this help and exit
REM       --indent int          Number of spaces per indentation (default 2)
REM   -M, --monochrome-output   Disable colored output
REM   -s, --slurp               Read entire stream into array, run script once
REM   -S, --sort-keys           Sort object keys
REM       --tab                 Use tabs instead of spaces
REM   -V, --version             Print program version and exit

echo "Test 0000 ******"
out\jx

REM Send a JSON string on stdin to jx for pretty-printing
echo "Test 0001 ******"
echo {"foo": 0} | out\jx "x"

REM Get the JSON from a file for pretty printing
echo "Test 0002 ******"
type tests\data\dragons.json | out\jx "x"

REM Use type (instead of cat) to get the data from a file
echo "Test 0003 ******"
type tests\data\dragons.json | out\jx "x"

REM Accept the name of a file as an argument to get the data
echo "Test 0004 ******"
out\jx "x" tests\data\dragons.json

REM 2 input files
echo "Test 0005 ******"
out\jx "x" tests\data\dragons.json tests\data\cats.json

REM --compact-output
REM Should look ~ like:
REM [["0","1","42"],{"baz":"1","foo":0,"bar":"42"}]
echo "Test 0006 ******"
type tests\data\compact_me.txt | out\jx --compact-output "x"

REM --exit-status (there's no direct equivalent for 'echo $?')
REM For Windows, %ERRORLEVEL% stores the exit status, so use conditional checks.
echo "Test 0007 ******"
echo "foo" | out\jx --exit-status "null"
echo Exit status: %ERRORLEVEL%

echo "foo" | out\jx --exit-status "{return false}"
echo Exit status: %ERRORLEVEL%

echo "foo" | out\jx --exit-status "{return true}"
echo Exit status: %ERRORLEVEL%

REM --from-file
REM 42
echo "Test 0008 ******"
echo {"z":1, "foo":2, "x51":3, "1":4, "a":32} | out\jx --from-file tests\add.js

REM --indent
REM {
REM        "foo": 0
REM }
echo "Test 0009 ******"
echo {"foo": 0} | out\jx --indent 7

REM --monochrome-output
REM Observe no color in output

REM --slurp
REM [
REM   "foo",
REM   "bar"
REM ]
echo "Test 0010 ******"
echo "foo" "bar" | out\jx --slurp

REM --sort-keys
REM {
REM    "bar": 2,
REM    "baz": 0,
REM    "foo": 40
REM }
echo "Test 0011 ******"
echo {"foo": 40, "baz": 0, "bar": 2} | out\jx --sort-keys "x"

REM --tab
echo "Test 0012 ******"
echo {"foo": 0} | out\jx --tab

REM --version
REM jx n.n.n
echo "Test 0013 ******"
out\jx --version

echo "Test 0014 ******"
type tests\data\dragons.json | out\jx "x.xyz()"
