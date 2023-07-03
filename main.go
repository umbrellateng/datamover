package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/xelabs/go-mydumper/common"
	"github.com/xelabs/go-mydumper/config"
	"github.com/xelabs/go-mysqlstack/xlog"
)

// 定义一个自定义类型，实现 flag.Value 接口
type dbSlice []string

// 实现 String 方法，返回参数的字符串表示
func (d *dbSlice) String() string {
	return fmt.Sprint(*d)
}

// 实现 Set 方法，将参数值追加到切片中
func (d *dbSlice) Set(value string) error {
	*d = append(*d, value)
	return nil
}

var (
	log *xlog.Log = xlog.NewStdLog(xlog.Level(xlog.INFO))

	// 定义命令行参数
	user     string
	password string
	host     string
	port     string
	//database string
	output   string
	input    string
	databases dbSlice

	all bool
	restore bool
	thread bool
)

func initFlags() {
	flag.StringVar(&user, "user", "root", "mysql user")
	flag.StringVar(&password, "password", "root", "mysql password")
	flag.StringVar(&host, "host", "127.0.0.1", "mysql host")
	flag.StringVar(&port, "port", "3306", "mysql port")
	flag.StringVar(&output, "output", "", "output file or directory")
	flag.StringVar(&input, "input", "", "input file or directory")
	flag.Var(&databases, "databases", "database name(s)")

	flag.BoolVar(&all, "all-databases", false, "mysql all databases")
	flag.BoolVar(&restore, "restore", false, "restore database from a sql file")
	flag.BoolVar(&thread, "thread", false, "use multi-threaded mode")

	// 定义短名称的参数，使用同一个变量地址
	flag.StringVar(&user, "u", "root", "mysql user (shorthand)")
	flag.StringVar(&password, "p", "root", "mysql password (shorthand)")
	flag.StringVar(&host, "h", "127.0.0.1", "mysql host (shorthand)")
	flag.StringVar(&port, "P", "3306", "mysql port (shorthand)")
	flag.StringVar(&input, "i", "", "input file or directory (shorthand)")
	flag.StringVar(&output, "o", "", "output file or directory (shorthand)")
	flag.Var(&databases, "d","database name(s) (shorthand)" )

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

func isDirectory(input string) bool {
	info, err := os.Stat(input)
	// 判断是否有错误发生
	if err != nil {
		log.Error("judge directory error: " + err.Error())
		os.Exit(4)
	}
	// 调用 IsDir 函数判断是否是目录
	if !info.IsDir() {
		return false
	}
	return true
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

func DumpDBToDirectory(outputDir string) error {
	dumperArgs := DefaultConfig()
	dumperArgs.User = user
	dumperArgs.Password = password
	dumperArgs.Address = fmt.Sprintf("%s:%s", host, port)
	if all {
		log.Info("dump all databases into file " + output + " at multi-threaded mode...")
		if len(outputDir) == 0 {
			outputDir = "all-databases"
		}
		fmt.Println()
	} else {
		dumperArgs.DatabaseRegexp = ""
		if len(databases) == 0 {
			return fmt.Errorf("%s","please provide at least one database name.")
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
		log.Info("dump database " + databasesStr + " into " + outputDir +  " directory... " )
		fmt.Println()
	}

	dumperArgs.Outdir = outputDir
	if _, err := os.Stat(dumperArgs.Outdir); os.IsNotExist(err) {
		x := os.MkdirAll(dumperArgs.Outdir, 0o777)
		common.AssertNil(x)
	}
	common.Dumper(log, dumperArgs)

	return nil
}

func DumpDBToSqlFile(outputSqlFile string) error {
	var cmd *exec.Cmd
	if all {
		if len(outputSqlFile) == 0 {
			outputSqlFile = "all-databases.sql"
		}
		log.Info("dump all databases into file " + outputSqlFile + " ...")
		fmt.Println()
		cmd = exec.Command("mysqldump", "-u", user, "-p"+password, "--host", host, "--port", port, "-A")
	} else {
		// 检查数据库名是否为空
		if  len(databases) == 0{
			return fmt.Errorf("Please provide at least one database name.")

		}
		if len(outputSqlFile) == 0 {
			outputSqlFile = databases[0] + ".sql"
		}
		log.Info("dump the certain database " + databases[0] +  " into file " + outputSqlFile + " ...")
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
		log.Info(cmd.String())
		return fmt.Errorf("cmd Run error: " + err.Error())
	}

	return nil
}

func RestoreDBFromDirectory(inputDir string) error {
	if !isDirectory(inputDir) {
		return fmt.Errorf("input is not a directory ,please specify the input directory with flag --input or -i")
	}

	restoreArgs := &config.Config{
		User:            user,
		Password:        password,
		Address:         fmt.Sprintf("%s:%s", host, port),
		Outdir:          inputDir,
		Threads:         16,
		IntervalMs:      10 * 1000,
		OverwriteTables: false,
	}
	log.Info("restore databases from the directory: " + input + " ...")
	fmt.Println()
	common.Loader(log, restoreArgs)

	return nil
}

func RestoreDBFromSqlFile(inputFile string) error {
	if len(inputFile) == 0 {
		return fmt.Errorf("please provide the certain input sql file with flag --input or -i")
	}

	if isDirectory(inputFile) {
		return fmt.Errorf("the input " + input + " is a directory, not a sql file, " +
			"please specify the sql file with flag --input or -i")
	}

	log.Info("restore the database from the certain sql file: " + inputFile + " ...")
	fmt.Println()

	cmd := exec.Command("mysql", "-u", user, "-p"+password, "--host", host, "--port", port)
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
		log.Info(cmd.String())
		return err
	}

	return nil
}

func main() {
	initFlags()

	var err error
	if thread {
		if restore {
			err = RestoreDBFromDirectory(input)
			if err != nil {
				log.Error("Restore DB from Directory " + input + " error: " + err.Error())
				return
			}

		} else {
			err = DumpDBToDirectory(output)
			if err != nil {
				log.Error("Dump DB to Directory " + output + " error: " + err.Error())
				return
			}
		}

	} else {
		if restore {
			err = RestoreDBFromSqlFile(input)
			if err != nil {
				log.Error("Restore DB from SqlFile error: " + err.Error())
				return
			}

		} else {
			err = DumpDBToSqlFile(output)
			if err != nil {
				log.Error("Dump DB to SqlFile error: " + err.Error())
				return
			}
		}
	}

	log.Info("Success!")
}
