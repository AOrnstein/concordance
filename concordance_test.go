package concordance

import (
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"
)

const sampleText = `Given an arbitrary text document written in English, write a program that will generate a
concordance, i.e. an alphabetical list of all word occurrences, labeled with word
frequencies.
Bonus: label each word with the sentence numbers in which each occurrence appeared.
`

var sampleResult = `a.    a              {2:1,1}
b.    all            {1:1}  
c.    alphabetical   {1:1}  
d.    an             {2:1,1}
e.    appeared       {1:2}  
f.    arbitrary      {1:1}  
g.    bonus          {1:2}  
h.    concordance    {1:1}  
i.    document       {1:1}  
j.    each           {2:2,2}
k.    english        {1:1}  
l.    frequencies    {1:1}  
m.    generate       {1:1}  
n.    given          {1:1}
o.    i.e.           {1:1}
p.    in             {2:1,2}
q.    label          {1:2}
r.    labeled        {1:1}
s.    list           {1:1}
t.    numbers        {1:2}
u.    occurrence     {1:2}
v.    occurrences    {1:1}
w.    of             {1:1}
x.    program        {1:1}
y.    sentence       {1:2}
z.    text           {1:1}
aa.   that           {1:1}
bb.   the            {1:2}
cc.   which          {1:2}
dd.   will           {1:1}
ee.   with           {2:1,2}
ff.   word           {3:1,1,2}
gg.   write          {1:1}
hh.   written        {1:1}
`

// getSampleResult with correct tabbing and cashes the correct one
func getSampleResult(t *testing.T) string {
	if strings.ContainsAny(sampleResult, "\t") {
		return sampleResult
	}
	sampleResultText, err := os.ReadFile("testing/sampleResults.txt")
	if err != nil {
		t.Fatal("unable to open sample result file: ", err)
	}
	sampleResult = string(sampleResultText)
	return sampleResult
}

var sampleConcordance = &Concordance{occurrences: map[string][]int{
	"a":            {1, 1},
	"all":          {1},
	"alphabetical": {1},
	"an":           {1, 1},
	"appeared":     {2},
	"arbitrary":    {1},
	"bonus":        {2},
	"concordance":  {1},
	"document":     {1},
	"each":         {2, 2},
	"english":      {1},
	"frequencies":  {1},
	"generate":     {1},
	"given":        {1},
	"i.e.":         {1},
	"in":           {1, 2},
	"label":        {2},
	"labeled":      {1},
	"list":         {1},
	"numbers":      {2},
	"occurrence":   {2},
	"occurrences":  {1},
	"of":           {1},
	"program":      {1},
	"sentence":     {2},
	"text":         {1},
	"that":         {1},
	"the":          {2},
	"which":        {2},
	"will":         {1},
	"with":         {1, 2},
	"word":         {1, 1, 2},
	"write":        {1},
	"written":      {1},
}}

func TestGenerateConcordance(t *testing.T) {
	document, err := os.Open("testing/sampleRequirements.txt")
	if err != nil {
		t.Fatal("unable to open sample file: ", err)
	}
	defer document.Close()

	sampleResultText := getSampleResult(t)

	var (
		doc    = strings.NewReader(sampleText)
		output strings.Builder
	)

	err = GenerateConcordance(doc, &output)
	if err != nil {
		t.Fatal("unable to generate concordance: ", err)
	}

	if got := output.String(); got != string(sampleResultText) {
		t.Errorf("Parse() = \n%v\n, want:\n%v", got, string(sampleResultText))
	}
}

func TestParse(t *testing.T) {
	type args struct {
		document io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *Concordance
		wantErr bool
	}{
		{
			name: "unableToParse",
			args: args{
				document: iotest.ErrReader(errors.New("custom error")),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "simple",
			args: args{
				document: strings.NewReader(sampleText),
			},
			want:    sampleConcordance,
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				document: strings.NewReader(""),
			},
			want:    &Concordance{occurrences: map[string][]int{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.document)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestConcordance_PrintTo(t *testing.T) {
	tests := []struct {
		name        string
		concordance *Concordance
		wantOutput  string
	}{
		{
			name: "emptyConcordance",
			concordance: &Concordance{
				occurrences: map[string][]int{},
			},
			wantOutput: "",
		},
		{
			name:        "filledConcordance",
			concordance: sampleConcordance,
			wantOutput:  getSampleResult(t),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			tt.concordance.PrintTo(output)

			gotOutput := output.String()
			gotOutputFields := strings.Fields(gotOutput)
			wantOutputFields := strings.Fields(tt.wantOutput)

			if !reflect.DeepEqual(gotOutputFields, wantOutputFields) {
				t.Errorf("strings.Fields(Concordance.PrintTo()) = \n%v\n, want \n%v", wantOutputFields, wantOutputFields)
			}
			if gotOutput := output.String(); gotOutput != tt.wantOutput {
				t.Errorf("Concordance.PrintTo() = \n%#v\n, want:\n%#v", gotOutput, tt.wantOutput)
			}
		})
	}
}
