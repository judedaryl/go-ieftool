package internal

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
)

var unwrapScalar = true
var outputFormat = "yaml"
var inputFormat = "yaml"
var exitStatus = false

var xmlAttributePrefix = "+"
var xmlContentName = "+content"
var xmlStrictMode = false
var xmlKeepNamespace = true
var xmlUseRawToken = true

var forceColor = false
var colorsEnabled = false
var indent = 2
var noDocSeparators = false
var splitFileExp = ""

func GetRequestedVariables(content string) []string {
	re := regexp.MustCompile(`(?s){{.*?}}`)
	matches := re.FindAllString(content, -1)
	cleanMatches := []string{}
	for _, match := range matches {
		match = strings.ReplaceAll(match, "{{", "")
		match = strings.ReplaceAll(match, "}}", "")
		match = strings.TrimSpace(match)
		cleanMatches = append(cleanMatches, match)
	}
	return cleanMatches
}

func GetVariable(expression string, configPath string) (value string, cmdError error) {
	fromEnv := os.Getenv("IEF_" + expression)
	if fromEnv != "" {
		return fromEnv, nil
	}
	_experssion := "." + expression
	var err error
	var b bytes.Buffer
	out := bufio.NewWriter(&b)

	if err != nil {
		return "", err
	}

	format, err := yqlib.OutputFormatFromString(outputFormat)
	if err != nil {
		return "", err
	}
	printerWriter, err := configurePrinterWriter(format, out)
	if err != nil {
		return "", err
	}
	encoder := configureEncoder(format)

	printer := yqlib.NewPrinter(encoder, printerWriter)

	decoder, err := configureDecoder()
	if err != nil {
		return "", err
	}
	streamEvaluator := yqlib.NewStreamEvaluator()
	streamEvaluator.EvaluateFiles(_experssion, []string{configPath}, printer, false, decoder)

	if err == nil && exitStatus && !printer.PrintedAnything() {
		return "", errors.New("no matches found")
	}

	return strings.TrimSpace(b.String()), nil
}
