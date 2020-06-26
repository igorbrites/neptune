package terraform

import (
	"bytes"
	"os/exec"
	"syscall"
)

const defaultFailedCode = 1

func RunCommand(name string, args ...string) (stdout string, stderr string, exitCode int) {
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	stdout = outbuf.String()
	stderr = errbuf.String()

	if err != nil {
			// try to get the exit code
			if exitError, ok := err.(*exec.ExitError); ok {
					ws := exitError.Sys().(syscall.WaitStatus)
					exitCode = ws.ExitStatus()
			} else {
					// This will happen (in OSX) if `name` is not available in $PATH,
					// in this situation, exit code could not be get, and stderr will be
					// empty string very likely, so we use the default fail code, and format err
					// to string and set to stderr
					exitCode = defaultFailedCode
					if stderr == "" {
							stderr = err.Error()
					}
			}
	} else {
			// success, exitCode should be 0 if go is ok
			ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
	}

	return
}

func IsCommandAvailable(path string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v " + path)

	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}
