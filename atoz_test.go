package main

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

/**
 * List Files
 * Get Groupings
 * Parse all Groups ( Definitions, then Objects, then Actions )
 * 		Get Line Type
 *   	Parse String
 *    	Parse KeyValue
 *     	Merge KeyValue
 */

type testParseLineTypeCase struct {
	line     string
	lineType string
	err      bool
}

var testParseLineTypeCases = []testParseLineTypeCase{
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

func TestParseLineType(t *testing.T) {
	var resultLineType string
	var resultErr error

	for _, test := range testParseLineTypeCases {
		resultLineType, resultErr = ParseLineType(test.line)

		if resultErr != nil {
			if !test.err {
				t.Errorf("TestParseLineType Unexpected error: %s", resultErr)
				return
			}
		} else {
			if test.err {
				t.Errorf("TestParseLineType - Should have errored out: %s", test.line)
				return
			}
			if resultLineType != test.lineType {
				t.Errorf("TestParseLineType Line Value Mismatch: %s\nExpected: %s\n  Actual: %s", test.line, test.lineType, resultLineType)
				return
			}
		}
	}
}

type testParseLineStringCase struct {
	line      string
	lineValue string
	err       bool
}

var testParseLineStringCases = []testParseLineStringCase{
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

func TestParseLineString(t *testing.T) {
	var resultLineValue string
	var resultErr error

	for _, test := range testParseLineStringCases {
		resultLineValue, resultErr = ParseLineString(test.line)

		if resultErr != nil {
			if !test.err {
				t.Errorf("TestParseLineString Unexpected error: %s", resultErr)
				return
			}
		} else {
			if test.err {
				t.Errorf("TestParseLineString - Should have errored out: %s", test.line)
				return
			}
			if resultLineValue != test.lineValue {
				t.Errorf("TestParseLineString Line Value Mismatch: %s\nExpected: %s\n  Actual: %s", test.line, test.lineValue, resultLineValue)
				return
			}
		}
	}
}

type testParseLineKeyValueCase struct {
	line            string
	lineType        string
	lineLimit       int64
	lineObjectspace string
	lineDescription string
	err             bool
}

var testParseLineKeyValueCases = []testParseLineKeyValueCase{
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

func TestParseLineKeyValue(t *testing.T) {
	var resultLineType string
	var resultLineLimit int64
	var resultLineObjectspace string
	var resultLineDescription string
	var resultErr error

	for _, test := range testParseLineKeyValueCases {
		resultLineType, resultLineLimit, resultLineObjectspace, resultLineDescription, resultErr = ParseLineKeyValue(test.line)

		if resultErr != nil {
			if !test.err {
				t.Errorf("TestParseLineKeyValue Unexpected error: %s", resultErr)
				return
			}
		} else {
			if test.err {
				t.Errorf("TestParseLineKeyValue - Should have errored out: %s", test.line)
				return
			}
			if resultLineType != test.lineType {
				t.Errorf("TestParseLineKeyValue Line Type Mismatch: %s\nExpected: %s\n  Actual: %s", test.line, test.lineType, resultLineType)
				return
			}
			if resultLineLimit != test.lineLimit {
				t.Errorf("TestParseLineKeyValue Line Limit Mismatch: %s\nExpected: %s\n  Actual: %s", test.line, test.lineLimit, resultLineLimit)
				return
			}
			if resultLineObjectspace != test.lineObjectspace {
				t.Errorf("TestParseLineKeyValue Line Objectspace Mismatch: %s\nExpected: %s\n  Actual: %s", test.line, test.lineObjectspace, resultLineObjectspace)
				return
			}
			if resultLineDescription != test.lineDescription {
				t.Errorf("TestParseLineKeyValue Line Description Mismatch: %s\nExpected: %s\n  Actual: %s", test.line, test.lineDescription, resultLineDescription)
				return
			}
		}

	}
}

type testParseGroupsCase struct {
	lines  string
	groups [][]string
	err    bool
}

var testParseGroupsCases = []testParseGroupsCase{
	{
		`

#include <stdio.h>
 
int main(void)
{
    printf("hello, world\n");
}

/**
 * ---ATOZDEF---
 * @ref /Defs/Authorization
 * @parameter {Object} auth 
 * @parameter {Integer} auth.id 
 * @parameter {String,64} auth.key 
 * ---ATOZEND---
 */

#include <stdio.h>
 
int main(void)
{
    printf("hello, world\n");
}

/**
 * ---ATOZDEF---
 * @name /Defs/BaseResult
 * @success {Boolean} success A boolean to show whether or not the request was successful.
 * @error {String} error An error message describing what went wrong.
 * ---ATOZEND---
 */

#include <stdio.h>
 
int main(void)
{
    printf("hello, world\n");
}

/**
 * ---ATOZAPI---
 * @name User Lookup
 * @ref /MyApp/User/Lookup
 * @uri /User/Lookup
 * @description Get the information for a user.
 * @include /Defs/Authorization
 * @parameter {Integer} id The ID of the user.
 * @include /Defs/BaseResult
 * @success {#/Application/User#} user
 * ---ATOZEND---
 */

#include <stdio.h>
 
int main(void)
{
    printf("hello, world\n");
}

/**
 * ---ATOZOBJ---
 * @name User
 * @ref /Application/User
 * @description A user in the application.
 * @property id INTEGER Unique ID of the user.
 * @property name STRING Name of the user.
 * @property email STRING Email address for the user.
 * ---ATOZEND---
 */

#include <stdio.h>
 
int main(void)
{
    printf("hello, world\n");
}
		`,
		[][]string{
			{
				" * ---ATOZDEF---",
				" * @ref /Defs/Authorization",
				" * @parameter {Object} auth ",
				" * @parameter {Integer} auth.id ",
				" * @parameter {String,64} auth.key ",
				" * ---ATOZEND---",
			},
			{
				" * ---ATOZDEF---",
				" * @name /Defs/BaseResult",
				" * @success {Boolean} success A boolean to show whether or not the request was successful.",
				" * @error {String} error An error message describing what went wrong.",
				" * ---ATOZEND---",
			},
			{
				" * ---ATOZAPI---",
				" * @name User Lookup",
				" * @ref /MyApp/User/Lookup",
				" * @uri /User/Lookup",
				" * @description Get the information for a user.",
				" * @include /Defs/Authorization",
				" * @parameter {Integer} id The ID of the user.",
				" * @include /Defs/BaseResult",
				" * @success {#/Application/User#} user",
				" * ---ATOZEND---",
			},
			{
				" * ---ATOZOBJ---",
				" * @name User",
				" * @ref /Application/User",
				" * @description A user in the application.",
				" * @property id INTEGER Unique ID of the user.",
				" * @property name STRING Name of the user.",
				" * @property email STRING Email address for the user.",
				" * ---ATOZEND---",
			},
		},
		false,
	},
	{
		`

#include <stdio.h>
 
int main(void)
{
    printf("hello, world\n");
}

/**
 * ---ATOZDEF---
 * @ref /Defs/Authorization
 * @parameter {Object} auth 
 * @parameter {Integer} auth.id 
 * @parameter {String,64} auth.key 
 */
		`,
		[][]string{
			{},
		},
		true,
	},
}

func TestParseGroups(t *testing.T) {
	var resultLineGroups [][]string
	var resultErr error

	for _, test := range testParseGroupsCases {
		buffer := bytes.NewBufferString(test.lines)
		reader := bufio.NewReader(buffer)
		resultLineGroups, resultErr = ParseGroups(reader)

		if resultErr != nil {
			if !test.err {
				t.Errorf("TestParseGroups Unexpected error: %s", resultErr)
				return
			}
		} else {
			for i, _ := range test.groups {
				if !reflect.DeepEqual(test.groups[i], resultLineGroups[i]) {
					t.Errorf("TestParseGroups Groups Mismatch:")
					t.Errorf("Expected:")
					for _, line := range test.groups[i] {
						t.Errorf("\t%s", line)
					}
					t.Errorf("Actual:")
					for _, line := range resultLineGroups[i] {
						t.Errorf("\t%s", line)
					}
				}
			}
		}
	}
}

type testParseGroupTypeCase struct {
	line      string
	groupType string
	err       bool
}

var testParseGroupTypeCases = []testParseGroupTypeCase{
	{
		" * ---ATOZDEF---",
		"definition",
		false,
	},
	{
		" * ---ATOZAPI---",
		"action",
		false,
	},
	{
		" * ---ATOZOBJ---",
		"object",
		false,
	},
	{
		" * ---ATOZEND---",
		"",
		true,
	},
}

func TestParseGroupType(t *testing.T) {
	var resultGroupType string
	var resultErr error

	for _, test := range testParseGroupTypeCases {
		resultGroupType, resultErr = ParseGroupType(test.line)

		if resultErr != nil {
			if !test.err {
				t.Errorf("TestParseGroupType Unexpected error: %s", resultErr)
				return
			}
		} else {
			if test.groupType != resultGroupType {
				t.Errorf("TestParseGroupType Group Type Mismatch: %s\nExpected: %s\n  Actual: %s", test.line, test.groupType, resultGroupType)
				return
			}
		}
	}
}

type testGenerateKeyValuesCase struct {
	keyValueType    string
	definition      string
	resultKeyValues [][]KeyValue
	resultErr       bool
}

var testGenerateKeyValuesCases = []testGenerateKeyValuesCase{
	{
		"property",
		`
/**
 * ---ATOZOBJ---
 * @name User
 * @ref /Application/User
 * @description A user in the application.
 * @property {Integer} id Unique ID of the user.
 * @property {String} name Name of the user.
 * @property {String} email Email address for the user.
 * ---ATOZEND---
 */
		`,
		[][]KeyValue{
			[]KeyValue{
				{
					"id",
					"integer",
					-1,
					"Unique ID of the user.",
					[]KeyValue{},
				},
				{
					"name",
					"string",
					0,
					"Name of the user.",
					[]KeyValue{},
				},
				{
					"email",
					"string",
					0,
					"Email address for the user.",
					[]KeyValue{},
				},
			},
		},
		false,
	},
	{
		"property",
		`
/**
 * ---ATOZOBJ---
 * @name User
 * @ref /Application/User
 * @description A user in the application.
 * @property {Object} user The user.
 * @property {Integer} user.id Unique ID of the user.
 * @property {String} user.name Name of the user.
 * @property {String} user.email Email address for the user.
 * @property {String} user.role The role of the user.
 * ---ATOZEND---
 */
		`,
		[][]KeyValue{
			[]KeyValue{
				KeyValue{
					"user",
					"object",
					-1,
					"The user.",
					[]KeyValue{
						KeyValue{
							"id",
							"integer",
							-1,
							"Unique ID of the user.",
							[]KeyValue{},
						},
						KeyValue{
							"name",
							"string",
							0,
							"Name of the user.",
							[]KeyValue{},
						},
						KeyValue{
							"email",
							"string",
							0,
							"Email address for the user.",
							[]KeyValue{},
						},
						KeyValue{
							"role",
							"string",
							0,
							"The role of the user.",
							[]KeyValue{},
						},
					},
				},
			},
		},
		false,
	},
	{
		"required",
		`
/**
 * ---ATOZDEF---
 * @ref /Application/Auth
 * @required {Object} auth Auth object.
 * @required {Object} auth.user User object.
 * @required {Integer} auth.user.id User ID.
 * @required {String} auth.user.name First and last ( or common ) name.
 * @required {String} auth.user.email Email address.
 * @required {Object} auth.token Token object.
 * @required {String} auth.token.key Token key.
 * @required {String} auth.token.secret Token secret.
 * ---ATOZEND---
 */
		`,
		[][]KeyValue{
			[]KeyValue{
				KeyValue{
					"auth",
					"object",
					-1,
					"Auth object.",
					[]KeyValue{
						KeyValue{
							"user",
							"object",
							-1,
							"User object.",
							[]KeyValue{
								KeyValue{
									"id",
									"integer",
									-1,
									"User ID.",
									[]KeyValue{},
								},
								KeyValue{
									"name",
									"string",
									0,
									"First and last ( or common ) name.",
									[]KeyValue{},
								},
								KeyValue{
									"email",
									"string",
									0,
									"Email address.",
									[]KeyValue{},
								},
							},
						},
						KeyValue{
							"token",
							"object",
							-1,
							"Token object.",
							[]KeyValue{
								KeyValue{
									"key",
									"string",
									0,
									"Token key.",
									[]KeyValue{},
								},
								KeyValue{
									"secret",
									"string",
									0,
									"Token secret.",
									[]KeyValue{},
								},
							},
						},
					},
				},
			},
		},
		false,
	},
}

func TestGenerateKeyValues(t *testing.T) {
	var resultKeyValues []KeyValue
	var resultGroups [][]string
	var resultErr error

	for _, test := range testGenerateKeyValuesCases {

		buffer := bytes.NewBufferString(test.definition)
		reader := bufio.NewReader(buffer)
		resultGroups, resultErr = ParseGroups(reader)

		if resultErr != nil {
			t.Errorf("Unexpected Error: Error parsing group definitions: %s", resultErr)
			return
		}

		for i, _ := range resultGroups {
			resultKeyValues, resultErr = GenerateKeyValues(test.keyValueType, resultGroups[i], "")

			if resultErr != nil {
				t.Errorf("Error generating keyvalues: %s", resultErr)
				return
			}

			if !reflect.DeepEqual(test.resultKeyValues[i], resultKeyValues) {
				t.Errorf("TestParseGroupType Group Type Mismatch: \n%s\nExpected: %s\n  Actual: %s", test.definition, test.resultKeyValues[i], resultKeyValues)
				return
			}
		}
	}
}