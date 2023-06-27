package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
)

func main() {

	// 定义命令行参数
	var user string
	var password string
	var host string
	var port string
	var database string
	var output string
	var sqlFile string

	var all bool
	var restore bool

	flag.StringVar(&user, "user", "root", "mysql user")
	flag.StringVar(&password, "password", "", "mysql password")
	flag.StringVar(&host, "host", "127.0.0.1", "mysql host")
	flag.StringVar(&port, "port", "3306", "mysql port")
	flag.StringVar(&database, "database", "", "mysql database")
	flag.StringVar(&output, "output", "default.sql", "output file")
	flag.StringVar(&sqlFile, "file", "", "certain sql file")

	flag.BoolVar(&all, "all-databases", false, "mysql all databases")
	flag.BoolVar(&restore, "restore", false, "restore database from a sql file")

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

	flag.Parse()


	// 检查是否所有的 flag 都已经解析
	if !flag.Parsed() {
		log.Println("error: invalid flag format. Please use --flag or -f.")
		os.Exit(3)
	}
	var cmd *exec.Cmd
	var err error

	if restore {
		if len(sqlFile) == 0 {
			log.Fatal("Please provide the certain sql file")
		}

		cmd = exec.Command("mysql", "-u", user, "-p"+password, "--host", host, "--port", port)
		// 打开源文件，用于读取SQL语句
		file, err := os.Open(sqlFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// 将命令的标准输入重定向到文件对象
		cmd.Stdin = file

	} else {
		if all {
			log.Println("enter in all")
			cmd = exec.Command("mysqldump", "-u", user, "-p"+password, "--host", host, "--port", port, "-A")
		} else {
			log.Println("enter in a certain database")
			// 检查数据库名是否为空
			if  len(database) == 0{
				log.Fatal("Please provide at least one database name.")
			}

			cmd = exec.Command("mysqldump", "-u", user, "-p"+password, "--host", host, "--port", port,"--databases", database)
		}

		f, err := os.Create(output)
		if err != nil {
			log.Fatal("create file error: " + err.Error())
		}
		defer f.Close()
		cmd.Stdout = f
	}

	err = cmd.Run()
	if err != nil {
		log.Println(cmd.String())
		log.Fatal("cmd Run error: " + err.Error())
	}

	log.Println("Complete!")
}
