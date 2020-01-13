package util

import (
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"sync"
)

var cache = struct {
	sync.RWMutex
	homedir string
}{}

// ExpandTilde returns the fully qualified path to a file in the user's home
// directory. I.E. it expands a path beginning with `~/`) and checks the file
// exists. ExpandTilde will cache the user's home directory to amortise the
// cost of the syscall
func ExpandTilde(path string) (*string, error) {
	if len(path) == 0 || path[0:2] != "~/" {
		return nil, errors.Errorf(`Path was empty or did not start with a tilde and a slash: "%s"`, path)
	}

	// discard ~/
	path = path[2:]

	var homedir string

	// Lock and read the cache to see if we already resolved the current user's
	// home directory
	cache.RLock()
	homedir = cache.homedir
	cache.RUnlock()
	if homedir == "" {
		// Take a write lock to update the cache
		cache.Lock()
		defer cache.Unlock()

		usr, err := user.Current()
		if err != nil {
			return nil, errors.Wrapf(err, "Error finding %s in home", path)
		}

		homedir = usr.HomeDir
		cache.homedir = homedir
	}

	p := filepath.Join(homedir, path)

	return &p, nil
}

// FileExists checks whether a path exists
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Slurp reads an entire file into a string in one operation and returns a
// pointer to the file's content or an error if any.  Similar to
// `ioutil.ReadFile` but it calls `filepath.Abs` first which cleans the path
// and resolves relative paths from the working directory.
//
// Note that this is slightly less efficient for zero-length files than
// `ioutil.Readfile` as it uses the default read buffer size of `bytes.MinRead`
// internally
func Slurp(path string) (*string, error) {
	fp, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Error resolving %s to slurp", path)
	}

	file, err := os.Open(fp)
	if err != nil {
		return nil, errors.Wrapf(err, "Error opening %s to slurp", path)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.Wrapf(err, "Error slurping %s", path)
	}

	output := string(b)
	return &output, nil
}

// MustSlurp is the panicky counterpart to Slurp. MustSlurp reads an entire
// file into a string in one operation and returns the contents or panics if it
// encouters and error
func MustSlurp(path string) string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return string(b)
}

// MustRemove removes a file or empty directory.  MustRemove will ignore an
// error if the path doesn't exist or panic for any other error
func MustRemove(path string) {
	err := os.Remove(path)
	if err != nil && os.IsNotExist(err) {
		return
	}
	if err != nil {
		panic(err)
	}
}

// EnsureFile checks a file exists and writes the supplied contents if not.
// returns a boolean indicating whether it wrote a file or not and any error
func EnsureFile(path, contents string) (bool, error) {
	exists, err := FileExists(path)
	if err != nil || exists {
		return false, err
	}

	file, err := os.Create(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	_, err = io.WriteString(file, contents)
	if err != nil {
		return false, err
	}

	err = file.Sync()
	if err != nil {
		return true, err
	}

	return true, nil
}

// OverwriteFile writes the supplied contents overwriting the path if it
// already exists.  It returns an error if any occurred
func OverwriteFile(path, contents string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, contents)
	if err != nil {
		return err
	}

	err = file.Sync()
	if err != nil {
		return err
	}

	return nil
}

// EnvOrDefault tries to read an environment variable with the supplied key and
// returns its value.  EnvOrDefault returns a default value if it is empty or
// unset
func EnvOrDefault(key, def string) string {
	var d = os.Getenv(key)
	if d == "" {
		d = def
	}

	return d
}
