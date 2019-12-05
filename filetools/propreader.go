package filetools

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	properties map[string]string
)

func ReadProperties(path string) map[string]string {
	paris := make(map[string]string)
	propertiesFile, err := os.Open(path)
	checkErr(err)
	propertiesReader := bufio.NewReader(propertiesFile)
	for {
		line, isPrefix, err := propertiesReader.ReadLine()
		if err == io.EOF {
			break
		}
		if !isPrefix {
			keyValue := strings.Split(string(line), "=")
			if len(keyValue) == 2 {
				paris[keyValue[0]] = keyValue[1]
			}
		} else {
			fmt.Println("配置文件单行定义过长")
		}
	}
	return paris
}
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
