package models

import (
	"MinecraftArchiveManagementServer/conf"
	"MinecraftArchiveManagementServer/filetools"
	"os"
	"path/filepath"
)

type Archive struct {
	Name string
	Active bool
	Description string
}

func GetArchive() []Archive {
	pairs:=filetools.ReadProperties(conf.GlobalConf.MinecraftDir+"/server.properties")
	var Archives []Archive
	filepath.Walk(conf.GlobalConf.MinecraftDir,func(path string,fileinfo os.FileInfo,err error) error{
		if fileinfo.IsDir(){
			fileName :=fileinfo.Name()
			if isMiecraftArchive(fileName){

				Archives=append(Archives,Archive{
					Name:        fileinfo.Name(),
					Active:      pairs["level"]==fileName,
					Description: "Minecraft archive",
				})
			}
		}
		return nil
	})
	return Archives
}
func isMiecraftArchive(fileName string) bool{
	otherFile :=[]string{"logs","crash-reports"}
	for _,value:=range otherFile{
		if value == fileName{
			return false
		}
	}
	return true
}
