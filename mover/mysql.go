/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/4 5:19 下午
 */
package mover

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"core.bank/datamover/log"
	"core.bank/datamover/utils"
	"github.com/xelabs/go-mydumper/common"
	"github.com/xelabs/go-mydumper/config"
)

type Mysql struct {
	BaseInfo
	conf *config.Config
	databases []string
	all bool
}

func NewMySql(user, pwd, h, p string, a bool, datas []string) *Mysql {
	info := BaseInfo{
		username: user,
		password: pwd,
		host: h,
		port: p,
	}
	m :=  &Mysql{
		BaseInfo: info,
		databases: datas,
		all: a,
	}

	return m
}

func (m *Mysql) Username() string {
	return m.username
}

func (m *Mysql) Password() string {
	return m.password
}

func (m *Mysql) Host() string {
	return m.host
}

func (m *Mysql) Port() string {
	return m.port
}

func (m *Mysql) DefaultConfig() {
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

	m.conf = args
}

func (m *Mysql) setConfig(config *config.Config) {
	m.conf = config
}

func (m *Mysql) Config() *config.Config {
	return m.conf
}

func (m *Mysql) UrlString() string {
	return m.username + ":" + m.password + "@tcp(" + m.host + ":" + m.port + ")"
}

func (m *Mysql) DumpToFile(filePath string) error {

	var cmd *exec.Cmd
	user, password, host, port := m.Username(), m.Password(), m.Host(), m.Port()
	outputSqlFile := filePath
	if m.all {
		if len(outputSqlFile) == 0 {
			outputSqlFile = "all-databases.sql"
		}
		log.Logger.Info("dump all databases into file " + outputSqlFile + " ...")
		fmt.Println()
		cmd = exec.Command("mysqldump", "-u", user, "-p"+password, "--host", host, "--port", port, "-A")
	} else {
		// 检查数据库名是否为空
		if  len(m.databases) == 0{
			return fmt.Errorf("please provide at least one database name with flag --databases or -d")

		}
		if len(outputSqlFile) == 0 {
			outputSqlFile = m.databases[0] + ".sql"
		}
		log.Logger.Info("dump the certain database " + m.databases[0] +  " into file " + outputSqlFile + " ...")
		fmt.Println()
		cmd = exec.Command("mysqldump", "-u", user, "-p"+password, "--host", host, "--port", port,"--databases", m.databases[0])
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

func (m *Mysql) DumpToDirectory(outputDir string) error {

	m.DefaultConfig()
	dumperArgs := m.Config()
	dumperArgs.User = m.Username()
	dumperArgs.Password = m.Password()
	dumperArgs.Address = fmt.Sprintf("%s:%s", m.Host(), m.Port())
	if m.all {
		log.Logger.Info("dump all databases into file " + outputDir + " at multi-threaded mode...")
		if len(outputDir) == 0 {
			outputDir = "all-databases"
		}
		fmt.Println()
	} else {
		dumperArgs.DatabaseRegexp = ""
		if len(m.databases) == 0 {
			return fmt.Errorf("%s","please provide at least one database name with flag --databases or -d.")
		}
		databasesStr := strings.Join(m.databases, ",")
		dumperArgs.Database = databasesStr
		if len(outputDir) == 0 {
			if len(m.databases) == 1 {
				outputDir = m.databases[0]
			} else {
				outputDir = strings.Join(m.databases, "_")
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

func (m *Mysql) RestoreFromFile(inputFile string) error {
	if len(inputFile) == 0 {
		return fmt.Errorf("please provide the certain input sql file with flag --input or -i")
	}

	if utils.IsDirectory(inputFile) {
		return fmt.Errorf("the input " + inputFile + " is a directory, not a sql file, " +
			"please specify the sql file with flag --input or -i")
	}

	log.Logger.Info("restore the database from the certain sql file: " + inputFile + " ...")
	fmt.Println()

	cmd := exec.Command("mysql", "-u", m.Username(), "-p"+m.Password(), "--host", m.Host(), "--port", m.Port())
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

func (m *Mysql) RestoreFromDirectory(inputDir string) error {

	if !utils.IsDirectory(inputDir) {
		return fmt.Errorf("input is not a directory ,please specify the input directory with flag --input or -i")
	}

	conf := &config.Config{
		User:            m.Username(),
		Password:        m.Password(),
		Address:         fmt.Sprintf("%s:%s", m.Host(), m.Port()),
		Outdir:          inputDir,
		Threads:         16,
		IntervalMs:      10 * 1000,
		OverwriteTables: false,
	}
	m.setConfig(conf)
	log.Logger.Info("restore databases from the directory: " + inputDir + " ...")
	fmt.Println()
	common.Loader(log.Logger, m.Config())

	return nil
}

func (m *Mysql) MoveOnline(infos []BaseInfo) error {

	if len(infos) == 0 {
		return fmt.Errorf("the target info is empty")
	}

	info := infos[0]

	target := NewMySql(info.username, info.password, info.host, info.port, false, nil)

	log.Logger.Info("source database connection string: " + m.UrlString())
	log.Logger.Info("target database connection string: " + target.UrlString())

	onlineTmpDir := "./tmpDir"

	err := m.DumpToDirectory(onlineTmpDir)
	if err != nil {
		return fmt.Errorf("dump source database error: " + err.Error())
	}

	err = target.RestoreFromDirectory(onlineTmpDir)
	if err != nil {
		_ = utils.DeleteDirAndFiles(onlineTmpDir)
		return fmt.Errorf("restore target database error: " + err.Error())
	}

	err = utils.DeleteDirAndFiles(onlineTmpDir)
	if err != nil {
		return fmt.Errorf("remove " + onlineTmpDir + " dir error: %s", err.Error())
	}
	fmt.Println()  // 空一行
	log.Logger.Info("move database online successfully!")

	return nil
}

