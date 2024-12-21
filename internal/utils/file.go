package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GenUniqueFilePath(filePath string) string {
	var	dir			string
	var	name		string
	var	fileBase	string
	var	ext			string
	var	counter		int
	var	err			error

	dir = filepath.Dir(filePath)
	ext = filepath.Ext(filePath)
	fileBase = filepath.Base(filePath)
	name = strings.TrimSuffix(fileBase, ext)
	counter = 1
	for {
		filePath = filepath.Join(dir, fileBase)
		if _, err = os.Stat(filePath); os.IsNotExist(err) {
			return dir + "/" + fileBase
		}
		fileBase = fmt.Sprintf("%s_%d%s", name, counter, ext)
		counter++
	}
}
