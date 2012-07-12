// © 2012 Jay Weisskopf

package pty

import "bufio"
import "fmt"
import "os/exec"
import "strings"
import "testing"

// Verify Open returns sane values.
func TestOpen(t *testing.T) {

	master, slaveName, err := Open()
	if err != nil {
		t.Error("Open returned an error")
		t.Fatal(err)
	}
	if master == nil {
		t.Error("Open returned a nil master device")
	} else {
		defer master.Close()
	}
	if slaveName == "" {
		t.Error("Open returned an empty slave device path")
	} else {
		fmt.Printf("Slave device path: '%s'\n", slaveName)
		tty, err := OpenTTY(slaveName)
		if err != nil {
			t.Error("OpenTTY returned an error")
			t.Fatal(err)
		}
		if tty == nil {
			t.Error("OpenTTY returned a nil slave device")
		} else {
			defer tty.Close()
		}
	}
}

// Verify OpenTTY returns an error and nil when given bad paths.
func TestOpenTTY_BadPaths(t *testing.T) {

	tty, err := OpenTTY("")
	if err == nil {
		t.Error("OpenTTY did not return an error when given an empty path")
	}
	if tty != nil {
		t.Error("OpenTTY returned a non-nil device when given an empty path")
	}

	tty, err = OpenTTY("/some/bull/shit/path")
	if err == nil {
		t.Error("OpenTTY did not return an error when given a fake path")
	}
	if tty != nil {
		t.Error("OpenTTY returned a non-nil device when given a fake path")
	}

	tty, err = OpenTTY("/bin/cat")
	if err == nil {
		t.Error("OpenTTY did not return an error when given a bad path")
	}
	if tty != nil {
		t.Error("OpenTTY returned a non-nil device when given a bad path")
	}
}

// Test data flow between a couple pseudoterminal device pairs.
func TestPty_DataFlow(t *testing.T) {

}

// Verify that concurrent commands have different terminal devices.
func TestSetCmdTTY_DiffDevs(t *testing.T) {

	thisTTY := "FIXME"

	thatCmd := exec.Command("tty")
	otherCmd := exec.Command("tty")
	thatMaster, err := SetCmdTTY(thatCmd, "")
	if err != nil {
		t.Fatal(err)
	}
	if thatMaster == nil {
		t.Fatal("SetCmdTTY returned a nil master device")
	} else {
		defer thatMaster.Close()
	}
	otherMaster, err := SetCmdTTY(otherCmd, "")
	if err != nil {
		t.Fatal(err)
	}
	if otherMaster == nil {
		t.Fatal("SetCmdTTY returned a nil master device")
	} else {
		defer otherMaster.Close()
	}
	err = thatCmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	err = otherCmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	err = thatCmd.Wait()
	if err != nil {
		t.Fatal(err)
	}
	err = otherCmd.Wait()
	if err != nil {
		t.Fatal(err)
	}

	thatTTY, err := bufio.NewReader(thatMaster).ReadString(10)
	if err != nil {
		t.Fatal(err)
	}
	thatTTY = strings.TrimSpace(thatTTY)
	if thatTTY == "" {
		t.Error("That TTY name is empty")
	} else {
		fmt.Printf("That TTY: '%s'\n", thatTTY)
	}
	otherTTY, err := bufio.NewReader(otherMaster).ReadString(10)
	if err != nil {
		t.Fatal(err)
	}
	otherTTY = strings.TrimSpace(otherTTY)
	if otherTTY == "" {
		t.Error("Other TTY name is empty")
	} else {
		fmt.Printf("Other TTY: '%s'\n", otherTTY)
	}

	if thisTTY == thatTTY {
		t.Error("This TTY and That TTY are the same")
	}
	if thisTTY == otherTTY {
		t.Error("This TTY and the Other TTY are the same")
	}
	if thatTTY == otherTTY {
		t.Error("That TTY and the Other TTY are the same")
	}
}

// Verify that the slave device File is closed after a command has been started.

// Verify that SetCmdTTY returns an error when given bad paths.
// Additionally, verify that TTY allocation "leaks" do not occur.

