package main

import (
	"backupdbs/lib"
	"fmt"
	"path"
)

func main() {
	conf := lib.CheckConfig()
	fmt.Println(conf)

	sftp := lib.SFTP{
		Username:conf.RemoteUser,
		Host:conf.RemoteHost,
		Port:conf.RemotePort,
		PrivateKeyFile:conf.RemoteKeyFile,
	}
	name := lib.MakeFileName(conf)
	if conf.Type == "get" {
		sftp.GetTo(conf.LocalBackupPath).From(path.Join(conf.RemoteBackupPath,name))
	}else if conf.Type == "put" {
		sftp.Put(path.Join(conf.LocalBackupPath,name)).To(conf.RemoteBackupPath)
	}


}


