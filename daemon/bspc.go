package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/nullishamy/bspwm-scratchpad/transport"
)

func runCommand(command ...string) (string, error) {
	cmd := exec.Command("bspc", command...)

	output, err := cmd.Output()
	if err != nil {
		casted := err.(*exec.ExitError)
		stderr := string(casted.Stderr)
		return "", errors.New("error running command 'bspc " + strings.Join(command, " ") + "':\n" + stderr)
	}

	return strings.TrimSuffix(string(output), "\n"), nil
}

func GetCurrentWindowId() (int64, error) {
	idStr, err := runCommand("query", "-N", "-n", "@focused:/#focused")
	if err != nil {
		return 0, err
	}

	id, err := strconv.ParseInt(idStr, 0, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetWindowDetails(id int64) (*transport.Window, error) {
	jsonStr, err := runCommand("query", "-T", "-n", fmt.Sprint(id))

	if err != nil {
		return nil, err
	}

	var window transport.Window
	err = json.Unmarshal([]byte(jsonStr), &window)
	if err != nil {
		return nil, err
	}

	return &window, nil
}

func GetCurrentWindow() (*transport.Window, error) {
	id, err := GetCurrentWindowId()
	if err != nil {
		return nil, err
	}

	return GetWindowDetails(id)
}

func HideWindow(window *transport.Window) error {
	_, err := runCommand("node", fmt.Sprint(window.ID), "--flag", "hidden=on")
	return err
}

func ShowWindow(window *transport.Window) error {
	_, err := runCommand("node", fmt.Sprint(window.ID), "--flag", "hidden=off")
	return err
}

func ToggleWindow(window *transport.Window) error {
	if window.Hidden {
		return ShowWindow(window)
	}

	return HideWindow(window)
}
