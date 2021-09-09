package matcher

import (
	"strings"
	"testing"
)

func TestKeywordCase(t *testing.T) {
	var m = NewMatcher(nil, nil, 0)

	for _, kw := range m.positiveKeywords {
		if strings.ToLower(kw) != kw {
			t.Errorf("Keyword %q is not lowercase in positiveKeywords, but should be", kw)
		}
	}
	for _, kw := range m.locationPositiveKeywors {
		if strings.ToLower(kw) != kw {
			t.Errorf("Keyword %q is not lowercase in locationPositiveKeywors, but should be", kw)
		}
	}
	for _, kw := range m.negativeKeywords {
		if strings.ToLower(kw) != kw {
			t.Errorf("Keyword %q is not lowercase in negativeKeywords, but should be", kw)
		}
	}
	for _, kw := range m.importantUsers {
		if strings.ToLower(kw) != kw {
			t.Errorf("Name %q is not lowercase in importantUsers, but should be", kw)
		}
	}
}
