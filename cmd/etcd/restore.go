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
	//initialClusterToken string
	name string

)

var restoreCmd = &cobra.Command{
	Use:   "restore [db_file_name]",
	Short: "etcd snapshot restore command",
	Args: cobra.ExactArgs(1),
	Run:   restoreCommandFunc,
}

func init() {
	restoreCmd.Flags().StringVar(&dataDir, "data-dir", "", "path to the data directory")
	restoreCmd.Flags().StringVar(&initialAdvertisePeerUrls, "initial-advertise-peer-urls","", "list of this member's peer URLs to advertise to the rest of the cluster")
	restoreCmd.Flags().StringVar(&initialCluster, "initial-cluster", "", "Initial cluster configuration for restore bootstrap")
	//restoreCmd.Flags().StringVar(&initialClusterToken, "initial-cluster-token", "initial cluster token for the etcd cluster during restore bootstrap", "path to the data directory")
	restoreCmd.Flags().StringVar(&name, "name", "", "human-readable name for this member")

	_ = restoreCmd.MarkFlagRequired("data-dir")
}

func restoreCommandFunc(cmd *cobra.Command, args []string) {
	err := restoreFromDBFile(args[0])
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
	var cmd *exec.Cmd
	if useTLS() && len(name) != 0 {
		log.Logger.Info("use tls and name")
		if err := checkURL(); err != nil {
			return err
		}
		cmd = exec.Command("etcdctl", "--endpoints", endpoints, "--cacert", cacert, "--cert", cert, "--key", key,
			"snapshot", "restore", dbFile, "--name", name, "--initial-cluster", initialCluster, "--initial-advertise-peer-urls",
			initialAdvertisePeerUrls, "--data-dir", dataDir)
	} else if len(name) != 0 && len(endpoints) != 0{
		log.Logger.Info("use name and endpoints")
		if err := checkAddrs(); err != nil {
			return err
		}
		cmd = exec.Command("etcdctl", "snapshot", "restore", dbFile, "--endpoints", endpoints, "--name", name, "--initial-cluster",
			initialCluster, "--initial-advertise-peer-urls", initialAdvertisePeerUrls, "--data-dir", dataDir)
	} else if len(name) != 0 {
		log.Logger.Info("use name only")
		if err := checkAddrs(); err != nil {
			return err
		}
		cmd = exec.Command("etcdctl", "snapshot", "restore", dbFile, "--name", name, "--initial-cluster",
			initialCluster, "--initial-advertise-peer-urls", initialAdvertisePeerUrls, "--data-dir", dataDir)
	} else if useTLS() {
		log.Logger.Info("use tls only")
		if err := checkURL(); err != nil {
			return err
		}
		cmd = exec.Command("etcdctl", "--endpoints", endpoints, "--cacert", cacert, "--cert", cert,
			"--key", key, "snapshot", "restore", dbFile, "--data-dir", dataDir)
	} else if len(endpoints) != 0 {
		log.Logger.Info("use endpoints only")
		cmd = exec.Command("etcdctl", "snapshot", "restore", dbFile, "--endpoints", endpoints, "--data-dir", dataDir)
	} else {
		log.Logger.Info("use nothing!")
		cmd = exec.Command("etcdctl", "snapshot", "restore", dbFile, "--data-dir", dataDir)
	}

	// 执行命令并获取输出
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("exec etcdctl command error: " + err.Error())
	}
	log.Logger.Info(string(out))

	return nil
}

func checkURL() error {
	if len(cacert) == 0 {
		return fmt.Errorf("cacert is empty, please specify it with the flag --cacert")
	}
	if len(cert) == 0 {
		return fmt.Errorf("cert is empty, please specify it with flag --cert")
	}
	if len(key) == 0 {
		return fmt.Errorf("key is empty, please specify it with flag --key")
	}

	return nil
}

func checkAddrs() error {
	if len(initialCluster) == 0 {
		return fmt.Errorf("initial-cluster is empty, please specify it with flag --initial-cluster")
	}
	if len(initialAdvertisePeerUrls) == 0 {
		return fmt.Errorf("initial-advertise-peer-urls is empty, please specify it with flag --initial-advertise-peer-urls")
	}
	return nil
}
