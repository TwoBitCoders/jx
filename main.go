package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bcicen/jstream"
	"github.com/dop251/goja"
	"github.com/tidwall/pretty"
	flag "github.com/spf13/pflag"
    colorable "github.com/mattn/go-colorable"
)

const (
    Version = "0.0.7"
)

type ResultCode int

const (
    NoError ResultCode = iota
    FalseOrNull
    UsageError
    CompileError
    NoValidResult
)

const template = `v = JSON.parse(json); f = x=>%s; v = f(v); JSON.stringify(v, null, null);`

// Test if stdout is hooked to a tty 
// We want to know, so we can make good decisions about when to colorize
func isTTY(f *os.File) (bool) {
    stat, _ := f.Stat()
    if (stat.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
       return true
    }
    return false
}

func formatJSON(buf []byte, opts Options) ([]byte, error) {
    if (opts.color || isTTY(os.Stdout)) && !opts.mono {
        buf = pretty.Color(buf, nil)
    }
	return buf, nil
}

func formatJSONIndent(buf []byte, indent string, opts Options) ([]byte, error) {
    var pOpts = &pretty.Options{Width: 80, Prefix:"", Indent: indent, SortKeys: opts.sortKeys}
    buf = pretty.PrettyOptions(buf, pOpts)
    if (opts.color || isTTY(os.Stdout)) && !opts.mono {
        buf = pretty.Color(buf, nil)
    }
	return buf, nil
}

func readUserFile(path string) ([]byte, error) {
        f, err := os.Open(path)
        if err != nil {
            return nil, fmt.Errorf("%s: Could not open %s: %s\n", os.Args[0], path, err) 
        }
        defer f.Close()

        buf, err := io.ReadAll(f)
        if err != nil {
            fmt.Printf("%s: Could not read %s: %s\n", os.Args[0], path, err)
            return nil, err
        }
        return buf, nil
}

func compileScript(script string) (*goja.Program, error) {
 
    vm := goja.New()
    vm.Set("json", vm.ToValue("{}"))
    js := fmt.Sprintf(template, script)

    p, err := goja.Compile("foo", js, false)
    if err != nil {
        return nil, err
    }

    _, err = vm.RunProgram(p)
    if err != nil {
        return nil, err
    }

    return p, nil
}

func runScript(vm *goja.Runtime, script string, jsonVal []byte) (goja.Value, error) {

    // create the script and put js value through it
    vm.Set("json", vm.ToValue(string(jsonVal)))
    js := fmt.Sprintf(template, script)

    // run it
    v, err := vm.RunString(js)
    if err != nil {
        return nil, err
    }

    return v, nil
}

func outputVal(v []byte, opts Options) error {
    var indentStr string

    var buf []byte
    var err error
    if opts.compact {
        buf, err = formatJSON(v, opts)
        if err != nil {
            return err
        }
    } else if opts.tabs {
        indentStr = "\t"
        buf, err = formatJSONIndent(v, indentStr, opts)
        if err != nil {
            return err
        }
    } else {
        indentStr = strings.Repeat(" ", opts.indent)
        buf, err = formatJSONIndent(v, indentStr, opts)
        if err != nil {
            return err
        }
    }

    out := colorable.NewColorableStdout()
    fmt.Fprintf(out, "%s", string(buf))
    return nil
}

func biteStream(vm *goja.Runtime, script string, d *jstream.Decoder, opts Options) (ResultCode, error) {
    var jv []byte
    var err error
    var buf []byte
    for mv := range d.Stream() {
        switch foo := mv.Value.(type) {
        case jstream.KVS:
            jv, err = foo.MarshalJSON()
        default:
            jv, err = json.Marshal(foo)
        }
        err = d.Err()
        if err != nil {
            return UsageError, fmt.Errorf("Error parsing inputs %s\n", err)
        } 

        rv, err := runScript(vm, script, jv)
        if err != nil {
            return CompileError, err
        }

        vm.ExportTo(rv, &buf)

        err = outputVal(buf, opts)
        if err != nil {
            return NoValidResult, fmt.Errorf("Error outputting value %s\n", err)
        }

    }

    // When exitStatus flag is set check for these specific conditions
    if opts.exitStatus {
        if string(buf) == "false" || string(buf) == "null" {
            return FalseOrNull, nil
        }
    }

    return NoError, nil
}

func slurpStream(vm *goja.Runtime, script string, d *jstream.Decoder, 
    opts Options) (ResultCode, error) {
    var a []string
    for mv := range d.Stream() {
        jv, err := json.Marshal(mv.Value)
        err = d.Err()
        if err != nil {
            return UsageError, fmt.Errorf("Error parsing inputs %s\n", err)
        }
        a = append(a, string(jv))
    }
    jv := fmt.Sprintf("[%s]", strings.Join(a, ","))
    rv, err := runScript(vm, script, []byte(jv))
    if err != nil {
        return CompileError, err
    }

    var buf []byte
    vm.ExportTo(rv, &buf)
    outputVal(buf, opts)

    return NoError, nil
}

