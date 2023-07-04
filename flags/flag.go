/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/4 11:02 上午
 */
package flags

import (
	"flag"
	"fmt"
	"os"

	"core.bank/datamover/log"
)

// 定义一个自定义类型，实现 flag.Value 接口
type DBSlice []string

// 实现 String 方法，返回参数的字符串表示
func (d *DBSlice) String() string {
	return fmt.Sprint(*d)
}

// 实现 Set 方法，将参数值追加到切片中
func (d *DBSlice) Set(value string) error {
	*d = append(*d, value)
	return nil
}

var (
	// 定义命令行参数
	User     string
	Password string
	Host     string
	Port     string
	//database string
	Output    string
	Input     string
	From      string
	To        string
	Databases DBSlice

	All     bool
	Restore bool
	Thread  bool
)

func InitFlags() {
	flag.StringVar(&User, "user", "root", "mysql User")
	flag.StringVar(&Password, "password", "root", "mysql Password")
	flag.StringVar(&Host, "host", "127.0.0.1", "mysql Host")
	flag.StringVar(&Port, "port", "3306", "mysql Port")
	flag.StringVar(&Output, "output", "", "Output file or directory")
	flag.StringVar(&Input, "input", "", "Input file or directory")
	flag.StringVar(&From, "from", "", "source database connection string （root:123456@tcp(localhost:3306)）")
	flag.StringVar(&To, "to", "", "target database connection string")

	flag.Var(&Databases, "databases", "database name(s)")

	flag.BoolVar(&All, "all-databases", false, "mysql All Databases")
	flag.BoolVar(&Restore, "restore", false, "Restore database From a sql file")
	flag.BoolVar(&Thread, "thread", false, "use multi-threaded mode")

	// 定义短名称的参数，使用同一个变量地址
	flag.StringVar(&User, "u", "root", "mysql User (shorthand)")
	flag.StringVar(&Password, "p", "root", "mysql Password (shorthand)")
	flag.StringVar(&Host, "h", "127.0.0.1", "mysql Host (shorthand)")
	flag.StringVar(&Port, "P", "3306", "mysql Port (shorthand)")
	flag.StringVar(&Input, "i", "", "Input file or directory (shorthand)")
	flag.StringVar(&Output, "o", "", "Output file or directory (shorthand)")
	flag.Var(&Databases, "d","database name(s) (shorthand)" )

	flag.BoolVar(&All, "a", false, "mysql All Databases")
	flag.BoolVar(&Restore, "r", false, "Restore database From a sql file")
	flag.BoolVar(&Thread, "t", false, "use multi-threaded mode")

	flag.Parse()

	// 检查是否所有的 flag 都已经解析
	if !flag.Parsed() {
		log.Logger.Error("error: invalid flag format. Please use --flag or -f.")
		os.Exit(3)
	}
}
