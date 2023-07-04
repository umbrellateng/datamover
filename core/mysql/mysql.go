/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/4 11:12 上午
 */
package mysql

import (
	"core.bank/datamover/log"
	"core.bank/datamover/utils"
	"fmt"
	"github.com/xelabs/go-mydumper/common"
	"github.com/xelabs/go-mydumper/config"
	"os"
	"os/exec"
	"strings"
)

// 定义一个结构体，存储解析后的信息
type DBInfo struct {
	Username string
	Password string
	Host     string
	Port     string
}

func DefaultConfig() *config.Config{
	args := &config.Config{
		User: "root",
		Password: "root",
		Address: "127.0.0.1:3306",
		Database: "",
		DatabaseRegexp: "^(mysql|sys|information_schema|performance_schema)$",
		DatabaseInvertRegexp: true,
		Table: "",
		Outdir: "",
		ChunksizeInMB: 128,
		SessionVars: "",
		Threads: 16,
		StmtSize: 1000000,
		IntervalMs: 10 * 1000,
		Wheres: make(map[string]string),
	}

	return args
}

func DumpDBToDirectory(dbInfo DBInfo, outputDir string, databases []string, all bool) error {

	dumperArgs := DefaultConfig()
	dumperArgs.User = dbInfo.Username
	dumperArgs.Password = dbInfo.Password
	dumperArgs.Address = fmt.Sprintf("%s:%s", dbInfo.Host, dbInfo.Port)
	if all {
		log.Logger.Info("dump all databases into file " + outputDir + " at multi-threaded mode...")
		if len(outputDir) == 0 {
			outputDir = "all-databases"
		}
		fmt.Println()
	} else {
		dumperArgs.DatabaseRegexp = ""
		if len(databases) == 0 {
			return fmt.Errorf("%s","please provide at least one database name with flag --databases or -d.")
		}
		databasesStr := strings.Join(databases, ",")
		dumperArgs.Database = databasesStr
		if len(outputDir) == 0 {
			if len(databases) == 1 {
				outputDir = databases[0]
			} else {
				outputDir = strings.Join(databases, "_")
			}
		}
		log.Logger.Info("dump database " + databasesStr + " into " + outputDir +  " directory... " )
		fmt.Println()
	}

	dumperArgs.Outdir = outputDir
	if _, err := os.Stat(dumperArgs.Outdir); os.IsNotExist(err) {
		x := os.MkdirAll(dumperArgs.Outdir, 0o777)
		common.AssertNil(x)
	}
	common.Dumper(log.Logger, dumperArgs)

	return nil
}

func DumpDBToSqlFile(dbInfo DBInfo, outputSqlFile string, databases []string, all bool) error {
	var cmd *exec.Cmd
	user, password, host, port := dbInfo.Username, dbInfo.Password, dbInfo.Host, dbInfo.Port
	if all {
		if len(outputSqlFile) == 0 {
			outputSqlFile = "all-databases.sql"
		}
		log.Logger.Info("dump all databases into file " + outputSqlFile + " ...")
		fmt.Println()
		cmd = exec.Command("mysqldump", "-u", user, "-p"+password, "--host", host, "--port", port, "-A")
	} else {
		// 检查数据库名是否为空
		if  len(databases) == 0{
			return fmt.Errorf("please provide at least one database name.")

		}
		if len(outputSqlFile) == 0 {
			outputSqlFile = databases[0] + ".sql"
		}
		log.Logger.Info("dump the certain database " + databases[0] +  " into file " + outputSqlFile + " ...")
		fmt.Println()
		cmd = exec.Command("mysqldump", "-u", user, "-p"+password, "--host", host, "--port", port,"--databases", databases[0])
	}

	f, err := os.Create(outputSqlFile)
	if err != nil {
		return fmt.Errorf("create file error: " + err.Error())
	}
	defer f.Close()
	cmd.Stdout = f
	err = cmd.Run()
	if err != nil {
		log.Logger.Info(cmd.String())
		return fmt.Errorf("cmd Run error: " + err.Error())
	}

	return nil
}

