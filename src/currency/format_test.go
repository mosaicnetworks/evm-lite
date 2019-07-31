package currency_test

import (
	"testing"

	"github.com/mosaicnetworks/evm-lite/src/currency"
)

type testRecord struct {
	input  string
	output string
}

//  1 000 000 000 000 000 000 000 000	yotta (Y)	1024
//  1 000 000 000 000 000 000 000		zetta (Z)	1021
//  1 000 000 000 000 000 000			exa (E)		1018
//  1 000 000 000 000 000				peta (P)	1015
//  1 000 000 000 000					tera (T)	1012
//  1 000 000 000						giga (G)	109
//  1 000 000							mega (M)	106
//  1 000								kilo (k)	103

func TestExpandCurrencyString(t *testing.T) {

	var tests = []testRecord{
		testRecord{input: "1K", output: "1000"},
		testRecord{input: "1.2M", output: "1200000"},
		testRecord{input: "1.23G", output: "1230000000"},
		testRecord{input: "1.2T", output: "1200000000000"},
		testRecord{input: "1.2P", output: "1200000000000000"},
		testRecord{input: "1.2E", output: "1200000000000000000"},
		testRecord{input: "1.2Z", output: "1200000000000000000000"},
		testRecord{input: "1.2Y", output: "1200000000000000000000000"},
		testRecord{input: "0.2K", output: "200"},
		testRecord{input: "0x122K", output: "0x122000"},
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
		testRecord{input: "1K", output: "1,000"},
		testRecord{input: "1.2M", output: "1,200,000"},
		testRecord{input: "1.23G", output: "1,230,000,000"},
		testRecord{input: "1.2T", output: "1,200,000,000,000"},
		testRecord{input: "1.2P", output: "1,200,000,000,000,000"},
		testRecord{input: "1.2E", output: "1,200,000,000,000,000,000"},
		testRecord{input: "1.2Z", output: "1,200,000,000,000,000,000,000"},
		testRecord{input: "1.2Y", output: "1,200,000,000,000,000,000,000,000"},
		testRecord{input: "0.2K", output: "200"},
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
		testRecord{input: "1000", output: "1K"},
		testRecord{input: "10000", output: "10K"},
		testRecord{input: "20000000", output: "20M"},
		testRecord{input: "1234", output: "1.234K"},
		testRecord{input: "1234000000000000000000000000", output: "1234Y"},
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
