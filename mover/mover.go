/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/4 5:17 下午
 */
package mover

type Mover interface {
	DumpToFile(filePath string) error
	DumpToDirectory(dirPath string) error
	RestoreFromFile(filePath string) error
	RestoreFromDirectory(dirPath string) error
	MoveOnline() error
}

type BaseInfo struct {
	username string
	password string
	host string
	port string
}

