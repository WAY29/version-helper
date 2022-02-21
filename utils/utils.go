package utils

import (
	"fmt"
	"os"
	"strings"
)

const (
	FLAG_CANCEL = iota
	FLAG_CHECK
	FLAG_CELEBRATION
	FLAG_GEAR
	FLAG_WARNING
	EMOJI_CANCEL      = "❌"
	EMOJI_CHECK       = "\033[1;32m✔️\033[0m"
	EMOJI_CELEBRATION = "🎉"
	EMOJI_GEAR        = "\033[1;36m⚙️\033[0m"
	EMOJI_WARNING     = "⚠️"
	COLOR_CANCEL      = "\033[1;31m"
	COLOR_CHECK       = "\033[1;32m"
	COLOR_CELEBRATION = "\033[1;33m"
	COLOR_GEAR        = "\033[1;36m"
	COLOR_WARNING     = "\033[1;33m"
)

var (
	SUPPORT_UNICODE = false
	EMOJI_DICT      = map[int][]string{
		FLAG_CANCEL:      {"[-]", EMOJI_CANCEL, COLOR_CANCEL},
		FLAG_CHECK:       {"[+]", EMOJI_CHECK, COLOR_CHECK},
		FLAG_CELEBRATION: {"[#]", EMOJI_CELEBRATION, COLOR_CELEBRATION},
		FLAG_GEAR:        {"[.]", EMOJI_GEAR, COLOR_GEAR},
		FLAG_WARNING:     {"[!]", EMOJI_WARNING, COLOR_WARNING},
	}
)

func init() {
	lang := strings.ToUpper(os.Getenv("LANG"))
	if strings.Contains(lang, "UTF") {
		SUPPORT_UNICODE = true
	}
}

func WrapRed(msg string) string {
	return "\033[1;31m" + msg + "\033[0m"
}

func WrapGreen(msg string) string {
	return "\033[1;32m" + msg + "\033[0m"
}

func WrapYellow(msg string) string {
	return "\033[1;33m" + msg + "\033[0m"
}

func WrapCyan(msg string) string {
	return "\033[1;36m" + msg + "\033[0m"
}

func printEmoji(msg string, level, flag int) {
	sign := EMOJI_DICT[flag][0]
	if SUPPORT_UNICODE {
		sign = EMOJI_DICT[flag][1]
		msg = EMOJI_DICT[flag][2] + msg + "\033[0m"
	}
	fmt.Printf("%s%s %s\n", strings.Repeat("  ", level), sign, msg)
}

func Error(msg string, level int) {
	printEmoji(msg, level, FLAG_CANCEL)
	os.Exit(1)
}

func Check(msg string, level int) {
	printEmoji(msg, level, FLAG_CHECK)
}

func Celebration(msg string, level int) {
	printEmoji(msg, level, FLAG_CELEBRATION)
}

func Work(msg string, level int) {
	printEmoji(msg, level, FLAG_GEAR)
}

func Warning(msg string, level int) {
	printEmoji(msg, level, FLAG_WARNING)
}
