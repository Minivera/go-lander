//go:build mage
// +build mage

package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/magefile/mage/mg"
)

// Updates the JavaScript glue code inside the example directory using the code provided in the $GOROOT.
// Will look up the JavaScript glue code needed to properly execute a WASM script in the `/misc/wasm/go_js_wasm_exec`
// directory of the $GOROOT and copy it to `./example` as `index.js`.
func UpdateGlue() error {
	goroot := os.Getenv("GOROOT")
	if goroot == "" {
		return fmt.Errorf("missing $GOROOT environement variable, it was either empty or missing")
	}

	filepath := fmt.Sprintf("%s/misc/wasm/wasm_exec.js", goroot)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("JavaScript glue code file was missing from GOROOT, make sure %s exists", filepath)
	}

	fmt.Printf("Copying %s to `./example/index.js`\n", filepath)

	cmd := exec.Command("cp", filepath, "./example/index.js")
	return cmd.Run()
}

// Runs the provided example by building it to WASM and running the serve tool in the example directory.
// RunExample will look up the example as a subdirectory in the example directory and build all go files it finds
// there to main.wasm in the example dir. It will then run the main command of serve.go to run the example
// as a web server on port 8080.
func RunExample(example string) error {
	mg.Deps(InstallDeps, mg.F(BuildExample, example))

	fmt.Printf("Executing the webserver and serving %s\n", example)

	exampleDir := fmt.Sprintf("example/%s", example)
	if _, err := os.Stat(exampleDir); os.IsNotExist(err) {
		return fmt.Errorf("%s is not a valid example, make sure a directory exists with the same name under `./example` and that it has a main function defined", example)
	}

	var errOutput bytes.Buffer
	cmd := exec.Command("go", "run", "./example/serve.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = &errOutput

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run build command %s, %s. Error was %w", cmd.String(), errOutput.String(), err)
	}

	return nil
}

// Builds the provided example as a wasm file at the root of the example directory.
// BuildExample will look up the example as a subdirectory in the example directory and build all go files
// it finds there to main.wasm at the root of the example directory.
func BuildExample(example string) error {
	mg.Deps(InstallDeps, UpdateGlue)

	fmt.Printf("Building %s to main.wasm\n", example)

	exampleDir := fmt.Sprintf("./example/%s", example)
	if _, err := os.Stat(exampleDir); os.IsNotExist(err) {
		return fmt.Errorf("%s is not a valid example, make sure a directory exists with the same name under `./example` and that it has a main function defined")
	}

	var errOutput bytes.Buffer
	cmd := exec.Command("go", "build", "-o", "example/main.wasm", exampleDir)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOOS=js", "GOARCH=wasm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = &errOutput

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run build command %s, %s. Error was %w", cmd.String(), errOutput.String(), err)
	}

	return nil
}

