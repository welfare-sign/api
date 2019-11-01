package wsgin

import "os"

// EnvSharpMode indicates environment name for sharp mode.
const EnvSharpMode = "SHARP_MODE"

const (
	// DevMode indicates sharp mode is dev.
	DevMode = "dev"
	// TestMode indicates sharp mode is test.
	TestMode = "test"
	// PrdMode indicates sharp mode is prd.
	PrdMode = "prd"
)
const (
	devCode = iota
	testCode
	prdCode
)

var (
	sharpMode = devCode
	modeName  = DevMode
)

func init() {
	mode := os.Getenv(EnvSharpMode)
	SetMode(mode)
}

// SetMode sets gin mode according to input string.
func SetMode(value string) {
	switch value {
	case DevMode, "":
		sharpMode = devCode
	case TestMode:
		sharpMode = testCode
	case PrdMode:
		sharpMode = prdCode
	default:
		panic("sharp mode unknown: " + value)
	}
	if value == "" {
		value = DevMode
	}
	modeName = value
}

// Mode returns currently gin mode.
func Mode() string {
	return modeName
}
