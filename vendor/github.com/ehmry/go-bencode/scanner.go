package bencode

import "strconv"

// checkValid verifies that data is valid bencode encoded data.
// scan is passed in for use by checkValid to avoid an allocation.
func checkValid(data []byte, scan *scanner) error {
	scan.reset()
	var c, op int
	for i := 0; i < len(data); i++ {
		c = int(data[i])
		scan.bytes++
		op = scan.step(scan, int(c))
		if op > 0 {
			scan.bytes += int64(op)
			i += op
		} else if op == scanError {
			return scan.err
		}
	}
	if scan.eof() == scanError {
		return scan.err
	}
	return nil
}

// nextValue splits data after the next whole bencode value,
// returning thet value and the bytes that follow it as seperate slices.
// scan is passed in for use by nextValue to avoid an allocation.
func nextValue(data []byte, scan *scanner) (value, rest []byte, err error) {
	scan.reset()
	var c, op int
	for i := 0; i < len(data); i++ {
		c = int(data[i])
		op = scan.step(scan, c)
		if op > 0 {
			i += op
		} else if op <= scanEnd {
			switch op {
			case scanError:
				return nil, nil, scan.err
			case scanEnd:
				return data[0:i], data[i:], nil
			}
		}
	}
	if scan.eof() == scanError {
		return nil, nil, scan.err
	}
	return data, nil, nil
}

// A SyntaxError is a description of a becode syntax error.
type SyntaxError struct {
	msg    string // description of error
	Offset int64  // error occurred after reading Offset bytes
}

func (e *SyntaxError) Error() string { return e.msg }

// These values are returned by the state transition functions
// assigned to scanner.state and the method scanner.eof.
// They give details about the current state of the scan that
// callers might be intererested to know about.
// It is ok to ignore the return value of any particular
// call to scanner.state: if one cal returns scanError,
// every subsequent call will retern scanError too.
const (
	// Continue.
	scanBeginInteger   = 0 - iota
	scanParseInteger   = 0 - iota
	scanEndInteger     = 0 - iota
	scanBeginStringLen = 0 - iota
	scanParseStringLen = 0 - iota
	scanEndStringLen   = 0 - iota
	scanParseString    = 0 - iota
	scanBeginList      = 0 - iota
	scanEndList        = 0 - iota
	scanEndValue       = 0 - iota
	scanBeginDict      = 0 - iota
	scanBeginDictKey   = 0 - iota
	scanDictKey        = 0 - iota
	scanBeginKeyLen    = 0 - iota
	scanParseKeyLen    = 0 - iota
	scanEndKeyLen      = 0 - iota
	scanParseKey       = 0 - iota
	scanDictValue      = 0 - iota
	scanEndDict        = 0 - iota

	// Stop.
	scanEnd   = 0 - iota
	scanError = 0 - iota // hit an error, scanner.err.
)

// These values are stored in the parseState stack.
// They give the current state of a composite value
// being scanned. If the parser is inside a nested value
// the parseState describes the nested state, outermost at entry 0.
const (
	parseInteger   = iota // parsing an integer
	parseString           // parsing a string
	parseDictKey          // parsing dict key
	parseDictValue        // parsing dict value
	parseListValue        // parsing list value
)

// reset prepares the scanner for use.
// It must be called before calling s.step.
func (s *scanner) reset() {
	s.step = stateBeginValue
	s.parseState = s.parseState[0:0]
	s.err = nil
	s.endTop = false
}

// stateError is the state after reaching a syntax error.
func stateError(s *scanner, c int) int {
	return scanError
}

// error records an error and switches to the error state.
func (s *scanner) error(c int, context string) int {
	s.step = stateError
	s.err = &SyntaxError{"invalid character " + strconv.Quote(string(c)) + " " + context, s.bytes}
	return scanError
}

// eof tells the scanner that the end of input has been reached.
// It returns a scan status just as s.step does.
func (s *scanner) eof() int {
	if s.err != nil {
		return scanError
	}
	if s.endTop {
		return scanEnd
	}
	s.step(s, 'e')
	if s.endTop {
		return scanEnd
	}
	if s.err == nil {
		s.err = &SyntaxError{"unexpected end of becode input", s.bytes}
	}
	return scanError
}

// pushParseState pushes a new parse state p onto the parse stack.
func (s *scanner) pushParseState(p int) {
	s.parseState = append(s.parseState, p)
}

// popParseState pops a parse state (alread obtained) off the stack
// and updates s.step accordingly.
func (s *scanner) popParseState() {
	n := len(s.parseState) - 1
	s.parseState = s.parseState[0:n]
	if n == 0 {
		s.endTop = true
	}
	s.step = stateEndValue
}


// stateParseInteger is the state after reading an `i`.
func stateParseInteger(s *scanner, c int) int {
	if c == 'e' {
		s.popParseState()
		if s.endTop {
			return scanEnd
		}
		return scanEndInteger
	}
	if (c >= '0' && c <= '9') || c == '-' {
		return scanParseInteger
	}
	return s.error(c, "in integer")
}

func stateBeginListValue(s *scanner, c int) int {
	if c == 'e' {
		s.popParseState()
		if s.endTop {
			return scanEnd
		}
		return scanEndList
	}
	return stateBeginValue(s, c)
}

// stateEndValue is the state after completing a value,
// such as after reading 'e' or finishing a string.
func stateEndValue(s *scanner, c int) int {
	n := len(s.parseState)
	if n == 0 {
		// Completed top-level before the current byte.
		s.step = stateBeginValue
		s.endTop = true
		return scanEnd
	}
	ps := s.parseState[n-1]
	switch ps {
	case parseDictKey:
		s.step = stateBeginValue
		return scanDictValue

	case parseDictValue:
		s.step = stateBeginDictKey
		return stateBeginDictKey(s, c)

	case parseListValue:
		if c == 'e' {
			s.popParseState()
			if s.endTop {
				return scanEnd
			}
			return scanEndList
		}
		s.step = stateBeginValue
		return stateBeginValue(s, c)
	}
	return s.error(c, "")
}