func processStream(script string, strm io.Reader, 
    opts Options) (ResultCode, error) {
    var f io.Reader
    var err error
    var rc ResultCode
    if strm == nil {
        f = os.Stdin
    } else {
        f = strm
    }

    vm := goja.New()
    decoder := jstream.NewDecoder(f, 0)

    // Required to get the decoder to preserve key order
    decoder = decoder.EmitKV()
    decoder = decoder.ObjectAsKVS()
    if opts.slurp {
        rc, err = slurpStream(vm, script, decoder, opts)
        if err != nil {
            return rc, err
        }
    } else {
        rc, err = biteStream(vm, script, decoder, opts)
        if err != nil {
            return rc, err
        }
    }
    return rc, nil
}

type Options struct {
    scriptFile string
    slurp bool
    compact bool
    indent int
    color bool
    mono bool
    tabs bool
    sortKeys bool
    exitStatus bool
    help bool
    printVersion bool
}

func main() {

    // Custom usage message like "Usage: jx [options] [user script] FILES\n"
    flag.Usage = func() {
        fmt.Fprintf(os.Stdout, 
            "Usage: %s [options] <user script> file...\n", 
            os.Args[0])
        flag.PrintDefaults()
    }

    opts := Options{}

    flag.StringVarP(&opts.scriptFile, "from-file", "f", "", "Read script from scriptfile")
    flag.BoolVarP(&opts.slurp, "slurp", "s", false, "Read entire stream into array, run script once")
    flag.IntVar(&opts.indent, "indent", 2, "Number of spaces per indentation")
    flag.BoolVarP(&opts.compact, "compact-output", "c", false, "Minimize output")
    flag.BoolVarP(&opts.color, "color-output", "C", false, "Colorize JSON output")
    flag.BoolVarP(&opts.mono, "monochrome-output", "M", false, "Disable colored output")
    flag.BoolVar(&opts.tabs, "tab", false, "Use tabs instead of spaces")
    flag.BoolVarP(&opts.sortKeys, "sort-keys", "S", false, "Sort object keys")
    flag.BoolVarP(&opts.exitStatus, "exit-status", "e", false, "Exit status based on user script return value")
    flag.BoolVarP(&opts.help, "help", "h", false, "Print this help and exit")
    flag.BoolVarP(&opts.printVersion, "version", "V", false, "Print program version and exit")
    flag.Parse()

    //
    // Pure opts stuff
    //

    if opts.help {
        flag.Usage()
        os.Exit(int(UsageError))
    }

    if opts.printVersion {
        fmt.Printf("%s %s\n", os.Args[0], Version)
        os.Exit(int(NoError))
    }

    // if the user set indent outside range let them know and exit
    if opts.indent < 0 || opts.indent > 7 {
        fmt.Printf("%s: --indent takes a number between -1 and 7\n", os.Args[0])
        flag.Usage()
        os.Exit(int(UsageError))
    }

	// 
    // Handle user script
	//
    var userScript string
    var err error
    userFileIndex := 1
    if opts.scriptFile != "" {
        buf, err := readUserFile(opts.scriptFile)
        if err != nil {
            fmt.Printf("%s\n", err)
            os.Exit(int(UsageError))
        }

        // if the user specified a script as an option userFile inputs move up 
        // one position
        userFileIndex = 0

        userScript = string(buf)
    } else if len(flag.Args()) > 0 {
        userScript = flag.Args()[0]
    }

    // Still no userScript and either stdout or stdin is not a terminal 
    //(i.e., it is redirected), default the userScript to "x"
    if(userScript == "" && (!isTTY(os.Stdout) || !isTTY(os.Stdin))) {
        userScript = "x"
    }

    if(userScript == "") {
        flag.Usage()
        os.Exit(int(UsageError))
    }

    // Try to compile the userScript
    _, err = compileScript(userScript)
    if err != nil {
        fmt.Printf("%s\n", err)
        os.Exit(int(CompileError))
    }

    // Scrabble up any remaining file paths
    userFiles := make([]string, 0)
    if len(flag.Args()) >= userFileIndex {
        userFiles = flag.Args()[userFileIndex:]
    }

	// If there are input files create a stream to read them
    var strm io.Reader
    if len(userFiles) > 0 {
        var filesBuf []byte 
        for i := range userFiles {
            userFile := userFiles[i]
            buf, err := readUserFile(userFile)
            if err != nil {
                fmt.Printf("%s: %s\n", os.Args[0], err)
                return
            }
            filesBuf = append(filesBuf, buf...)
        }
        if filesBuf != nil {
            strm = strings.NewReader(string(filesBuf))
        }
    }

    // finally get down to business
    // run the provided script either against the
    // specified files or stdin
	rc, err := processStream(userScript, strm, opts)
    if err != nil {
        fmt.Printf("%s: %s\n", os.Args[0], err)
        return
    }

    os.Exit(int(rc))
}

