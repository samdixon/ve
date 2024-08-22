package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	venvPath := findVenv(currentDir)

	if venvPath != "" {
		printActivationCommand(venvPath)
	} else {
		venvPaths := findVenvsDownTree(currentDir)
		if len(venvPaths) > 0 {
			// Print paths for fzf
			for _, path := range venvPaths {
				fmt.Println(path)
			}
			// output for the corresponding bash function
			fmt.Println("__FZF_SELECTION_REQUIRED__")
		} else {
			fmt.Print("No virtual environment found. Create one? (y/n): ")
			var answer string
			fmt.Scanln(&answer)
			if strings.ToLower(answer) == "y" {
				venvPath = createVenv(currentDir)
				if venvPath != "" {
					printActivationCommand(venvPath)
				}
			}
		}
	}
}

func findVenv(dir string) string {
	venvPath := filepath.Join(dir, ".venv")
	if _, err := os.Stat(venvPath); err == nil {
		return venvPath
	}
	return ""
}

func findVenvsDownTree(dir string) []string {
	var venvPaths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".venv" {
			venvPaths = append(venvPaths, filepath.Dir(path))
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error searching for virtual environments:", err)
	}
	return venvPaths
}

func createVenv(dir string) string {
	venvPath := filepath.Join(dir, ".venv")
	cmd := exec.Command("python", "-m", "venv", venvPath)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error creating virtual environment:", err)
		return ""
	}
	fmt.Printf("Virtual environment created at %s\n", venvPath)
	return venvPath
}

func printActivationCommand(venvPath string) {
	fmt.Printf("source %s/bin/activate\n", venvPath)
}
