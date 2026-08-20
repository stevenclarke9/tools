// create a string of the format "yyyymmdd-hhmm" for the time input paameter.
package file

import (
	// "fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	YearFirstDateLayout = "20060102"
	HHMMTimeLayout = "1504"
	FileDateTimeLayout = "20060102-1504"
)

func PrintYearFirstDate(t time.Time) string {
	return t.Format(YearFirstDateLayout)
}

func PrintHHMMTime(t time.Time) string {
	return t.Format(HHMMTimeLayout)
}


func FormatTime24Hour(t time.Time) (string) {

//	year := t.Year()
//	month := t.Month()
//	day := t.Day()
//	hour24 := 10
//	minute := 10

//	layout := "20060102-1504"
//	s := t.Format(layout)
//	return s
	
	return t.Format(FileDateTimeLayout)
}

func Create(suffix string) (*os.File, error) {
	filename := FormatTime24Hour(time.Now()) + "-" + suffix
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fileNamePath := filepath.Join(pwd,filename)
	filePtr, err := os.Create(fileNamePath)
	return filePtr, err
}