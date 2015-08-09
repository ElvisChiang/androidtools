package main

import (
	"flag"
	"fmt"
	"os"
)

func usageAsCheck() {
	fmt.Printf("Usage: %s -f apkfile [-c check certification] [-l check locale] [-v check version]\n",
		os.Args[0])
}

func usageAsPackZip() {
	fmt.Printf("Usage: %s -d directory [-l filelist]\n",
		os.Args[0])
}

/*
	asPackzip let app acts as packzip
	Usage:
	$ packzip -d directory [-l filelist]
*/
func asPackzip() {
	flag.Usage = usageAsPackZip
	directory := flag.String("d", "", "directory to pack")
	filelist := flag.String("l", "", "filelist")
	flag.Parse()

	args := flag.Args()
	if len(args) != 0 {
		fmt.Printf("Unknown argument: %v\n", args)
		os.Exit(1)
	}

	// there's no -d argument or empty string directory
	if *directory == "" {
		flag.Usage()
		os.Exit(1)
	}

	if _, err := os.Stat(*directory); os.IsNotExist(err) {
		fmt.Println("directory doesn't exist")
		os.Exit(1)
	}

	if *filelist != "" {
		if _, err := os.Stat(*filelist); os.IsNotExist(err) {
			fmt.Println("filelist doesn't exist")
			os.Exit(1)
		}
		PackZip(*directory, *filelist)
	}
}

func asCheck() {
	/*
		Usage:
		$ check -c -l -v -f file.apk
		-f apk filename
		-c certification check
		-l locale check
		-v version format check
	*/
	flag.Usage = usageAsCheck
	apkfile := flag.String("f", "", "apk filename")
	checkcert := flag.Bool("c", false, "enable certification check")
	checklocale := flag.Bool("l", false, "enable locale check")
	checkversion := flag.Bool("v", false, "enable version check")
	flag.Parse()

	args := flag.Args()
	if len(args) != 0 {
		fmt.Printf("Unknown argument: %v\n", args)
		os.Exit(1)
	}

	// there's no -d argument or empty string directory
	if *apkfile == "" {
		flag.Usage()
		os.Exit(1)
	}

	if _, err := os.Stat(*apkfile); os.IsNotExist(err) {
		fmt.Println("apk file doesn't exist")
		os.Exit(1)
	}

	if *checkcert {
		Checkcert("", "", "", "")
	}

	if *checklocale {
		Checklocale()
	}

	if *checkversion {
		Checkver()
	}

	if !*checkcert && !*checklocale && !*checkversion {
		flag.Usage()
		os.Exit(0)
	}
}

func main() {
	if os.Args[0] == "packzip" {
		// if true {
		asPackzip()
	} else {
		asCheck()
	}
}
