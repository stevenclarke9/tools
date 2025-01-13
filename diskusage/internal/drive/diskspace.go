package drive

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/sys/windows"
)

// DiskSpaceStatus contains the size of the space on the requested disk drive.
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
	// var usedPercentage float64
	usedPercentage := float64((float64(d.Used) / float64(d.All)) * 100)
	fmt.Printf("usedPercentage: %.2f\n", usedPercentage)
	freePercentage := 100 - usedPercentage
	
	return p.Sprintf("All Space %d bytes %d%%\nUsed Space %d bytes %.2f%%\nFree Space %d bytes %.2f%%",
		d.All, 100,
		d.Used, usedPercentage,
		d.Free, freePercentage)
}

func GetDiskSpaceStatus(diskPath string) (DiskSpaceStatus, error) {

	// test windows package function first:
	//var directoryName string = diskdrive
	var freeBytesAvailableToCaller uint64
	var totalNumberOfBytes uint64
	var totalNumberOfFreeBytes uint64


	ptr, convertError := windows.UTF16PtrFromString(diskPath)
	if convertError != nil {
		// there is a convertError value returned
		fmt.Println("convertError:", convertError)
		return DiskSpaceStatus{}, convertError
	}
	err := windows.GetDiskFreeSpaceEx(ptr,&freeBytesAvailableToCaller,&totalNumberOfBytes,&totalNumberOfFreeBytes)
	if err != nil {
		fmt.Println("windows.GetDiskFreeSpaceEx error:", err)
		return DiskSpaceStatus{}, err
	}
	d := DiskSpaceStatus{
		All: totalNumberOfBytes,
		Free: totalNumberOfFreeBytes,
		Used: totalNumberOfBytes - totalNumberOfFreeBytes,
	}
	return d, nil
}