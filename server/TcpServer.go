package server

import (
	"MinecraftArchiveManagementServer/conf"
	"MinecraftArchiveManagementServer/entity"
	"MinecraftArchiveManagementServer/filetools"
	"MinecraftArchiveManagementServer/shell"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func CreateServer() {
	addr := ":" + strconv.Itoa(conf.GlobalConf.Port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	checkError(err)
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := tcpListener.AcceptTCP()
		checkError(err)
		fmt.Println("Server listening at port:" + addr)
		handleArrive(conn)
	}
}
func handleArrive(conn net.Conn) {
	var authorized = conf.GlobalConf.Authentication
	if authorized {
		conn.Write([]byte("authority"))
		success := auth(conn, authorized)
		if success {
			for handleData(conn) {
			}
		}
	} else {
		conn.Write([]byte("noAuthority"))
		for handleData(conn) {
		}
	}

}
func handleData(conn net.Conn) (isConnect bool) {
	var data []byte = make([]byte, 4)

	_, err := conn.Read(data) // 读取控制信息
	checkError(err)
	cmd := parseBytesToInt32(data[0:]) //控制信息
	fmt.Println(cmd)
	if err != nil {
		conn.Close()
		fmt.Println("Connect lose")
		return false
	}
	switch cmd {
	case 1:
		/**读取参数长度，每次上传文件会先生成文件参数（名称，大小，校验码等）参数以json格式字符串生成并计算转换后byte数组的长度，将计算结果先发送过来*/
		_, err = conn.Read(data)
		checkError(err)
		paramLength := parseBytesToInt32(data[0:])
		var param []byte = make([]byte, paramLength)
		_, err = conn.Read(param) //读取文件参数
		checkError(err)
		var params entity.Params
		err = json.Unmarshal(param, &params)
		checkError(err)
		file := filetools.InitFile(params.Name + ".zip")
		var byteCount float64 = 0
		buffer := make([]byte, 1024)
		for byteCount < params.Length { //读取文件
			len, err := conn.Read(buffer[0:])
			if len == 0 {
				break
			} else {
				byteCount += float64(len)
				filetools.GenerateZipFile(buffer[0:len], file)
				checkError(err)
			}

		}
		checkError(err)
		err = file.Close()
		checkError(err)
		shell.Minecraft("stop")           //关闭minecraft服务器
		filetools.UnCompress(file.Name()) //解压读取到的文件
		_, err = conn.Write([]byte("File uploade done"))
		shell.Minecraft("start") //启动minecraft服务器

	case 2:
		sendFile(conn)
	}
	return true
}
func sendFile(conn net.Conn) {

	fileInfo, _ := os.Lstat(conf.GlobalConf.CurrentArchive)
	nameHasExtend := fileInfo.Name()
	nameHasNotExtend := strings.Split(nameHasExtend, ".")[0]
	destFile:= "tempzip/"+nameHasNotExtend+".zip"
	err:=filetools.Compress(conf.GlobalConf.CurrentArchive, destFile)
	zipFileInfo,_:=os.Lstat("tempzip/"+nameHasNotExtend+".zip")
	fileLength := float64(zipFileInfo.Size())
	sendTime := time.Now().Format("2006-01-02 03:04:05")
	fileType := "file"
	if fileInfo.IsDir() {
		fileType = "Directory"
	}
	params := entity.Params{
		Name:   nameHasExtend,
		Type:   fileType,
		Length: fileLength,
		Time:   sendTime,
	}
	paramsBinary, _ := json.Marshal(&params)
	fmt.Println(string(paramsBinary))
	paramLength := int32(len(paramsBinary))
	_, _ = conn.Write(parseInt32ToBytes(paramLength))
	_, _ = conn.Write(paramsBinary)
	file, err := os.Open(destFile)
	checkError(err)
	buffer := make([]byte, 1024)
	count := 0
	for {
		read, err := file.Read(buffer)
		fmt.Println(err)
		if read != 0 {
			conn.Write(buffer[0:read])
		} else {
			break
		}
	}

	fmt.Println(count)
}

/**
权限验证
params:
	conn: 当前socket连接
	authentication:是否进行权限验证
*/
func auth(conn net.Conn, authentication bool) bool {
	param := make([]byte, 1024)
	len, err := conn.Read(param)
	if err != nil {
		conn.Close()
		return false
	}
	fmt.Println(len)
	if authentication {
		checkError(err)
		fmt.Println(string(param[0:len]))
		params := strings.Split(string(param[0:len]), "&")
		users := conf.GlobalConf.User
		for _, user := range users {
			if params[0] == user.Username && params[1] == user.Password {
				conn.Write([]byte("success"))
				return true
			} else {
				conn.Write([]byte("fail"))
			}
		}
	} else {
		return true
	}
	conn.Close()
	return false
}
func parseBytesToInt32(args []byte) int32 {
	var param int32
	buffer := bytes.NewBuffer(args)
	binary.Read(buffer, binary.BigEndian, &param)
	return param

}
func parseInt32ToBytes(args int32) []byte {
	i := int32(args)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, i)
	return bytesBuffer.Bytes()
}
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
