package util

import (
	"path/filepath"
	"runtime"
)

func getRoot() string {
	_, b, _, _ := runtime.Caller(0)

	return filepath.Join(filepath.Dir(b), "../..")

}

var Root = getRoot()
