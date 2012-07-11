// © 2012 Jay Weisskopf

// Package pty implements utility routines for psuedoterminal devices.
package pty

import "os"
import "os/exec"
import "syscall"

// SetCmdTTY attaches a command's standard streams to the named terminal device
// and returns nil if successful.
// If a terminal device is not named (empty string), a new one will be opened
// and the corresponding master device will be returned.
func SetCmdTTY(cmd *exec.Cmd, ttyName string) (master *os.File, err error) {

	if ttyName == "" {
		master, ttyName, err = Open()
		if err != nil {
			return nil, err
		}
		defer func() {
			if err != nil && master != nil {
				master.Close()
			}
		}()
	}

	tty, err := OpenTTY(ttyName)
	if err != nil {
		return nil, err
	}
	(*cmd).Stdin = tty
	(*cmd).Stdout = tty
	(*cmd).Stderr = tty
	(*cmd).SysProcAttr = &syscall.SysProcAttr{
		Setsid:  true,
		Setctty: true,
	}
	// tty.Close() is not necessary because Cmd.Start handles it.

	return master, nil
}

// OpenTTY opens the named file for reading and writing as a terminal device.
func OpenTTY(ttyName string) (tty *os.File, err error) {

	flags := syscall.O_RDWR | syscall.O_NOCTTY
	tty, err = os.OpenFile(ttyName, flags, 0)
	if err != nil {
		return nil, err
	}

	return tty, nil
}
