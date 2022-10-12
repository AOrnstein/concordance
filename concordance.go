// package concordance creates a alphabetical list of words and the sentences in which they occur
package concordance

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
)

// GenerateConcordance creates a concordance from the given document and writes it to the given output
func GenerateConcordance(document io.Reader, output io.Writer) error {
	concordance, err := Parse(document)
	if err != nil {
		return err
	}
	concordance.PrintTo(output)
	return nil
}

// Parse constructs a Concordance from a input
func Parse(document io.Reader) (*Concordance, error) {
	scanner := bufio.NewScanner(document)
	parser := newParser()
	for scanner.Scan() {
		parser.parseLine(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return parser.concordance, nil
}

// Concordance is a list of words and the sentences in which they occur
type Concordance struct {
	// occurrences is a map of words to the sentence in which the word was found, duplicates allowed
	occurrences map[string][]int
}

// Add a word and the sentence in which it was found
func (concordance *Concordance) Add(word string, sentence int) {
	concordance.occurrences[word] = append(concordance.occurrences[word], sentence)
}

func (concordance *Concordance) String() string {
	var sb strings.Builder
	concordance.PrintTo(&sb)
	return sb.String()
}

// PrintTo pretty prints the concordance to the output
func (concordance *Concordance) PrintTo(output io.Writer) {
	var (
		longestWord    int
		wordList       = make([]string, 0, len(concordance.occurrences))
		idxEnumeration num2MultiLetter
	)
	for word := range concordance.occurrences {
		wordList = append(wordList, word)
		if wordLen := len(word); wordLen >= longestWord {
			longestWord = wordLen
		}
	}
	sort.Strings(wordList)

	indexFmt := "%-" + fmt.Sprintf("%d", idxEnumeration.maxChars(len(wordList)+1)) + "s\t"
	wordFmt := "%-" + fmt.Sprintf("%d", longestWord) + "s\t"
	for i, word := range wordList {
		fmt.Fprintf(output, indexFmt, idxEnumeration.formatIndex(i)+".")
		fmt.Fprintf(output, wordFmt, word)

		// print sentence occurrences count and sentences in which the word appears
		occurrences := concordance.occurrences[word]
		fmt.Fprint(output, "{", len(occurrences), ":", occurrences[0]) // guaranteed to be at least one occurrence
		for _, sentence := range occurrences[1:] {
			fmt.Fprint(output, ",", sentence)
		}
		fmt.Fprint(output, "}\n")
	}
}

// num2MultiLetter helps format a index (0 indexed) to a string using repeating letters.
// i.e. 0 -> a , 26 -> aa
type num2MultiLetter struct{}

// maxChars returns the maximum width of the index
func (num2MultiLetter) maxChars(numWords int) int { return (numWords / 26) + 1 }

// formatIndex converts the given index (0 based) into a formatted string
func (num2MultiLetter) formatIndex(num int) string {
	letter := string([]byte{'a' + byte(num%26)})
	return strings.Repeat(letter, (num/26)+1)
}
