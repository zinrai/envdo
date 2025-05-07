package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// stripQuotes removes surrounding quotes from a string if present
func stripQuotes(s string) string {
	if len(s) < 2 {
		return s
	}

	firstChar := s[0]
	lastChar := s[len(s)-1]

	if (firstChar == '"' && lastChar == '"') ||
		(firstChar == '\'' && lastChar == '\'') {
		return s[1 : len(s)-1]
	}

	return s
}

// parseEnvLine parses a line from an env file and returns key-value pair if valid
func parseEnvLine(line string) (string, string, bool) {
	// Skip empty lines and comments
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false
	}

	// Split by first equals sign
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	// Remove quotes if present
	value = stripQuotes(value)

	return key, value, true
}

// loadEnvFile loads environment variables from the specified file
func loadEnvFile(filename string) (map[string]string, error) {
	envVars := make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key, value, valid := parseEnvLine(scanner.Text())
		if valid {
			envVars[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envVars, nil
}

// findEnvFile tries to find an env file in the current directory
func findEnvFile(specifiedFile string) string {
	if specifiedFile != "" {
		return specifiedFile
	}

	// Try to find .env or .envrc in the current directory
	if _, err := os.Stat(".env"); err == nil {
		return ".env"
	} else if _, err := os.Stat(".envrc"); err == nil {
		return ".envrc"
	}

	return ""
}

// executeCommand executes a command with given environment variables
func executeCommand(args []string, envVars map[string]string) error {
	// Create a shell command
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	// Create the command line to execute
	cmdLine := strings.Join(args, " ")

	// Use the shell to execute the command
	cmd := exec.Command(shell, "-c", cmdLine)

	// Set up environment variables
	cmd.Env = os.Environ() // Start with current environment

	// Add variables from the env file
	for key, value := range envVars {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	// Connect standard I/O
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	return cmd.Run()
}

func main() {
	// Define flags
	envFile := flag.String("f", "", "Specify an environment file to use")
	flag.Parse()

	// Get the remaining arguments (the command to execute)
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: envdo [-f envfile] command [args...]")
		os.Exit(1)
	}

	// Find the appropriate env file
	envFilePath := findEnvFile(*envFile)

	// Load environment variables if a file was found
	var envVars map[string]string
	var err error
	if envFilePath != "" {
		envVars, err = loadEnvFile(envFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading environment file %s: %v\n", envFilePath, err)
			os.Exit(1)
		}
	} else {
		envVars = make(map[string]string)
	}

	// Execute the command
	err = executeCommand(args, envVars)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		} else {
			fmt.Fprintf(os.Stderr, "Failed to execute command: %v\n", err)
			os.Exit(1)
		}
	}
}
