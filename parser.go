package concordance

import (
	"strings"
	"unicode"
)

// parser lexes, parses, normalizes, tokenizes, and processes a document to construct a concordance.
type parser struct {
	sentence       int
	insideSentence bool
	foundNewline   bool

	concordance *Concordance
}

// newParser creates a new parser.
func newParser() *parser {
	return &parser{
		concordance: &Concordance{occurrences: map[string][]int{}},
	}
}

// parseLine parses and processes a single line
func (parser *parser) parseLine(line string) {
	parser.lineStart()
	words := strings.Fields(line)
	for _, word := range words {
		parser.parseWord(word)
	}
}

// lineStart processes a new line character.
func (parser *parser) lineStart() {
	// only end a sentence from newlines if two or more newlines in a row
	// this is to handle potential formatting, this may not be wanted
	if parser.insideSentence && parser.foundNewline {
		parser.insideSentence = false
	}
	parser.foundNewline = true
}

// parseWord parses and processes a given word.
func (parser *parser) parseWord(word string) {
	normalized, hasEOSPunctuation, ignore := normalize(word)
	if ignore {
		return
	}
	if len(normalized) > 0 {
		if !parser.insideSentence {
			parser.sentence++
			parser.insideSentence = true
		}
		parser.foundNewline = false
		parser.concordance.Add(normalized, parser.sentence)
	}
	if hasEOSPunctuation {
		parser.insideSentence = false
	}
}

// normalize converts a raw character white-spaced word to a likely word token.
// A more advanced algorithm would return a slice of tokens.
func normalize(word string) (normalized string, hasEOSPunctuation bool, ignore bool) {
	if len(word) == 0 {
		ignore = true
		return
	}
	normalized = strings.ToLower(word)
	normalized = strings.TrimRightFunc(normalized, unicode.IsPunct)

	if len(normalized) != len(word) {

		// word was only punctuation
		if len(normalized) == 0 {
			ignore = true
			return
		}

		// find ending punctuation assuming correct simple grammar (e.g. period inside of quotes)
		punct := word[len(normalized)]

		normalized = strings.TrimLeftFunc(normalized, unicode.IsPunct)

		// check if word is likely an abbreviation (with multiple periods like "i.e.") which should keep a period and is not an EOS
		// doesn't catch abbreviations like "mr."" or "etc."
		if punct == '.' && strings.Count(normalized, ".") > 0 {
			normalized += "."
			return
		}
		switch punct {
		case '.', '?', '!':
			hasEOSPunctuation = true
		case '\'':
			// check if word is likely possessive
			if normalized[len(normalized)-1] == 's' {
				normalized += "'"
			}
		}
	}

	// strip punctuation like quotes
	normalized = strings.TrimLeftFunc(normalized, unicode.IsPunct)

	return
}
