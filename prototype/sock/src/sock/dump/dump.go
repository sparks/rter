package dump

import (
	r "reflect"
	"fmt"
	"strconv"
	"io"
	"os"
)

var emptyString = ""

// Prints to the writer the value with indentation.
func Fdump(out io.Writer, v_ interface{}) {
	// forward decl
	var dump0 func(r.Value, int)
	var dump func(r.Value, int, *string, *string)

	done := make(map[string]bool)

	dump = func(v r.Value, d int, prefix *string, suffix *string) {
		pad := func() {
			res := ""
			for i := 0; i < d; i++ {
				res += "  "
			}
			fmt.Fprintf(out, res)
		}

		padprefix := func() {
			if prefix != nil {
				fmt.Fprintf(out, *prefix)
			} else {
				res := ""
				for i := 0; i < d; i++ {
					res += "  "
				}
				fmt.Fprintf(out, res)
			}
		}

		printv := func(o interface{}) { fmt.Fprintf(out, "%v", o) }

		printf := func(s string, args ...interface{}) { fmt.Fprintf(out, s, args) }

		// prevent circular for composite types
		switch o := v.(type) {
		case nil:
			// do nothing
		case *r.ArrayValue, *r.SliceValue, *r.MapValue, *r.PtrValue, *r.StructValue, *r.InterfaceValue:
			addr := v.Addr()
			key := fmt.Sprintf("%x %v", addr, v.Type())
			if _, exists := done[key]; exists {
				padprefix()
				printf("<%s>", key)
				return
			} else {
				done[key] = true
			}
		default:
			// do nothing
		}

		switch o := v.(type) {
		case *r.ArrayValue:
			padprefix()
			printf("[%d]%s {\n", o.Len(), o.Type().(*r.ArrayType).Elem())
			for i := 0; i < o.Len(); i++ {
				dump0(o.Elem(i), d+1)
				if i != o.Len()-1 {
					printf(",\n")
				}
			}
			print("\n")
			pad()
			print("}")

		case *r.SliceValue:
			padprefix()
			printf("[]%s (len=%d) {\n", o.Type().(*r.SliceType).Elem(), o.Len())
			for i := 0; i < o.Len(); i++ {
				dump0(o.Elem(i), d+1)
				if i != o.Len()-1 {
					printf(",\n")
				}
			}
			print("\n")
			pad()
			print("}")

		case *r.MapValue:
			padprefix()
			t := o.Type().(*r.MapType)
			printf("map[%s]%s {\n", t.Key(), t.Elem())
			for i, k := range o.Keys() {
				dump0(k, d+1)
				printf(": ")
				dump(o.Elem(k), d+1, &emptyString, nil)
				if i != o.Len()-1 {
					printf(",\n")
				}
			}
			print("\n")
			pad()
			print("}")

		case *r.PtrValue:
			padprefix()
			if o.Elem() == nil {
				printf("(*%s) nil", o.Type().(*r.PtrType).Elem())
			} else {
				print("&")
				dump(o.Elem(), d, &emptyString, nil)
			}

		case *r.StructValue:
			padprefix()
			t := o.Type().(*r.StructType)
			printf("%s {\n", t)
			d += 1
			for i := 0; i < o.NumField(); i++ {
				pad()
				printv(t.Field(i).Name)
				printv(": ")
				dump(o.Field(i), d, &emptyString, nil)
				if i != o.NumField()-1 {
					printf(",\n")
				}
			}
			d -= 1
			print("\n")
			pad()
			print("}")

		case *r.InterfaceValue:
			padprefix()
			t := o.Type().(*r.InterfaceType)
			printf("(%s) ", t)
			dump(o.Elem(), d, &emptyString, nil)

		case *r.StringValue:
			padprefix()
			printv(strconv.Quote(o.Get()))

		case *r.BoolValue,
			*r.IntValue, *r.Int8Value, *r.Int16Value, *r.Int32Value, *r.Int64Value,
			*r.UintValue, *r.Uint8Value, *r.Uint16Value, *r.Uint32Value, *r.Uint64Value,
			*r.FloatValue, *r.Float32Value, *r.Float64Value:
			padprefix()
			//printv(o.Interface());
			i := o.Interface()
			if stringer, ok := i.(interface {
				String() string
			}); ok {
				printf("(%v) %s", o.Type(), stringer.String())
			} else {
				printv(i)
			}

		case nil:
			padprefix()
			printv("nil")

		default:
			padprefix()
			printf("(%v) %v", o.Type(), o.Interface())
		}
	}

	dump0 = func(v r.Value, d int) { dump(v, d, nil, nil) }

	v := r.NewValue(v_)
	dump0(v, 0)
	fmt.Fprintf(out, "\n")
}

// Print to standard out the value that is passed as the argument with indentation.
// Pointers are dereferenced.
func Dump(v_ interface{}) { Fdump(os.Stdout, v_) }
