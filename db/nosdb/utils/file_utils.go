/*
 * @Author: sjhuang
 * @Date: 2022-04-22 21:15:07
 * @LastEditTime: 2022-04-24 09:40:42
 * @FilePath: /nosdb/utils/file_utils.go
 */
package utils

import (
	"errors"
	"fmt"
	"os"
)

const (
	FILE_PERM = 0664 // 0110 0110 0100 , rwe rwe rwe
	PATH_PERM = 0777
)

// 重新命名一个文件
func ReNameFile(path string, oldName string, newName string) error {
	oldfilePath := fmt.Sprintf("%s%c%s", path, os.PathSeparator, oldName)
	newfilePath := fmt.Sprintf("%s%c%s", path, os.PathSeparator, newName)
	return os.Rename(oldfilePath, newfilePath)
}

func FileIsExists(path string, fileName string) bool {
	filePath := fmt.Sprintf("%s%c%s", path, os.PathSeparator, fileName)
	file, err := os.Stat(filePath)
	if err != nil {
		return os.IsExist(err)
	}
	return !file.IsDir()
}

func CheckDir(path string) error {
	d, err := os.Stat(path)
	if err != nil {
		if !os.IsExist(err) {
			// 如果不存在，则新建一个文件夹
			return os.Mkdir(path, PATH_PERM)
		}
		return nil
	}
	if !d.IsDir() {
		err = errors.New(" 该文件不是文件夹")
		return err
	}
	return nil
}

// 新建一个文件
func Open(path string, fileName string) (*os.File, error) {
	if err := CheckDir(path); err != nil {
		return nil, err
	}
	filePath := fmt.Sprintf("%s%c%s", path, os.PathSeparator, fileName)
	return os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, FILE_PERM)
}
