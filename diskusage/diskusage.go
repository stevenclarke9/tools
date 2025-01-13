package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"unsafe"

	"golang.org/x/sys/windows"
	// "golang.org/x/text/language"
	// "golang.org/x/text/message"
	
	"github.com/stevenclarke9/tools/diskusage/internal/drive"
)


func DiskUsage(path string) (drive.DiskSpaceStatus, error, error) {
	h := windows.MustLoadDLL("kernel32.dll")
	c := h.MustFindProc("GetDiskFreeSpaceExW")
	lpFreeBytesAvailable := uint64(0)
	lpTotalNumberOfBytes := uint64(0)
	lpTotalNumberOfFreeBytes := uint64(0)

	// r1, r2, err := c.Call(uintptr(unsafe.Pointer(windows.StringToUTF16Ptr("C:"))),
	r1, r2, err := c.Call(uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)))

	fmt.Println("first result:", r1)
	fmt.Println("second result:", r2)
	fmt.Println("return messsge:",fmt.Sprintf("|%s|",err))
	errno := windows.GetLastError()
	fmt.Println("errno:", errno)

	if fmt.Sprintf("%s",err) != "The operation completed successfully." {
		return drive.DiskSpaceStatus{}, errno, err
	}

	return drive.DiskSpaceStatus{
		All: lpTotalNumberOfBytes,
		Free: lpTotalNumberOfFreeBytes,
		Used: lpTotalNumberOfBytes - lpFreeBytesAvailable,
	}, nil, nil
}

func main() {
	driveFlag := flag.String("d", "", "diskusage for drive letter")
	flag.Parse()

	diskdrive := "C:"
	if len(*driveFlag) == 1 {
		if drives := drive.GetAllDrives(); slices.Contains(drives,*driveFlag) {
			diskdrive = *driveFlag + ":"
		} else {
			fmt.Println("Available drives: ", drives)
			os.Exit(1)
		}
	}
	fmt.Println("disk space status for Drive ", diskdrive)
	
	d, err := drive.GetDiskSpaceStatus(diskdrive)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println(d)
	}
	/*
	ptr, convertError := windows.UTF16PtrFromString(directoryName)
	if convertError == nil {
		err := windows.GetDiskFreeSpaceEx(ptr,&freeBytesAvailableToCaller,&totalNumberOfBytes,&totalNumberOfFreeBytes)
		if err != nil {
			fmt.Println("windows.GetDiskFreeSpaceEx error:", err)
			dss, windowsGetLastError, callError := DiskUsage(diskdrive)
			if callError != nil {
				fmt.Println("callError: ", callError, "windowsGetLastError: ", windowsGetLastError)
			} else {
				fmt.Println(dss)
			}
		} else {
			d := DiskSpaceStatus{
				All: totalNumberOfBytes,
				Free: totalNumberOfFreeBytes,
				Used: totalNumberOfBytes - totalNumberOfFreeBytes,
			}
			fmt.Println("windows.GetDiskFreeSpaceEx err is nil")
			//			
			// fmt.Println("freeBytesAvailableToCaller: ",freeBytesAvailableToCaller)
			// fmt.Println("totalNumberOfBytes: ",totalNumberOfBytes)
			// fmt.Println("totalNumberOfFreeBytes: ",totalNumberOfFreeBytes)
			//
			fmt.Println(d)
		}
	} else {
		// there is a convertError value returned
		fmt.Println("convertError:", convertError)
	}
	*/
}
