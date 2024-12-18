package main // require main package name

import "os"

func main() { // check only calls in main function
	os.Exit(0)   // want "use return instead of stop process with code 0"
	os.Exit(1)   // want "use panic instead of stop process with code 1"
	os.Exit(125) // want "use panic instead of stop process with code 125"

	// nested calls
	f1()
	f2()
}

func f1() {
	os.Exit(0) // want "use return instead of stop process with code 0"
	os.Exit(1) // want "use panic instead of stop process with code 1"
}

func f2() {
	os.Exit(0) // want "use return instead of stop process with code 0"
	os.Exit(1) // want "use panic instead of stop process with code 1"
	f3()
}

func f3() {
	os.Exit(0) // want "use return instead of stop process with code 0"
	os.Exit(1) // want "use panic instead of stop process with code 1"
}

func unreachable() { // no usage in main function
	os.Exit(0)
	os.Exit(1)
}
