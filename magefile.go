//go:build mage
// +build mage

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

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
		return fmt.Errorf("%s is not a valid example, make sure a directory exists with the same name under `./example` and that it has a main function defined")
	}

	var errOutput bytes.Buffer
	cmd := exec.Command("go", "run", "./example/serve.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = &errOutput

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Failed to run build command %s, %s. Error was %w.", cmd.String(), errOutput.String(), err)
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
		return fmt.Errorf("Failed to run build command %s, %s. Error was %w.", cmd.String(), errOutput.String(), err)
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
