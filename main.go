package main

import (
	"os"
	"os/exec"
	"strings"
)

const (
	exitSuccess = iota
	exitFailure
)

// dmenuErorr show error message in dmenu.
func dmenuErorr(err string) {
	var cmd = exec.Command("dmenu", "-p", "error: "+err)
	cmd.Run()
	os.Exit(exitFailure)
}

// trimSpaces trim spaces given value.
func trimSpaces(value []byte) string {
	return strings.TrimSpace(string(value))
}

func main() {
	var (
		exitEnter          = "press enter to exit"
		invalidInputError  = "password cannot empty. " + exitEnter
		inccorectPassError = "vault: password incorrect. " + exitEnter
		emptyVaultError    = "vault: no password available"
	)

	// Ask vault credentials in dmenu.
	var cmd = exec.Command("dmenu", "-p", "Enter vault password: ", "-P")
	out, err := cmd.Output()
	if err != nil {
		dmenuErorr("input: " + err.Error())
	}

	var secret = trimSpaces(out)
	if secret == "" {
		dmenuErorr(invalidInputError)
	}

	// Get available passwords.
	out, err = exec.Command("skunk", "list", "--pass", secret).Output()
	if err != nil {
		dmenuErorr(inccorectPassError)
	}
	var passwords = strings.ReplaceAll(trimSpaces(out), "- ", "")

	if passwords == emptyVaultError {
		dmenuErorr(emptyVaultError + ". " + exitEnter)
	}

	// Ask to select password in vault.
	cmd = exec.Command("dmenu", "-p", "Select a password: ")
	cmd.Stdin = strings.NewReader(passwords)
	out, err = cmd.Output()
	if err != nil {
		dmenuErorr("select: " + err.Error())
	}
	var selectedPass = trimSpaces(out)

	// Copy into clipboard.
	cmd = exec.Command("skunk", "show", "--name", selectedPass, "--pass", secret, "--copy")
	if err = cmd.Run(); err != nil {
		dmenuErorr("copy: " + err.Error())
	}
	os.Exit(exitSuccess)
}
