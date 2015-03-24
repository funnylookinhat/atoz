package main

import (
	"testing"
)

/**
 * List Files
 * Get Groupings
 * Parse all Definitions First
 * Parse all Objects
 * Parse all Actions
 * Get Line Type
 * Parse String
 * Parse KeyValue
 * Merge KeyValue
 */

/*
func TestGetLineType(t *testing.T) {
	t.Errorf("TestGetLineType: Not built yet.")
}

func TestParseString(t *testing.T) {
	t.Errorf("TestParseString: Not built yet.")
}
*/

type testParseTypeCase struct {
	line     string
	lineType string
	err      bool
}

var testParseTypeCases = []testParseTypeCase{
	{
		"@name Namespace",
		"name",
		false,
	},
	{
		"@ref Namespace",
		"ref",
		false,
	},
	{
		"@uri Namespace",
		"uri",
		false,
	},
	{
		"@description Description",
		"description",
		false,
	},
	{
		"@include Namespace",
		"include",
		false,
	},
	{
		"@parameter {Type,Limit} Objectspace Description",
		"parameter",
		false,
	},
	{
		"@required {Type,Limit} Objectspace Description",
		"required",
		false,
	},
	{
		"@optional {Type,Limit} Objectspace Description",
		"optional",
		false,
	},
	{
		"@return {Type,Limit} Objectspace Description",
		"return",
		false,
	},
	{
		"@success {Type,Limit} Objectspace Description",
		"success",
		false,
	},
	{
		"@failure {Type,Limit} Objectspace Description",
		"failure",
		false,
	},
}

func TestParseType(t *testing.T) {
	var resultLineType string
	var resultErr error

	for _, test := range testParseTypeCases {
		resultLineType, resultErr = ParseType(test.line)

		if resultErr != nil {
			if !test.err {
				t.Errorf("TestParseType Unexpected error: %s", resultErr)
				return
			}
		} else {
			if test.err {
				t.Errorf("TestParseType - Should have errored out: %s", test.line)
				return
			}
			if resultLineType != test.lineType {
				t.Errorf("TestParseType Line Value Mismatch: %s\nExpected: %s\n  Actual: %s", test.line, test.lineType, resultLineType)
				return
			}
		}
	}
}

type testParseStringCase struct {
	line      string
	lineValue string
	err       bool
}

var testParseStringCases = []testParseStringCase{
	{
		"@ref /Defs/Authorization",
		"/Defs/Authorization",
		false,
	},
	{
		"@name /Defs/BaseResult",
		"/Defs/BaseResult",
		false,
	},
	{
		"@uri Namespace",
		"Namespace",
		false,
	},
	{
		"@description This is a really short description.",
		"This is a really short description.",
		false,
	},
	{
		"@include /Some/Path/To/Something",
		"/Some/Path/To/Something",
		false,
	},
	{
		"@include",
		"",
		true,
	},
	{
		"@include ",
		"",
		true,
	},
}

func TestParseString(t *testing.T) {
	var resultLineValue string
	var resultErr error

	for _, test := range testParseStringCases {
		resultLineValue, resultErr = ParseString(test.line)

		if resultErr != nil {
			if !test.err {
				t.Errorf("TestParseString Unexpected error: %s", resultErr)
				return
			}
		} else {
			if test.err {
				t.Errorf("TestParseString - Should have errored out: %s", test.line)
				return
			}
			if resultLineValue != test.lineValue {
				t.Errorf("TestParseString Line Value Mismatch: %s\nExpected: %s\n  Actual: %s", test.line, test.lineValue, resultLineValue)
				return
			}
		}
	}
}

type testParseKeyValueCase struct {
	line            string
	lineType        string
	lineLimit       int64
	lineObjectspace string
	lineDescription string
	err             bool
}

