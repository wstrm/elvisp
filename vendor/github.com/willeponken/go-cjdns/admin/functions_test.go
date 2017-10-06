// +build functions

package admin

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
	"unicode"
	"unicode/utf8"
)

// TODO(inhies): Cleanup this test.
type support struct {
	Cjdns   bool
	Package bool
}

// Compare functions that Conn exports versus what cjdns tells us is available.
func Test_FunctionCoverage(t *testing.T) {
	conn, err := Connect(nil)

	if err != nil {
		t.Fatal(err)
	}

	cjdnsFuncs, err := conn.Admin_availableFunctions()
	if err != nil {
		t.Fatal(err)
	}

	v := reflect.TypeOf(&Conn{})
	var ourFuncs = make(map[string]support)

	// Loop over each of Conn's methods.
	for i := 0; i < v.NumMethod(); i++ {
		ourFunc := v.Method(i).Name

		// Get the first character as a rune.
		r, n := utf8.DecodeRuneInString(ourFunc)
		// Only check Conn's exported methods.
		if !unicode.IsLower(r) {
			// Check to see if cjdns supports our function.
			if cjdnsFuncs[ourFunc] != nil {
				ourFuncs[ourFunc] = support{true, true}
				delete(cjdnsFuncs, ourFunc)
			} else {
				// Convert first letter to lowercase, then check for support
				// from cjdns. This is for the 'ping' and 'memory' functions.
				lowerFunc := string(unicode.ToLower(r)) + ourFunc[n:]
				if cjdnsFuncs[lowerFunc] != nil {
					ourFuncs[lowerFunc] = support{true, true}
					delete(cjdnsFuncs, lowerFunc)
				} else {
					// cjdns does not support this.
					ourFuncs[ourFunc] = support{false, true}
				}
			}
		}
	}

	funcs := make([]string, 0)
	for k, _ := range cjdnsFuncs {
		funcs = append(funcs, k)
	}
	if len(funcs) > 0 {
		fmt.Println("Functions supported by cjdns but not this package:")
		sort.Strings(funcs)
		for _, f := range funcs {
			fmt.Println("   ", f)
		}
		fmt.Println()
	}

	funcs = make([]string, 0)
	for k, v := range ourFuncs {
		if v.Cjdns == false && v.Package {
			funcs = append(funcs, k)
		}
	}
	if len(funcs) > 0 {
		fmt.Println("Functions supported by this package but not cjdns:")
		sort.Strings(funcs)
		for _, f := range funcs {
			fmt.Println("   ", f)
		}
		fmt.Println()
	}
}
