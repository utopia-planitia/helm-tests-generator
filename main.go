package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed job.yaml.gtpl
var yamlTemplate string

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

type TestDetails struct {
	TestName string
	Image    string
	Command  []string
}

func run() error {
	testFiles, err := listTestfiles()
	if err != nil {
		return err
	}

	tmpl, err := template.New("test").Parse(yamlTemplate)

	tests := []TestDetails{}

	for _, testFile := range testFiles {
		extension := filepath.Ext(testFile)

		image, err := image(extension)
		if err != nil {
			return err
		}

		command, err := command(filepath.Base(testFile))
		if err != nil {
			return err
		}

		testDetails := TestDetails{
			TestName: strings.TrimSuffix(filepath.Base(testFile), extension),
			Command:  command,
			Image:    image,
		}

		tests = append(tests, testDetails)
	}

	err = tmpl.Execute(os.Stdout, tests)
	if err != nil {
		return err
	}

	return nil
}

func command(testFile string) ([]string, error) {
	extension := filepath.Ext(testFile)
	commands := make(map[string][]string)

	commands[".sh"] = []string{
		"/bin/bash",
		"-o=pipefail",
		"-eu",
		"/test/" + testFile,
	}

	command, ok := commands[extension]
	if !ok {
		return []string{}, fmt.Errorf("no command is defined for %s scripts", extension)
	}

	return command, nil
}

func image(extension string) (string, error) {
	images := make(map[string]string)

	// TODO update via renovate bot
	images[".sh"] = "utopiaplanitia/helm-tools:v1.0.2"

	image, ok := images[extension]
	if !ok {
		return "", fmt.Errorf("no image is defined for %s scripts", extension)
	}

	return image, nil
}

func listTestfiles() ([]string, error) {
	files := []string{}

	files, err := filepath.Glob(filepath.Join("test-scripts", "*.sh"))
	if err != nil {
		return []string{}, err
	}

	return files, nil
}
