/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/26 3:45 下午
 */
package redis

import (
	"core.bank/datamover/log"
	"core.bank/datamover/utils"
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
)

var (
	url string
)

var saveCmd = &cobra.Command{
	Use: "save",
	Short: "redis generates rdb snapshot files and outputs them to the specified directory",
	Args: cobra.MaximumNArgs(1),
	Run: saveCommandFunc,
}

func init() {
	saveCmd.Flags().StringVarP(&url, "url", "u", "redis://127.0.0.1:6379", "redis server url")
}

func saveCommandFunc(cmd *cobra.Command, args []string) {
	var output string
	if len(args) == 0 {
		output = utils.GenFilenameByDate("dump")
	} else {
		output = args[0]
	}

	err := saveRDBFile(output)
	if err != nil {
		log.Logger.Error("redis save rdb file error: " + err.Error())
		return
	}
	fmt.Println()
	log.Logger.Info("redis save rdb file to " + output + " on success!")
}


func saveRDBFile(outputPath string) error {

	execCmd := exec.Command("redis-cli", "-u", url, "save")
	out, err := execCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("redis save: " + err.Error())
	}
	log.Logger.Info("redis-cli save output: " + string(out))

	execCmd = exec.Command("redis-cli", "-u", url, "--rdb", outputPath)
	out, err = execCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("redis-cli rdb: " + err.Error())
	}
	log.Logger.Info("redis-cli rdb output: " + string(out))

	return nil
}
