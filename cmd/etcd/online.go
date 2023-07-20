/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/20 11:22 上午
 */
package etcd

import (
	"core.bank/datamover/log"
	"github.com/spf13/cobra"
)

var (
	libraryUse bool
)

var onlineCmd = &cobra.Command{
	Use:   "online",
	Short: "etcd migrate online command",
	Run:   onlineCommandFunc,
}

func init() {
	onlineCmd.Flags().BoolVarP(&libraryUse, "library-use", "l", false, "whether use etcd library to realize the etcd migrate")
}

func onlineCommandFunc(cmd *cobra.Command, args []string) {

}

func onlineMove(from, to string) error {
	if libraryUse {
		log.Logger.Info("use the third library of etcd to realize migrate...")
	} else {

	}

	return nil
}

func onlineUseLibrary(oldEndpoints, newEndpoints string) error {
	return nil
}
