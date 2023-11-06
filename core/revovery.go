package core

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

type Result interface{}

var (
	// ShowShortName 显示短的名字
	ShowShortName = true
)

// ErrorWW 输出异常 Error with where
func ErrorWithStack1(err error) string {
	if err == nil {
		return ""
	}
	return string(Stack3(2, 5, ShowShortName)) + " -> " + err.Error()
}

// ErrorWW 输出异常 Error with where
func ErrorWithStack3(err error, skip int, max int) string {
	if err == nil {
		return ""
	}
	return string(Stack3(skip, max, ShowShortName)) + " -> " + err.Error()
}

//===============
// gin.Revovery
//===============

// Stack returns a nicely formatted stack frame, skipping skip frames.
func Stack(skip int) []byte {
	return Stack2(skip, 0)
}

// Whereis where is
func Whereis(short bool) []byte {
	return Stack3(2, 3, short)
}

// Stack2 returns a nicely formatted stack frame, skipping skip frames.
func Stack2(skip int, max int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; max <= 0 || i < max; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x); ", file, line, pc)
		if file != lastFile {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// Stack3 returns a nicely formatted stack frame, skipping skip frames.
func Stack3(skip, max int, short bool) []byte {
	buf := new(bytes.Buffer)
	for i := skip; max <= 0 || i < max; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		if short {
			if j1 := strings.LastIndex(file, "/"); j1 > 0 {
				if j2 := strings.LastIndex(file[:j1], "/"); j2 > 0 {
					file = file[j2+1:]
				} else {
					file = file[j1+1:]
				}
			}
		}
		fmt.Fprintf(buf, "%s:%d (0x%x); ", file, line, pc)
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
