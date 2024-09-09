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
out\jx 

REM Send a json string on stdin to jx for pretty-printing
echo '{"foo": 0}' | out\jx 'x'

REM Get the json from a file for pretty printing
out\jx 'x' < data\dragons.json

REM Use cat to get the data from a file
cat data\dragons.json | out\jx 'x'

REM Accept the name of a file as an argument to get the data
out\jx 'x' data\dragons.json

REM 2 input files
out\jx 'x' data\dragons.json data\cats.json

REM --compact-output
REM Should look ~ like:
REM [["0","1","42"],{"baz":"1","foo":0,"bar":"42"}]
cat data\compact_me.txt | out\jx --compact-output 'x' 

REM No idea how to print exit status on Windows
REM --exit-status 
REM null 1
REM false 1
REM true 0
echo '"foo"' | out\jx --exit-status 'null' && echo %ERRORLEVEL% 
echo '"foo"' | out\jx --exit-status '{return false}' && echo %ERRORLEVEL%
echo '"foo"' | out\jx --exit-status '{return true}' && echo %ERRORLEVEL%

REM --from-file
REM 42
echo '{"z":1, "foo":2, "x51":3, "1":4, "a":32}' | out\jx --from-file tests\add.js 

REM --help (see top of file)

REM --indent 
REM {
REM        "foo": 0
REM }
echo '{"foo":0}' | out\jx --indent 7

REM --monochrome-output
REM NA observe no color in output

REM --slurp
REM [
REM   "foo",
REM   "bar"
REM ]
echo '"foo" "bar"' | out\jx --slurp

REM --sort-keys
REM {
REM    "bar": 2,
REM    "baz": 0,
REM    "foo": 40
REM }
echo '{"foo":40,"baz":0,"bar":2}' | out\jx -sort-keys 'x' 

REM --tab 
echo '{"foo":0}' | out\jx --tab 

REM --version
REM jx n.n.n
out\jx --version

