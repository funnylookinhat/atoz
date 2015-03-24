package main

import (
	//"encoding/json"
	"bufio"
	"fmt"
	"os"
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

func (a Action) String() string {
	return "\tName: " + a.Name + "\n" +
		"\tRef: " + a.Ref + "\n" +
		"\tUri: " + a.Uri + "\n" +
		"\tDescription: " + a.Description + "\n" +
		"\tParameterChildren: " + strconv.Itoa(len(a.ParameterChildren)) + "\n" +
		"\tSuccessChildren: " + strconv.Itoa(len(a.SuccessChildren)) + "\n" +
		"\tErrorChildren: " + strconv.Itoa(len(a.ErrorChildren))
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
	Limit       string     `json:"limit"`
	Description string     `json:"description"`
	Children    []KeyValue `json:"children"`
}

const (
	startDefinition = "---ATOZDEF---"
	startAction     = "---ATOZAPI---"
	startObject     = "---ATOZAPI---"
	endDefinition   = "---ATOZEND---"
)

func ParseLineGroups(path string) ([][]string, error) {
	groups := make([][]string, 0)

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
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
		} else if strings.ContainsAny(line, endDefinition) {
			group = append(group, line)
			groups = append(groups, group)
			group = make([]string, 0)
		} else if len(group) > 0 {
			group = append(group, line)
		}
	}

	if len(group) > 0 {
		return nil, fmt.Errorf("Unclosed definition found in %s", path)
	}

	return groups, nil
}

func ParseAction(lines []string) (Action, error) {
	action := Action{}
	return action, fmt.Errorf("This function is not built yet.")
}

func ParseLineType(line string) (string, error) {

	if strings.Contains(line, "@name") {
		return "name", nil
	}

	return "", fmt.Errorf("Not built yet.")
}
