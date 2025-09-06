package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/urfave/cli/v3"
)

//go:embed job.yaml.gtpl
var yamlTemplate string

func main() {
	err := run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	app := &cli.App{
		Name:   "helm tests generator",
		Usage:  "generate helm test jobs based on test scripts",
		Action: render,
	}

	err := app.Run(os.Args)
	if err != nil {
		return err
	}

	return nil
}

func render(c *cli.Context) error {
	testspath := "test-scripts"

	if c.Args().Present() {
		testspath = c.Args().First()
	}

	if len(c.Args().Tail()) != 0 {
		return fmt.Errorf("please provide only one path")
	}

	tmpl, err := template.New("test").Parse(yamlTemplate)
	if err != nil {
		return err
	}

	tests, err := allTests(testspath)
	if err != nil {
		return err
	}

	err = tmpl.Execute(os.Stdout, tests)
	if err != nil {
		return err
	}

	return nil
}

type Testcase interface {
	Name() string
	Image() string
	Command() []string
}

func allTests(testspath string) ([]Testcase, error) {
	tests := []Testcase{}

	shelltests, err := shellTests(testspath)
	if err != nil {
		return []Testcase{}, err
	}

	tests = append(tests, shelltests...)

	batstests, err := batsTests(testspath)
	if err != nil {
		return []Testcase{}, err
	}

	tests = append(tests, batstests...)

	return tests, nil
}

func shellTests(testspath string) ([]Testcase, error) {
	files, err := filepath.Glob(filepath.Join(testspath, "*.sh"))
	if err != nil {
		return []Testcase{}, err
	}

	testcases := []Testcase{}

	for _, file := range files {
		testcase := SHTest{
			basename: filepath.Base(file),
		}
		testcases = append(testcases, testcase)
	}

	return testcases, nil
}

func batsTests(testspath string) ([]Testcase, error) {
	files, err := filepath.Glob(filepath.Join(testspath, "*.bats"))
	if err != nil {
		return []Testcase{}, err
	}

	testcases := []Testcase{}

	for _, file := range files {
		testcase := BatsTest{
			basename: filepath.Base(file),
		}
		testcases = append(testcases, testcase)
	}

	return testcases, nil
}

type SHTest struct {
	basename string
}

func (s SHTest) Name() string {
	extension := filepath.Ext(s.basename)
	name := strings.TrimSuffix(s.basename, extension)
	return name + "-" + strings.ReplaceAll(extension, ".", "")
}

func (s SHTest) Image() string {
	return "ghcr.io/utopia-planitia/helm-tests-image:latest"
}

func (s SHTest) Command() []string {
	return []string{
		"bash",
		s.basename,
	}
}

type BatsTest struct {
	basename string
}

func (s BatsTest) Name() string {
	extension := filepath.Ext(s.basename)
	name := strings.TrimSuffix(s.basename, extension)
	return name + "-" + strings.ReplaceAll(extension, ".", "")
}

func (s BatsTest) Image() string {
	return "ghcr.io/utopia-planitia/helm-tests-image:latest"
}

func (s BatsTest) Command() []string {
	return []string{
		"bats",
		s.basename,
	}
}
