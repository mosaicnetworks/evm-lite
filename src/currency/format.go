// Package currency provides library functions for manipulating token balances
//
// Units
//
// Taking inspiration from the SI units, we have suitable multiples:
//
// 1/ 1 000 000 000 000 000 000			atto		(a)	10^-18
// 1/ 1 000 000 000 000 000				femto 		(f)	10^-15
// 1/ 1 000 000 000 000					pico		(p)	10^-12
// 1/ 1 000 000 000						nano		(n)	10^-9
// 1/ 1 000 000							micro		(u)	10^-6
// 1/ 1 000								milli		(m)	10^-3
// 1									Token		(T)	1
// All letters are lowercase except for T for Token
//
package currency

import (
	"fmt"
	"strconv"
	"strings"
)

const tokenLetters = "afpnumT"

var thouSeparator = ","
var decSeparator = "."

//ExpandCurrencyString takes a string with a token suffix and expands it with
//the appropriate number of zeroes. The input string may contain a decimal
//point, and it will expand it, standard form style. By conventtion hex
//balances have a leading 0x, decimals are bare. THis function will work
//correctly with either. This function returns Attoms.
func ExpandCurrencyString(input string) string {

	// trim whitespace as it would mess with the place counting further on
	cleanInput := strings.TrimSpace(input)

	if cleanInput == "" {
		return ""
	}

	//TODO check for and strip thouSeparator

	// token is the last character in the string
	token := cleanInput[len(cleanInput)-1:]

	// search for token in tokenLetters. If not found, Index() returns -1,
	// and the expression evaluates to zero. Otherwise the list is ordered
	// in ascending powers in multiples of three.
	tokenPower := (strings.Index(tokenLetters, token) + 1) * 3

	// If token not found, there is nothing to do, bar the TrimSpace() we have
	// already done.
	if tokenPower == 0 {
		return cleanInput
	}

	tokenPower -= 3

	// Remove the token from the input string
	last := len(cleanInput) - 1
	cleanInput = cleanInput[:last] // trim token from string.

	// Check for a decimal point
	idx := strings.Index(cleanInput, decSeparator)
	if idx >= 0 {
		// Reduce out zero count (tokenPower), by the number of characters
		// after the decimal point. Remove dot from the string
		pre := cleanInput[:idx]
		fix := cleanInput[idx+1:]
		tokenPower -= len(fix)

		cleanInput = pre + fix
	}

	// Add the requisite number of zeroes to the end of the number
	format := "%0" + strconv.Itoa(tokenPower) + "d"
	cleanInput = cleanInput + fmt.Sprintf(format, 0)

	// Loop cleans leading zeroes from the string. Coded to preserve leading 0[xX]
	// despite it not being an expected value
	for len(cleanInput) > 1 && cleanInput[0:1] == "0" && cleanInput[1:2] != "x" && cleanInput[1:2] != "X" {
		cleanInput = cleanInput[1:]
	}

	return cleanInput
}

// ExpandAndSeparateCurrencyString expands the input string, then applies
// comma separators
func ExpandAndSeparateCurrencyString(input string) string {
	expanded := ExpandCurrencyString(input)
	l := len(expanded)

	for l > 3 {
		l -= 3
		expanded = expanded[:l] + thouSeparator + expanded[l:]
	}

	return expanded
}

// FormatTenomString is a wrapper to FormatUnitString
func FormatTenomString(input string) string {
	return FormatUnitString(input, 18)
}

// FormatUnitString takes an atomic input and returns
func FormatUnitString(input string, power int) string {

	// clean input. Guarantees no whitespace or tokens.
	cleanInput := ExpandCurrencyString(input)

	if cleanInput == "0" {
		return cleanInput
	}

	if len(cleanInput) == 18 { // We need some zero padding
		return strings.TrimRight("0."+cleanInput, "0")
	}

	if len(cleanInput) < 18 { // We need some zero padding
		format := "%0" + strconv.Itoa(power-len(cleanInput)) + "d"
		cleanInput = "0." + fmt.Sprintf(format, 0) + cleanInput
		return strings.TrimRight(cleanInput, "0")
	}

	cleanSuffix := strings.TrimRight(cleanInput[len(cleanInput)-power:], "0")
	if len(cleanSuffix) != 0 {
		cleanSuffix = "." + cleanSuffix
	}
	return cleanInput[:len(cleanInput)-power] + cleanSuffix
}

// FormatCurrencyString ...
func FormatCurrencyString(input string) string {

	// clean input. Guarantees no whitespace or tokens.
	cleanInput := ExpandCurrencyString(input)

	//If we have less than a thousand there is nothing to do
	if len(cleanInput) < 4 {
		return cleanInput
	}

	strpos := (len(cleanInput) / 3)
	if strpos >= len(tokenLetters) {
		strpos = len(tokenLetters) - 1
	}

	tokenLetter := string([]byte{tokenLetters[strpos]})

	zeroplaces := strpos * 3

	for cleanInput[len(cleanInput)-1:] == "0" && zeroplaces > 0 {
		last := len(cleanInput) - 1
		cleanInput = cleanInput[:last]
		zeroplaces--
	}

	if zeroplaces < 1 {
		return cleanInput + tokenLetter
	}

	idx := len(cleanInput) - zeroplaces

	return cleanInput[:idx] + decSeparator + cleanInput[idx:] + tokenLetter

}
