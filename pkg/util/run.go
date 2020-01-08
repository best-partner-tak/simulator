package util

import (
	"bufio"
	"bytes"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// MustResolve returns an absolute path or panics if the underlying sys call
// fails
func MustResolve(wd string) string {
	absPath, err := filepath.Abs(wd)
	if err != nil {
		panic(err)
	}

	return absPath
}

// Run runs a child process and returns its buffer stdout.  Run also tees the
// output to stdout of this process, `env` will be appended to the current
// environment.  `wd` is the working directory for the child
func Run(wd string, env []string, cmd string, args ...string) (*string, error) {
	child := exec.Command(cmd, args...)

	child.Env = append(os.Environ(), env...)

	childIn, _ := child.StdinPipe()
	childErr, _ := child.StderrPipe()
	childOut, _ := child.StdoutPipe()
	defer childIn.Close()
	defer childErr.Close()
	defer childOut.Close()

	dir := MustResolve(wd)

	child.Dir = dir

	// Copy child stdout to stdout but also into a buffer to be returned
	var buf bytes.Buffer
	tee := io.TeeReader(childOut, &buf)

	err := child.Start()
	if err != nil {
		return nil, errors.Wrapf(err, "Error starting child process %s", cmd)
	}

	if _, err = io.Copy(os.Stdout, tee); err != nil {
		return nil, err
	}
	if _, err = io.Copy(os.Stderr, childErr); err != nil {
		return nil, err
	}

	err = child.Wait()
	// TODO: (rem) make this parameterisable?
	if err != nil && err.Error() != "exit status 127" {
		return nil, errors.Wrapf(err, "Error waiting for child process %s", cmd)
	}

	out := buf.String()
	return &out, nil
}

// RunSilently runs a sub command silently
func RunSilently(wd string, env []string, cmd string, args ...string) (*string, *string, error) {
	child := exec.Command(cmd, args...)

	child.Env = append(os.Environ(), env...)

	var outBuf, errBuf bytes.Buffer
	child.Stdout = bufio.NewWriter(&outBuf)
	child.Stderr = bufio.NewWriter(&errBuf)
	dir := MustResolve(wd)

	child.Dir = dir

	err := child.Start()
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Error starting child process %s", cmd)
	}

	err = child.Wait()
	// TODO: (rem) make this parameterisable?
	if err != nil && err.Error() != "exit status 127" {
		//Debug("Error waiting for child process", err)
		childErr := errBuf.String()
		return nil, &childErr, err
	}

	childErr := errBuf.String()
	childOut := outBuf.String()

	return &childOut, &childErr, nil
}
