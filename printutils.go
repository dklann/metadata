package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/coreos/go-systemd/journal"
)

// verbosePrint adds a bit of "template" around verbose print
// statements. Output directs to STDERR if attached to a terminal,
// otherwise to the system journal.
func verbosePrint(message string) {
	pcs := make([]uintptr, 12)
	n := runtime.Callers(2, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	frame, _ := frames.Next()
	// No newline here: message includes it, or not.
	if cmdline.Verbose {
		if fileInfo, _ := os.Stderr.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
			fmt.Fprintf(os.Stderr, "\t%s:%s():%d: %s",
				path.Base(frame.File), frame.Function, frame.Line, message)
		} else {
			if journal.Enabled() {
				_ = journal.Print(journal.PriInfo,
					"%s():%d: %s",
					frame.Function, frame.Line, message)
			}
		}
	}
}

// debugPrint adds a bit of "template" around debugging print
// statements. Output directs to STDERR if attached to a terminal,
// otherwise to the system journal.
func debugPrint(message string) {
	pcs := make([]uintptr, 12)
	n := runtime.Callers(2, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	frame, _ := frames.Next()
	// No newline here: message includes it, or not.
	if cmdline.Debug {
		if fileInfo, _ := os.Stderr.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
			fmt.Fprintf(os.Stderr, "\t[D] %s:%s():%d: %s",
				path.Base(frame.File), frame.Function, frame.Line, message)
		} else {
			if journal.Enabled() {
				_ = journal.Print(journal.PriDebug,
					"%s():%d: %s",
					frame.Function, frame.Line, message)
			}
		}
	}
}

// logSprintf wraps fmt.Sprintf() with a consistent template.
func logSprintf(format string, a ...interface{}) string {
	pcs := make([]uintptr, 12)
	n := runtime.Callers(2, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	frame, _ := frames.Next()

	// No newline here: format includes it, or not.
	augmentedFormat := fmt.Sprintf("%s:%s():%d: %s", path.Base(frame.File), frame.Function, frame.Line, format)

	if len(a) > 0 {
		return fmt.Sprintf(augmentedFormat, a...)
	}
	return augmentedFormat
}
