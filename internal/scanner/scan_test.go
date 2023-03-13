package scanner

import (
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestScanLine(t *testing.T) {
	defaultRegEx := regexp.MustCompile(`(?P<key>((private_key|public_key)\w+)|(private_key|public_key))`)

	tests := []struct {
		name  string
		regEx *regexp.Regexp
		line  string
		want  []string
	}{
		{
			name:  "find key",
			regEx: defaultRegEx,
			line:  "private_key public_key private_key_test",
			want:  []string{"private_key", "public_key", "private_key_test"},
		},
		{
			name:  "find combine name",
			regEx: defaultRegEx,
			line:  "private_key_public_key",
			want:  []string{"private_key_public_key"},
		},
		{
			name:  "find end of line",
			regEx: defaultRegEx,
			line:  "private_key public_key_test",
			want:  []string{"private_key", "public_key_test"},
		},
	}

	for _, test := range tests {
		t.Log(test.name)

		got := scanLine(test.regEx, test.line)
		if !cmp.Equal(got, test.want) {
			t.Fatalf("scanLine mismatch want: `%v`, got: `%v`", test.want, got)
		}
	}
}
