package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/ConvertAPI/convertapi-go"
	"github.com/ConvertAPI/convertapi-go/config"
	"github.com/ConvertAPI/convertapi-go/param"
	"os"
	"sync"
)

const Version = 1
const Name = "convertapi"
const HelpFlagName = "help"

func main() {
	infF := flag.String("inf", "", "input format e.g. docx, pdf, jpg")
	outfF := flag.String("outf", "", "output format e.g. pdf, jpg, zip")
	secretF := flag.String("secret", "", "ConvertAPI user secret. Get your secret at https://www.convertapi.com/a")
	proxyF := flag.String("proxy", "", "HTTP proxy server")
	outF := flag.String("out", "url", "place where to output converted files. Allowed values: [ url | @<path> | stdout ]. Save to directory example: --out=\"@/path/to/dir\" ")
	paramsF := flag.String("params", "", "conversion parameters, see full list of available parameters at https://www.convertapi.com .  Allowed values: [ value | @<path> | @< | < ]. Usage example: --params=\"file:@/path/to/file.doc, pdftitle:My title\"")
	verF := flag.Bool("version", false, "output version information and exit")
	helpF := flag.Bool(HelpFlagName, false, "display this help and exit")
	flag.Parse()

	if *verF {
		printVersion()
	}
	if *helpF {
		printHelp()
	}
	if *infF == "" {
		printError(errors.New("Input format is not set. Please set --inf"), 1)
	}
	if *outfF == "" {
		printError(errors.New("Output format is not set. Please set --outf"), 1)
	}
	if *proxyF != "" { /*TODO: set proxy*/
	}

	if *secretF == "" {
		printError(errors.New("ConvertAPI user secret is not set. Please set --secret parameter. Get your secret at https://www.convertapi.com/a"), 1)
	} else {
		config.Default.Secret = *secretF
	}

	if *paramsF == "" {
		printError(errors.New("Conversion parameters are not set. Please set --params parameter."), 1)
	} else {
		if paramsets, err := parseParams(*paramsF, *infF); err == nil {
			wg := &sync.WaitGroup{}
			//fmt.Printf("%+v\n", paramsets)
			for _, set := range paramsets {
				//fmt.Println("Uploading ", set)
				if err := prepare(set); err != nil {
					printError(err, 1)
				}
				//fmt.Println("Converting ", set)

				go func(set []param.IParam) {
					wg.Add(1)
					defer wg.Done()
					res := convertapi.Convert(*infF, *outfF, set, nil)
					if files, err := res.Files(); err == nil {
						//fmt.Println("Downloading ", set)
						for _, file := range files {
							output(file, *outF)
						}
					} else {
						printError(err, 1)
					}
				}(set)
			}
			wg.Wait()
		} else {
			printError(fmt.Errorf("Conversion parameters are invalid. %s", err), 1)
		}
	}
}

func printError(err error, exitCode int) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", Name, err)
	fmt.Fprintf(os.Stderr, "Try '%s --%s' for more information.\n", Name, HelpFlagName)
	os.Exit(exitCode)
}

func printHelp() {
	flag.PrintDefaults()
	os.Exit(0)
}

func printVersion() {
	fmt.Printf("%s %d\n", Name, Version)
	os.Exit(0)
}
