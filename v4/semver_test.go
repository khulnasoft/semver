package semver

import (
	"testing"
)

func prstr(s string) PRVersion {
	return PRVersion{s, 0, false}
}

func prnum(i uint64) PRVersion {
	return PRVersion{"", i, true}
}

type formatTest struct {
	v      Version
	result string
}

var formatTests = []formatTest{
	{Version{1, 2, 3, -1, nil, nil}, "1.2.3"},
	{Version{0, 0, 1, -1, nil, nil}, "0.0.1"},
	{Version{0, 0, 1, -1, []PRVersion{prstr("alpha"), prstr("preview")}, []string{"123", "456"}}, "0.0.1-alpha.preview+123.456"},
	{Version{1, 2, 3, -1, []PRVersion{prstr("alpha"), prnum(1)}, []string{"123", "456"}}, "1.2.3-alpha.1+123.456"},
	{Version{1, 2, 3, -1, []PRVersion{prstr("alpha"), prnum(1)}, nil}, "1.2.3-alpha.1"},
	{Version{1, 2, 3, -1, nil, []string{"123", "456"}}, "1.2.3+123.456"},
	{Version{1, 2, 3, 4, nil, nil}, "1.2.3.4"},
	{Version{0, 0, 1, 2, nil, nil}, "0.0.1.2"},
	{Version{0, 0, 1, 2, []PRVersion{prstr("alpha"), prstr("preview")}, []string{"123", "456"}}, "0.0.1.2-alpha.preview+123.456"},
	{Version{1, 2, 3, 4, []PRVersion{prstr("alpha"), prnum(1)}, []string{"123", "456"}}, "1.2.3.4-alpha.1+123.456"},
	{Version{1, 2, 3, 4, []PRVersion{prstr("alpha"), prnum(1)}, nil}, "1.2.3.4-alpha.1"},
	{Version{1, 2, 3, 4, nil, []string{"123", "456"}}, "1.2.3.4+123.456"},
	// Prereleases and build metadata hyphens
	{Version{1, 2, 3, -1, []PRVersion{prstr("alpha"), prstr("b-eta")}, []string{"123", "b-uild"}}, "1.2.3-alpha.b-eta+123.b-uild"},
	{Version{1, 2, 3, -1, nil, []string{"123", "b-uild"}}, "1.2.3+123.b-uild"},
	{Version{1, 2, 3, -1, []PRVersion{prstr("alpha"), prstr("b-eta")}, nil}, "1.2.3-alpha.b-eta"},
	{Version{1, 2, 3, 4, []PRVersion{prstr("alpha"), prstr("b-eta")}, []string{"123", "b-uild"}}, "1.2.3.4-alpha.b-eta+123.b-uild"},
	{Version{1, 2, 3, 4, nil, []string{"123", "b-uild"}}, "1.2.3.4+123.b-uild"},
	{Version{1, 2, 3, 4, []PRVersion{prstr("alpha"), prstr("b-eta")}, nil}, "1.2.3.4-alpha.b-eta"},
}

var tolerantFormatTests = []formatTest{
	{Version{1, 2, 3, -1, nil, nil}, "v1.2.3"},
	{Version{1, 2, 0, -1, []PRVersion{prstr("alpha")}, nil}, "1.2.0-alpha"},
	{Version{1, 2, 0, -1, nil, nil}, "1.2.00"},
	{Version{1, 2, 3, -1, nil, nil}, "	1.2.3 "},
	{Version{1, 2, 3, -1, nil, nil}, "01.02.03"},
	{Version{0, 0, 3, -1, nil, nil}, "00.0.03"},
	{Version{0, 0, 3, -1, nil, nil}, "000.0.03"},
	{Version{1, 2, 0, -1, nil, nil}, "1.2"},
	{Version{1, 0, 0, -1, nil, nil}, "1"},
	{Version{1, 2, 3, 4, nil, nil}, "v1.2.3.4"},
	{Version{1, 2, 3, 4, nil, nil}, "    v1.2.3.4"},
	{Version{1, 2, 3, 4, nil, nil}, "    1.2.3.4"},
	{Version{1, 2, 0, 0, []PRVersion{prstr("alpha")}, nil}, "1.2.0.0-alpha"},
	{Version{1, 2, 0, 0, nil, nil}, "1.2.00.00"},
	{Version{1, 2, 3, 4, nil, nil}, "	1.2.3.4 "},
	{Version{1, 2, 3, 4, nil, nil}, "01.02.03.04"},
	{Version{0, 0, 3, 1, nil, nil}, "00.0.03.01"},
	{Version{0, 0, 3, 1, nil, nil}, "000.0.03.01"},
}

