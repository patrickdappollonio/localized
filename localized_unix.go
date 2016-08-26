// +build darwin freebsd linux netbsd openbsd
package localized

const (
	LC_ALL = "LC_ALL"
	LANG   = "LANG"
)

func (d *defaults) GetLanguage() (*Language, error) {
	// Find the lang code using environment variables
	// and also check if the code matches the regular
	// expression
	lang := findCode()

	if lang == "" {
		if d.Set {
			return &Language{
				Code:     d.Code,
				Encoding: d.Encoding,
			}, nil
		}

		return nil, ErrNoLangDetected
	}
}

func findCode() string {
	langcode := envdef(LC_ALL, envdef(LANG, ""))

	if !reHasFormat.MatchString(langcode) {
		return ""
	}

	return langcode
}
