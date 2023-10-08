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

// exit show error message in dmenu.
func exit(msg string) {
	_ = exec.Command("dmenu", "-p", msg).Run()
	os.Exit(exitFailure)
}

func main() {
	var (
		data, secret, selected, passwords string
	)

	// Ask vault credentials.
	out, err := exec.Command("dmenu", "-p", "Enter vault passphrase: ", "-P").Output()
	if err != nil {
		os.Exit(exitFailure)
	}

	if secret = strings.TrimSpace(string(out)); secret == "" {
		exit("Error: Passphrase cannot empty")
	}

	// Retrieve password names.
	out, err = exec.Command("skunk", "list", "--pass", secret).Output()
	if err != nil {
		exit("Error: Passphrase is incorrect")
	}
	data = string(out)

	if !strings.Contains(data, "- ") {
		exit("Vault: No password available")
	}
	passwords = strings.ReplaceAll(data, "- ", "")

	// Select a password.
	var cmd = exec.Command("dmenu", "-p", "Select a password: ")
	cmd.Stdin = strings.NewReader(passwords)
	out, err = cmd.Output()
	if err != nil {
		os.Exit(exitFailure)
	}
	selected = strings.TrimSpace(string(out))

	// Copy into clipboard.
	err = exec.Command("skunk", "show", "-name", selected, "-pass", secret, "-copy").Run()
	if err != nil {
		exit("Error: Unable to copy password")
	}
	os.Exit(exitSuccess)
}