type exampleData struct {
	Examples []struct {
		Path   string
		Active bool
		Name   string
	}
	Path string
	Name string
	Url  string
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSentenceCase(str string) string {
	sentence := matchFirstCap.ReplaceAllString(str, "${1} ${2}")
	sentence = matchAllCap.ReplaceAllString(sentence, "${1} ${2}")

	return strings.ToTitle(strings.ToLower(sentence))
}

// Builds the all the examples from the example directory into a special viewer.
// BuildExampleViewer will look up the examples as subdirectories in the example directory
// and build all go files it finds there to the build directory of the `example/.exampleSwitcher`
// directory. It will generate the necessary code for the switcher to work as intended.
func BuildExampleViewer() error {
	mg.Deps(InstallDeps, UpdateGlue)

	fmt.Println("Creating the build directory for the example viewer")
	if err := os.RemoveAll("./example/.exampleSwitcher/build"); err != nil {
		return err
	}

	if err := os.MkdirAll("./example/.exampleSwitcher/build", os.ModePerm); err != nil {
		return err
	}

	rootTemplate, err := template.ParseFiles("./example/.exampleSwitcher/root.html.tpl")
	if err != nil {
		return fmt.Errorf("failed to read index.html remplate. Error was %w", err)
	}

	indexTemplate, err := template.ParseFiles("./example/.exampleSwitcher/index.html.tpl")
	if err != nil {
		return fmt.Errorf("failed to read example index.html remplate. Error was %w", err)
	}

	fmt.Printf("Listing all directories of the example directory, except for the viewer\n")
	examples, _ := os.ReadDir("./example")

	var exampleDefs []exampleData
	for _, dir := range examples {
		if dir.IsDir() && dir.Name() != ".exampleSwitcher" {
			example := dir.Name()

			exampleDefs = append(exampleDefs, exampleData{
				Examples: []struct {
					Path   string
					Active bool
					Name   string
				}{},
				Path: example,
				Name: toSentenceCase(example),
				Url:  fmt.Sprintf("https://github.com/Minivera/go-lander/tree/main/example/%s/main.go", example),
			})
		}
	}

	fmt.Printf("Executing all found examples and creating their build files\n")
	for _, exampleDef := range exampleDefs {
		example := exampleDef.Path
		exampleBuildDir := fmt.Sprintf("./example/.exampleSwitcher/build/%s", example)
		if err := os.MkdirAll(exampleBuildDir, os.ModePerm); err != nil {
			return err
		}

		fmt.Printf("Copying glue to example build directory %s\n", exampleBuildDir)
		cmd := exec.Command("cp", "./example/index.js", exampleBuildDir+"/index.js")
		if err := cmd.Run(); err != nil {
			return err
		}

		fmt.Printf("Building %s to main.wasm in the build directory %s\n", example, exampleBuildDir)
		exampleDir := fmt.Sprintf("./example/%s", example)
		if _, err := os.Stat(exampleDir); os.IsNotExist(err) {
			return fmt.Errorf("%s is not a valid example, make sure a directory exists with the same name under `./example` and that it has a main function defined", example)
		}

		var errOutput bytes.Buffer
		cmd = exec.Command("go", "build", "-o", exampleBuildDir+"/main.wasm", exampleDir)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "GOOS=js", "GOARCH=wasm")
		cmd.Stdout = os.Stdout
		cmd.Stderr = &errOutput

		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to run build command %s, %s. Error was %w", cmd.String(), errOutput.String(), err)
		}

		fmt.Printf("Writing index.html in the build directory %s\n", exampleBuildDir)
		indexFile, err := os.Create(exampleBuildDir + "/index.html")
		if err != nil {
			return fmt.Errorf("failed to create index file for example %s build. Error was %w", example, err)
		}

		for _, subDef := range exampleDefs {
			exampleDef.Examples = append(exampleDef.Examples, struct {
				Path   string
				Active bool
				Name   string
			}{
				Path:   subDef.Path,
				Active: subDef.Path == exampleDef.Path,
				Name:   subDef.Name,
			})
		}

		err = indexTemplate.Execute(indexFile, exampleDef)
		if err != nil {
			return fmt.Errorf("failed to create index file for example %s build. Error was %w", example, err)
		}

		err = indexFile.Close()
		if err != nil {
			return err
		}
	}

	fmt.Printf("Writing root index.html in the build directory\n")
	rootFile, err := os.Create("./example/.exampleSwitcher/build/index.html")
	if err != nil {
		return fmt.Errorf("failed to create index file for root. Error was %w", err)
	}

	rootDef := exampleData{
		Examples: []struct {
			Path   string
			Active bool
			Name   string
		}{},
	}
	for _, subDef := range exampleDefs {
		rootDef.Examples = append(rootDef.Examples, struct {
			Path   string
			Active bool
			Name   string
		}{
			Path:   subDef.Path,
			Active: subDef.Path == rootDef.Path,
			Name:   subDef.Name,
		})
	}

	err = rootTemplate.Execute(rootFile, rootDef)
	if err != nil {
		return fmt.Errorf("failed to create root index file for build. Error was %w", err)
	}

	return rootFile.Close()
}

// Runs the examples by building them to WASM in a common viewer directory and then serving them as HTML files.
// RunExampleViewer will look up the examples as subdirectories in the example directory and build all go
// example it finds there to their respective index.html and main.wasm files. It will then run the main
// command of serve.go to run the .exampleSwitcher directory as a web server on port 8080.
func RunExampleViewer() error {
	mg.Deps(InstallDeps, BuildExampleViewer)

	fmt.Println("Executing the webserver and serving")

	var errOutput bytes.Buffer
	cmd := exec.Command("go", "run", "./example/.exampleSwitcher/serve.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = &errOutput

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run build command %s, %s. Error was %w", cmd.String(), errOutput.String(), err)
	}

	return nil
}

// Install all dependencies of the project in one command
func InstallDeps() error {
	fmt.Println("Installing all dependencies")

	cmd := exec.Command("go", "get", "./...")
	return cmd.Run()
}

// Cleans the WASm and build artifacts from the repository
func Clean() {
	fmt.Println("Cleaning WASM and build artifacts")

	os.RemoveAll("vendor")
	os.Remove("example/main.wasm")
}
