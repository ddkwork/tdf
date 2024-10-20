package tdf

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"math"
	"strings"
	"time"

	"github.com/ddkwork/golibrary/mylog"
)

// PrintEncoder is a TdfEncoder that writes TDF data to an io.Writer in a human-readable format.
type PrintEncoder struct {
	w          io.Writer
	indent     int
	indentStr  string
	pretty     bool
	escapeHTML bool
}

// NewPrintEncoder creates a new PrintEncoder that writes to the given io.Writer.
func NewPrintEncoder(w io.Writer) *PrintEncoder {
	return &PrintEncoder{
		w:          w,
		indentStr:  "  ",
		pretty:     true,
		escapeHTML: true,
	}
}

// SetIndent sets the number of spaces to use for each level of indentation.
func (e *PrintEncoder) SetIndent(indent int) {
	e.indent = indent
}

// SetPretty sets whether or not to pretty-print the output.
func (e *PrintEncoder) SetPretty(pretty bool) {
	e.pretty = pretty
}

// SetEscapeHTML sets whether or not to escape HTML characters in string values.
func (e *PrintEncoder) SetEscapeHTML(escapeHTML bool) {
	e.escapeHTML = escapeHTML
}

// Encode writes the given Tdf to the underlying io.Writer in a human-readable format.
func (e *PrintEncoder) Encode(tdf any) error {
	return e.encodeTdf(tdf, 0)
}

func (e *PrintEncoder) encodeTdf(tdf any, level int) error {
	//switch tdf.Type {
	//case TdfTypeNull:
	//	return e.write("null")
	//case TdfTypeBoolean:
	//	return e.writeBool(tdf.BoolValue)
	//case TdfTypeByte:
	//	return e.writeInt(int64(tdf.ByteValue))
	//case TdfTypeShort:
	//	return e.writeInt(int64(tdf.ShortValue))
	//case TdfTypeInteger:
	//	return e.writeInt(int64(tdf.IntValue))
	//case TdfTypeLong:
	//	return e.(tdf.LongValue)
	//case TdfTypeFloat:
	//	return e.writeFloat(tdf.FloatValue)
	//case TdfTypeString:
	//	return e.writeString(tdf.StringValue)
	//case TdfTypeBlob:
	//	return e.writeBlob(tdf.BlobValue)
	//case TdfTypeTime:
	//	return e.writeTime(tdf.TimeValue)
	//case TdfTypeObject:
	//	return e.writeObject(tdf.ObjectValue)
	//case TdfTypeArray:
	//	return e.writeArray(tdf.ArrayValue, level)
	//case TdfTypeMap:
	//	return e.writeMap(tdf.MapValue, level)
	//case TdfTypeUnion:
	//	return e.writeUnion(tdf.UnionValue, level)
	//case TdfTypeVariableContainer:
	//	return e.writeVariableContainer(tdf.VariableContainerValue, level)
	//default:
	//	return fmt.Errorf("unknown TdfType: %v", tdf.Type)
	//}
	return nil
}

func (e *PrintEncoder) writeBool(b bool) error {
	if b {
		return e.write("true")
	}
	return e.write("false")
}

func (e *PrintEncoder) writeInt(i int64) error {
	return e.write(fmt.Sprintf("%d", i))
}

func (e *PrintEncoder) writeFloat(f float32) error {
	if math.IsNaN(float64(f)) {
		return e.write("NaN")
	}
	if math.IsInf(float64(f), 1) {
		return e.write("Infinity")
	}
	if math.IsInf(float64(f), -1) {
		// return e.writeInfinity
	}
	return e.write(fmt.Sprintf("%g", f))
}

func (e *PrintEncoder) writeString(s string) error {
	if e.escapeHTML {
		s = escapeHTML(s)
	}
	return e.write(fmt.Sprintf("\"%s\"", s))
}

func (e *PrintEncoder) writeBlob(b []byte) error {
	return e.write(fmt.Sprintf("blob(%d)", len(b)))
}

func (e *PrintEncoder) writeTime(t time.Time) error {
	return e.write(fmt.Sprintf("time(%d)", t.Unix()))
}

func (e *PrintEncoder) writeObject(obj *BlazeObjectId) error {
	// return e.write(fmt.Sprintf("object(%d,%d)", obj.Type, obj.Id))
	return nil
}