var testParseKeyValueCases = []testParseKeyValueCase{
	// Integer
	{
		"@required {Integer} Some.Integer This is an integer.",
		"integer",
		-1,
		"some.integer",
		"This is an integer.",
		false,
	},
	// Integer -Error - Integers don't take a limit
	{
		"@required {Integer,0} Some.Integer This is an integer.",
		"integer",
		0,
		"some.integer",
		"This is an integer.",
		true,
	},
	// Integer - Error - no {} around type
	{
		"@required Integer Some.Integer This is an integer.",
		"integer",
		0,
		"some.integer",
		"This is an integer.",
		true,
	},
	// Integer - Error - invalid type
	{
		"@required {BLARG} Some.Fake.Type This should error.",
		"blarg",
		0,
		"some.fake.type",
		"This should error.",
		true,
	},
	// Boolean
	{
		"@required {Boolean} Some.Path.To.Boolean This is a bool.",
		"boolean",
		-1,
		"some.path.to.boolean",
		"This is a bool.",
		false,
	},
	// Boolean - Error - Booleans don't take a limit.
	{
		"@required {Boolean,0} Some.Path.To.Boolean This is a bool.",
		"boolean",
		-1,
		"some.path.to.boolean",
		"This is a bool.",
		true,
	},
	// Decimal
	{
		"@required {Decimal} some.path.to.decimal This is a decimal.",
		"decimal",
		0,
		"some.path.to.decimal",
		"This is a decimal.",
		false,
	},
	// Decimal - Limit 1
	{
		"@required {Decimal,1} some.path.to.decimal This is a decimal.",
		"decimal",
		1,
		"some.path.to.decimal",
		"This is a decimal.",
		false,
	},
	// Decimal - Limit 5
	{
		"@required {Decimal,5} some.path.to.decimal This is a decimal.",
		"decimal",
		5,
		"some.path.to.decimal",
		"This is a decimal.",
		false,
	},
	// Decimal - Error - Invalid limit
	{
		"@required {Decimal,ABC} some.path.to.decimal This is a decimal.",
		"decimal",
		-1,
		"some.path.to.decimal",
		"This is a decimal.",
		true,
	},
	// String
	{
		"@required {String} some.path.to.string This is a string.",
		"string",
		0,
		"some.path.to.string",
		"This is a string.",
		false,
	},
	// String - Limit 1
	{
		"@required {String,1} some.path.to.string This is a string.",
		"string",
		1,
		"some.path.to.string",
		"This is a string.",
		false,
	},
	// String - Limit 5
	{
		"@required {String,5} some.path.to.string This is a string.",
		"string",
		5,
		"some.path.to.string",
		"This is a string.",
		false,
	},
	// String - Error - Invalid limit
	{
		"@required {String,ABC} some.path.to.string This is a string.",
		"string",
		-1,
		"some.path.to.string",
		"This is a string.",
		true,
	},
	// Array
	{
		"@required {Array} some.path.to.array This is an array.",
		"array",
		0,
		"some.path.to.array",
		"This is an array.",
		false,
	},
	// Object
	{
		"@required {Object} some.path.to.object This is an object.",
		"object",
		-1,
		"some.path.to.object",
		"This is an object.",
		false,
	},
}

func TestParseKeyValue(t *testing.T) {
	var resultLineType string
	var resultLineLimit int64
	var resultLineObjectspace string
	var resultLineDescription string
	var resultErr error

	for _, test := range testParseKeyValueCases {
		resultLineType, resultLineLimit, resultLineObjectspace, resultLineDescription, resultErr = ParseKeyValue(test.line)

		if resultErr != nil {
			if !test.err {
				t.Errorf("TestParseKeyValue Unexpected error: %s", resultErr)
				return
			}
		} else {
			if test.err {
				t.Errorf("TestParseKeyValue - Should have errored out: %s", test.line)
				return
			}
			if resultLineType != test.lineType {
				t.Errorf("TestParseKeyValue Line Type Mismatch: %s\nExpected: %s\n  Actual: %s", test.line, test.lineType, resultLineType)
				return
			}
			if resultLineLimit != test.lineLimit {
				t.Errorf("TestParseKeyValue Line Limit Mismatch: %s\nExpected: %s\n  Actual: %s", test.line, test.lineLimit, resultLineLimit)
				return
			}
			if resultLineObjectspace != test.lineObjectspace {
				t.Errorf("TestParseKeyValue Line Objectspace Mismatch: %s\nExpected: %s\n  Actual: %s", test.line, test.lineObjectspace, resultLineObjectspace)
				return
			}
			if resultLineDescription != test.lineDescription {
				t.Errorf("TestParseKeyValue Line Description Mismatch: %s\nExpected: %s\n  Actual: %s", test.line, test.lineDescription, resultLineDescription)
				return
			}
		}

	}
}
