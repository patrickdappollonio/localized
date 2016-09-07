// +build windows

package localized

import (
	"regexp"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

const kernel32DLL = "kernel32.dll"

// This validates according to the docs at https://msdn.microsoft.com/en-us/library/windows/desktop/dd373814(v=vs.85).aspx
// This won't accept "replacement locales" or user-defined locales. It won't even allow language without a region.
var reValidateLang = regexp.MustCompile(`[a-z]{2}\-[A-Z0-9]{2,3}`)

func (d *Detector) Detect() error {
	// Load Kernel32.dll, check if it fails
	k32dll, err := windows.LoadDLL(kernel32DLL)

	// Check if it fails, and if so, return
	if err != nil {
		return setDefaultNotFound(d)
	}

	// Once we're done with it, just release the DLL
	defer k32dll.Release()

	// Holder for the language
	var detectedLanguage string

	// Check first if we can get the user language
	if userlang, found := callLang(k32dll, "GetUserDefaultLocaleName"); found {
		d.Detected = true
		detectedLanguage = userlang
	}

	// If it wasn't found, try getting the system language
	if syslang, found := callLang(k32dll, "GetSystemDefaultLocaleName"); detectedLanguage == "" && found {
		d.Detected = true
		detectedLanguage = syslang
	}

	// Check if the string we got matches at least the format
	// we're expecting.
	if !reValidateLang.MatchString(detectedLanguage) {
		return setDefaultNotFound(d)
	}

	// If the language was actually found, parse it. This will be
	// assumed to work since it uses parts of the original regular expression
	// used to handle these stuff.
	d.Lang, d.Region = parselang(detectedLanguage)
	d.Detected = true

	return nil
}

func setDefaultNotFound(d *Detector) error {
	// If we had a default language set, then set detected to false
	// and don't return an error, since default language doesn't give any error
	if d.defaultSet {
		d.Detected = false
		return nil
	}

	return ErrNoLangDetected
}

func callLang(k *windows.DLL, processName string) (string, bool) {
	// Check if we did passed kernel32
	if k == nil {
		return "", false
	}

	// Let's create a buffer that holds the response from the system call.
	// According to Windows, the buffer max length is 85, taken from LOCALE_NAME_MAX_LENGTH
	// so we will use that number.
	buf := make([]uint16, 85)

	// Now let's find the process that will return the output we need
	proc, err := k.FindProc(processName)

	// Check if we did find that process
	if err != nil {
		return "", false
	}

	// Get the first position of the buffer
	firstpos := uintptr(unsafe.Pointer(&buf[0]))

	// Try calling now the process with the information we have.
	// We need to give the pointer to the location of the first value of
	// our buffer, so Windows can store the answer there. We also need to pass
	// the length of the buffer, which is always 85 according to the docs.
	// This will return "r" which is a status number, where 0 means it failed.
	if r, _, _ := proc.Call(firstpos, uintptr(85)); r == 0 {
		return "", false
	}

	// If we reach this part, this means we successfully got the information as
	// a UTF16 byte array we need to convert to string and return.
	return strings.TrimSpace(windows.UTF16ToString(buf)), true
}

func parselang(lang string) (string, string) {
	// We find the language part, as well as the region
	return reLowerLang.FindString(lang), reUpperRegion.FindString(lang)
}
