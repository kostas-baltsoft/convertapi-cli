package main

import (
	"github.com/ConvertAPI/convertapi-go"
	"github.com/ConvertAPI/convertapi-go/config"
	"github.com/ConvertAPI/convertapi-go/param"
	"sync"
)

func convert(inFormat string, outFormat string, paramsets [][]param.IParam, out string) {
	config.Default.HttpClient = newHttpClient()
	wg := &sync.WaitGroup{}
	//	debug(paramsets)
	for _, set := range paramsets {
		//		debug("Uploading ", set)
		if err := prepare(set); err != nil {
			printError(err, 1)
		}
		//		debug("Converting ", set)

		wg.Add(1)
		go func(set []param.IParam) {
			defer wg.Done()
			res := convertapi.Convert(inFormat, outFormat, set, nil)
			if files, err := res.Files(); err == nil {
				//				debug("Downloading ", set)
				for _, file := range files {
					output(file, out)
				}
			} else {
				printError(err, 1)
			}
		}(set)
	}
	wg.Wait()
}
