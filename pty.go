// © 2012 Jay Weisskopf

package pty

import "os"
import "os/exec"
import "syscall"

func Start(c *exec.Cmd) (ptm *os.File, err error) {

	// Open PTY master/slave pair
	ptm, ptsname, err := Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil && ptm != nil {
			ptm.Close()
		}
	}()

	// Open slave
	oflags := syscall.O_RDWR | syscall.O_NOCTTY
	pts, err := os.OpenFile(ptsname, oflags, 0)
	if err != nil {
		return nil, err
	}
	defer pts.Close()

	// Use pseudo-terminal slave for command's TTY
	(*c).Stdin = pts
	(*c).Stdout = pts
	(*c).Stderr = pts
	(*c).SysProcAttr = &syscall.SysProcAttr{
		Setsid:  true,
		Setctty: true,
	}

	err = (*c).Start()
	if err != nil {
		return nil, err
	}

	return ptm, nil
}
