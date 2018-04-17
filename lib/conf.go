package lib

import (
	"path"
	"fmt"
	"os"
	"bufio"
	"io"
	"bytes"
	"strings"
	"strconv"
)

type Config struct {
	RemoteHost string
	RemotePort int
	RemoteUser string
	RemotePass string
	RemoteKeyFile string
	RemoteBackupPath string
	BackupRule string
	LocalBackupPath string
	Type string
	TimeFormat string
	ConfigPath string
	Divider string
	Prefix string
	Suffix string
}


func CheckConfig() (conf Config) {
	currentWorkingPath := GetCurrentWorkPath()
	configDir := path.Join(currentWorkingPath,"conf")
	configFilePath := path.Join(configDir,"setting.conf")
	conf.ConfigPath = configDir
	fmt.Println(configFilePath)
	stat,err := os.Stat(configFilePath)
	if err != nil{
		panic("检查配置文件路径出错")
		panic(err)
	}
	if !stat.IsDir() {
		lines,err := ReadConfigLineByLine(configFilePath)
		if err != nil{
			panic("无法读取配置文件")
			panic(err.Error())
		}
		for _,line := range lines {
			if !strings.HasPrefix(line,"#") {

				strs := strings.Split(line,"=")
				confKey := GetConfigKey(strs)
				switch confKey {
				case "RemoteHost":
					PanicNoValueErr(strs)
					conf.RemoteHost = GetConfigValue(strs)
					break
				case "RemotePort":
					value := GetConfigValue(strs)
					if value == "" {
						 conf.RemotePort = 22
					}else {
						port,err := strconv.Atoi(value)
						if err != nil {
							panic("端口格式有误")
						}
						conf.RemotePort = port
					}
					break
				case "RemoteUser":
					PanicNoValueErr(strs)
					conf.RemoteUser = GetConfigValue(strs)
					break
				case "RemotePass":
					conf.RemotePass = GetConfigValue(strs)

					break
				case "RemoteKeyFile":
					conf.RemoteKeyFile = GetConfigValue(strs)
					// 下面将相对路径组装成绝对路径，未判断windows 下的，windows 下请放在conf目录下作为绝对路径
					if conf.RemoteKeyFile != "" {
						if !strings.HasPrefix(conf.RemoteKeyFile,"/") {
							conf.RemoteKeyFile = path.Join(configDir,conf.RemoteKeyFile)
						}
					}
					break
				case "RemoteBackupPath":
					PanicNoValueErr(strs)
					conf.RemoteBackupPath = GetConfigValue(strs)
					break
				case "BackupRule":
					PanicNoValueErr(strs)
					conf.BackupRule = GetConfigValue(strs)
					break
				case "LocalBackupPath":
					PanicNoValueErr(strs)
					conf.LocalBackupPath = GetConfigValue(strs)
					break
				case "Type":
					PanicNoValueErr(strs)
					conf.Type = GetConfigValue(strs)
					break
				case "TimeFormat":
					conf.TimeFormat = GetConfigValue(strs)
					break
				case "Divider":
					conf.Divider = GetConfigValue(strs)
					break
				case "Prefix":
					conf.Prefix =GetConfigValue(strs)
					break
				case "Suffix":
					conf.Suffix = GetConfigValue(strs)
					break
				}
			}
		}
	}
	if conf.RemotePass == "" && conf.RemoteKeyFile == "" {
		panic("配置文件的两种验证方式都不存在，请添加RemotePass或者RemoteKeyFile至少一项")
	}
	return
}

func ReadConfigLineByLine(path string) (lines []string,err error) {
	var (
		file *os.File
		part [] byte
		prefix bool
	)

	if file, err = os.Open(path); err != nil {
		return
	}

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte,1024))

	for {
		if part, prefix, err = reader.ReadLine();err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines,buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

func PanicNoValueErr(strs []string)  {
	if len(strs)<=1 {
		panic("不存在配置："+strs[0]+"的值，请检查配置文件")
	}
}

func GetConfigKey(strs []string) (key string) {
	if len(strs) <= 0 {
		panic("这不是一行有效的配置key")
	}else {
		key = strings.TrimSpace(strs[0])
	}
	return
}

func GetConfigValue(strs []string) (value string) {
	if len(strs)<=1 {
		value = ""
	}else {
		value = strings.TrimSpace(strs[1])
	}
	return
}



