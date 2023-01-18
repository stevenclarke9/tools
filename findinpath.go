package main

import (
	//"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	verboseFlag *bool = flag.Bool("v", false, "verbose messages")
	quietFlag *bool = flag.Bool("q", false, "quiet mode")
)

func findFileinPath(filename string, pathlist []string) (string, bool, error) {
	var err error
	found := false
	foundPath := ""
	fullpath := ""
	for _, path := range pathlist {
		if path == "" {
			continue
		}
		pathlen := len(path)
		if path[0] == '%' {
			// get the environment value between the pair of %.
			envname := ""
			for j := 1; path[j] != '%' && j < pathlen; j++ {
				envname = envname + string(path[j])
			}
			newpath := os.Getenv(envname)
			if newpath == "" {
				continue
			}
			fmt.Println("the env name",envname,"has the value",newpath,".")
			path = newpath
		}
		if path[pathlen-1] == '\\' {
			fullpath = path + filename
		} else {
			fullpath = path + "\\" + filename
		}
		if _, err = os.Lstat(fullpath); err == nil {
			if (*verboseFlag) {
				fmt.Println("the file",fullpath,"is found in a directory within the PATH environment variable.")
			}
			found = true
			foundPath = fullpath
			break
		} else {
			if os.IsNotExist(err) {
				if (*verboseFlag) {
					fmt.Println("the file", fullpath, "does not exist on the file system.")
				}
				err = nil
			//}
			//var e *os.PathError
			//if errors.As(err, &e) {
			//	err = nil
			} else {
				return "",false,err
			}
		}
	}
	return foundPath,found,err
}

func main() {
	// parse the command line options.
	flag.Parse()

	arg := flag.Args()
	
	// get the MSYSTEM value from the environment.	
	msystem := os.Getenv("MSYSTEM")
	
	// get the path value from the environment.
	envpath := os.Getenv("PATH")
	
	// set the path seperator character to the sep variable.
	sep := ";"
	
	// split the string into a slice of strings. The string is seperated by a ";".
	envpatharray := strings.Split(envpath,sep)

	if (*verboseFlag) {
		fmt.Println("The MSYSTEM value is",msystem)
		fmt.Println("the raw PATH separater is ",sep)
		fmt.Printf("%s\n%s\n","The value of the raw PATH environment variable is",envpath)
		fmt.Println("The list of elements in the PATH is")
		for i, path := range envpatharray {
			fmt.Println(i,path)
		}
		fmt.Println("args is",arg)
	}
	
	if len(arg) == 1 {
		fullpath := ""
		found := false
		var err error
		fullpath, found, err = findFileinPath(arg[0],envpatharray)
		if (err == nil) && (found == true) {
			if (!*quietFlag) {
				fmt.Println("the full path to the file is",fullpath)
			}
			os.Exit(0)
		} else {
			if err != nil {
				if (!*quietFlag) {
					fmt.Println("ERROR:",err)
				}
			} else {
				if (!*quietFlag) {
					fmt.Println("the file",arg[0],"is not found in the system PATH environment variable")
				}
			}
			os.Exit(1)
		}
	} else {
		if (!*quietFlag) {
			fmt.Println("one filename only")
		}
		os.Exit(1)
	}

}
