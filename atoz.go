package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

type ApiSpec struct {
	Actions []Action `json:"actions"`
	Objects []Object `json:"objects"`
}

type Action struct {
	Name        string     `json:"name"`
	Ref         string     `json:"ref"`
	Uri         string     `json:"uri"`
	Description string     `json:"description"`
	Parameters  []KeyValue `json:"parameters"`
	Returns     []KeyValue `json:"returns"`
}

type Object struct {
	Name        string     `json:"name"`
	Ref         string     `json:"ref"`
	Description string     `json:"description"`
	Properties  []KeyValue `json:"properties"`
}

type KeyValue struct {
	Name        string     `json:"name"`
	Flag        string     `json:"name"`
	Type        string     `json:"type"`
	Limit       int64      `json:"limit"`
	Description string     `json:"description"`
	Children    []KeyValue `json:"children"`
}

const (
	startDefinition = "---ATOZDEF---"
	startAction     = "---ATOZAPI---"
	startObject     = "---ATOZOBJ---"
	endDefinition   = "---ATOZEND---"
)

func GenerateApiSpec(files []string) (ApiSpec, error) {
	var err error

	groups := make([][]string, 0)

	var definitionGroups map[string][]string
	var actionGroups map[string][]string
	var objectGroups map[string][]string

	for _, path := range files {
		file, err := os.Open(path)
		reader := bufio.NewReader(file)

		if err != nil {
			return ApiSpec{}, err
		}

		parseGroupsFiles, parseGroupsErr := ParseGroups(reader)

		if parseGroupsErr != nil {
			return ApiSpec{}, parseGroupsErr
		}

		for _, parseGroupFile := range parseGroupsFiles {
			groups = append(groups, parseGroupFile)
		}
	}

	definitionGroups, err = GetDefinitionGroups(groups)

	if err != nil {
		return ApiSpec{}, err
	}

	actionGroups, err = GetActionGroups(groups)

	if err != nil {
		return ApiSpec{}, err
	}

	objectGroups, err = GetObjectGroups(groups)

	if err != nil {
		return ApiSpec{}, err
	}

	fmt.Printf("%v", definitionGroups)
	fmt.Printf("%v", actionGroups)
	fmt.Printf("%v", objectGroups)

	/**
	 * List Files
	 * Get Groupings
	 * Get Definitions Hashmap
	 * 		Remove First and Last Line
	 * Get Objects Hashmap
	 * 		Remove First and Last Line
	 * Get Actions Hashmap
	 * 		Remove First and Last Line
	 * Parse all Groups ( Definitions, then Objects, then Actions )
	 * 		Get Line Type
	 *   	Parse String
	 *    	Parse KeyValue
	 *     	Merge KeyValue
	 */

	return ApiSpec{}, nil
}

func ParseGroups(r *bufio.Reader) ([][]string, error) {
	groups := make([][]string, 0)

	scanner := bufio.NewScanner(r)

	var line string
	group := make([]string, 0)

	for scanner.Scan() {
		line = scanner.Text()

		if !utf8.ValidString(line) {
			return make([][]string, 0), nil
		}

		if strings.Contains(line, startDefinition) ||
			strings.Contains(line, startAction) ||
			strings.Contains(line, startObject) {
			group = append(group, line)
		} else if strings.Contains(line, endDefinition) {
			group = append(group, line)
			groups = append(groups, group)
			group = make([]string, 0)
		} else if len(group) > 0 {
			group = append(group, line)
		}
	}

	if len(group) > 0 {
		return nil, fmt.Errorf("Unclosed definition found.")
	}

	return groups, nil
}

// Receive
// @name Something something
// Return "name"
func ParseLineType(line string) (string, error) {
	var lineTypes = map[string]string{
		"@name":        "name",
		"@ref":         "ref",
		"@uri":         "uri",
		"@description": "description",
		"@include":     "include",
		"@parameter":   "parameter",
		"@required":    "parameter",
		"@optional":    "parameter",
		"@return":      "return",
		"@success":     "return",
		"@failure":     "return",
		"@property":    "property",
	}

	var returnValue string
	var ok bool

	atIndex := strings.Index(line, "@")

	if atIndex < 0 {
		return "", fmt.Errorf("Invalid line - missing @declaration.")
	}

	line = line[atIndex:]

	lineParts := strings.Split(line, " ")

	if len(lineParts) < 1 {
		return "", fmt.Errorf("Invalid line - missing @declaration.")
	}

	if returnValue, ok = lineTypes[lineParts[0]]; !ok {
		return "", fmt.Errorf("Invalid line - unknown @declaration type. " + lineParts[0])
	}

	return returnValue, nil
}

