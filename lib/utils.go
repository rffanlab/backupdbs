package lib

import (
	"path/filepath"
	"os"
	"strings"
	"fmt"
	"time"
)

func GetCurrentWorkPath() string {
	dir,err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil{
		fmt.Println(err.Error())
	}
	return strings.Replace(dir,"\\","/",-1)
}

func IsRelativePath(str string) (isOrNot bool) {

	return
}

/**
	当前时间：精确到小时，格式为2006-01-02_15
 */
func TimeNowForHour() string {
	return time.Now().Format("2006-01-02_15")
}

/**
	格式化时间
 */
func FormattedTime(format string) string  {
	return time.Now().Format(format)
}

/**
	当前时间：精确到秒，格式为2006-01-02 15:04:05
 */
func TimeNowForSec() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func MakeFileName(config Config) (filename string) {
	var prefix string
	var suffix string
	if config.Prefix != "" {
		prefix =  config.Prefix+config.Divider
	}
	if config.Suffix != "" {
		suffix = "."+config.Suffix
	}
	filename = prefix+FormattedTime(config.TimeFormat)+suffix
	return filename
}