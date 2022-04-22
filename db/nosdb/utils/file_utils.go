package utils

import (
	"fmt"
	"os"
)

const (
	PERM = 0664 // 0110 0110 0100 , rwe rwe rwe
)

// 重新命名一个文件
func ReNameFile(path string, oldName string, newName string) error {
	oldfilePath := fmt.Sprintf("%s%d%s", path, os.PathSeparator, oldName)
	newfilePath := fmt.Sprintf("%s%d%s", path, os.PathSeparator, newName)
	return os.Rename(oldfilePath, newfilePath)
}

// 新建一个文件
func Open(path string, fileName string) (*os.File, error) {
	filePath := fmt.Sprintf("%s%d%s", path, os.PathSeparator, fileName)
	return os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, PERM)
}
