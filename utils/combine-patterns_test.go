package utils

import "testing"

func TestCombinePatterns_CompilesAndMatchesAnyAlternative(t *testing.T) {
	re := CombinePatterns([]string{"foo", "bar"})

	cases := []struct {
		in   string
		want bool
	}{
		{"foo", true},
		{"bar", true},
		{"xxfooyy", true},
		{"xxbaryy", true},
		{"baz", false},
		{"", false},
	}

	for _, tc := range cases {
		if got := re.MatchString(tc.in); got != tc.want {
			t.Fatalf("MatchString(%q)=%v, want %v (regex=%q)", tc.in, got, tc.want, re.String())
		}
	}
}

func TestCombinePatterns_WorksWithNamedSliceType(t *testing.T) {
	type ArrayFlag []string

	re := CombinePatterns(ArrayFlag{"cat", "dog"})

	if !re.MatchString("hotdog") {
		t.Fatalf("expected regex %q to match input %q", re.String(), "hotdog")
	}
	if re.MatchString("mouse") {
		t.Fatalf("expected regex %q to NOT match input %q", re.String(), "mouse")
	}
}

func TestCombinePatterns_PreservesRegexMetacharacters(t *testing.T) {
	// This test documents current behavior: patterns are treated as raw regex,
	// not literals (no escaping is applied).
	re := CombinePatterns([]string{`a.*b`, `^hello$`})

	if !re.MatchString("axxxb") {
		t.Fatalf("expected %q to match %q", re.String(), "axxxb")
	}
	if !re.MatchString("hello") {
		t.Fatalf("expected %q to match %q", re.String(), "hello")
	}
	if re.MatchString("hello!") {
		t.Fatalf("expected %q to NOT match %q", re.String(), "hello!")
	}
}

func TestCombinePatterns_EmptyPatternsProducesEmptyRegexWhichMatchesEverything(t *testing.T) {
	// strings.Join on an empty slice => ""
	// regexp "" matches at position 0 in any string.
	re := CombinePatterns([]string{})

	cases := []string{"", "anything", "123"}
	for _, in := range cases {
		if !re.MatchString(in) {
			t.Fatalf("expected empty regex to match %q", in)
		}
	}
}

func TestCombinePatterns_PanicsOnInvalidRegex(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic on invalid regex, got none")
		}
	}()

	_ = CombinePatterns([]string{"("}) // invalid regex => MustCompile panics
}

func TestCombinePatterns_ElementContainsAlternative(t *testing.T) {
	re := CombinePatterns([]string{"foo|bar", "baz"})

	cases := []struct {
		in   string
		want bool
	}{
		{"foo", true},
		{"bar", true},
		{"baz", true},
		{"qux", false},
	}

	for _, tc := range cases {
		if got := re.MatchString(tc.in); got != tc.want {
			t.Fatalf(
				"MatchString(%q)=%v, want %v (regex=%q)",
				tc.in, got, tc.want, re.String(),
			)
		}
	}
}

func TestCombinePatterns_MultipleElementsWithAlternatives(t *testing.T) {
	re := CombinePatterns([]string{
		"foo|bar",
		"baz|qux",
	})

	cases := []struct {
		in   string
		want bool
	}{
		{"foo", true},
		{"bar", true},
		{"baz", true},
		{"qux", true},
		{"nope", false},
	}

	for _, tc := range cases {
		if got := re.MatchString(tc.in); got != tc.want {
			t.Fatalf(
				"MatchString(%q)=%v, want %v (regex=%q)",
				tc.in, got, tc.want, re.String(),
			)
		}
	}
}

func TestCombinePatterns_AlternativePrecedenceIsNotIsolated(t *testing.T) {
	re := CombinePatterns([]string{
		"ab|cd",
		"efg",
	})

	if !re.MatchString("cd") {
		t.Fatalf("expected %q to match %q", re.String(), "cd")
	}
	if !re.MatchString("efg") {
		t.Fatalf("expected %q to match %q", re.String(), "efg")
	}
	if re.MatchString("def") {
		t.Fatalf("did not expect %q to match %q", re.String(), "cdef")
	}
}

func TestCombinePatterns_AlternativeWithAnchors(t *testing.T) {
	re := CombinePatterns([]string{
		"^foo|bar$",
		"baz",
	})

	cases := []struct {
		in   string
		want bool
	}{
		{"foo", true}, // ^foo
		{"bar", true}, // bar$
		{"baz", true},
		{"foobar", true},   // matches ^foo
		{"xxbar", true},    // matches bar$
		{"xxbazxx", true},  // contains baz
		{"xxbarxx", false}, // bar not at the end
	}

	for _, tc := range cases {
		if got := re.MatchString(tc.in); got != tc.want {
			t.Fatalf(
				"MatchString(%q)=%v, want %v (regex=%q)",
				tc.in, got, tc.want, re.String(),
			)
		}
	}
}
