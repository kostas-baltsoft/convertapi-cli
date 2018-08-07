package main

import (
	"../convertapi-go"
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func parseParams(paramString string, ext string) (paramsets [][]*convertapi.Param, err error) {
	var newParams []*convertapi.Param
	var parallel bool
	for _, p := range strings.Split(paramString, ",") {
		kv := strings.Split(p, ":")
		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])
		if newParams, parallel, err = newCaParams(k, v, ext); err != nil {
			return
		}
		paramsets = mergeParams(paramsets, newParams, parallel)
	}
	return
}

func mergeParams(paramsets [][]*convertapi.Param, params []*convertapi.Param, parallel bool) (res [][]*convertapi.Param) {
	if params == nil || len(params) == 0 {
		return paramsets
	}
	if parallel {
		for _, param := range params {
			if paramsets == nil || len(paramsets) == 0 {
				res = append(res, []*convertapi.Param{param})
			}
			for _, set := range paramsets {
				mergedSet := append(set, param)
				res = append(res, mergedSet)
			}
		}
	} else {
		if paramsets == nil || len(paramsets) == 0 {
			return [][]*convertapi.Param{params}
		}
		for i, set := range paramsets {
			res[i] = append(set, params...)
		}
	}
	return
}

func newCaParams(k string, v string, ext string) (caParams []*convertapi.Param, parallel bool, err error) {
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

		if paths, err = flattenPaths(paths, ext); err != nil {
			return
		}

		for i, p := range paths {
			name := k
			if !parallel {
				name = fmt.Sprintf("%s[%d]", k, i)
			}
			caParam := convertapi.NewFilePathParam(name, p, nil)
			caParams = append(caParams, caParam)
		}
	} else if strings.HasPrefix(v, "<<") {
		caParam := convertapi.NewReaderParam(k, os.Stdin, "file."+ext, nil)
		caParams = append(caParams, caParam)
	} else {
		var vals []string
		if strings.HasPrefix(v, "<") {
			vals, err = stdinLines()
		} else {
			vals = []string{v}
		}

		for i, val := range vals {
			name := k
			if !parallel {
				name = fmt.Sprintf("%s[%d]", k, i)
			}
			caParam := convertapi.NewStringParam(name, val)
			caParams = append(caParams, caParam)
		}
	}
	return
}

func stdinLines() (lines []string, err error) {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	return
}

func flattenPaths(paths []string, ext string) (res []string, err error) {
	flat := []string{}
	for _, p := range paths {
		if flat, err = dirToFiles(p, ext); err != nil {
			return
		}
		res = append(res, flat...)
	}
	return
}

func dirToFiles(path string, ext string) (paths []string, err error) {
	dir, err := isDir(path)
	if err == nil {
		if dir {
			paths = []string{}
			wildcardPath := filepath.Join(path, "*."+ext)
			files, err := filepath.Glob(wildcardPath)
			if err == nil {
				for _, f := range files {
					paths = append(paths, f)
				}
			}
			sort.Strings(paths)
		} else {
			if strings.EqualFold(filepath.Ext(path), "."+ext) {
				paths = []string{path}
			} else {
				err = errors.New(fmt.Sprintf("File %s is not %s format.", path, ext))
			}
		}
	}
	return
}

func isDir(path string) (isDir bool, err error) {
	info, err := os.Stat(path)
	if err == nil {
		isDir = info.IsDir()
	}
	return
}
