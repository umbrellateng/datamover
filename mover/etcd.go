/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/4 5:31 下午
 */
package mover

type Etcd struct {
	BaseInfo
}

func (e *Etcd) Username() string {
	return e.username
}

func (e *Etcd) Password() string {
	return e.password
}

func (e *Etcd) Host() string {
	return e.host
}

func (e *Etcd) Port() string {
	return e.port
}


func (e *Etcd) DumpToFile() error {
	panic("implement me")
}

func (e *Etcd) DumpToDirectory() error {
	panic("implement me")
}

func (e *Etcd) RestoreFromFile() error {
	panic("implement me")
}

func (e *Etcd) RestoreFromDirectory() error {
	panic("implement me")
}

func (e *Etcd) MoveOnline() error {
	panic("implement me")
}
