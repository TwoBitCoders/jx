package main

// forked from here - https://github.com/TylerBrock/colorjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
    "github.com/iancoleman/orderedmap"
)

const initialDepth = 0
const valueSep = ","
const null = "null"
const startMap = "{"
const endMap = "}"
const startArray = "["
const endArray = "]"

const emptyMap = startMap + endMap
const emptyArray = startArray + endArray

type Formatter struct {
	KeyColor        *color.Color
	StringColor     *color.Color
	BoolColor       *color.Color
	NumberColor     *color.Color
	NullColor       *color.Color
	StringMaxLength int
	Indent          string
	DisabledColor   bool
	RawStrings      bool
    SortKeys        bool
}

func ColorJsonFormatterNew() *Formatter {
	return &Formatter{
		KeyColor:        color.New(color.FgWhite),
		StringColor:     color.New(color.FgGreen),
		BoolColor:       color.New(color.FgYellow),
		NumberColor:     color.New(color.FgCyan),
		NullColor:       color.New(color.FgMagenta),
		StringMaxLength: 0,
		DisabledColor:   false,
		Indent:          "",
		RawStrings:      false,
	}
}

func (f *Formatter) sprintfColor(c *color.Color, format string, args ...interface{}) string {
	if f.DisabledColor || c == nil {
		return fmt.Sprintf(format, args...)
	}
	return c.SprintfFunc()(format, args...)
}

func (f *Formatter) sprintColor(c *color.Color, s string) string {
	if f.DisabledColor || c == nil {
		return fmt.Sprint(s)
	}
	return c.SprintFunc()(s)
}

func (f *Formatter) writeIndent(buf *bytes.Buffer, depth int) {
	buf.WriteString(strings.Repeat(f.Indent, depth))
}

func (f *Formatter) writeObjSep(buf *bytes.Buffer) {
	if f.Indent != "" {
		buf.WriteByte('\n')
	}
}

func (f *Formatter) Marshal(jsonObj interface{}) ([]byte, error) {
	buffer := bytes.Buffer{}
	f.marshalValue(jsonObj, &buffer, initialDepth)
	return buffer.Bytes(), nil
}

func (f *Formatter) MarshalIndent(jsonObj interface{}, indent string) ([]byte, error) {
    f.Indent = indent
	buffer := bytes.Buffer{}
	f.marshalValue(jsonObj, &buffer, initialDepth)
	return buffer.Bytes(), nil
}

func (f *Formatter) marshalOrderedMap(o orderedmap.OrderedMap, buf *bytes.Buffer, depth int) {
	remaining := len(o.Keys())
	if remaining == 0 {
		buf.WriteString(emptyMap)
		return
	}

    if f.SortKeys {
        o.SortKeys(sort.Strings)
    }

	buf.WriteString(startMap)
	f.writeObjSep(buf)

	for _, k := range o.Keys() {
		f.writeIndent(buf, depth+1)
        if f.Indent != "" {
            buf.WriteString(f.sprintfColor(f.KeyColor, "\"%s\": ", k))
        } else {
            buf.WriteString(f.sprintfColor(f.KeyColor, "\"%s\":", k))
        }
        v, _ := o.Get(k)
		f.marshalValue(v, buf, depth+1)
		remaining--
		if remaining != 0 {
			buf.WriteString(valueSep)
		}
		f.writeObjSep(buf)
    }
	f.writeIndent(buf, depth)
	buf.WriteString(endMap)
}

func (f *Formatter) marshalMap(m map[string]interface{}, buf *bytes.Buffer, depth int) {
	remaining := len(m)

	if remaining == 0 {
		buf.WriteString(emptyMap)
		return
	}

	keys := make([]string, 0)
	for key := range m {
		keys = append(keys, key)
	}

    if f.SortKeys {
    	sort.Strings(keys)
    }

	buf.WriteString(startMap)
	f.writeObjSep(buf)

	for _, key := range keys {
		f.writeIndent(buf, depth+1)
		buf.WriteString(f.KeyColor.Sprintf("\"%s\": ", key))
		f.marshalValue(m[key], buf, depth+1)
		remaining--
		if remaining != 0 {
			buf.WriteString(valueSep)
		}
		f.writeObjSep(buf)
	}
	f.writeIndent(buf, depth)
	buf.WriteString(endMap)
}

func (f *Formatter) marshalArray(a []interface{}, buf *bytes.Buffer, depth int) {
	if len(a) == 0 {
		buf.WriteString(emptyArray)
		return
	}

	buf.WriteString(startArray)
	f.writeObjSep(buf)

	for i, v := range a {
		f.writeIndent(buf, depth+1)
		f.marshalValue(v, buf, depth+1)
		if i < len(a)-1 {
			buf.WriteString(valueSep)
		}
		f.writeObjSep(buf)
	}
	f.writeIndent(buf, depth)
	buf.WriteString(endArray)
}

func (f *Formatter) marshalValue(val interface{}, buf *bytes.Buffer, depth int) {
	switch v := val.(type) {
	case orderedmap.OrderedMap:
		f.marshalOrderedMap(v, buf, depth)
	case map[string]interface{}:
		f.marshalMap(v, buf, depth)
	case []interface{}:
		f.marshalArray(v, buf, depth)
	case string:
		f.marshalString(v, buf)
	case float64:
		buf.WriteString(f.sprintColor(f.NumberColor, strconv.FormatFloat(v, 'f', -1, 64)))
	case bool:
		buf.WriteString(f.sprintColor(f.BoolColor, (strconv.FormatBool(v))))
	case nil:
		buf.WriteString(f.sprintColor(f.NullColor, null))
	case json.Number:
		buf.WriteString(f.sprintColor(f.NumberColor, v.String()))
	}
}

func (f *Formatter) marshalString(str string, buf *bytes.Buffer) {
	if !f.RawStrings {
		strBytes, _ := json.Marshal(str)
		str = string(strBytes)
	}

	if f.StringMaxLength != 0 && len(str) >= f.StringMaxLength {
		str = fmt.Sprintf("%s...", str[0:f.StringMaxLength])
	}

	buf.WriteString(f.sprintColor(f.StringColor, str))
}

