package ui

import (
	"fmt"
	"os"
)

func colorsEnabled() bool {
	return os.Getenv("NO_COLOR") == "" && os.Getenv("TERM") != "dumb"
}

func applyColor(code string, format string, a ...interface{}) string {
	if !colorsEnabled() {
		return fmt.Sprintf(format, a...)
	}
	return fmt.Sprintf("\033["+code+"m"+format+"\033[0m", a...)
}

// RedString wraps the input string with ANSI codes for red color
func RedString(format string, a ...interface{}) string {
	return applyColor("31", format, a...)
}

// GreenString wraps the input string with ANSI codes for green color
func GreenString(format string, a ...interface{}) string {
	return applyColor("32", format, a...)
}

// YellowString wraps the input string with ANSI codes for yellow color
func YellowString(format string, a ...interface{}) string {
	return applyColor("33", format, a...)
}

// BlueString wraps the input string with ANSI codes for blue color
func BlueString(format string, a ...interface{}) string {
	return applyColor("34", format, a...)
}

// CyanString wraps the input string with ANSI codes for cyan color
func CyanString(format string, a ...interface{}) string {
	return applyColor("36", format, a...)
}

// GrayString wraps the input string with ANSI codes for gray color
func GrayString(format string, a ...interface{}) string {
	return applyColor("90", format, a...)
}