func RestoreDBFromDirectory(dbInfo DBInfo, inputDir string) error {
	if !utils.IsDirectory(inputDir) {
		return fmt.Errorf("input is not a directory ,please specify the input directory with flag --input or -i")
	}

	restoreArgs := &config.Config{
		User:            dbInfo.Username,
		Password:        dbInfo.Password,
		Address:         fmt.Sprintf("%s:%s", dbInfo.Host, dbInfo.Port),
		Outdir:          inputDir,
		Threads:         16,
		IntervalMs:      10 * 1000,
		OverwriteTables: false,
	}
	log.Logger.Info("restore databases from the directory: " + inputDir + " ...")
	fmt.Println()
	common.Loader(log.Logger, restoreArgs)

	return nil
}

func RestoreDBFromSqlFile(dbInfo DBInfo, inputFile string) error {
	if len(inputFile) == 0 {
		return fmt.Errorf("please provide the certain input sql file with flag --input or -i")
	}

	if utils.IsDirectory(inputFile) {
		return fmt.Errorf("the input " + inputFile + " is a directory, not a sql file, " +
			"please specify the sql file with flag --input or -i")
	}

	log.Logger.Info("restore the database from the certain sql file: " + inputFile + " ...")
	fmt.Println()

	cmd := exec.Command("mysql", "-u", dbInfo.Username, "-p"+dbInfo.Password, "--host", dbInfo.Host, "--port", dbInfo.Port)
	// 打开源文件，用于读取SQL语句
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("open input file error: " + err.Error())
	}
	defer file.Close()

	// 将命令的标准输入重定向到文件对象
	cmd.Stdin = file
	err = cmd.Run()
	if err != nil {
		log.Logger.Info(cmd.String())
		return err
	}

	return nil
}



func IsOnlineMode(from, to string) bool {
	return len(from) != 0 && len(to) != 0
}

func OnlineMove(dbInfo DBInfo, from, to string, databases []string, all bool) error {
	log.Logger.Info("source database connection string: " + from)
	log.Logger.Info("target database connection string: " + to)
	fromInfo, err := ParseDBStringWithoutDB(from)
	if err != nil {
		return fmt.Errorf("parse source database connection error: " + err.Error())
	}
	toInfo, err := ParseDBStringWithoutDB(to)
	if err != nil {
		return fmt.Errorf("parse target database connection error: " + err.Error())
	}

	onlineTmpDir := "./tmpDir"

	err = DumpDBToDirectory(fromInfo, onlineTmpDir, databases, all)
	if err != nil {
		return fmt.Errorf("dump source database error: " + err.Error())
	}

	err = RestoreDBFromDirectory(toInfo, onlineTmpDir)
	if err != nil {
		_ = utils.DeleteDirAndFiles(onlineTmpDir)
		return fmt.Errorf("restore target database error: " + err.Error())
	}

	err = utils.DeleteDirAndFiles(onlineTmpDir)
	if err != nil {
		return fmt.Errorf("remove " + onlineTmpDir + " dir error: ", err.Error())
	}
	fmt.Println()  // 空一行
	log.Logger.Info("move database online successfully!")
	return nil
}

// 定义一个函数，接受一个字符串参数，返回一个 DBInfo 结构体和一个错误值  user:Password@tcp(localhost:3306)
func ParseDBStringWithoutDB(s string) (DBInfo, error) {
	// 定义一个空的 DBInfo 结构体
	var info DBInfo

	// 按照 @ 符号分割字符串，得到用户名和密码部分和主机和端口部分
	parts := strings.Split(s, "@")
	if len(parts) != 2 {
		return info, fmt.Errorf("invalid format1")
	}

	// 按照 : 符号分割用户名和密码部分，得到用户名和密码
	userpass := strings.Split(parts[0], ":")
	if len(userpass) != 2 {
		return info, fmt.Errorf("invalid format2")
	}
	info.Username = userpass[0]
	info.Password = userpass[1]

	// 按照 ( 符号去掉主机和端口部分的 tcp 前缀，得到主机和端口
	hostport := strings.TrimPrefix(parts[1], "tcp(")
	hostport = strings.TrimSuffix(hostport, ")")
	hostports := strings.Split(hostport, ":")
	if len(hostports) != 2 {
		return info, fmt.Errorf("invalid format4")
	}
	info.Host = hostports[0]
	info.Port = hostports[1]

	// 返回解析后的结构体和 nil 错误值
	return info, nil
}
