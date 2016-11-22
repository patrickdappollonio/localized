# localized

[![GoDoc](https://godoc.org/github.com/patrickdappollonio/localized?status.svg)](https://godoc.org/github.com/patrickdappollonio/localized)

`localized` is a simple package to detect the user language based on the environment. On Windows, it works by performing a Win API call. On *NIX, it works by simply retrieving `LC_ALL` environment variable first, then `LANG` if the first one didn't contain any information.

### Usage

In order to use it, just create an instance of `Detector` with a given configuration (or no configuration at all), and then call `Detect()`. 

Depending on the configuration passed, you will receive the detected language and region. If no configuration is passed and the language couldn't be detected, then `Detect()` will return `ErrNoLangDetected`. If a configuration is passed, then two outcomes may happen: a) if the language was detected, then the configuration is ommited, and the `Detected` flag in the `Detector` will be set as `true`, effectively implying that the language was detected, or; b) if the language wasn't detected, then the configuration is returned (which contains a default language and region), and the value of `Detected` is set to `false`.

### Examples

**Default configuration**
```go
// we create a new instance of the detector
lang := localized.New() 

// check if there was an error because the language couldn't be detected
if err := lang.Detect(); err != nil && err == localized.ErrNoLangDetected {
  // the language wasn't detected. `localized` doesn't return any other error
}

// print output if appropiate
fmt.Printf("Language: %q, Region: %q", lang.Lang, lang.Region)
```

**With default settings**
```go
// we create a new instance of the detector with default settings
lang := localized.New(localized.Config{ DefaultLanguage: "en", DefaultRegion: "US" })

// we call the detection process, which even when it returns an error, since
// we passed a default configuration, it will never return one of those
lang.Detect()

// since it doesn't return an error, the only way to know if the language
// was indeed detected is by checking `lang.Detected`. If true, then 
// the language was detected, if false, then it means we couldn't find a 
// language so we just returned the default values from the configuration
fmt.Printf("Detected?: %s, Language: %q, Region: %q", lang.Detected, lang.Lang, lang.Region)
```
