package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func removeDupes(sarray []string) ([]string) {
	stringmap := make(map[string]bool)
	out := []string{}

	for _, path := range sarray {
		// check that the 'path' value starts with a '/' and path[0:1] is not "C:". If not prepend a '/' to the 'path' value.
		// fmt.Println(path[0:2])
		if path[0] != '/' && ( path[0:2] != "C:" ) {
			path = "/" + path
		}
		// check if we have seen the path before.
		if _, ok := stringmap[path]; !ok {
			stringmap[path] = true
			out = append(out,path)
		}
	}
	return out
}

func main() {

	verboseFlag := flag.Bool("v", false, "verbose messages")

	// parse the command line options.
	flag.Parse()

	// get the MSYSTEM value from the environment.	
	var msystem string = os.Getenv("MSYSTEM")
	
	// get the path value from the environment.
	var envpath string = os.Getenv("PATH")
	
	// set the path seperator character to the sep variable.
	var sep string = ";"
	
	// split the string into a slice of strings. The string is seperated by a ";".
	envpatharray := strings.Split(envpath,sep)

	if (*verboseFlag == true) {
		fmt.Println("The MSYSTEM value is",msystem)
		fmt.Println("the raw PATH separater is ",sep)
		fmt.Printf("%s\n%s\n","The value of the raw PATH environment variable is",envpath)
		fmt.Println("The list of elements in the PATH is")
		for i, path := range envpatharray {
			fmt.Println(i,path)
		}
	}
	
	// remove the duplicates from the slice of strings
	noDupes := removeDupes(envpatharray)
	
	if (*verboseFlag == true) {
		fmt.Println("\nThe list of elements in the PATH with the duplicates removed is")
		for i, path := range noDupes {
			fmt.Println(i,path)
		}
		// print an empty line after the elements of the noDupes slice are printed.
		fmt.Println("")
	}	
	// rejoin the split path elements. The slice elements are joined together by a ";".
	pathWithNoDupes := strings.Join(noDupes,sep)
	if msystem == "" {
		fmt.Println(pathWithNoDupes)
	} else {
		// replace the "C:" with "/c"
		re := strings.NewReplacer("\\","/")
		pathInMsys2Format := re.Replace(pathWithNoDupes)
		fmt.Println(pathInMsys2Format)
	}
	os.Exit(0)
}
