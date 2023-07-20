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
	"github.com/xelabs/go-mydumper/config"
	"os"
	"os/exec"
)

var (
	input string
)

var restoreCmd = &cobra.Command{
	Use: "restore",
	Short: "restore data target mysql database",
	Run: restoreCommandFunc,
	Args: cobra.NoArgs,
}

func init() {
	restoreCmd.Flags().StringVarP(&input, "input","i","", "the input sql file or directory for mysql restore")
}

func restoreCommandFunc(cmd *cobra.Command, args []string) {

	username, password, host, port, err := utils.ParseDBStringWithoutDB(target)
	if err != nil {
		log.Logger.Error("parse mysql connection string error: " + err.Error())
		return
	}

	if thread {
		err := restoreFromDirectory(username, password, host, port, input)
		if err != nil {
			log.Logger.Error("restore mysql database in multi-threaded mode error: " + err.Error())
			return
		}
	} else {
		err := restoreFromSqlfile(username, password, host, port, input)
		if err != nil {
			log.Logger.Error("restore mysql database in single-threaded mode error: " + err.Error())
			return
		}
	}
	fmt.Println()
	log.Logger.Info("restore database on success!")
}

func restoreFromDirectory(username, password, host, port, inputDir string) error {
	if !utils.IsDirectory(inputDir) {
		return fmt.Errorf("input is not a directory ,please specify the input directory with flag --input or -i")
	}

	conf := &config.Config{
		User:            username,
		Password:        password,
		Address:         fmt.Sprintf("%s:%s", host, port),
		Outdir:          inputDir,
		Threads:         16,
		IntervalMs:      10 * 1000,
		OverwriteTables: false,
	}

	log.Logger.Info("restore databases from the directory: " + inputDir + " ...")
	fmt.Println()
	common.Loader(log.Logger, conf)

	return nil
}

func restoreFromSqlfile(username, password, host, port, inputFile string) error {
	if len(inputFile) == 0 {
		return fmt.Errorf("please provide the certain input sql file with flag --input or -i")
	}

	if utils.IsDirectory(inputFile) {
		return fmt.Errorf("the input " + inputFile + " is a directory, not a sql file, " +
			"please specify the sql file with flag --input or -i")
	}

	log.Logger.Info("restore the database from the certain sql file: " + inputFile + " ...")
	fmt.Println()

	cmd := exec.Command("mysql", "-u", username, "-p" + password, "--host", host, "--port", port)
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
