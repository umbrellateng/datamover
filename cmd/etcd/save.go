/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/20 11:22 上午
 */
package etcd

import (
	"core.bank/datamover/log"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"time"
)

var (
	endpoints string
	output string
)

var saveCmd = &cobra.Command{
	Use: "save",
	Short: "etcd snapshot save command",
	Args: cobra.NoArgs,
	Run: saveCommandFunc,
}

func init() {
	saveCmd.Flags().StringVarP(&output, "output", "o", "", "the location file to save the etcd snapshot")
}

func saveCommandFunc(cmd *cobra.Command, args []string) {

	err := saveSnapShot(output)
	if err != nil {
		log.Logger.Error("etcd snapshot save error: " + err.Error())
		return
	}
	fmt.Println()
	log.Logger.Info("etcd snapshot save on success!")

}

func saveSnapShot(outputPath string) error {
	// 设置ETCDCTL_API环境变量为3
	_ = os.Setenv("ETCDCTL_API", "3")

	// 生成快照文件的路径和名称
	if len(outputPath) == 0 {
		outputPath = fmt.Sprintf("etcd-snapshot-%s.db", time.Now().Format("2006-01-02 15:04:05"))
	}

	// 构造etcdctl命令
	execCmd := exec.Command("etcdctl", "--endpoints", endpoints, "--cacert", cacert, "--cert", cert, "--key", key, "snapshot", "save", outputPath)

	// 执行命令并获取输出
	out, err := execCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("exec etcdctl command error: " + err.Error())
	}
	log.Logger.Info(string(out))

	return nil
}