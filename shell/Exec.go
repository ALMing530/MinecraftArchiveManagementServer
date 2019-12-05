package shell

import (
	"MinecraftArchiveManagementServer/conf"
	"os/exec"
)

var cmd *exec.Cmd

func Minecraft(op string) {
	var command string
	switch op {
	case "start":
		command = `./Start_Minecraft.sh`
		cmd := exec.Command(conf.GlobalConf.ShellBash, "-c", command)
		cmd.Start()
	case "stop":
		command = `./Stop_Minecraft.sh`
		cmd := exec.Command(conf.GlobalConf.ShellBash, "-c", command)
		cmd.Run()
	}

}
