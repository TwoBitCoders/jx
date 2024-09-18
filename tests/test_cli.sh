#!/bin/sh

function usage {
    echo "Usage: test_cli.sh OS ARCH"
}

if [ -z "$1" ]; then 
    echo "Parameter OS is required"
    usage
    exit
fi

if [ -z "$2" ]; then 
    echo "Usage: parameter ARCH is required"
    usage
    exit
fi

OS="$1"
ARCH="$2"

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
out/"$OS"/"$ARCH"/jx 

# Send a json string on stdin to jx for pretty-printing
echo "Test 0001 ****** echo stdin"
echo '{"foo": 0}' | out/"$OS"/"$ARCH"/jx x

# Get the json from a file for pretty printing
echo "Test 0002 ****** redirection stdin"
out/"$OS"/"$ARCH"/jx x < tests/data/dragons.json

# Use cat to get the data from a file
echo "Test 0003 ****** cat stdin"
cat tests/data/dragons.json | out/"$OS"/"$ARCH"/jx x

# Accept the name of a file as an argument to get the data
echo "Test 0004 ****** filename"
out/"$OS"/"$ARCH"/jx x tests/data/dragons.json

# 2 input files
echo "Test 0005 ****** two input files"
out/"$OS"/"$ARCH"/jx x tests/data/dragons.json tests/data/cats.json

# --compact-output
# Should look ~ like:
# [["0","1","42"],{"baz":"1","foo":0,"bar":"42"}]
echo "Test 0006 ******"
cat tests/data/compact_me.txt | out/"$OS"/"$ARCH"/jx --compact-output x  
echo ""

# No idea how to print exit status on Windows
# --exit-status 
# null 1
# false 1
# true 0
echo "Test 0007 ******"
echo '"foo"' | out/"$OS"/"$ARCH"/jx --exit-status 'null'; echo $?
echo '"foo"' | out/"$OS"/"$ARCH"/jx --exit-status '{return false}'; echo $?
echo '"foo"' | out/"$OS"/"$ARCH"/jx --exit-status '{return true}'; echo $?

# --from-file
# 42
echo "Test 0008 ****** script file"
echo '{"z":1, "foo":2, "x51":3, "1":4, "a":32}' | out/"$OS"/"$ARCH"/jx --from-file tests/add.js 

# --help (see top of file)

# --indent 
# {
#        "foo": 0
# }
echo "Test 0009 ******"
echo '{"foo":0}' | out/"$OS"/"$ARCH"/jx --indent 7

# --monochrome-output
# NA observe no color in output

# --slurp
# [
#   "foo",
#   "bar"
# ]
echo "Test 0010 ******"
echo '"foo" "bar"' | out/"$OS"/"$ARCH"/jx --slurp

# --sort-keys
# {
#    "bar": 2,
#    "baz": 0,
#    "foo": 40
# }
echo "Test 0011 ******"
echo '{"foo":40,"baz":0,"bar":2}' | out/"$OS"/"$ARCH"/jx --sort-keys 'x' 

# --tab 
echo "Test 0012 ******"
echo '{"foo":0}' | out/"$OS"/"$ARCH"/jx --tab 

# --version
# jx n.n.n
echo "Test 0013 ******"
out/"$OS"/"$ARCH"/jx --version

echo "Test 0014 ******"
cat tests/data/dragons.json | out/"$OS"/"$ARCH"/jx 'x.xyz()'

