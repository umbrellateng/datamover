/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/20 3:44 下午
 */
package redis

import (
	"core.bank/datamover/log"
	"github.com/spf13/cobra"
)

func NewRedisCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "redis",
		Short: "Realize data migration commands between redis",
		Long:  "Realize data migration commands between redis, support move from one target another",
		Run:   redisCommandFunc,
	}

	cmd.AddCommand(saveCmd)
	cmd.AddCommand(onlineCmd)

	return cmd
}

func redisCommandFunc(cmd *cobra.Command, args []string) {
	log.Logger.Info("redis command")
}