func TestStringer(t *testing.T) {
	for _, test := range formatTests {
		if res := test.v.String(); res != test.result {
			t.Errorf("Stringer, expected %q but got %q", test.result, res)
		}
	}
}

func TestParse(t *testing.T) {
	for _, test := range formatTests {
		if v, err := Parse(test.result); err != nil {
			t.Errorf("Error parsing %q: %q", test.result, err)
		} else if comp := v.Compare(test.v); comp != 0 {
			t.Errorf("Parsing, expected %q but got %q, comp: %d ", test.v, v, comp)
		} else if err := v.Validate(); err != nil {
			t.Errorf("Error validating parsed version %q: %q", test.v, err)
		}
	}
}

func TestParseTolerant(t *testing.T) {
	for _, test := range tolerantFormatTests {
		if v, err := ParseTolerant(test.result); err != nil {
			t.Errorf("Error parsing %q: %q", test.result, err)
		} else if comp := v.Compare(test.v); comp != 0 {
			t.Errorf("Parsing, expected %q but got %q, comp: %d ", test.v, v, comp)
		} else if err := v.Validate(); err != nil {
			t.Errorf("Error validating parsed version %q: %q", test.v, err)
		}
	}
}

func TestMustParse(t *testing.T) {
	_ = MustParse("32.2.1-alpha")
}

func TestMustParse_panic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Errorf("Should have panicked")
		}
	}()
	_ = MustParse("invalid version")
}

func TestValidate(t *testing.T) {
	for _, test := range formatTests {
		if err := test.v.Validate(); err != nil {
			t.Errorf("Error validating %q: %q", test.v, err)
		}
	}
}

var finalizeVersionMethod = []formatTest{
	{Version{1, 2, 3, -1, nil, nil}, "1.2.3"},
	{Version{0, 0, 1, -1, nil, nil}, "0.0.1"},
	{Version{0, 0, 1, -1, []PRVersion{prstr("alpha"), prstr("preview")}, []string{"123", "456"}}, "0.0.1"},
	{Version{1, 2, 3, -1, []PRVersion{prstr("alpha"), prnum(1)}, []string{"123", "456"}}, "1.2.3"},
	{Version{1, 2, 3, -1, []PRVersion{prstr("alpha"), prnum(1)}, nil}, "1.2.3"},
	{Version{1, 2, 3, -1, nil, []string{"123", "456"}}, "1.2.3"},
	{Version{1, 2, 3, 4, nil, nil}, "1.2.3.4"},
	{Version{0, 0, 1, 0, nil, nil}, "0.0.1.0"},
	{Version{0, 0, 1, 0, []PRVersion{prstr("alpha"), prstr("preview")}, []string{"123", "456"}}, "0.0.1.0"},
	{Version{1, 2, 3, 4, []PRVersion{prstr("alpha"), prnum(1)}, []string{"123", "456"}}, "1.2.3.4"},
	{Version{1, 2, 3, 4, []PRVersion{prstr("alpha"), prnum(1)}, nil}, "1.2.3.4"},
	{Version{1, 2, 3, 4, nil, []string{"123", "456"}}, "1.2.3.4"},
	// Prereleases and build metadata hyphens
	{Version{1, 2, 3, -1, []PRVersion{prstr("alpha"), prstr("b-eta")}, []string{"123", "b-uild"}}, "1.2.3"},
	{Version{1, 2, 3, -1, nil, []string{"123", "b-uild"}}, "1.2.3"},
	{Version{1, 2, 3, -1, []PRVersion{prstr("alpha"), prstr("b-eta")}, nil}, "1.2.3"},
	{Version{1, 2, 3, 4, []PRVersion{prstr("alpha"), prstr("b-eta")}, []string{"123", "b-uild"}}, "1.2.3.4"},
	{Version{1, 2, 3, 4, nil, []string{"123", "b-uild"}}, "1.2.3.4"},
	{Version{1, 2, 3, 4, []PRVersion{prstr("alpha"), prstr("b-eta")}, nil}, "1.2.3.4"},
}

func TestFinalizeVersionMethod(t *testing.T) {
	for _, test := range finalizeVersionMethod {
		out := test.v.FinalizeVersion()
		if out != test.result {
			t.Errorf("Finalized version error, expected %q but got %q", test.result, out)
		}
	}
}

