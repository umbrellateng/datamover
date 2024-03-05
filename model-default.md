## 部署模板

#### 1、基本信息
+ ServiceName:                  服务名称            
+ Replicas:                     副本数              
+ Cpu:                          CPU 核数            
+ Memory:                       内存

#### 2、网络信息
+ CniType:                      网络类型
+ IpSetName:                    IpSet 名称
+ Protocol:                     网络传输协议
+ ContainerPort:                容器端口号
+ NodePort:                     节点端口号

#### 3、环境变量列表
+ Envs:                         map[string]string

#### 4、配置信息列表
+ ConfigGroupName:              配置组名称
+ ConfigFileNames:              配置文件名称列表
+ MountDir:                     挂载路径

#### 5、卷存储信息列表
+ VolumnName:                   卷名称
+ VolumnType:                   卷类型 (cbs | cfs | local)
+ VolumnSize:                   卷大小
+ MountDir                      挂载路径      

