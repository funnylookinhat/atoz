package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Definition struct {
	Ref      string     `json:"ref"`
	Children []KeyValue `json:"children:"`
}

type Action struct {
	Name              string     `json:"name"`
	Ref               string     `json:"ref"`
	Uri               string     `json:"uri"`
	Description       string     `json:"description"`
	ParameterChildren []KeyValue `json:"parameters"`
	SuccessChildren   []KeyValue `json:"returnSuccess"`
	ErrorChildren     []KeyValue `json:"returnError"`
}

type Object struct {
	Name        string     `json:"name"`
	Ref         string     `json:"ref"`
	Description string     `json:"description"`
	Children    []KeyValue `json:"children"`
}

type KeyValue struct {
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	Limit       int        `json:"limit"`
	Description string     `json:"description"`
	Children    []KeyValue `json:"children"`
}

const (
	startDefinition = "---ATOZDEF---"
	startAction     = "---ATOZAPI---"
	startObject     = "---ATOZOBJ---"
	endDefinition   = "---ATOZEND---"
)

func GetLineType(line string) (string, error) {
	return "", nil
}

// Receive
// @name Something something
// Return "name"
func ParseLineType(line string) (string, error) {
	var lineTypes = map[string]bool{
		"@name":        true,
		"@ref":         true,
		"@uri":         true,
		"@description": true,
		"@include":     true,
		"@parameter":   true,
		"@required":    true,
		"@optional":    true,
		"@return":      true,
		"@success":     true,
		"@failure":     true,
	}

	var returnValue string

	atIndex := strings.Index(line, "@")

	if atIndex < 0 {
		return "", fmt.Errorf("Invalid line - missing @declaration.")
	}

	line = line[atIndex:]

	lineParts := strings.Split(line, " ")

	if len(lineParts) < 1 {
		return "", fmt.Errorf("Invalid line - missing @declaration.")
	}

	if _, ok := lineTypes[lineParts[0]]; !ok {
		return "", fmt.Errorf("Invalid line - unknown @declaration type.")
	}

	returnValue = lineParts[0]

	return returnValue[1:len(returnValue)], nil
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
