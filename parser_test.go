package concordance

import "testing"

func Test_normalize(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name                  string
		args                  args
		wantNormalized        string
		wantHasEOSPunctuation bool
		wantIgnore            bool
	}{
		{
			name: "empty",
			args: args{
				word: "",
			},
			wantNormalized:        "",
			wantHasEOSPunctuation: false,
			wantIgnore:            true,
		},
		{
			name: "word",
			args: args{
				word: "word",
			},
			wantNormalized:        "word",
			wantHasEOSPunctuation: false,
			wantIgnore:            false,
		},
		{
			name: "case",
			args: args{
				word: "Title-Case",
			},
			wantNormalized:        "title-case",
			wantHasEOSPunctuation: false,
			wantIgnore:            false,
		},
		{
			name: "period",
			args: args{
				word: "period.",
			},
			wantNormalized:        "period",
			wantHasEOSPunctuation: true,
			wantIgnore:            false,
		},
		{
			name: "abbreviationWithMultiplePeriods",
			args: args{
				word: "i.e.",
			},
			wantNormalized:        "i.e.",
			wantHasEOSPunctuation: false,
			wantIgnore:            false,
		},
		{
			name: "quotes",
			args: args{
				word: `"Hey,"`,
			},
			wantNormalized:        "hey",
			wantHasEOSPunctuation: false,
			wantIgnore:            false,
		},
		{
			name: "quotesEOS",
			args: args{
				word: `"Hey!"`,
			},
			wantNormalized:        "hey",
			wantHasEOSPunctuation: true,
			wantIgnore:            false,
		},
		{
			name: "onlyPunctuation",
			args: args{
				word: `"!"`,
			},
			wantNormalized:        "",
			wantHasEOSPunctuation: false,
			wantIgnore:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNormalized, gotHasEOSPunctuation, gotIgnore := normalize(tt.args.word)
			if gotNormalized != tt.wantNormalized {
				t.Errorf("normalize() gotNormalized = %v, want %v", gotNormalized, tt.wantNormalized)
			}
			if gotHasEOSPunctuation != tt.wantHasEOSPunctuation {
				t.Errorf("normalize() gotHasEOSPunctuation = %v, want %v", gotHasEOSPunctuation, tt.wantHasEOSPunctuation)
			}
			if gotIgnore != tt.wantIgnore {
				t.Errorf("normalize() gotIgnore = %v, want %v", gotIgnore, tt.wantIgnore)
			}
		})
	}
}
