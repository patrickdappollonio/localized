// +build linux darwin freebsd netbsd openbsd
package localized

import (
	"os"
	"testing"
)

func TestParsingRegex(t *testing.T) {
	var langs = []struct {
		Code, ELang, ERegion string
		DLang, DRegion       string
		Detected             bool
		EError               error
		Switch               int
	}{
		{
			Switch:   1,
			Code:     "en_US.UTF_8",
			ELang:    "en",
			ERegion:  "US",
			Detected: true,
		}, {
			Switch:   2,
			Code:     "es_CL.ISO_8859_1",
			ELang:    "es",
			ERegion:  "CL",
			Detected: true,
		}, {
			Switch:   1,
			Code:     "co_CR",
			ELang:    "co",
			ERegion:  "CR",
			Detected: true,
		}, {
			Switch:   2,
			Code:     "es_419.iso_8859_15",
			ELang:    "es",
			ERegion:  "419",
			Detected: true,
		}, {
			Switch:   1,
			Code:     "this_is_not_valid1",
			Detected: false,
			EError:   ErrNoLangDetected,
		}, {
			Switch:   1,
			Code:     "this_is_not_valid2",
			ELang:    "en",
			ERegion:  "US",
			DLang:    "en",
			DRegion:  "US",
			Detected: false,
			EError:   nil,
		},
	}

	for _, tv := range langs {
		validateLang(t, tv.Switch, tv.Code, tv.DLang, tv.DRegion, tv.ELang, tv.ERegion, tv.Detected, tv.EError)
	}
}

func validateLang(t *testing.T, ver int, code, defLang, defRegion, expectedLang, expectedRegion string, detected bool, expectedErr error) {
	t.Logf("Setting computer language to: %q", code)

	// Clean env vars
	os.Setenv(LC_ALL, "")
	os.Setenv(LANG, "")

	// Set the language to the code passed
	if ver == 1 {
		if err := os.Setenv(LC_ALL, code); err != nil {
			t.Fatalf("Error while setting computer language at %q: %s", LC_ALL, err.Error())
		}
	} else {
		if err := os.Setenv(LANG, code); err != nil {
			t.Fatalf("Error while setting computer language at %q: %s", LANG, err.Error())
		}
	}

	// Create a new default detector
	d := New()

	// Check if we passed default data
	if defLang != "" && defRegion != "" {
		d = New(Config{DefaultLanguage: defLang, DefaultRegion: defRegion})
	}

	// Perform the detection and see if it fails
	if err := d.Detect(); err != nil {
		if expectedErr != nil {
			if err != expectedErr {
				t.Fatalf("For lang %q -- Expected error %q, got %q instead", code, expectedErr.Error(), err.Error())
			}
		} else {
			t.Fatalf("For lang %q -- Error received but not expected: %s", code, err.Error())
		}
	}

	// Now check if we did detect a language
	if d.Detected != detected {
		t.Fatalf("For lang %q -- Language was detected? %v Expected detection? %v", code, d.Detected, detected)
	}

	// Now check if the language detected was indeed the one we needed
	if d.Lang != expectedLang {
		t.Fatalf("For lang %q -- Expected language: %q, but we got %q:%q instead", code, expectedLang, d.Lang, d.Region)
	}

	// Now check if the region code was indeed the one we needed
	if d.Region != expectedRegion {
		t.Fatalf("For lang %q -- Expected region: %q, but we got %q:%q instead", code, expectedRegion, d.Lang, d.Region)
	}

}
