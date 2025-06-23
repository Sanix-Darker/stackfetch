package ui

import "fmt"

// RedString wraps the input string with ANSI codes for red color
func RedString(format string, a ...interface{}) string {
	return fmt.Sprintf("\033[31m"+format+"\033[0m", a...)
}

// GreenString wraps the input string with ANSI codes for green color
func GreenString(format string, a ...interface{}) string {
	return fmt.Sprintf("\033[32m"+format+"\033[0m", a...)
}

// YellowString wraps the input string with ANSI codes for yellow color
func YellowString(format string, a ...interface{}) string {
	return fmt.Sprintf("\033[33m"+format+"\033[0m", a...)
}
