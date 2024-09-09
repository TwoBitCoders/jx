package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bcicen/jstream"
	"github.com/dop251/goja"
    "github.com/iancoleman/orderedmap"
	flag "github.com/spf13/pflag"
)

const (
    Version = "0.0.2"
    NoError = 0 
    FalseOrNull = 1 
    UsageError = 2
    CompileError = 3
    NoValidResult = 4
)

const template = `
    v = JSON.parse(json);
    f = x=>%s;
    r = {};
    r.v = f(v);
    JSON.stringify(r, null, null);
`
func formatJSON(v interface{}, opts Options) ([]byte, error) {
    f := ColorJsonFormatterNew()
    f.SortKeys = opts.sortKeys
    f.DisabledColor = opts.mono
    s, err := f.Marshal(v)
    if err != nil {
        fmt.Printf("err:%s\n", err)
    }
	return s, nil
}

func formatJSONIndent(v interface{}, indent string, opts Options) ([]byte, error) {
    f := ColorJsonFormatterNew()
    f.SortKeys = opts.sortKeys
    f.DisabledColor = opts.mono
    s, err := f.MarshalIndent(v, indent)
    if err != nil {
        fmt.Printf("err:%s\n", err)
    }
	return s, nil
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

func runScript(vm *goja.Runtime, script string, jsonVal []byte) (goja.Value, error) {

    // create the script and put js value through it
    vm.Set("json", vm.ToValue(string(jsonVal)))
    js := fmt.Sprintf(template, script)

    // run it
    v, err := vm.RunString(js)
    if err != nil {
        panic(err)
    }

    return v, nil
}

func outputVal(v interface{}, opts Options) {
    var indentStr string

    var buf []byte
    var err error
    if opts.compact {
        buf, err = formatJSON(v, opts)
    } else if opts.tabs {
        indentStr = "\t"
        buf, err = formatJSONIndent(v, indentStr, opts)
        if err != nil {
            fmt.Printf("%s\n", err)
        }
    } else {
        indentStr = strings.Repeat(" ", opts.indent)
        buf, err = formatJSONIndent(v, indentStr, opts)
        if err != nil {
            fmt.Printf("%s\n", err)
        }
    }
    fmt.Printf("%s\n", string(buf))
}

func biteStream(vm *goja.Runtime, script string, d *jstream.Decoder,
    stdoutTTY bool, opts Options) (bool, error) {
    var jv []byte
    var err error
    for mv := range d.Stream() {
        switch foo := mv.Value.(type) {
        case jstream.KVS:
            jv, err = foo.MarshalJSON()
        default:
            jv, err = json.Marshal(foo)
        }
        err = d.Err()
        if err != nil {
            return false, err
        } 

        rv, err := runScript(vm, script, jv)
        if err != nil {
            return false, fmt.Errorf("Error runing user script %s\n", err)
        }

        var x *string
        vm.ExportTo(rv, &x)
        r := orderedmap.New()
        json.Unmarshal([]byte(*x), &r)
        v, _ := r.Get("v")

        outputVal(v, opts)
    }

    return true, nil
}

func slurpStream(vm *goja.Runtime, script string, d *jstream.Decoder, stdoutTTY bool, opts Options) error {
    var a []string
    for mv := range d.Stream() {
        jv, err := json.Marshal(mv.Value)
        err = d.Err()
        if err != nil {
            return err
        }
        a = append(a, string(jv))
    }
    jv := fmt.Sprintf("[%s]", strings.Join(a, ","))
    rv, err := runScript(vm, script, []byte(jv))
    if err != nil {
        return fmt.Errorf("Error runing user script %s\n", err)
    }

    var x *string
    vm.ExportTo(rv, &x)
    r := orderedmap.New()
    json.Unmarshal([]byte(*x), &r)
    v, _ := r.Get("v")

    outputVal(v, opts)

    return nil
}

func runScriptForFile(script string, path *string, 
    stdoutTTY bool, opts Options) (bool, error) {
    var f *os.File
    var err error
    var x bool
    if path == nil {
        f = os.Stdin
    } else {
        f, err = os.Open(*path)
        if err != nil {
            return x, fmt.Errorf("%s: Could not open %s: %s\n", os.Args[0], *path, err) 
        }
        defer f.Close()
    }

    vm := goja.New()
    decoder := jstream.NewDecoder(f, 0)

    // This nonsense is required to get the decoder to preserve key order
    // It does at least seem to work
    decoder = decoder.EmitKV()
    decoder = decoder.ObjectAsKVS()
    if opts.slurp {
        err = slurpStream(vm, script, decoder, stdoutTTY, opts)
        if err != nil {
            return x, err
        }
    } else {
        x, err = biteStream(vm, script, decoder, stdoutTTY, opts)
        if err != nil {
            return x, err
        }
    }
    return x, nil
}

type Options struct {
    scriptFile string
    slurp bool
    compact bool
    indent int
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
    flag.BoolVarP(&opts.mono, "monochrome-output", "M", false, "Disable colored output")
    flag.BoolVar(&opts.tabs, "tab", false, "Use tabs instead of spaces")
    flag.BoolVarP(&opts.sortKeys, "sort-keys", "S", false, "Sort object keys")
    flag.BoolVarP(&opts.exitStatus, "exit-status", "e", false, "Exit status based on user script return value")
    flag.BoolVarP(&opts.help, "help", "h", false, "Print this help and exit")
    flag.BoolVarP(&opts.printVersion, "version", "V", false, "Print program version and exit")
    flag.Parse()

    if opts.help {
        flag.Usage()
        os.Exit(UsageError)
    }

    if opts.printVersion {
        fmt.Printf("%s %s\n", os.Args[0], Version)
        os.Exit(NoError)
    }

    // test if stdin has stuff to read, so we can make more decisions later
    stdinData := false 
    stat, _ := os.Stdin.Stat()
    if (stat.Mode() & os.ModeCharDevice) == 0 {
        stdinData = true
    }

    // No data to read and no args at all, give the user a hand and exit
    if !stdinData && (len(flag.Args()) < 1) {
        flag.Usage()
        os.Exit(UsageError)
    } 

    // test if the user specified a script as an option, if 
    // they did, userFile inputs move up one position
    userFileIndex := 1
    if opts.scriptFile != "" {
        userFileIndex = 0
    }

    userFiles := make([]string, 0)
    if len(flag.Args()) >= userFileIndex {
        userFiles = flag.Args()[userFileIndex:]
    }

    // test if stdout is a tty
    stdoutTTY := false 
    stat, _ = os.Stdout.Stat()
    if (stat.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
       stdoutTTY = true
    }

    // if the user set indent outside range let them know and exit
    if opts.indent < 0 || opts.indent > 7 {
        fmt.Printf("%s: --indent takes a number between -1 and 7\n", os.Args[0])
        flag.Usage()
        os.Exit(UsageError)
    }

    // get the user script
    var userScript string
    var err error
    if opts.scriptFile != "" {
        buf, err := readUserFile(opts.scriptFile)
        if err != nil {
            fmt.Printf("%s\n", err)
            os.Exit(UsageError)
        }
        userScript = string(buf)
    } else if opts.scriptFile == "" && len(flag.Args()) < 1 {
        userScript = "x"
    } else {
        userScript = flag.Args()[0]
    }

    // finally get down to business
    // run the provided script either against the
    // specified files or stdin
    var x interface{}
    if len(userFiles) > 0 {
        for i := range userFiles {
            userFile := userFiles[i]
            x, err = runScriptForFile(userScript, &userFile, stdoutTTY, opts)
            if err != nil {
                fmt.Printf("%s: Error - %s\n", os.Args[0], err)
                return
            }
        }
    } else {
        x, err = runScriptForFile(userScript, nil, stdoutTTY, opts)
        if err != nil {
            fmt.Printf("%s: Error - %s\n", os.Args[0], err)
            return
        }
    }

    // When exitStatus flag is set check for these specific conditions
    r := 0
    if opts.exitStatus && x == false {
        r = FalseOrNull
    }

    os.Exit(r)
}


