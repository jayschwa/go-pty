// © 2012 Jay Weisskopf

package pty

import "fmt"
import "io"
import "os"
import "os/exec"
import "strings"
import "testing"
import "time"

func TestOpen(t *testing.T) {
	master, slave, err := Open()
	if err != nil {
		t.Error(err)
	}
	if master == nil {
		t.Error("No master device")
	}
	t.Logf("Slave device path: '%s'\n", slave)
	if !strings.HasPrefix(slave, "/dev/") {
		t.Fail()
	}
}

func TestFork(t *testing.T) {
	cmd := exec.Command("sh")
	ptm, err := SetCmdTTY(cmd, "")
	if err != nil {
		t.Fatal(err)
	}
	cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	go io.Copy(os.Stdout, ptm)
	time.Sleep(100 * time.Millisecond)
	fmt.Fprintln(ptm, "ps T")
	time.Sleep(100 * time.Millisecond)
	fmt.Fprintln(ptm, "tty")
	time.Sleep(100 * time.Millisecond)
	fmt.Fprintln(ptm, "who")
	time.Sleep(100 * time.Millisecond)
	fmt.Fprintln(ptm, "exit")
	cmd.Wait()
}
