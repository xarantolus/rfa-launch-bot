package matcher

import "testing"

func Test_anyWordStartsWith(t *testing.T) {
	var searchedPrefixes = []string{"test", "best", "rest", "more than one word"}

	tests := []struct {
		argText string
		want    bool
	}{
		{"testing is nice", true},
		{"wrongprefixtesting is nice", false},
		{"the test keyword can be at any point in the string", true},
		{"we want to support more than one word", true},
		{"we want to support less than one word", false},

		{"the #test hashtag should still be recognized", true},
		{"also @test should work", true},
		{"#test at the beginning", true},
		{"@test should work", true},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := anyWordStartsWith(tt.argText, searchedPrefixes...); got != tt.want {
				t.Errorf("containsAny(%q) = %v, want %v", tt.argText, got, tt.want)
			}
		})
	}
}

func Test_containsAnyGeneric(t *testing.T) {
	var searchedInfixes = []string{"test", "best", "rest", "more than one word"}

	tests := []struct {
		argText string
		want    bool
	}{
		{"testing is nice", true},
		{"wrongprefixtesting is nice", true},
		{"the test keyword can be at any point in the string", true},
		{"we want to support more than one word", true},
		{"we want to support less than one word", false},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := textContainsAny(tt.argText, searchedInfixes...); got != tt.want {
				t.Errorf("containsAny(%q) = %v, want %v", tt.argText, got, tt.want)
			}
		})
	}
}

func Test_containsStringCaseInsensitive(t *testing.T) {
	var slice = []string{"test", "best", "rest", "more than one word"}

	var contains = []string{
		"test", "Test", "tEst",
		"rest", "REST",
	}

	var notContains = []string{
		"not test",
		"one word",
		"anything else",
	}

	for _, positive := range contains {
		if !containsStringCaseInsensitive(slice, positive) {
			t.Errorf("Slice %v contains %q, but wasn't detected", slice, positive)
		}
	}
	for _, negative := range notContains {
		if containsStringCaseInsensitive(slice, negative) {
			t.Errorf("Slice %v doesn't contain %q, but was detected", slice, negative)
		}
	}
}
