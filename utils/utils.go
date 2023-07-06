/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/4 11:08 上午
 */
package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"core.bank/datamover/log"
)

func IsDirectory(input string) bool {
	info, err := os.Stat(input)
	// 判断是否有错误发生
	if err != nil {
		log.Logger.Error("judge directory error: " + err.Error())
		os.Exit(4)
	}
	// 调用 IsDir 函数判断是否是目录
	if !info.IsDir() {
		return false
	}
	return true
}

func DeleteDirAndFiles(dir string) error {

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		tmpPath := filepath.Join(dir, entry.Name())
		if entry.Type().IsRegular() {
			err = os.Remove(tmpPath)
			if err != nil {
				log.Logger.Error("remove file " + tmpPath + " error: ", err.Error())
				continue
			}
		} else {
			_ = DeleteDirAndFiles(tmpPath)
		}
	}

	err = os.Remove(dir)
	if err != nil {
		return err
	}

	log.Logger.Info("remove dir " + dir + " on success!")
	return nil
}

// 定义一个函数，接受一个字符串参数，返回4个字符串和一个错误值  user:Password@tcp(localhost:3306)
func ParseDBStringWithoutDB(s string) (string, string, string, string, error) {
	// 按照 @ 符号分割字符串，得到用户名和密码部分和主机和端口部分
	var username, password, host, port string
	parts := strings.Split(s, "@")
	if len(parts) != 2 {
		return "", "", "", "", fmt.Errorf("invalid input, not found \"@\" character")
	}

	// 按照 : 符号分割用户名和密码部分，得到用户名和密码
	userpass := strings.Split(parts[0], ":")
	if len(userpass) != 2 {
		return "", "", "", "", fmt.Errorf("invalid input, not found \":\" character")
	}
	username = userpass[0]
	password = userpass[1]

	// 按照 ( 符号去掉主机和端口部分的 tcp 前缀，得到主机和端口
	hostport := strings.TrimPrefix(parts[1], "tcp(")
	hostport = strings.TrimSuffix(hostport, ")")
	hostports := strings.Split(hostport, ":")
	if len(hostports) != 2 {
		return "", "", "", "", fmt.Errorf("invalid input, can not find host and port")
	}
	host = hostports[0]
	port = hostports[1]

	return username, password, host, port, nil
}

func OnlineMode(from, to string) bool {
	j0 := len(from) != 0 && len(to) != 0
	j1 := strings.Contains(from, "@tcp(")
	j2 := strings.Contains(to, "@tcp(")
	j3 := strings.Contains(from, ":")
	j4 := strings.Contains(to, ":")
	return j0 && j1 && j2 && j3 && j4
}

//// 定义一个函数，接受一个字符串参数，返回一个 DBInfo 结构体和一个错误值
//func ParseDBString(s string) (DBInfo, error) {
//	// 定义一个空的 DBInfo 结构体
//	var info DBInfo
//
//	// 按照 @ 符号分割字符串，得到用户名和密码部分和主机和端口部分
//	parts := strings.Split(s, "@")
//	if len(parts) != 2 {
//		return info, fmt.Errorf("invalid format1")
//	}
//
//	// 按照 : 符号分割用户名和密码部分，得到用户名和密码
//	userpass := strings.Split(parts[0], ":")
//	if len(userpass) != 2 {
//		return info, fmt.Errorf("invalid format2")
//	}
//	info.username = userpass[0]
//	info.password = userpass[1]
//
//	// 按照 / 符号分割主机和端口部分，得到主机和端口和数据库名
//	hostportdb := strings.Split(parts[1], "/")
//	if len(hostportdb) != 2 {
//		return info, fmt.Errorf("invalid format3")
//	}
//
//	// 按照 ( 符号去掉主机和端口部分的 tcp 前缀，得到主机和端口
//	hostport := strings.TrimPrefix(hostportdb[0], "tcp(")
//	hostport = strings.TrimSuffix(hostport, ")")
//	hostports := strings.Split(hostport, ":")
//	if len(hostports) != 2 {
//		return info, fmt.Errorf("invalid format4")
//	}
//	info.host = hostports[0]
//	info.port = hostports[1]
//
//	// 得到数据库名
//	info.database = hostportdb[1]
//
//	// 返回解析后的结构体和 nil 错误值
//	return info, nil
//}
//
//func PrintDBInfo(s string) {
//
//	info, err := parseDBStringWithoutDB(s)
//	if err != nil {
//		log.Error("parse db string error: " + err.Error())
//		return
//	}
//	log.Info("%s %s %s %s %s",info.username, info.password, info.host, info.port, info.database)
//}

