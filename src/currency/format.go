// Package currency provides library functions for manipulating token balances
//
// Units
//
// Taking inspiration from the SI units, we have suitable multiples:
//
//  1 000 000 000 000 000 000 000 000	yotta	(Y)	10^24
//  1 000 000 000 000 000 000 000		zetta	(Z)	10^21
//  1 000 000 000 000 000 000			exa		(E)	10^18
//  1 000 000 000 000 000				peta	(P)	10^15
//  1 000 000 000 000					tera	(T)	10^12
//  1 000 000 000						giga	(G)	10^9
//  1 000 000							mega	(M)	10^6
//  1 000								kilo	(K)	10^3
//
// NB we use a capital K for kilo, so all letters are capital.
//
// Capital E as the last digit is treated as an exponential, lower case e is
// a hex number.
//
package currency

import (
	"fmt"
	"strconv"
	"strings"
)

const tokenLetters = "KMGTPEZY"

var thouSeparator = ","

//ExpandCurrencyString takes a string with a token suffix and expands it with
//the appropriate number of zeroes. The input string may contain a decimal
//point, and it will expand it, standard form style. By conventtion hex
//balances have a leading 0x, decimals are bare. THis function will work
//correctly with either.
func ExpandCurrencyString(input string) string {

	// trim whitespace as it would mess with the place counting further on
	cleanInput := strings.TrimSpace(input)

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

	// Remove the token from the input string
	last := len(cleanInput) - 1
	cleanInput = cleanInput[:last] // trim token from string.

	// Check for a decimal point
	idx := strings.Index(cleanInput, ".")
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
