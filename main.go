package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/ConvertAPI/convertapi-go"
	"os"
)

const Version = 1
const Name = "convertapi"
const HelpFlagName = "help"

func main() {
	infF := flag.String("inf", "", "input format e.g. docx, pdf, jpg")
	outfF := flag.String("outf", "", "output format e.g. pdf, jpg, zip")
	secretF := flag.String("secret", "", "ConvertAPI user secret. Get your secret at https://www.convertapi.com/a")
	proxyF := flag.String("proxy", "", "HTTP proxy server")
	//saveF := flag.String("save", "", "path for saving files after conversion")
	paramsF := flag.String("params", "", "conversion parameters, see full list of available parameters at https://www.convertapi.com")
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
	if *paramsF == "" {
		printError(errors.New("Conversion parameters are not set. Please set --params parameter."), 1)
	} else {
		params := parseParams(*paramsF)
		fmt.Println(params["name"])
	}
	if *secretF == "" {
		printError(errors.New("ConvertAPI user secret is not set. Please set --secret parameter. Get your secret at https://www.convertapi.com/a"), 1)
	} else {
		convertapi.Default.Secret = *secretF
	}

	if *proxyF != "" {
		//TODO: set proxy
	}

	printError(errors.New("No arguments set"), 1)
}

func printError(err error, exitCode int) {
	fmt.Printf("%s: %s\n", Name, err)
	fmt.Printf("Try '%s --%s' for more information.", Name, HelpFlagName)
	os.Exit(exitCode)
}

func printHelp() {
	flag.PrintDefaults()
	os.Exit(0)
}

func printVersion() {
	fmt.Sprintf("%s %d", Name, Version)
	os.Exit(0)
}

/*

convertapi doc pdf --file=@/tmp/src.doc --save=/tmp/dst.pdf
convertapi doc jpg --file=@/tmp/src.doc --save=/tmp
convertapi pdf merge --files[]=@/tmp/pdfs/* --save=/tmp/dst.pdf
convertapi pdf merge --files[]=@/tmp/pdfs/pirst.pdf;@/tmp/pdfs/second.pdf --save=/tmp/dst.pdf
cat /tmp/src.pdf | convertapi pdf compress --file=@< --save=/tmp/dst.pdf
convertapi doc jpg --file=@/tmp/src.doc | convertapi any zip --files[]=< --save=/tmp/dst.zip
cat /tmp/src.pdf | convertapi any zip --files[]=@<myfile.pdf --save=/tmp/dst.zip

convertapi --inf=pdf --outf=merge --save=/tmp/dst.pdf --params={"files[]":"@/tmp/pdfs/*"}


*/
