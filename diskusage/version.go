package main

import (
	"fmt"
	"runtime/debug"
)

// printVersion prints the application version
func printVersion() {
	buildinfo, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("Unable to determine version information.")
		return
	}
	if buildinfo.Main.Version != "" {
		fmt.Printf("Version: %s\n", buildinfo.Main.Version)
	} else {
		fmt.Println("Version: unknown")
	}
}
