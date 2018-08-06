package main

import (
	"errors"
	"fmt"
	"github.com/ConvertAPI/convertapi-go"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type caParams []*convertapi.Param
type caConversions []*caParams

func parseParams(paramString string) (paramMap map[string]string) {
	paramMap = make(map[string]string)
	for _, p := range strings.Split(paramString, ",") {
		kv := strings.Split(p, ":")
		paramMap[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	return
}

func newCaParams(k string, v string, inf string) (caParams []*convertapi.Param, parallel bool, err error) {
	parallel = !strings.HasSuffix(k, "[]")
	if !parallel {
		k = strings.TrimSuffix(k, "[]")
	}

	if strings.HasPrefix(v, "@") {
		v = strings.TrimPrefix(v, "@")

		paths := []string{}
		if strings.HasPrefix(v, "<") {
			paths, err = stdinLines()
		} else {
			paths = strings.Split(v, ";")
		}

		paths, err = inFlatten(paths, inf)

		for i, p := range paths {
			name := k
			if !parallel {
				name = fmt.Sprintf("%s[%d]", k, i)
			}
			caParam := convertapi.NewFilePathParam(name, p, nil)
			caParams = append(caParams, caParam)
		}
	} else if strings.HasPrefix(v, "<") {
		if strings.HasPrefix(v, "<<") {
			caParam := convertapi.NewReaderParam(k, os.Stdin, "file."+inf, nil)
			caParams = append(caParams, caParam)
		} else {
			var urls []string
			urls, err = stdinLines()
			for i, url := range urls {
				name := k
				if !parallel {
					name = fmt.Sprintf("%s[%d]", k, i)
				}
				caParam := convertapi.NewStringParam(name, url)
				caParams = append(caParams, caParam)
			}
		}
	}
	return
}

func stdinLines() (lines []string, err error) {
	b := []byte{}
	if b, err = ioutil.ReadAll(os.Stdin); err == nil {
		lines = strings.Split(string(b), "\n")
	}
	return
}

func inFlatten(paths []string, inf string) (res []string, err error) {
	var flat []string
	for _, p := range paths {
		if flat, err = inDirToFiles(p, inf); err == nil {
			res = append(res, flat...)
		} else {
			return
		}

	}
}

func inDirToFiles(path string, inf string) (paths []string, err error) {
	dir, err := isDir(path)
	if err == nil {
		if dir {
			paths = []string{}
			wildcardPath := filepath.Join(path, "*."+inf)
			files, err := filepath.Glob(wildcardPath)
			if err == nil {
				for _, f := range files {
					paths = append(paths, f)
				}
			}
			sort.Strings(paths)
		} else {
			if strings.EqualFold(filepath.Ext(path), "."+inf) {
				paths = []string{path}
			} else {
				err = errors.New(fmt.Sprintf("File %s is not %s format.", path, inf))
			}
		}
	}
	return
}

func isDir(path string) (bool, error) {
	info, err := os.Stat(path)
	return info.IsDir(), err
}
