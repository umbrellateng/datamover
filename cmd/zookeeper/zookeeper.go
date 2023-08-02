/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/20 2:55 下午
 */
package zookeeper

import (
	"core.bank/datamover/log"
	"github.com/spf13/cobra"
)


func NewZKCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "zookeeper",
		Short: "Realize data migration commands between different zookeepers",
		Long:  "Realize data migration commands between different zookeepers, support move from one to another",
		Run:   zkCommandFunc,
	}

	cmd.AddCommand(onlineCmd)

	return cmd
}

func zkCommandFunc(cmd *cobra.Command, args []string) {
	log.Logger.Info("zookeeper command")
}
