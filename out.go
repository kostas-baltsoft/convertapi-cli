package main

import (
	"fmt"
	"github.com/ConvertAPI/convertapi-go"
	"io"
	"os"
	"strings"
	"sync/atomic"
)

var outCnt uintptr = 0

func output(file *convertapi.ResFile, dst string) {
	dst = strings.TrimSpace(dst)
	if strings.EqualFold(dst, "stdout") {
		fd := outIdx()
		if fd == 2 {
			fd = outIdx()
		} // Skipping stderr
		out := os.NewFile(fd, file.FileName)
		io.Copy(out, file)
		file.Delete()
	} else if strings.HasPrefix(dst, "@") {
		dst = strings.TrimPrefix(dst, "@")
		paths := strings.Split(dst, ";")
		i := int(outIdx())
		if len(paths) < i {
			i = len(paths)
		}
		file.ToPath(paths[i-1])
		file.Delete()
	} else {
		fmt.Println(file.Url)
	}
}

func outIdx() uintptr {
	return atomic.AddUintptr(&outCnt, 1)
}
