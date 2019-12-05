package controllers

import (
	"MinecraftArchiveManagementServer/models"
	"fmt"
	"github.com/astaxie/beego"
)

type ArchiveController struct {
	beego.Controller
}
// @router / [get]
func (a *ArchiveController) Get()  {
	a.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", a.Ctx.Request.Header.Get("Origin"))
	archives:=models.GetArchive()
	fmt.Println(archives)
	a.Data["json"] = archives
	a.ServeJSON()
}
