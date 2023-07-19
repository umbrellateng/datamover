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
	MoveOnline(info []BaseInfo) error
}

type BaseInfo struct {
	username string
	password string
	host string
	port string
}

func NewBaseInfo(u, pwd, h, p string) BaseInfo {
	return BaseInfo{
		username: u,
		password: pwd,
		host: h,
		port: p,
	}
}