type compareTest struct {
	v1     Version
	v2     Version
	result int
}

var compareTests = []compareTest{
	{Version{1, 0, 0, -1, nil, nil}, Version{1, 0, 0, -1, nil, nil}, 0},
	{Version{2, 0, 0, -1, nil, nil}, Version{1, 0, 0, -1, nil, nil}, 1},
	{Version{0, 1, 0, -1, nil, nil}, Version{0, 1, 0, -1, nil, nil}, 0},
	{Version{0, 2, 0, -1, nil, nil}, Version{0, 1, 0, -1, nil, nil}, 1},
	{Version{0, 0, 1, -1, nil, nil}, Version{0, 0, 1, -1, nil, nil}, 0},
	{Version{0, 0, 2, -1, nil, nil}, Version{0, 0, 1, -1, nil, nil}, 1},
	{Version{1, 2, 3, -1, nil, nil}, Version{1, 2, 3, -1, nil, nil}, 0},
	{Version{2, 2, 4, -1, nil, nil}, Version{1, 2, 4, -1, nil, nil}, 1},
	{Version{1, 3, 3, -1, nil, nil}, Version{1, 2, 3, -1, nil, nil}, 1},
	{Version{1, 2, 4, -1, nil, nil}, Version{1, 2, 3, -1, nil, nil}, 1},

	// Spec Examples #11
	{Version{1, 0, 0, -1, nil, nil}, Version{2, 0, 0, -1, nil, nil}, -1},
	{Version{2, 0, 0, -1, nil, nil}, Version{2, 1, 0, -1, nil, nil}, -1},
	{Version{2, 1, 0, -1, nil, nil}, Version{2, 1, 1, -1, nil, nil}, -1},

	// Spec Examples #9
	{Version{1, 0, 0, -1, nil, nil}, Version{1, 0, 0, -1, []PRVersion{prstr("alpha")}, nil}, 1},
	{Version{1, 0, 0, -1, []PRVersion{prstr("alpha")}, nil}, Version{1, 0, 0, -1, []PRVersion{prstr("alpha"), prnum(1)}, nil}, -1},
	{Version{1, 0, 0, -1, []PRVersion{prstr("alpha"), prnum(1)}, nil}, Version{1, 0, 0, -1, []PRVersion{prstr("alpha"), prstr("beta")}, nil}, -1},
	{Version{1, 0, 0, -1, []PRVersion{prstr("alpha"), prstr("beta")}, nil}, Version{1, 0, 0, -1, []PRVersion{prstr("beta")}, nil}, -1},
	{Version{1, 0, 0, -1, []PRVersion{prstr("beta")}, nil}, Version{1, 0, 0, -1, []PRVersion{prstr("beta"), prnum(2)}, nil}, -1},
	{Version{1, 0, 0, -1, []PRVersion{prstr("beta"), prnum(2)}, nil}, Version{1, 0, 0, -1, []PRVersion{prstr("beta"), prnum(11)}, nil}, -1},
	{Version{1, 0, 0, -1, []PRVersion{prstr("beta"), prnum(11)}, nil}, Version{1, 0, 0, -1, []PRVersion{prstr("rc"), prnum(1)}, nil}, -1},
	{Version{1, 0, 0, -1, []PRVersion{prstr("rc"), prnum(1)}, nil}, Version{1, 0, 0, -1, nil, nil}, -1},

	// Ignore Build metadata
	{Version{1, 0, 0, -1, nil, []string{"1", "2", "3"}}, Version{1, 0, 0, -1, nil, nil}, 0},
}

func TestCompare(t *testing.T) {
	for _, test := range compareTests {
		if res := test.v1.Compare(test.v2); res != test.result {
			t.Errorf("Comparing %q : %q, expected %d but got %d", test.v1, test.v2, test.result, res)
		}
		// Test counterpart
		if res := test.v2.Compare(test.v1); res != -test.result {
			t.Errorf("Comparing %q : %q, expected %d but got %d", test.v2, test.v1, -test.result, res)
		}
	}
}

type wrongformatTest struct {
	v   *Version
	str string
}

