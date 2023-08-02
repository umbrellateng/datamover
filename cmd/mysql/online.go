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
)


var onlineCmd = &cobra.Command{
	Use: "online",
	Short: "migrate mysql database online from one mysql to another",
	Run: onlineCommandFunc,
	Args: cobra.NoArgs,
}

func init() {
	onlineCmd.Flags().StringArrayVarP(&databases, "databases", "d", nil, "the dump databases of mysql")
	onlineCmd.Flags().BoolVarP(&all, "all-databases", "a", false, "all mysql databases except(mysql|sys|performance_schema|information_schema)")
}

func onlineCommandFunc(cmd *cobra.Command, args []string) {

	onlineTmpDir := "./tmpDir"

	defer func() {
		if r := recover(); r != nil {
			_ = utils.DeleteDirAndFiles(onlineTmpDir)
			log.Logger.Error("online move something wrong, received from panic: %v", r)
		}
	}()

	srcUsername, srcPassword, srcHost, srcPort, err := utils.ParseDBStringWithoutDB(from)
	if err != nil {
		log.Logger.Error("parse from mysql connection string error: " + err.Error())
		return
	}

	dstUsername, dstPassword, dstHost, dstPort, err := utils.ParseDBStringWithoutDB(target)
	if err != nil {
		log.Logger.Error("parse target mysql connection string error: " + err.Error())
		return
	}



	err = dumpToDirectory(srcUsername, srcPassword, srcHost, srcPort, onlineTmpDir)
	if err != nil {
		log.Logger.Error("dump mysql database online error: " + err.Error())
		return
	}

	err = restoreFromDirectory(dstUsername, dstPassword, dstHost, dstPort, onlineTmpDir)
	if err != nil {
		log.Logger.Error("restore mysql database online error: " + err.Error())
		return
	}

	err = utils.DeleteDirAndFiles(onlineTmpDir)
	if err != nil {
		log.Logger.Warning("remove " + onlineTmpDir + " dir error: %s", err.Error())
	}
	fmt.Println()  // 空一行
	log.Logger.Info("move mysql database online successfully!")
}