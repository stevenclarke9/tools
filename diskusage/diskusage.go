package main

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// diskSpaceStatus contains the size of the space on the requested disk drive.
// All is the total space size in bytes of all the drive.
// Used is the used space in bytes on the drive.
// Free is the free space in bytes on the drive
type DiskSpaceStatus struct {
	All  uint64
	Used uint64
	Free uint64
}

func (d DiskSpaceStatus) String() string {
	p := message.NewPrinter(language.English)
	// return fmt.Sprintf("All Space %d\nUsed Space %d\nFreeSpace %d\n", d.All, d.Used, d.Free)
	return p.Sprintf("All Space %d bytes\nUsed Space %d bytes\nFreeSpace %d bytes", d.All, d.Used, d.Free)
}

func DiskUsage(path string) (DiskSpaceStatus, error, error) {
	h := windows.MustLoadDLL("kernel32.dll")
	c := h.MustFindProc("GetDiskFreeSpaceExW")
	lpFreeBytesAvailable := uint64(0)
	lpTotalNumberOfBytes := uint64(0)
	lpTotalNumberOfFreeBytes := uint64(0)

	// r1, r2, err := c.Call(uintptr(unsafe.Pointer(windows.StringToUTF16Ptr("C:"))),
	_, _, err := c.Call(uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)))
	fmt.Println(fmt.Sprintf("%s",err))
	if fmt.Sprintf("%s",err) != "The operation completed successfully." {
		errno := windows.GetLastError()
		return DiskSpaceStatus{}, errno, err
	}

	return DiskSpaceStatus{
		All: lpTotalNumberOfBytes,
		Free: lpTotalNumberOfFreeBytes,
		Used: lpTotalNumberOfBytes - lpFreeBytesAvailable,
	}, nil, nil
}

func main() {
	drive := "C:"
	fmt.Println("disk space status for Drive ", drive)
	dss, windowsGetLastError, callError := DiskUsage(drive)
	if windowsGetLastError != nil {
		fmt.Println("callError: ", callError, "windowsGetLastError: ", windowsGetLastError)
	} else {
		fmt.Println(dss)
	}
}
