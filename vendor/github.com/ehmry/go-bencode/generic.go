// +build 386 amd64
// build flags for little-endian architecture

package bencode

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
	strLen int

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
		s.strLen = (c & 0xcf)
		s.step = stateParseStringLen
		s.pushParseState(parseString)
		return scanBeginStringLen
	}
	return s.error(c, "looking for beginning of value")
}

func stateParseStringLen(s *scanner, c int) int {
	if c == ':' {
		// decoder should read this string as a slice
		s.popParseState()
		// BUG(emery): undefined behavior with top level strings

		// this is a problem, if this string is a top-level object,
		// the fact that this scanner has reached the end isn't communicated.
		// I guess I could shift the string length and set scanEnd bits
		return s.strLen
	}
	if c >= '0' && c <= '9' {
		s.strLen *= 10
		s.strLen += (c & 0xcf)
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
		s.strLen = (c & 0xcf)
		s.step = stateParseKeyLen
		return scanBeginKeyLen
	}
	return s.error(c, "in start of dictionary key length")
}

func stateParseKeyLen(s *scanner, c int) int {
	if c == ':' {
		s.step = stateBeginValue

		// decoder should read this key chunk at once
		return s.strLen
	}
	if c >= '0' && c <= '9' {
		s.strLen *= 10
		s.strLen += (c & 0xcf)
		return scanParseKeyLen
	}
	return s.error(c, "in dictionary key length")
}
