/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/19 3:25 下午
 */
package mysql

import (
	"core.bank/datamover/log"
	"core.bank/datamover/utils"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xelabs/go-mydumper/common"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	output string
	databases []string
	tables []string

	all bool
	withoutCreateDatabase bool
)

var dumpCmd = &cobra.Command{
	Use: "dump",
	Short: "dump data from mysql database",
	Run: dumpCommandFunc,
	Args: cobra.NoArgs,
}

func init() {
	dumpCmd.Flags().StringArrayVarP(&databases, "databases", "d", nil, "the dump databases of mysql")
	dumpCmd.Flags().StringArrayVarP(&tables, "tables", "t", nil, "the table name of some database")
	dumpCmd.Flags().StringVarP(&output, "output", "o", "", "the location that save the dump file or directory ")

	dumpCmd.Flags().BoolVarP(&all, "all-databases", "a", false, "all mysql databases except(mysql|sys|performance_schema|information_schema)")
	dumpCmd.Flags().BoolVarP(&withoutCreateDatabase, "without-create-database", "w", false, "if true the create-database.sql will be removed from the output directory")
}

func dumpCommandFunc(cmd *cobra.Command, args []string) {
	username, password, host, port, err := utils.ParseDBStringWithoutDB(from)
	if err != nil {
		log.Logger.Error("parse mysql connection string error: " + err.Error())
		return
	}

	if thread {
		outputDir, err := dumpToDirectory(username, password, host, port, output)
		if err != nil {
			log.Logger.Error("dump mysql databases in multi-threaded mode error: " + err.Error())
			return
		}

		if withoutCreateDatabase && len(databases) == 1{
			fileName := filepath.Join(outputDir, fmt.Sprintf("%s-schema-create.sql", databases[0]))
			err := os.Remove(fileName)
			log.Logger.Info("remove the sql file: %s", fileName)
			if err != nil {
				log.Logger.Warning("remove the %s file error: %s", fileName, err.Error())
			}
		}

	} else {
		err := dumpToSqlFile(username, password, host, port, output)
		if err != nil {
			log.Logger.Error("dump mysql database in single-threaded mode error: " + err.Error())
			return
		}
	}
	fmt.Println()
	log.Logger.Info("dump database on success!")
}

func dumpToDirectory(username, password, host, port, outputDir string) (string, error) {
	dumperArgs := defaultConfig()
	dumperArgs.User = username
	dumperArgs.Password = password
	dumperArgs.Address = fmt.Sprintf("%s:%s", host, port)
	if all {
		if len(outputDir) == 0 {
			outputDir = "all-databases"
		}
		log.Logger.Info("dump all databases into file " + outputDir + " at multi-threaded mode...")
		fmt.Println()
	} else {
		dumperArgs.DatabaseRegexp = ""
		if len(databases) == 0 {
			return "", fmt.Errorf(" please provide at least one database name with flag --databases or -d.")
		}
		databasesStr := strings.Join(databases, ",")
		dumperArgs.Database = databasesStr
		if len(tables) != 0 {
			if len(databases) != 1 {
				return "", fmt.Errorf("if you specify the tables, the database can only specify one")
			}
			dumperArgs.Table = strings.Join(tables, ",")
		}
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

	return outputDir, nil
}

func dumpToSqlFile(username, password, host, port, outputFile string) error {
	var execCmd *exec.Cmd
	if all {
		if len(outputFile) == 0 {
			outputFile = "all-databases.sql"
		}
		log.Logger.Info("dump all databases into file " + outputFile + " ...")
		fmt.Println()
		execCmd = exec.Command("mysqldump", "-u", username, "-p"+password, "--host", host, "--port", port, "-A")
	} else {
		// 检查数据库名是否为空
		if  len(databases) == 0{
			return fmt.Errorf("please provide the mysql database name with flag --databases or -d")
		}

		if len(databases) > 1 {
			return fmt.Errorf("only one mysql database dump is supported in single-threaded mode")
		}

		if len(outputFile) == 0 {
			outputFile = databases[0] + ".sql"
		}
		log.Logger.Info("dump the certain database " + databases[0] +  " into file " + outputFile + " ...")
		fmt.Println()
		execCmd = exec.Command("mysqldump", "-u", username, "-p"+password, "--host", host, "--port", port,"--databases", databases[0])
	}

	f, err := os.Create(outputFile)
	if err != nil {
		log.Logger.Error("create " +  outputFile + " file error: " + err.Error())
	}
	defer f.Close()

	execCmd.Stdout = f
	err = execCmd.Run()
	if err != nil {
		log.Logger.Info(execCmd.String())
		return fmt.Errorf("mysqldump command run error: " + err.Error())
	}

	return nil
}
