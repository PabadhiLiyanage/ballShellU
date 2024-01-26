package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	// Get the absolute path to the script
	scriptPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		fmt.Println("Error getting script path:", err)
		os.Exit(1)
	}

	javaCommand := "java"
	scriptPath = "/usr/lib/ballerina/bin/bal"
	scriptDir := filepath.Dir(scriptPath)

	switch osType := runtime.GOOS; osType {
	case "darwin":
		// path set
	default:
		scriptDir, err = filepath.EvalSymlinks(scriptDir)
		if err != nil {
			fmt.Println("Error resolving symlinks:", err)
			os.Exit(1)
		}
		scriptDir, err = filepath.Abs(scriptDir)
		if err != nil {
			fmt.Println("Error getting absolute path:", err)
			os.Exit(1)
		}
	}

	// Check for specific Java versions
	javaVersions := []string{"jdk-17.0.7+7-jre", "jdk-11.0.18+10-jre", "jdk-11.0.15+10-jre", "jdk-11.0.8+10-jre", "jdk8u265-b01-jre"}
	for _, version := range javaVersions {
		javaDir := filepath.Join(scriptDir, "..", "dependencies", version)
		if _, err := os.Stat(javaDir); err == nil {
			javaCommand = filepath.Join(javaDir, "bin", "java")
			break
		}
	}
	fmt.Println("ScriptPath :", scriptPath)
	fmt.Println("ScriptDir :", scriptDir)
	fmt.Println("Java Command:", javaCommand)

	//bal completion bash or bal completion zsh commands implementation
	if len(os.Args) >= 3 && os.Args[1] == "completion" {
		completionScriptPath := filepath.Join(scriptDir, "..", "scripts", "bal_completion.bash")

		if _, err := os.Stat(completionScriptPath); err == nil {
			if os.Args[2] == "zsh" {
				fmt.Printf("autoload -U +X bashcompinit && bashcompinit\n")
				fmt.Printf("autoload -U +X compinit && compinit\n\n")
				fmt.Printf("#!/usr/bin/env bash\n\n")
				printCompletionScript(completionScriptPath)
			} else if os.Args[2] == "bash" {
				fmt.Printf("#!/usr/bin/env bash\n\n")
				printCompletionScript(completionScriptPath)
			} else {
				fmt.Printf("ballerina: unknown command '%s'\n", os.Args[2])
				os.Exit(1)
			}
		} else {
			fmt.Println("Completion scripts not found")
		}
	}
	//
	runCommand := false
	runBallerina := true

	if os.Args[1] == "dist" || os.Args[1] == "update" || (os.Args[1] == "dist" && os.Args[2] == "update") {
		runCommand = true
		runBallerina = false
	}
	if os.Args[1] == "build" {
		runCommand = true
	}
}

func printCompletionScript(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Failed to generate the completion script", err)
		os.Exit(1)
	}
	fmt.Print(string(content))
}
