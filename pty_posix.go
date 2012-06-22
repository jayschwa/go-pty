// © 2012 Jay Weisskopf

package pty

// #include <stdlib.h>
// #include <fcntl.h>
import "C"

import "os"

func Open() (ptm *os.File, ptsName string, err error) {

	ptmFd, err := C.posix_openpt(C.O_RDWR | C.O_NOCTTY)
	if err != nil {
		return nil, "", err
	}
	ptm = os.NewFile(uintptr(ptmFd), "")
	defer func() {
		if err != nil && ptm != nil {
			ptm.Close()
			ptm = nil
		}
	}()

	_, err = C.grantpt(ptmFd)
	if err != nil {
		return nil, "", err
	}

	_, err = C.unlockpt(ptmFd)
	if err != nil {
		return nil, "", err
	}

	ptsNameCstr, err := C.ptsname(ptmFd)
	if err != nil {
		return nil, "", err
	}
	ptsName = C.GoString(ptsNameCstr)

	return ptm, ptsName, nil
}
