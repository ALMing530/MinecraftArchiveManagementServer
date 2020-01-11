package filetools

import (
	"MinecraftArchiveManagementServer/conf"
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

/**
解压文件
params:
	fileDir:目的文件的文件目录
*/
func UnCompress(fileDir string) {
	minecraftDir := conf.GlobalConf.MinecraftDir
	r, err := zip.OpenReader(fileDir)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, file := range r.Reader.File {
		if file.FileInfo().IsDir() {
			err := os.MkdirAll(minecraftDir+file.Name, os.ModeDir)
			checkError(err)
			//err = os.Chmod(minecraftDir+file.Name, 0755)
			checkError(err)
		} else {
			parentPath := getDir(file.Name)
			if parentPath != "" {
				err := os.MkdirAll(minecraftDir+parentPath, os.ModeDir)
				checkError(err)
				//err = os.Chmod(minecraftDir+parentPath, 0755)
				checkError(err)
			}
			openFile, err := file.Open()
			checkError(err)
			newFile, err := os.Create(minecraftDir + file.Name)
			fmt.Println(minecraftDir + file.Name)
			checkError(err)
			//err = os.Chmod(minecraftDir+file.Name, 0755)
			checkError(err)
			_, err = io.Copy(newFile, openFile)
			err = openFile.Close()
			checkError(err)
			err = newFile.Close()
			checkError(err)
			//panic("err")
		}
	}
	err = r.Close()
	checkError(err)
}

/**
压缩文件
params:
	srcFile:源文件目录
	destZip:压缩后文件目录
*/
func Compress(srcFile string, destZip string) error {
	srcFileInfo, _ := os.Lstat(srcFile)
	srcFileName := srcFileInfo.Name()
	zipFile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		index := strings.Index(path, srcFileName)
		header.Name = path[index:]
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	return err
}

/**
获取文件父目录
params:
	fileName:文件目录
*/
func getDir(fileName string) (dir string) {
	index := strings.LastIndex(fileName, "/")
	if index == -1 {
		return
	}
	runeStr := []rune(fileName)
	dir = string(runeStr[:index])
	return
}
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
