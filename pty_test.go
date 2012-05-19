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
	cmd := exec.Command("login", "-f", "jay")
	ptm, err := Start(cmd)
	if err != nil {
		t.Fatal(err)
	}
	go io.Copy(os.Stdout, ptm)
	time.Sleep(100 * time.Millisecond) // Avoid premature echo
	fmt.Fprintln(ptm, "ps T")
	fmt.Fprintln(ptm, "tty")
	fmt.Fprintln(ptm, "who")
	fmt.Fprintln(ptm, "exit")
	cmd.Wait()
}