var wrongformatTests = []wrongformatTest{
	{nil, ""},
	{nil, "."},
	{nil, "1."},
	{nil, ".1"},
	{nil, "a.b.c"},
	{nil, "1.a.b"},
	{nil, "1.1.a"},
	{nil, "1.a.1"},
	{nil, "a.1.1"},
	{nil, ".."},
	{nil, "1.."},
	{nil, "1.1."},
	{nil, "1..1"},
	{nil, "1.1.+123"},
	{nil, "1.1.-beta"},
	{nil, "-1.1.1"},
	{nil, "1.-1.1"},
	{nil, "1.1.-1"},
	{nil, "1.1.1.a"},
	{nil, "a.1.1.1"},
	{nil, "1.1.1."},
	{nil, "1.1.1.1.+123"},
	{nil, "1.1.1.1.+-beta"},
	// giant numbers
	{nil, "20000000000000000000.1.1"},
	{nil, "1.20000000000000000000.1"},
	{nil, "1.1.20000000000000000000"},
	{nil, "1.1.1-20000000000000000000"},
	{nil, "1.1.1.20000000000000000000"},
	{nil, "1.1.1.20000000000000000000-alpha"},
	// Leading zeroes
	{nil, "01.1.1"},
	{nil, "001.1.1"},
	{nil, "1.01.1"},
	{nil, "1.001.1"},
	{nil, "1.1.01"},
	{nil, "1.1.001"},
	{nil, "1.1.1-01"},
	{nil, "1.1.1-001"},
	{nil, "1.1.1-beta.01"},
	{nil, "1.1.1-beta.001"},
	{nil, "1.1.0.01"},
	{nil, "1.1.0.01-1"},
	{nil, "1.1.0.01-001"},
	{nil, "1.1.1.1-beta.001"},
	{&Version{0, 0, 0, -1, []PRVersion{prstr("!")}, nil}, "0.0.0-!"},
	{&Version{0, 0, 0, -1, nil, []string{"!"}}, "0.0.0+!"},
	// empty prversion
	{&Version{0, 0, 0, -1, []PRVersion{prstr(""), prstr("alpha")}, nil}, "0.0.0-.alpha"},
	// empty build meta data
	{&Version{0, 0, 0, -1, []PRVersion{prstr("alpha")}, []string{""}}, "0.0.0-alpha+"},
	{&Version{0, 0, 0, -1, []PRVersion{prstr("alpha")}, []string{"test", ""}}, "0.0.0-alpha+test."},
}

func TestWrongFormat(t *testing.T) {
	for _, test := range wrongformatTests {

		if res, err := Parse(test.str); err == nil {
			t.Errorf("Parsing wrong format version %q, expected error but got %q", test.str, res)
		}

		if test.v != nil {
			if err := test.v.Validate(); err == nil {
				t.Errorf("Validating wrong format version %q (%q), expected error", test.v, test.str)
			}
		}
	}
}

var wrongTolerantFormatTests = []wrongformatTest{
	{nil, "1.0+abc"},
	{nil, "1.0-rc.1"},
}

func TestWrongTolerantFormat(t *testing.T) {
	for _, test := range wrongTolerantFormatTests {
		if res, err := ParseTolerant(test.str); err == nil {
			t.Errorf("Parsing wrong format version %q, expected error but got %q", test.str, res)
		}
	}
}

func TestCompareHelper(t *testing.T) {
	v := Version{1, 0, 0, -1, []PRVersion{prstr("alpha")}, nil}
	v1 := Version{1, 0, 0, -1, nil, nil}
	if !v.EQ(v) {
		t.Errorf("%q should be equal to %q", v, v)
	}
	if !v.Equals(v) {
		t.Errorf("%q should be equal to %q", v, v)
	}
	if !v1.NE(v) {
		t.Errorf("%q should not be equal to %q", v1, v)
	}
	if !v.GTE(v) {
		t.Errorf("%q should be greater than or equal to %q", v, v)
	}
	if !v.LTE(v) {
		t.Errorf("%q should be less than or equal to %q", v, v)
	}
	if !v.LT(v1) {
		t.Errorf("%q should be less than %q", v, v1)
	}
	if !v.LTE(v1) {
		t.Errorf("%q should be less than or equal %q", v, v1)
	}
	if !v.LE(v1) {
		t.Errorf("%q should be less than or equal %q", v, v1)
	}
	if !v1.GT(v) {
		t.Errorf("%q should be greater than %q", v1, v)
	}
	if !v1.GTE(v) {
		t.Errorf("%q should be greater than or equal %q", v1, v)
	}
	if !v1.GE(v) {
		t.Errorf("%q should be greater than or equal %q", v1, v)
	}
}

const (
	MAJOR = iota
	MINOR
	PATCH
	REVISION
)

