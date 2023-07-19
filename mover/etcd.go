/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/4 5:31 下午
 */
package mover

import (
	"context"
	"fmt"
	"time"

	"core.bank/datamover/log"
	etcdclient "go.etcd.io/etcd/client/v3"
)

type Etcd struct {
	BaseInfo
}

func NewEtcd(user, pwd, h, p string) *Etcd {
	info := BaseInfo{
		username: user,
		password: p,
		host: h,
		port: p,
	}

	return &Etcd{
		BaseInfo:info,
	}
}

func NewEtcdFromBaseInfo(info BaseInfo) *Etcd {
	return &Etcd{
		BaseInfo: info,
	}
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

func (e *Etcd) UrlString() string {
	return "http://" + e.host + ":" + e.port
}

type EtcdCluster struct {
	entries []*Etcd
	username string
	password string
}



func NewEtcdCluster() *EtcdCluster {
	return &EtcdCluster{
		entries: make([]*Etcd, 0),
	}
}

func NewEtcdClusterFromBaseInfos(infos []BaseInfo) *EtcdCluster {
	var etcd *Etcd
	ec := NewEtcdCluster()
	for _, info := range infos {
		etcd = NewEtcdFromBaseInfo(info)
		ec.addEntry(etcd)
	}
	if len(infos) != 0 {
		ec.username = infos[0].username
		ec.password = infos[0].password
	}
	return ec
}

func (ec *EtcdCluster) addEntry(etcd *Etcd) {
	ec.entries = append(ec.entries, etcd)
}

func (ec *EtcdCluster) EndPoints() []string {
	var endPoints []string
	for _, entry := range ec.entries {
		endPoints = append(endPoints, entry.UrlString())
	}

	return endPoints
}

func (ec *EtcdCluster) Empty() bool {
	return len(ec.entries) == 0
}

func (ec *EtcdCluster) Username() string {
	return ec.username
}

func (ec *EtcdCluster) Password() string {
	return ec.password
}


func (ec *EtcdCluster) DumpToFile(filePath string) error {
	panic("implement me")
}

func (ec *EtcdCluster) DumpToDirectory(dirPath string) error {
	panic("implement me")
}

func (ec *EtcdCluster) RestoreFromFile(filePath string) error {
	panic("implement me")
}

func (ec *EtcdCluster) RestoreFromDirectory(dirPath string) error {
	panic("implement me")
}

func (ec *EtcdCluster) MoveOnline(infos []BaseInfo) error {
	if ec.Empty() {
		return fmt.Errorf("source etcd cluster is empty")
	}

	if len(infos) == 0 {
		return fmt.Errorf("target etcd cluster info is emtpy")
	}
	// 创建一个源 etcd 的客户端
	srcClient, err := etcdclient.New(etcdclient.Config{
		Endpoints: ec.EndPoints(),
		Username:  ec.entries[0].username,
		Password:  ec.entries[0].password,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("new source etcd client error: %s", err.Error())
	}
	defer srcClient.Close()

	targetEtcdCluster := NewEtcdClusterFromBaseInfos(infos)

	// 创建一个目标 etcd 的客户端
	dstClient, err := etcdclient.New(etcdclient.Config{
		Endpoints: targetEtcdCluster.EndPoints(),
		Username:  targetEtcdCluster.Username(),
		Password:  targetEtcdCluster.Password(),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("new target etcd client error: %s", err.Error())
	}
	defer dstClient.Close()

	// 创建一个上下文，用于控制超时或取消操作
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// 使用源客户端获取所有的键值对
	resp, err := srcClient.Get(ctx, "", etcdclient.WithPrefix())
	if err != nil {
		return fmt.Errorf("get source etcd cluster all key value error: ", err)
	}

	// 遍历所有的键值对，使用目标客户端将它们写入目标集群
	for _, kv := range resp.Kvs {
		key := string(kv.Key)
		value := string(kv.Value)
		fmt.Printf("Migrating key: %s, value: %s\n", key, value)
		_, err := dstClient.Put(ctx, key, value)
		if err != nil {
			log.Logger.Error("put key value  %s : %s into target etcd cluster error: ", key, value)
			continue
		}
	}

	log.Logger.Info("etcd cluster migration completed!")

	return nil
}
