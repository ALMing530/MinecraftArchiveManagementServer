package main

import (
	"MinecraftArchiveManagementServer/conf"
	_ "MinecraftArchiveManagementServer/routers"
	"MinecraftArchiveManagementServer/server"
	"MinecraftArchiveManagementServer/shell"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"os"
	"strconv"
)

func main() {
	go initArchiveManage()
		beego.Run()
	//maps:=filetools.ReadProperties("Oracle.properties")
	//fmt.Println(maps["user"])
}
func initArchiveManage() {
	initGlobalConf()
	initMinecraft()
	server.CreateServer()
}
func initGlobalConf() {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&conf.GlobalConf)
	if err != nil {
		fmt.Println("Parse Fail")
	}
	fmt.Println("loade minecraft directory:\t" + conf.GlobalConf.MinecraftDir)
	fmt.Println("loade server port:        \t" + strconv.Itoa(conf.GlobalConf.Port))
	fmt.Println("loade shell command       \t" + conf.GlobalConf.ShellBash)
	fmt.Println("loade Authentication      \t" + strconv.FormatBool(conf.GlobalConf.Authentication))
}
func initMinecraft() {
	shell.Minecraft("start")
}
