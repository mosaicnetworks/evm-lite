package currency_test

import (
	"testing"

	"github.com/mosaicnetworks/evm-lite/src/currency"
)

type testRecord struct {
	input  string
	output string
}

// 1/ 1 000 000 000 000 000 000			atto		(a)	10^-18
// 1/ 1 000 000 000 000 000				femto 		(f)	10^-15
// 1/ 1 000 000 000 000					pico		(p)	10^-12
// 1/ 1 000 000 000						nano		(n)	10^-9
// 1/ 1 000 000							micro		(u)	10^-6
// 1/ 1 000								milli		(m)	10^-3
// 1									Token		(T)	1

func TestExpandCurrencyString(t *testing.T) {

	var tests = []testRecord{
		testRecord{input: "1f", output: "1000"},
		testRecord{input: "1.2p", output: "1200000"},
		testRecord{input: "1.23n", output: "1230000000"},
		testRecord{input: "1.2u", output: "1200000000000"},
		testRecord{input: "1.2m", output: "1200000000000000"},
		testRecord{input: "1.2T", output: "1200000000000000000"},
		testRecord{input: "1200T", output: "1200000000000000000000"},
		testRecord{input: "1200000T", output: "1200000000000000000000000"},
		testRecord{input: "0.2f", output: "200"},
		testRecord{input: "0x122f", output: "0x122000"},
	}

	for _, test := range tests {
		ret := currency.ExpandCurrencyString(test.input)
		if ret != test.output {
			t.Errorf("\nWrong Answer: %s\nGot: %s\nExpected: %s\n", test.input, ret, test.output)
		} else {
			t.Logf("%s => %s", test.input, test.output)
		}
	}
}

func TestExpandAndSeparateCurrencyString(t *testing.T) {

	var tests = []testRecord{
		testRecord{input: "1f", output: "1,000"},
		testRecord{input: "1.2p", output: "1,200,000"},
		testRecord{input: "1.23n", output: "1,230,000,000"},
		testRecord{input: "1.2u", output: "1,200,000,000,000"},
		testRecord{input: "1.2m", output: "1,200,000,000,000,000"},
		testRecord{input: "1.2T", output: "1,200,000,000,000,000,000"},
		testRecord{input: "1200T", output: "1,200,000,000,000,000,000,000"},
		testRecord{input: "1200000T", output: "1,200,000,000,000,000,000,000,000"},
		testRecord{input: "0.2f", output: "200"},
	}

	for _, test := range tests {
		ret := currency.ExpandAndSeparateCurrencyString(test.input)
		if ret != test.output {
			t.Errorf("\nWrong Answer: %s\nGot: %s\nExpected: %s\n", test.input, ret, test.output)
		} else {
			t.Logf("%s => %s", test.input, test.output)
		}
	}
}

func TestFormatCurrencyString(t *testing.T) {

	var tests = []testRecord{
		testRecord{input: "1000", output: "1f"},
		testRecord{input: "10000", output: "10f"},
		testRecord{input: "20000000", output: "20p"},
		testRecord{input: "1234", output: "1.234f"},
		testRecord{input: "1234000000000000000000000", output: "1234000T"},
	}

	for _, test := range tests {
		ret := currency.FormatCurrencyString(test.input)
		if ret != test.output {
			t.Errorf("\nWrong Answer: %s\nGot: %s\nExpected: %s\n", test.input, ret, test.output)
		} else {
			t.Logf("%s => %s", test.input, test.output)
		}
	}
}

func TestFormatTenomString(t *testing.T) {

	var tests = []testRecord{
		testRecord{input: "1m", output: "0.001"},
		testRecord{input: "900000000000000000", output: "0.9"},
		testRecord{input: "1000000000000000000", output: "1"},
		testRecord{input: "2000000000000000000", output: "2"},
		testRecord{input: "1234000000000000000", output: "1.234"},
		testRecord{input: "123400000000000000000", output: "123.4"},
	}

	for _, test := range tests {
		ret := currency.FormatTenomString(test.input)
		if ret != test.output {
			t.Errorf("\nWrong Answer: %s\nGot: %s\nExpected: %s\n", test.input, ret, test.output)
		} else {
			t.Logf("%s => %s", test.input, test.output)
		}
	}
}