func ParseLineFlag(line string) (string, error) {
	var lineFlags = map[string]string{
		"@required": "required",
		"@optional": "optional",
		"@success":  "success",
		"@error":    "error",
	}

	var returnValue string
	var ok bool

	atIndex := strings.Index(line, "@")

	if atIndex < 0 {
		return "", fmt.Errorf("Invalid line - missing @declaration.")
	}

	line = line[atIndex:]

	lineParts := strings.Split(line, " ")

	if len(lineParts) < 1 {
		return "", fmt.Errorf("Invalid line - missing @declaration.")
	}

	// If no flag, return a blank string.
	if returnValue, ok = lineFlags[lineParts[0]]; !ok {
		return "", nil
	}

	return returnValue, nil
}

// Receive
// @name Some string value
// Return Value
func ParseLineString(line string) (string, error) {
	var returnValue string

	atIndex := strings.Index(line, "@")

	if atIndex < 0 {
		return "", fmt.Errorf("Invalid line - missing @declaration.")
	}

	line = line[atIndex:]

	lineParts := strings.Split(line, " ")

	if len(lineParts) < 2 {
		return "", fmt.Errorf("Invalid line - missing value.")
	}

	returnValue = strings.TrimSpace(strings.Join(lineParts[1:], " "))

	if len(returnValue) < 1 {
		return "", fmt.Errorf("Invalid line - missing value.")
	}

	return returnValue, nil
}

// Receive
// @returns {Type,Limit} Objectspace Description
// Return Type, Limit, Objectspace, Description
func ParseLineKeyValue(line string) (string, int64, string, string, error) {
	lineTypeLimits := map[string]bool{
		"boolean": false,
		"integer": false,
		"decimal": true,
		"string":  true,
		"array":   true,
		"object":  false,
	}

	var returnType string
	var returnLimit int64
	var returnObjectspace string
	var returnDescription string
	var err error

	atIndex := strings.Index(line, "@")

	if atIndex < 0 {
		return "", -1, "", "", fmt.Errorf("Invalid line - missing @declaration.")
	}

	line = line[atIndex:]

	lineParts := strings.Split(line, " ")

	if len(lineParts) < 4 {
		return "", -1, "", "", fmt.Errorf("Invalid line - missing one or more statements.")
	}

	lineType := lineParts[1]

	if strings.Index(lineType, "{") != 0 || strings.Index(lineType, "}") != (len(lineType)-1) {
		return "", -1, "", "", fmt.Errorf("Invalid line - missing {} type.")
	}

	lineType = lineType[1 : len(lineType)-1]

	lineTypeParts := strings.Split(lineType, ",")

	returnType = strings.ToLower(lineTypeParts[0])

	var returnTypeHasLimit bool
	var ok bool

	if returnTypeHasLimit, ok = lineTypeLimits[returnType]; !ok {
		return "", -1, "", "", fmt.Errorf("Invalid type: %s", returnType)
	}

	if len(lineTypeParts) > 2 {
		return "", -1, "", "", fmt.Errorf("Invalid {} type - must be in format {Type,Limit}")
	} else if len(lineTypeParts) == 2 {
		if _, err := strconv.Atoi(lineTypeParts[1]); err != nil {
			return "", -1, "", "", fmt.Errorf("Invalid Type Limit - must be an integer.")
		}

		returnLimit, err = strconv.ParseInt(lineTypeParts[1], 10, 64)

		if err != nil {
			return "", -1, "", "", err
		}
	} else {
		if returnTypeHasLimit {
			returnLimit = 0
		} else {
			returnLimit = -1
		}
	}

	if returnLimit >= 0 && !returnTypeHasLimit {
		return "", -1, "", "", fmt.Errorf("Invalid limit: %s does not accept a limit.", returnType)
	}

	if len(lineTypeParts) > 2 {
		return "", -1, "", "", fmt.Errorf("Invalid {} type - must be in format {Type,Limit}.")
	}

	returnObjectspace = strings.ToLower(lineParts[2])

	returnDescription = strings.Join(lineParts[3:], " ")

	return returnType, returnLimit, returnObjectspace, returnDescription, nil
}

func ParseGroupType(line string) (string, error) {
	if strings.Contains(line, startDefinition) {
		return "definition", nil
	}
	if strings.Contains(line, startAction) {
		return "action", nil
	}
	if strings.Contains(line, startObject) {
		return "object", nil
	}

	return "", fmt.Errorf("Invalid line: no starting group identifier found.")
}

func ParseGroupRef(group []string) (string, error) {
	var lineType string
	var err error

	for _, line := range group {
		lineType, err = ParseLineType(line)

		if err != nil {
			return "", err
		}

		if lineType == "ref" {
			var lineValue, err = ParseLineString(line)

			if err != nil {
				return "", err
			}

			return lineValue, nil
		}
	}

	return "", fmt.Errorf("No line type found.")
}

