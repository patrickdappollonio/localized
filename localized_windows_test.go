// +build windows

package localized

import "testing"

func TestLanguageDetection(t *testing.T) {
	// Create an instance of the detector
	d := New()

	// We expect that every computer out there has at least an user-defined language
	if err := d.Detect(); err != nil {
		t.Fatalf("Unable to detect user language: %s", err.Error())
	}

	// Check if the detection was successful
	if !d.Detected {
		t.Fatalf("Language was not detected")
	}
}
