/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/20 4:01 下午
 */
package kafka

import (
	"core.bank/datamover/log"
	"github.com/spf13/cobra"
)

func NewKafkaCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "kafka",
		Short: "Realize data migration commands between kafkas",
		Long:  "Realize data migration commands between kafkas, support move from one target another",
		Run:   kafkaCommandFunc,
	}

	cmd.AddCommand(onlineCmd)

	return cmd
}

func kafkaCommandFunc(cmd *cobra.Command, args []string) {
	log.Logger.Info("kafka command")
}
