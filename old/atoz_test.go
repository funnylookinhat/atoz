package main

import (
	"reflect"
	"testing"
)

type parseLineTypeTest struct {
	n        int
	line     string
	expected string
}

var parseLineTypeTests = []parseLineTypeTest{
	{
		1,
		" * @name User Lookup",
		"name",
	},
	{
		1,
		"* @name User Lookup",
		"name",
	},
	{
		1,
		" @name User Lookup",
		"name",
	},
	{
		1,
		"@name User Lookup",
		"name",
	},
}

func TestParseLineType(t *testing.T) {
	for _, test := range parseLineTypeTests {
		result, err := ParseLineType(test.line)

		if err != nil {
			t.Errorf("ParseLineType Test %d Failed with an Error: %v", test.n, err)
		}

		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("ParseLineType Test %d Failed \nExpected: \n%v\nResult: \n%v", test.n, test.expected, result)
		}
	}
}

/*
var parseActionTests = []parseActionTest{
	{
		1,
		[]string{
			" * @name User Lookup",
			" * @ref /MyApp/User/Lookup",
			" * @uri /User/Lookup",
			" * @description Get the information for a user.",
		},
		Action{
			Name:        "User Lookup",
			Ref:         "/MyApp/User/Lookup",
			Uri:         "/User/Lookup",
			Description: "Get the information for a user.",
		},
	},
}

func TestParseAction(t *testing.T) {
	for _, test := range parseActionTests {
		result, err := ParseAction(test.lines)

		if err != nil {
			t.Errorf("ParseAction Test %d Failed with an Error: %v", test.n, err)
		}

		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("ParseAction Test %d Failed \nExpected: \n%v\nResult: \n%v", test.n, test.expected, result)
		}
	}
}

func TestParseLine(line string) whatever {

}

*/
