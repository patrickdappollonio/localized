package localized

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Language struct {
	Code     string
	Encoding string
}

type defaults struct {
	Language
	Set bool
}

var (
	Lang defaults

	ErrNoLangDetected = fmt.Errorf("localized: no language detected")

	reHasFormat = regexp.MustCompile(`[A-Za-z]{2}\_[A-Za-z0-9]{2,3}\.[\w\-\_]+`)
)

func (l *defaults) SetDefault(def Language) *defaults {
	l.Code = def.Code
	l.Encoding = def.Encoding
	l.Set = true
	return l
}

func envdef(key, defvalue string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}

	return defvalue
}
