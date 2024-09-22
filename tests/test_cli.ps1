#
# Test command-line interface 
#
# Nothing, should display usage
# Usage: jx [options] <user script> file...
#   -c, --compact-output      Minimize output
#   -e, --exit-status         Exit status based on user script return value
#   -f, --from-file string    Read script from scriptfile
#   -h, --help                Print this help and exit
#       --indent int          Number of spaces per indentation (default 2)
#   -M, --monochrome-output   Disable colored output
#   -s, --slurp               Read entire stream into array, run script once
#   -S, --sort-keys           Sort object keys
#       --tab                 Use tabs instead of spaces
#   -V, --version             Print program version and exit
echo "Test 0000 ****** usage msg"
out/windows/amd64/jx

# Send a json string on stdin to jx for pretty-printing
echo "Test 0001 ****** echo stdin"
echo '{"foo": 0}' | out/windows/amd64/jx x

# Get the json from a file for pretty printing
echo "Test 0002 ****** redirection stdin not supported in ps"
#out/windows/amd64/jx x < tests/data/dragons.json

# Use cat to get the data from a file
echo "Test 0003 ****** cat stdin"
cat tests/data/dragons.json | out/windows/amd64/jx x

# Accept the name of a file as an argument to get the data
echo "Test 0004 ****** filename"
out/windows/amd64/jx x tests/data/dragons.json

# 2 input files
echo "Test 0005 ****** two input files"
out/windows/amd64/jx x tests/data/dragons.json tests/data/cats.json

# --compact-output
# Should look ~ like:
# [["0","1","42"],{"baz":"1","foo":0,"bar":"42"}]
echo "Test 0006 ******"
cat tests/data/compact_me.txt | out/windows/amd64/jx --compact-output x  
echo ""

# No idea how to print exit status on windows
# --exit-status 
# null 
# 1
# false 
# 1
# true 
# 0
echo "Test 0007 ******"
echo 'null' | out/windows/amd64/jx --exit-status 'x'; echo $LastExitCode
echo 'false' | out/windows/amd64/jx --exit-status 'x'; echo $LastExitCode
echo 'true' | out/windows/amd64/jx --exit-status; echo $LastExitCode

# --from-file
# 42
echo "Test 0008 ****** script file"
echo '{"z":1, "foo":2, "x51":3, "1":4, "a":32}' | out/windows/amd64/jx --from-file tests/add.js 

# --help (see top of file)

# --indent 
# {
#        "foo": 0
# }
echo "Test 0009 ******"
echo '{"foo":0}' | out/windows/amd64/jx --indent 7

# --monochrome-output
# NA observe no color in output

# --slurp
# [
#   "foo",
#   "bar"
# ]
echo "Test 0010 ******"
echo '"foo" "bar"' | out/windows/amd64/jx --slurp

# --sort-keys
# {
#    "bar": 2,
#    "baz": 0,
#    "foo": 40
# }
echo "Test 0011 ******"
echo '{"foo":40,"baz":0,"bar":2}' | out/windows/amd64/jx --sort-keys 'x' 

# --tab 
echo "Test 0012 ******"
echo '{"foo":0}' | out/windows/amd64/jx --tab 

# --version
# jx n.n.n
echo "Test 0013 ******"
out/windows/amd64/jx --version

echo "Test 0014 ******"
cat tests/data/dragons.json | out/windows/amd64/jx 'x.xyz()'
