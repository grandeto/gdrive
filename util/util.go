package util

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/grandeto/gdrive/constants"
)

func GetDefaultConfigDir() string {
	return filepath.Join(Homedir(), ".gdrive")
}

func ConfigFilePath(basePath, name string) string {
	return filepath.Join(basePath, name)
}

func GetAuthFileNamePath() string {
	return filepath.Join(Homedir(), constants.AuthFileName)
}

func Homedir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("APPDATA")
	}

	homedir := os.Getenv("HOME")

	if homedir != "" {
		return homedir
	}

	if osUser := os.Getenv("OS_USER"); osUser != "" {
		return filepath.Join(constants.HomeDir, osUser)
	}

	return homedir
}

func Equal(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func ExitF(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Println("")
	os.Exit(1)
}

func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func WriteJson(path string, data interface{}) error {
	tmpFile := path + ".tmp"
	f, err := os.Create(tmpFile)
	if err != nil {
		return err
	}

	err = json.NewEncoder(f).Encode(data)
	f.Close()
	if err != nil {
		os.Remove(tmpFile)
		return err
	}

	return os.Rename(tmpFile, path)
}

func Md5sum(path string) string {
	h := md5.New()
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	io.Copy(h, f)
	return fmt.Sprintf("%x", h.Sum(nil))
}
