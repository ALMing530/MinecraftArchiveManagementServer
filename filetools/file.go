package filetools

import (
	"fmt"
	"os"
)
/**
生成zip文件
params:
	data:将要写入的数据、
	file:目的zip文件
 */
func GenerateZipFile(data []byte,file *os.File)  {
	file.Write(data)
}
/**
初始化zip文件
params:
	name:文件名称
 */
func InitFile(name string) *os.File {
	var file *os.File
	path:="zip/"
	filePath:=path+name
	_,err:=os.Stat(path)
	if err==nil{
		file, err = os.Create(filePath)
		if err !=nil{
			fmt.Println(err)
		}
	}
	if os.IsNotExist(err){
		err = os.MkdirAll(path, os.ModeDir)
		if err!=nil{
			panic("create folder fail")
		}
		file, _ = os.Create(filePath)
	}
	return file
}