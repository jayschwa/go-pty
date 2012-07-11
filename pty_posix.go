// © 2012 Jay Weisskopf

package pty

/*
#define _XOPEN_SOURCE 600
#include <stdlib.h>
#include <fcntl.h>
*/
import "C"
import "os"

// Open returns a new pair of master and slave pseudoterminal devices.
//
// When successful, the master device is ready for reading and writing; the
// slave device can be opened internally with OpenTTY, control another process
// with SetCmdTTY, or be passed off to some other part of the system.
func Open() (master *os.File, slaveName string, err error) {

	mfd, err := C.posix_openpt(C.O_RDWR | C.O_NOCTTY)
	if err != nil {
		return nil, "", err
	}
	master = os.NewFile(uintptr(mfd), "")
	defer func() {
		if err != nil && master != nil {
			master.Close()
			master = nil
		}
	}()

	_, err = C.grantpt(mfd)
	if err != nil {
		return nil, "", err
	}

	_, err = C.unlockpt(mfd)
	if err != nil {
		return nil, "", err
	}

	C_slaveName, err := C.ptsname(mfd)
	if err != nil {
		return nil, "", err
	}
	slaveName = C.GoString(C_slaveName)

	return master, slaveName, nil
}
