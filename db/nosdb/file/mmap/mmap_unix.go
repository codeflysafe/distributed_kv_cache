package mmap

import (
	"golang.org/x/sys/unix"
)

// See https://github.com/edsrzf/mmap-go/blob/37b9a46b37a6638c1da75bfbceba03d066ff777e/mmap_unix.go#L13
func mmap(len int, inProt, fd int, off int64) ([]byte, error) {

	flags := unix.MAP_SHARED
	prot := unix.PROT_READ
	switch {
	case inProt&COPY != 0:
		prot |= unix.PROT_WRITE
		flags = unix.MAP_PRIVATE
	case inProt&RDWR != 0:
		prot |= unix.PROT_WRITE
	}
	if inProt&EXEC != 0 {
		prot |= unix.PROT_EXEC
	}
	b, err := unix.Mmap(fd, off, len, prot, flags)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func msync(b []byte) error {
	return unix.Msync(b, unix.MS_SYNC)
}

func munmap(b []byte) error {
	return unix.Munmap(b)
}