/*
func GenerateDefinition(lines []string) (Definition, error) {
	return Definition{}, nil
}
*/

func GenerateObject(group []string, definitions map[string][]string) (Object, error) {
	returnObject := Object{}

	var err error
	var lineType string

	for _, line := range group {
		lineType, err = ParseLineType(line)

		if err != nil {
			return returnObject, err
		}

		if lineType == "name" {
			returnObject.Name, err = ParseLineString(line)

			if err != nil {
				return returnObject, err
			}
		}
		else if lineType == "ref" {
			returnObject.Ref, err = ParseLineString(line)

			if err != nil {
				return returnObject, err
			}
		}
		else if lineType == "description" {
			returnObject.Description, err = ParseLineString(line)

			if err != nil {
				return returnObject, err
			}
		}
	}

	returnObject.Properties, err = GenerateKeyValues("property", group, "")

	if err != nil {
		return returnObject, err
	}

	return returnObject, nil
}

func GenerateAction(group []string, definitions map[string][]string) (Action, error) {
	return Action{}, nil
}

func GenerateKeyValues(keyValueType string, lines []string, objectspace string) ([]KeyValue, error) {
	keyValues := make([]KeyValue, 0)

	// Unpacking each line each iteration will be a bit more inefficient,
	// but should provide a nice proof-of-concept
	var lineKeyValueType string
	var lineKeyValueLimit int64
	var lineKeyValueObjectspace string
	var lineKeyValueDescription string
	var lineKeyValueError error

	var lineType string
	var lineTypeError error

	var lineFlag string
	// var lineFlagError error

	var lineKeyValue KeyValue

	for _, line := range lines {
		if !strings.Contains(line, startDefinition) &&
			!strings.Contains(line, startAction) &&
			!strings.Contains(line, startObject) &&
			!strings.Contains(line, endDefinition) {

			lineType, lineTypeError = ParseLineType(line)

			lineFlag = ""
			// lineFlagError = nil

			if lineTypeError != nil {
				return make([]KeyValue, 0), lineTypeError
			}

			if lineType == keyValueType {

				lineKeyValueType, lineKeyValueLimit, lineKeyValueObjectspace, lineKeyValueDescription, lineKeyValueError = ParseLineKeyValue(line)

				if lineKeyValueError != nil {
					return keyValues, lineKeyValueError
				}

				if strings.Contains(lineKeyValueObjectspace, objectspace) &&
					strings.Index(strings.Replace(lineKeyValueObjectspace, objectspace, "", 1), ".") < 0 {
					lineKeyValue = KeyValue{
						strings.Replace(lineKeyValueObjectspace, objectspace, "", 1),
						lineFlag,
						lineKeyValueType,
						lineKeyValueLimit,
						lineKeyValueDescription,
						make([]KeyValue, 0),
					}

					lineKeyValue.Children, lineKeyValueError = GenerateKeyValues(keyValueType, lines, lineKeyValueObjectspace+".")

					if lineKeyValueError != nil {
						return make([]KeyValue, 0), lineKeyValueError
					}

					keyValues = append(keyValues, lineKeyValue)
				}

			}
		}
	}

	return keyValues, nil
}

func GetDefinitionGroups(groups [][]string) (map[string][]string, error) {
	definitionGroups := make(map[string][]string, 0)

	for _, group := range groups {
		if groupType, err := ParseGroupType(group[0]); err != nil {
			return definitionGroups, err
		} else {
			if groupType == "definition" {
				group = group[1:]
				group = group[0 : len(group)-1]
				if groupRef, err := ParseGroupRef(group); err != nil {
					return definitionGroups, err
				} else {
					definitionGroups[groupRef] = group
				}
			}
		}
	}

	return definitionGroups, nil
}

func GetObjectGroups(groups [][]string) (map[string][]string, error) {
	objectGroups := make(map[string][]string, 0)

	for _, group := range groups {
		if groupType, err := ParseGroupType(group[0]); err != nil {
			return objectGroups, err
		} else {
			if groupType == "object" {
				group = group[1:]
				group = group[0 : len(group)-1]
				if groupRef, err := ParseGroupRef(group); err != nil {
					return objectGroups, err
				} else {
					objectGroups[groupRef] = group
				}
			}
		}
	}

	return objectGroups, nil
}

func GetActionGroups(groups [][]string) (map[string][]string, error) {
	actionGroups := make(map[string][]string, 0)

	for _, group := range groups {
		if groupType, err := ParseGroupType(group[0]); err != nil {
			return actionGroups, err
		} else {
			if groupType == "action" {
				group = group[1:]
				group = group[0 : len(group)-1]
				if groupRef, err := ParseGroupRef(group); err != nil {
					return actionGroups, err
				} else {
					actionGroups[groupRef] = group
				}
			}
		}
	}

	return actionGroups, nil
}
