/*
 * @Author: sjhuang
 * @Date: 2022-04-22 21:15:07
 * @LastEditTime: 2022-04-24 11:10:59
 * @FilePath: /nosdb/utils/file_utils.go
 */
package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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

func PrefixPath(path string, prefixName string) (fileName string, err error) {
	if err = CheckDir(path); err != nil {
		return
	}
	var files []os.FileInfo
	if files, err = ioutil.ReadDir(path); err != nil {
		return
	}
	var fileNames []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), prefixName) {
			fileNames = append(fileNames, file.Name())
		}
	}
	if len(fileNames) == 0 {
		err = fmt.Errorf(" %d  file match the %s ", len(fileNames), prefixName)
		return
	} else {
		// 如果不止一个，就返回最近的一个
		fileName = fileNames[len(fileNames)-1]
		return
	}
}
