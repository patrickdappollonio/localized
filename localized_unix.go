// +build linux darwin freebsd netbsd openbsd
package localized

import "regexp"

const (
	LC_ALL = "LC_ALL"
	LANG   = "LANG"
)

var (
	// Any letter lowercase two times, then an underscore, then
	// any letter or number 2 or 3 times, then an optional dot
	// then zero or more characters.
	reUnixLang = regexp.MustCompile(`([a-z]{2})\_([A-Z0-9]{2,3})\.?.*`)

	// Remove all the suffix containing the encoding
	reEncoding = regexp.MustCompile(`\..*$`)

	// Find a two-letter language code
	reLowerLang = regexp.MustCompile(`[a-z]{2}`)

	// Find the two or three letter-number region
	reUpperRegion = regexp.MustCompile(`[A-Z0-9]{2,3}`)
)

func (d *Detector) Detect() error {
	// Find the lang code using environment variables
	// and also check if the code matches the regular
	// expression
	lang := evars()

	// If we didn't find any language then return the defaults
	// if they were set, but also set Detected as false
	if lang == "" {
		// If there was a default value, then return it
		// and don't fail, otherwise, return NoLangDetected
		if d.defaultSet {
			d.Detected = false
			return nil
		} else {
			return ErrNoLangDetected
		}
	}

	// If the language was actually found, parse it. This will be
	// assumed to work since it uses parts of the original regular expression
	// used to handle these stuff.
	d.Lang, d.Region = parselang(lang)
	d.Detected = true

	return nil
}

func evars() string {
	// Try finding first in $LC_ALL, then in $LANG
	langcode := envdef(LC_ALL, envdef(LANG, ""))

	// Check if what we found doesn't look like a language code
	// from a UNIX env
	if !reUnixLang.MatchString(langcode) {
		return ""
	}

	// Else return the language code
	return langcode
}

func parselang(lang string) (string, string) {
	// First, we remove the encoding from the string
	lang = reEncoding.ReplaceAllString(lang, "")

	// Then, we find the language part, as well as the region
	return reLowerLang.FindString(lang), reUpperRegion.FindString(lang)
}