type incrementTest struct {
	version         Version
	incrementType   int
	expectingError  bool
	expectedVersion Version
}

var incrementTests = []incrementTest{
	{Version{1, 2, 3, -1, nil, nil}, PATCH, false, Version{1, 2, 4, -1, nil, nil}},
	{Version{1, 2, 3, -1, nil, nil}, MINOR, false, Version{1, 3, 0, -1, nil, nil}},
	{Version{1, 2, 3, -1, nil, nil}, MAJOR, false, Version{2, 0, 0, -1, nil, nil}},
	{Version{0, 1, 2, -1, nil, nil}, PATCH, false, Version{0, 1, 3, -1, nil, nil}},
	{Version{0, 1, 2, -1, nil, nil}, MINOR, false, Version{0, 2, 0, -1, nil, nil}},
	{Version{0, 1, 2, -1, nil, nil}, MAJOR, false, Version{1, 0, 0, -1, nil, nil}},
	{Version{1, 2, 3, 1, nil, nil}, REVISION, false, Version{1, 2, 3, 2, nil, nil}},
	{Version{0, 1, 2, 3, nil, nil}, REVISION, false, Version{0, 1, 2, 4, nil, nil}},
	{Version{0, 1, 2, -1, nil, nil}, REVISION, false, Version{0, 1, 2, -1, nil, nil}},
	{Version{1, 2, 3, 1, nil, nil}, PATCH, false, Version{1, 2, 4, 0, nil, nil}},
	{Version{1, 2, 3, 1, nil, nil}, MINOR, false, Version{1, 3, 0, 0, nil, nil}},
	{Version{1, 2, 3, 1, nil, nil}, MAJOR, false, Version{2, 0, 0, 0, nil, nil}},
}

func TestIncrements(t *testing.T) {
	for _, test := range incrementTests {
		var originalVersion = Version{
			test.version.Major,
			test.version.Minor,
			test.version.Patch,
			test.version.Revision,
			test.version.Pre,
			test.version.Build,
		}
		var err error
		switch test.incrementType {
		case PATCH:
			err = test.version.IncrementPatch()
		case MINOR:
			err = test.version.IncrementMinor()
		case MAJOR:
			err = test.version.IncrementMajor()
		case REVISION:
			err = test.version.IncrementRevision()
		}
		if test.expectingError {
			if err != nil {
				t.Errorf("Increment version, expecting %q, got error %q", test.expectedVersion, err)
			}
			if test.version.EQ(originalVersion) {
				t.Errorf("Increment version, expecting %q, got %q", test.expectedVersion, test.version)
			}
		} else {
			if (err != nil) && !test.expectingError {
				t.Errorf("Increment version %q, not expecting error, got %q", test.version, err)
			}
			if test.version.NE(test.expectedVersion) {
				t.Errorf("Increment version, expecting %q, got %q", test.expectedVersion, test.version)
			}
		}
	}
}

func TestPreReleaseVersions(t *testing.T) {
	p1, err := NewPRVersion("123")
	if !p1.IsNumeric() {
		t.Errorf("Expected numeric prversion, got %q", p1)
	}
	if p1.VersionNum != 123 {
		t.Error("Wrong prversion number")
	}
	if err != nil {
		t.Errorf("Not expected error %q", err)
	}
	p2, err := NewPRVersion("alpha")
	if p2.IsNumeric() {
		t.Errorf("Expected non-numeric prversion, got %q", p2)
	}
	if p2.VersionStr != "alpha" {
		t.Error("Wrong prversion string")
	}
	if err != nil {
		t.Errorf("Not expected error %q", err)
	}
}

func TestBuildMetaDataVersions(t *testing.T) {
	_, err := NewBuildVersion("123")
	if err != nil {
		t.Errorf("Unexpected error %q", err)
	}

	_, err = NewBuildVersion("build")
	if err != nil {
		t.Errorf("Unexpected error %q", err)
	}

	_, err = NewBuildVersion("test?")
	if err == nil {
		t.Error("Expected error, got none")
	}

	_, err = NewBuildVersion("")
	if err == nil {
		t.Error("Expected error, got none")
	}
}

func TestNewHelper(t *testing.T) {
	v, err := New("1.2.3")
	if err != nil {
		t.Fatalf("Unexpected error %q", err)
	}

	// New returns pointer
	if v == nil {
		t.Fatal("Version is nil")
	}
	if v.Compare(Version{1, 2, 3, -1, nil, nil}) != 0 {
		t.Fatal("Unexpected comparison problem")
	}
}

