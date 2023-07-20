/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/20 11:22 上午
 */
package etcd

import (
	"core.bank/datamover/log"
	"core.bank/datamover/utils"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var (
	dataDir string
	initialAdvertisePeerUrls string
	initialCluster string
	initialClusterToken string
	name string
	input string

)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "etcd snapshot restore command",
	Args: cobra.NoArgs,
	Run:   restoreCommandFunc,
}

func init() {
	restoreCmd.Flags().StringVar(&dataDir, "data-dir", "", "path to the data directory")
	restoreCmd.Flags().StringVar(&initialAdvertisePeerUrls, "initial-advertise-peer-urls","", "list of this member's peer URLs to advertise to the rest of the cluster")
	restoreCmd.Flags().StringVar(&initialCluster, "initial-cluster", "", "Initial cluster configuration for restore bootstrap")
	restoreCmd.Flags().StringVar(&initialClusterToken, "initial-cluster-token", "initial cluster token for the etcd cluster during restore bootstrap", "path to the data directory")
	restoreCmd.Flags().StringVar(&name, "name", "", "human-readable name for this member")
	restoreCmd.Flags().StringVarP(&input, "input", "i", "", "the snapshot db file for input")
}

func restoreCommandFunc(cmd *cobra.Command, args []string) {
	err := restoreFromDBFile(input)
	if err != nil {
		log.Logger.Error("restore etcd error: " + err.Error())
		return
	}

	fmt.Println()
	log.Logger.Info("etcd restore on success!")
}

func restoreFromDBFile(dbFile string) error {

	if len(dbFile) == 0 {
		return fmt.Errorf("please specify the snapshot db file")
	}

	if utils.IsDirectory(dbFile) {
		return fmt.Errorf("not directory, should be specify the snapshot db file")
	}

	// 设置ETCDCTL_API环境变量为3
	_ =os.Setenv("ETCDCTL_API", "3")

	// 构造etcdctl命令
	cmd := exec.Command("etcdctl", "--endpoints", endpoints, "--cacert", cacert, "--cert", cert, "--key", key,
		"snapshot", "restore", dbFile, "--name", name, "--initial-cluster", initialCluster, "--initial-advertise-peer-urls",
		initialAdvertisePeerUrls, "--initial-cluster-token", initialClusterToken, "--data-dir", dataDir)

	// 执行命令并获取输出
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("exec etcdctl command error: " + err.Error())
	}
	log.Logger.Info(string(out))

	return nil
}
