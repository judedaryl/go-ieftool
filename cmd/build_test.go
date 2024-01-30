package cmd

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExecuteBuild(t *testing.T) {
	actual := new(bytes.Buffer)
	cmd := rootCmd
	cmd.SetOut(actual)
	cmd.SetArgs([]string{"build", "--config", "../test/fixtures/config.yaml", "--source", "../test/fixtures/base", "--destination", "../build"})
	cmd.Execute()

	log.Println(actual.String())
	assert.Equal(t, "", actual.String(), "error is expected")
}

func Test_ExecuteBuildErrorOnMissingVariables(t *testing.T) {
	actual := new(bytes.Buffer)
	cmd := rootCmd
	cmd.SetErr(actual)
	cmd.SetArgs([]string{"build", "--config", "../test/fixtures/config.yaml", "--source", "../test/fixtures/missingVariables", "--destination", "../build"})
	cmd.Execute()

	log.Println(actual.String())
	assert.NotEmpty(t, actual.String(), "error is expected")
}
