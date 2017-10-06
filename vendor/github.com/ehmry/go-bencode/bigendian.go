// +build !386,!amd64
// build flags for unknown endian architecture

package bencode

import "strconv"

type scanner struct {
	// The step is a func to be called to execute the next transition.
	step func(*scanner, int) int

	// Reached end of top-level value.
	endTop bool

	// Stack of what we're in the middle of.
	parseState []int

	// Error that happened, if any.
	err error

	// storage for string length numeral bytes
	strLenB []byte

	// total bytes consumed, updated by decoder.Decode
	bytes int64
}

// stateBeginValue is the state at the beginning of the input.
func stateBeginValue(s *scanner, c int) int {
	switch c {
	case 'i':
		s.step = stateParseInteger
		s.pushParseState(parseInteger)
		return scanBeginInteger
	case 'l':
		s.step = stateBeginListValue
		s.pushParseState(parseListValue)
		return scanBeginList
	case 'd':
		s.step = stateBeginDictKey
		s.pushParseState(parseDictValue)
		return scanBeginDict
	}

	if c >= '0' && c <= '9' {
		s.strLenB = append(s.strLenB[0:0], byte(c))
		s.step = stateParseStringLen
		s.pushParseState(parseString)
		return scanBeginStringLen
	}
	return s.error(c, "looking for beginning of value")
}

func stateParseStringLen(s *scanner, c int) int {
	if c == ':' {
		l, err := strconv.Atoi(string(s.strLenB))
		if err != nil {
			s.err = err
			return scanError
		}
		// decoder should read this string as a slice
		s.popParseState()
		// BUG(emery): undefined behavior with top level strings

		// this is a problem, if this string is a top-level object,
		// the fact that this scanner has reached the end isn't communicated.
		// I guess I could shift the string length and set scanEnd bits
		return l
	}
	if c >= '0' && c <= '9' {
		s.strLenB = append(s.strLenB, byte(c))
		return scanParseStringLen
	}
	return s.error(c, "in string length")
}

func stateBeginDictKey(s *scanner, c int) int {
	if c == 'e' {
		s.popParseState() // pop parseDictValue
		if s.endTop {
			return scanEnd
		}
		return scanEndDict
	}
	if c >= '0' && c <= '9' {
		s.strLenB = append(s.strLenB[0:0], byte(c))
		s.step = stateParseKeyLen
		return scanBeginKeyLen
	}
	return s.error(c, "in start of dictionary key length")
}

func stateParseKeyLen(s *scanner, c int) int {
	if c == ':' {
		l, err := strconv.Atoi(string(s.strLenB))
		if err != nil {
			s.err = err
			return scanError
		}
		// decoder should read this chunk at once
		s.step = stateBeginValue
		return l
	}
	if c >= '0' && c <= '9' {
		s.strLenB = append(s.strLenB, byte(c))
		return scanParseKeyLen
	}
	return s.error(c, "in dictionary key length")
}
