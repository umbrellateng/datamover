/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/19 2:50 下午
 */
package mysql

import (
	"core.bank/datamover/log"
	"github.com/spf13/cobra"
	"github.com/xelabs/go-mydumper/config"
)

var (
	thread bool

	from   string
	target string
)

func NewMysqlCommand() *cobra.Command {
	// mysqlCmd represents the mysql command
	cmd := &cobra.Command{
		Use:   "mysql",
		Short: "Realize data migration commands between isomorphic mysql",
		Long: "Realize data migration commands between isomorphic mysql, support single-threaded mode and multi-threaded mode",
		Run: mysqlCommandFunc,
	}

	cmd.AddCommand(dumpCmd)
	cmd.AddCommand(restoreCmd)
	cmd.AddCommand(onlineCmd)

	cmd.PersistentFlags().BoolVarP(&thread, "thread", "T", false, "whether target enable multi-threaded mode")
	cmd.PersistentFlags().StringVarP(&from, "from", "f", "root:root@tcp(localhost:3306)", "from mysql connection string")
	cmd.PersistentFlags().StringVarP(&target, "target", "t", "root:root@tcp(localhost:3306)", "target mysql connection string")

	return cmd
}

func mysqlCommandFunc(cmd *cobra.Command, args []string) {
	log.Logger.Info("mysql command")
}


func defaultConfig() *config.Config{
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
