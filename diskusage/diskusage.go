package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	//"time"
	"unsafe"

	"golang.org/x/sys/windows"
	// "golang.org/x/text/language"
	// "golang.org/x/text/message"
	
	"github.com/stevenclarke9/tools/diskusage/internal/drive"
	"github.com/stevenclarke9/tools/diskusage/internal/file"
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
//	now := time.Now()
//	filetime := file.FormatTime24Hour(now)

	outputFile := flag.String("o", "", "write diskusage to a file")

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

	var filePtr *os.File
	var fileErr error

	if (outputFile != nil) && (*outputFile != "") {
		//fmt.Println("outputFile:", *outputFile)
		filePtr, fileErr = file.Create(*outputFile)
		if fileErr != nil {
			fmt.Println("file create error:", fileErr)
			os.Exit(2)
		}
	}
	
	d, err := drive.GetDiskSpaceStatus(diskdrive)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		if filePtr == nil {
			fmt.Println("disk space status for Drive ", diskdrive)
			fmt.Println(d)
		} else {
			filePtr.WriteString(fmt.Sprintf("disk space status for Drive %s\n", diskdrive))
			filePtr.WriteString(fmt.Sprint(d))
			filePtr.Close()
		}
	}
}