func (e *PrintEncoder) writeArray(arr []any, level int) error {
	if len(arr) == 0 {
		return e.write("[]")
	}
	if e.pretty {
		mylog.Check(e.write("[\n"))
		for _, tdf := range arr {
			mylog.Check(e.writeIndent(level + 1))
			mylog.Check(e.encodeTdf(tdf, level+1))
			mylog.Check(e.write(",\n"))
		}
		mylog.Check(e.writeIndent(level))
		return e.write("]")
	}
	var buf bytes.Buffer
	for _, tdf := range arr {
		mylog.Check(e.encodeTdf(tdf, level+1))
		mylog.Check(e.write(","))
	}
	return e.write(fmt.Sprintf("[%s]", strings.TrimRight(buf.String(), ",")))
}

//func (e *PrintEncoder) writeMap(m *TagInfoMap, level int) error {
//	if m == nil || len(m.Tags) == 0 {
//		return e.write("{}")
//	}
//	if e.pretty {
//		if err := e.write("{\n"); err != nil {
//			return err
//		}
//		keys := make([]string, 0, len(m.Tags))
//		for key := range m.Tags {
//			keys = append(keys, key)
//		}
//		sort.Strings(keys)
//		_, key := range keys{
//			tag, := m.Tags[key]
//			if err := e.writeIndent(level + 1); err != nil{
//			return err
//		}
//			if err := e.writeString(key); err != nil{
//			return err
//		}
//			if err := e.write(": "); err != nil{
//			return err
//		}
//			if err := e.encodeTdf(tag.Value, level+1); err != nil{
//			return err
//		}
//			if err := e.write(",\n"); err != nil{
//			return err
//		}
//		}
//		if err := e.writeIndent(level); err != nil {
//			return err
//		}
//		return e.write("}")
//	}
//	var buf bytes.Buffer
//	keys := make([]string, 0, len(m.Tags))
//	for key := range m.Tags {
//		keys = append(keys, key)
//	}
//	sort.Strings(keys)
//	for _, key := range keys {
//		tag := m.Tags[key]
//		if err := e.writeString(key); err != nil {
//			return err
//		}
//		if err := e.write(":"); err != nil {
//			return err
//		}
//		if err := e.encodeTdf(tag.Value, level+1); err != nil {
//			return err
//		}
//		if err := e.write(","); err != nil {
//			return err
//		}
//	}
//	return e.write(fmt.Sprintf("{%s}", strings.TrimRight(buf.String(), ",")))
//}

func (e *PrintEncoder) writeUnion(u *Union, level int) error {
	if u == nil {
		return e.write("null")
	}
	if e.pretty {
		mylog.Check(e.write("{\n"))
		mylog.Check(e.writeIndent(level + 1))
		//if err := e.writeString(u.Type); err != nil {
		//	return err
		//}
		mylog.Check(e.write(": "))
		//if err := e.encodeTdf(u.Value, level+1); err != nil {
		//	return err
		//}
		mylog.Check(e.write("\n"))
		mylog.Check(e.writeIndent(level))
		return e.write("}")
	}
	//if err := e.writeString(u.Type); err != nil {
	//	return err
	//}
	mylog.Check(e.write(":"))
	// return e.encodeTdf(u.Value, level+1)
	return nil
}

//func (e *PrintEncoder) writeVariableContainer(vc *VariableTdfContainer, level int) error {
//	if vc == nil || len(vc.Tags) == 0 {
//		return e.write("{}")
//	}
//	if e.pretty {
//		if err := e.write("{\n"); err != nil {
//			return err
//		}
//		for _, tag := range vc.Tags {
//			if err := e.writeIndent(level + 1); err != nil {
//				return err
//			}
//			if err := e.writeString(tag.Name); err != nil {
//				return err
//			}
//			if err := e.write(": "); err != nil {
//				return err
//			}
//			if err := e.encodeTdf(tag.Value, level+1); err != nil {
//				return err
//			}
//			if err := e.write(",\n"); err != nil {
//				return err
//			}
//		}
//		if err := e.writeIndent(level); err != nil {
//			return err
//		}
//		return e.write("}")
//	}
//	var buf bytes.Buffer
//	for _, tag := range vc.Tags {
//		if err := e.writeString(tag.Name); err != nil {
//			return err
//		}
//		if err := e.write(":"); err != nil {
//			return err
//		}
//		if err := e.encodeTdf(tag.Value, level+1); err != nil {
//			return err
//		}
//		if err := e.write(","); err != nil {
//			return err
//		}
//	}
//	return e.write(fmt.Sprintf("{%s}", strings.TrimRight(buf.String(), ",")))
//}

func (e *PrintEncoder) writeIndent(level int) error {
	if !e.pretty {
		return nil
	}
	return e.write(strings.Repeat(e.indentStr, level))
}

func (e *PrintEncoder) write(s string) error {
	mylog.Check2(io.WriteString(e.w, s))
	return nil
}

func escapeHTML(s string) string {
	return html.EscapeString(s)
}
