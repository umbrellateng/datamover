/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/20 3:44 下午
 */
package redis

import (
	"core.bank/datamover/log"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
)

var (
	from string
	to string
	db int
)

var moveCmd = &cobra.Command{
	Use: "move",
	Short: "move redis data from source cluster to the target cluster",
	Args: cobra.NoArgs,
	Run: moveCommandFunc,
}

func init() {
	moveCmd.Flags().StringVarP(&from, "from", "f", "", "source redis cluster url")
	moveCmd.Flags().StringVarP(&to, "to", "t", "", "target redis cluster url")
	moveCmd.Flags().IntVarP(&db, "db", "d", 0, "redis db number url")

	_ = moveCmd.MarkFlagRequired("from")
	_ = moveCmd.MarkFlagRequired("to")
}

func moveCommandFunc(cmd *cobra.Command, args []string) {

	err := redisMove(from, to, db)
	if err != nil {
		log.Logger.Error("redis data migration error: " + err.Error())
		return
	}

	fmt.Println()
	log.Logger.Info("redis data migration on success!")
}

func redisMove(sourceURL, targetURL string, db int) error {
	// 创建redis-cli命令对象
	cliCmd := exec.Command("redis-cli", "-u", sourceURL, "-n", fmt.Sprint(db), "keys", "*")
	// 执行redis-cli命令并获取所有键的列表
	keys, err := cliCmd.Output()
	if err != nil {
		return fmt.Errorf("redis cliCmd.Output error: " + err.Error())
	}

	// 创建xargs命令对象
	xargsCmd := exec.Command("xargs", "-I", "{}", "redis-cli", "-u", sourceURL, "-n", fmt.Sprint(db),
		"migrate", targetURL, "{}", fmt.Sprint(db), "10000", "COPY")
	// 创建一个读写管道
	stdin, err := xargsCmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("redis xargsCmd.StdinPipe error: " + err.Error())
	}
	stdout, err := xargsCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("redis xargsCmd.StdoutPipe error: " + err.Error())
	}
	// 将键的列表写入到管道中
	go func() {
		defer stdin.Close()
		_, _ = io.WriteString(stdin, string(keys))
	}()

	// 从管道中读取输出信息，并打印到标准输出中
	go func() {
		defer stdout.Close()
		_, _ = io.Copy(os.Stdout, stdout)
	}()
	// 执行xargs命令并等待完成
	err = xargsCmd.Run()
	if err != nil {
		return fmt.Errorf("xargsCmd.Run error: " + err.Error())
	}

	return nil
}
