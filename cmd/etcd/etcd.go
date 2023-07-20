/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/20 10:45 上午
 */
package etcd

import (
	"core.bank/datamover/log"
	"github.com/spf13/cobra"
)

var (
	cacert string
	cert string
	key string
)


func NewETCDCommand() *cobra.Command {
	// mysqlCmd represents the mysql command
	cmd := &cobra.Command{
		Use:   "etcd",
		Short: "Realize data migration commands between isomorphic etcd",
		Long:  "Realize data migration commands between isomorphic etcd, support save, restore and online move",
		Run:   etcdCommandFunc,
	}

	cmd.PersistentFlags().StringVar(&cacert, "cacert", "","the cacert path of the etcd endpoints")
	cmd.PersistentFlags().StringVar(&cert, "cert", "","the cert path of the etcd endpoints")
	cmd.PersistentFlags().StringVar(&key, "key", "","the key path of the etcd endpoints")
	cmd.PersistentFlags().StringVar(&endpoints,"endpoints", "", "the endpoints of the etcd cluster")

	cmd.AddCommand(saveCmd)
	cmd.AddCommand(restoreCmd)
	//cmd.AddCommand(onlineCmd)

	return cmd
}

func etcdCommandFunc(cmd *cobra.Command, args []string) {
	log.Logger.Info("etcd command")
}