func TestMakeHelper(t *testing.T) {
	v, err := Make("1.2.3")
	if err != nil {
		t.Fatalf("Unexpected error %q", err)
	}
	if v.Compare(Version{1, 2, 3, -1, nil, nil}) != 0 {
		t.Fatal("Unexpected comparison problem")
	}
}

type finalizeTest struct {
	input  string
	output string
}

var finalizeTests = []finalizeTest{
	{"", ""},
	{"1.2.3", "1.2.3"},
	{"0.0.1", "0.0.1"},
	{"0.0.1-alpha.preview+123.456", "0.0.1"},
	{"1.2.3-alpha.1+123.456", "1.2.3"},
	{"1.2.3-alpha.1", "1.2.3"},
	{"1.2.3+123.456", "1.2.3"},
	{"1.2.3-alpha.b-eta+123.b-uild", "1.2.3"},
	{"1.2.3+123.b-uild", "1.2.3"},
	{"1.2.3-alpha.b-eta", "1.2.3"},
	{"1.2-alpha", ""},
	{"1.2.3.4", "1.2.3.4"},
	{"0.0.1.0", "0.0.1.0"},
	{"0.0.1.0-alpha.preview+123.456", "0.0.1.0"},
	{"1.2.3.4-alpha.1+123.456", "1.2.3.4"},
	{"1.2.3.4-alpha.1", "1.2.3.4"},
	{"1.2.3.4+123.456", "1.2.3.4"},
	{"1.2.3.4-alpha.b-eta+123.b-uild", "1.2.3.4"},
	{"1.2.3.4+123.b-uild", "1.2.3.4"},
	{"1.2.3.4-alpha.b-eta", "1.2.3.4"},
}

func TestFinalizeVersion(t *testing.T) {
	for _, test := range finalizeTests {
		finalVer, err := FinalizeVersion(test.input)
		if finalVer == "" {
			if err == nil {
				t.Errorf("Finalize Version error, expected error but got nil")
			}
		} else if finalVer != test.output && err != nil {
			t.Errorf("Finalize Version error expected %q but got %q", test.output, finalVer)
		}
	}
}

func BenchmarkParseSimple(b *testing.B) {
	const VERSION = "0.0.1"
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = Parse(VERSION)
	}
}

func BenchmarkParseComplex(b *testing.B) {
	const VERSION = "0.0.1-alpha.preview+123.456"
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = Parse(VERSION)
	}
}

func BenchmarkParseAverage(b *testing.B) {
	l := len(formatTests)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = Parse(formatTests[n%l].result)
	}
}

func BenchmarkParseTolerantAverage(b *testing.B) {
	l := len(tolerantFormatTests)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = ParseTolerant(tolerantFormatTests[n%l].result)
	}
}

func BenchmarkStringSimple(b *testing.B) {
	const VERSION = "0.0.1"
	v, _ := Parse(VERSION)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.String()
	}
}

func BenchmarkStringLarger(b *testing.B) {
	const VERSION = "11.15.2012"
	v, _ := Parse(VERSION)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.String()
	}
}

func BenchmarkStringComplex(b *testing.B) {
	const VERSION = "0.0.1-alpha.preview+123.456"
	v, _ := Parse(VERSION)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.String()
	}
}

func BenchmarkStringAverage(b *testing.B) {
	l := len(formatTests)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = formatTests[n%l].v.String()
	}
}

func BenchmarkValidateSimple(b *testing.B) {
	const VERSION = "0.0.1"
	v, _ := Parse(VERSION)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.Validate()
	}
}

func BenchmarkValidateComplex(b *testing.B) {
	const VERSION = "0.0.1-alpha.preview+123.456"
	v, _ := Parse(VERSION)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.Validate()
	}
}

func BenchmarkValidateAverage(b *testing.B) {
	l := len(formatTests)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = formatTests[n%l].v.Validate()
	}
}

func BenchmarkCompareSimple(b *testing.B) {
	const VERSION = "0.0.1"
	v, _ := Parse(VERSION)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.Compare(v)
	}
}

func BenchmarkCompareComplex(b *testing.B) {
	const VERSION = "0.0.1-alpha.preview+123.456"
	v, _ := Parse(VERSION)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.Compare(v)
	}
}

func BenchmarkCompareAverage(b *testing.B) {
	l := len(compareTests)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		compareTests[n%l].v1.Compare((compareTests[n%l].v2))
	}
}
