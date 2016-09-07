package localized

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Detector is the instance of the detection mechanism
type Detector struct {
	Lang     string // Lang holds a two-letter key of the language (like `en`, `es` or similar)
	Region   string // Region holds a key for the region code (like `US` for USA, `CL` for Chile, or `419` for Latin-america)
	Detected bool   // Detected will be true if it was originally detected from the environment, false if we're returning the defaults

	defaultSet bool
}

// Config holds a default configuration for the detector
// like a default language and a default region, if nothing is found
type Config struct {
	DefaultLanguage string
	DefaultRegion   string
}

var (
	ErrNoLangDetected = fmt.Errorf("localized: no language detected") // Returned when no language was detected and there are no defaults

	reLangWithRegion = regexp.MustCompile(`[A-Za-z]{2}\_[A-Za-z0-9]{2,3}`)
	reLowerLang      = regexp.MustCompile(`[a-z]{2}`)      // Find a two-letter language code
	reUpperRegion    = regexp.MustCompile(`[A-Z0-9]{2,3}`) // Find the two or three letter-number region
)

// New takes an optional configuration with a default language and a default region
// and creates an instance of the detector. Subsequent calls to Detect() will return
// ErrNoLangDetected if no config was passed and it was impossible to detect the language.
// When multiple configs are passed, only the first one gets used, all others are silently
// discarded.
func New(c ...Config) *Detector {
	var defaults Config

	if len(c) >= 1 {
		defaults = c[0]
	}

	var d Detector

	if defaults.DefaultLanguage != "" {
		d.Lang = defaults.DefaultLanguage
		d.defaultSet = true
	}

	if defaults.DefaultRegion != "" {
		d.Region = defaults.DefaultRegion
		d.defaultSet = true
	}

	return &d
}

func envdef(key, defvalue string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}

	return defvalue
}
