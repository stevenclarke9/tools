package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	
	duw "github.com/stevenclarke9/tools/diskusage/drives/windows"
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
	var usedPercentage float64
	usedPercentage = float64((float64(d.Used) / float64(d.All)) * 100)
	fmt.Printf("usedPercentage: %.2f\n", usedPercentage)
	freePercentage := 100 - usedPercentage
	
	return p.Sprintf("All Space %d bytes %d%%\nUsed Space %d bytes %.2f%%\nFree Space %d bytes %.2f%%",
		d.All, 100,
		d.Used, usedPercentage,
		d.Free, freePercentage)
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
	driveFlag := flag.String("d", "", "diskusage for drive letter")
	flag.Parse()

	drive := "C:"
	fmt.Println("len(*driveFlag) = ", len(*driveFlag))
	if len(*driveFlag) == 1 {
		if drives := duw.GetAllDrives(); slices.Contains(drives,*driveFlag) {
			drive = *driveFlag + ":"
		} else {
			fmt.Println("Available drives: ", drives)
			os.Exit(1)
		}
	}
	fmt.Println("disk space status for Drive ", drive)
	dss, windowsGetLastError, callError := DiskUsage(drive)
	if windowsGetLastError != nil {
		fmt.Println("callError: ", callError, "windowsGetLastError: ", windowsGetLastError)
	} else {
		fmt.Println(dss)
	}
}
