package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/xelabs/go-mydumper/common"
	"github.com/xelabs/go-mydumper/config"
	"github.com/xelabs/go-mysqlstack/xlog"
)

var (
	log *xlog.Log = xlog.NewStdLog(xlog.Level(xlog.INFO))

	// 定义命令行参数
	user string
	password string
	host string
	port string
	database string
	output string
	sqlFile string

	all bool
	restore bool
	thread bool
)

func initFlags() {
	flag.StringVar(&user, "user", "root", "mysql user")
	flag.StringVar(&password, "password", "", "mysql password")
	flag.StringVar(&host, "host", "127.0.0.1", "mysql host")
	flag.StringVar(&port, "port", "3306", "mysql port")
	flag.StringVar(&database, "database", "", "mysql database")
	flag.StringVar(&output, "output", "default.sql", "output file")
	flag.StringVar(&sqlFile, "file", "", "certain sql file")

	flag.BoolVar(&all, "all-databases", false, "mysql all databases")
	flag.BoolVar(&restore, "restore", false, "restore database from a sql file")
	flag.BoolVar(&thread, "thread", false, "use multi-threaded mode")

	// 定义短名称的参数，使用同一个变量地址
	flag.StringVar(&user, "u", "root", "mysql user (shorthand)")
	flag.StringVar(&password, "p", "", "mysql password (shorthand)")
	flag.StringVar(&host, "h", "127.0.0.1", "mysql host (shorthand)")
	flag.StringVar(&port, "P", "3306", "mysql port (shorthand)")
	flag.StringVar(&database, "d", "", "mysql database (shorthand)")
	flag.StringVar(&sqlFile, "f", "", "certain sql file (shorthand)")
	flag.StringVar(&output, "o", "default.sql", "output file (shorthand)")

	flag.BoolVar(&all, "a", false, "mysql all databases")
	flag.BoolVar(&restore, "r", false, "restore database from a sql file")
	flag.BoolVar(&thread, "t", false, "use multi-threaded mode")

	flag.Parse()

	// 检查是否所有的 flag 都已经解析
	if !flag.Parsed() {
		log.Error("error: invalid flag format. Please use --flag or -f.")
		os.Exit(3)
	}
}

func main() {

	initFlags()

	var cmd *exec.Cmd
	var err error

	if thread {
		if restore {
			restoreArgs := &config.Config{
				User:            user,
				Password:        password,
				Address:         fmt.Sprintf("%s:%s", host, port),
				Outdir:          output,
				Threads:         16,
				IntervalMs:      10 * 1000,
				OverwriteTables: false,
			}
			log.Info("restore databases from the directory: " + output + " ...")
			fmt.Println()
			common.Loader(log, restoreArgs)

		} else {
			// TODO: do not dependent config file
			filePath := "./datamover.ini"
			dumperArgs, err := config.ParseDumperConfig(filePath)
			if err != nil {
				log.Error("read the datamover ini config error!")
				return
			}

			if dumperArgs == nil {
				log.Error("read the datamover ini config nil!")
				return
			}

			dumperArgs.User = user
			dumperArgs.Password = password
			dumperArgs.Address = fmt.Sprintf("%s:%s", host, port)
			if all {
				log.Info("dump all databases into file " + output + " at multi-threaded mode...")
				fmt.Println()
			} else {
				dumperArgs.DatabaseRegexp = ""
				// TODO: multi databases
				if len(database) == 0 {
					log.Error("Please provide at least one database name.")
					return
				}
				dumperArgs.Database = database
				log.Info("dump database " + database + " into " + output +  " directory... " )
				fmt.Println()
			}

			dumperArgs.Outdir = output
			if _, err := os.Stat(dumperArgs.Outdir); os.IsNotExist(err) {
				x := os.MkdirAll(dumperArgs.Outdir, 0o777)
				common.AssertNil(x)
			}
			common.Dumper(log, dumperArgs)
		}

	} else {
		if restore {
			if len(sqlFile) == 0 {
				log.Error("Please provide the certain sql file")
				return
			}

			log.Info("restore the database from the certain sql file: " + sqlFile + " ...")
			fmt.Println()

			cmd = exec.Command("mysql", "-u", user, "-p"+password, "--host", host, "--port", port)
			// 打开源文件，用于读取SQL语句
			file, err := os.Open(sqlFile)
			if err != nil {
				log.Fatal(err.Error())
			}
			defer file.Close()

			// 将命令的标准输入重定向到文件对象
			cmd.Stdin = file

		} else {
			if all {
				log.Info("dump all databases into file " + output + " ...")
				fmt.Println()
				cmd = exec.Command("mysqldump", "-u", user, "-p"+password, "--host", host, "--port", port, "-A")
			} else {
				// 检查数据库名是否为空
				if  len(database) == 0{
					log.Error("Please provide at least one database name.")
					return
				}

				log.Info("dump the certain database " + database +  " into file " + output + " ...")
				fmt.Println()
				cmd = exec.Command("mysqldump", "-u", user, "-p"+password, "--host", host, "--port", port,"--databases", database)
			}

			f, err := os.Create(output)
			if err != nil {
				log.Error("create file error: " + err.Error())
				return
			}
			defer f.Close()
			cmd.Stdout = f
		}

		err = cmd.Run()
		if err != nil {
			log.Info(cmd.String())
			log.Error("cmd Run error: " + err.Error())
			return
		}
	}

	log.Info("Success!")
}
