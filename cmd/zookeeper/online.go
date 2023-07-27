/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/20 2:55 下午
 */
package zookeeper

import (
	"core.bank/datamover/log"
	"fmt"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/spf13/cobra"
)

var (
	from   string
	target string
)

var onlineCmd = &cobra.Command{
	Use:   "online",
	Short: "move zookeeper data from source cluster to the target cluster",
	Args:  cobra.NoArgs,
	Run:   onlineCommandFunc,
}

func init() {
	onlineCmd.Flags().StringVarP(&from, "from", "f", "", "source zookeeper cluster url")
	onlineCmd.Flags().StringVarP(&target, "target", "t", "", "target zookeeper cluster url")

}

func onlineCommandFunc(cmd *cobra.Command, args []string) {

	if len(from) == 0 {
		log.Logger.Error("source zookeeper url is empty, please specify it with --from or -f like \"127.0.0.1:2181\"")
		return
	}

	if len(target) == 0 {
		log.Logger.Error("target zookeeper url is empty, please specify it with --target or -t like \"127.0.0.1:2181\"")
		return
	}
	err := zkOnlineMove(from, target)
	if err != nil {
		log.Logger.Error("zookeeper data migration error: " + err.Error())
		return
	}

	fmt.Println()
	log.Logger.Info("zookeeper data migration on success!")
}

func zkOnlineMove(source, target string) error {
	// 连接源集群
	sourceConn, _, err := zk.Connect([]string{source}, 15 * time.Second)
	if err != nil {
		return fmt.Errorf("connect source zookeeper cluster error: " + err.Error())
	}
	defer sourceConn.Close()

	// 连接目标集群
	targetConn, _, err := zk.Connect([]string{target}, 15 * time.Second)
	if err != nil {
		return fmt.Errorf("connect target zookeeper cluster error: " + err.Error())
	}
	defer targetConn.Close()

	// 递归遍历源集群的根节点，获取所有子节点的路径和数据
	var paths []string
	var data [][]byte
	err = walk(sourceConn, "/", &paths, &data)
	if err != nil {
		return err
	}

	log.Logger.Info("path length: %d", len(paths))
	log.Logger.Info("datas length: %d", len(data))
	// 在目标集群创建相同的节点路径和数据
	for i, path := range paths {
		log.Logger.Info("path: " + path)
		log.Logger.Info("data: %v", data[i])
		_, err := targetConn.Create(path, data[i], 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return fmt.Errorf("targetConn.Create error: " + err.Error())
		}
	}

	return nil
}

// walk 遍历给定节点及其所有子节点，将路径和数据追加到切片中
func walk(conn *zk.Conn, path string, paths *[]string, data *[][]byte) error {
	children, _, err := conn.Children(path)
	if err != nil {
		return fmt.Errorf("conn.Children error: " + err.Error())
	}

	if path != "/" { // 跳过根节点，因为无法创建
		log.Logger.Info("Found node:" + path)
		datum, _, err := conn.Get(path)
		if err != nil {
			return fmt.Errorf("conn.Get error: " + err.Error())
		}
		log.Logger.Info("Node data: %s", string(datum))
		fmt.Println()

		// 将路径和数据追加到切片中
		*paths = append(*paths, path)
		*data = append(*data, datum)
	}

	for _, child := range children {
		childPath := path + "/" + child
		if path == "/" {
			childPath = "/" + child
		}

		err = walk(conn, childPath, paths, data)
		if err != nil {
			return err
		}
	}

	return nil
}