# DataMover
## 数据迁移工具，支持同构 mysql、etcd、redis、zookeeper、kafka 之间的迁移。

## 一、前置条件
+ 1、go 1.16 以上（包括1.16） 
+ 2、确保mysql、mysqldump、etcdctl 等在系统$PATH中

## 二、编译生成 datamover
+ 1、make or make build，编译出来的可执行文件 datamover 依赖于系统
+ 2、支持 mac 和 linux
+ 3、linux编译： make linux 
+ 4、mac编译： make mac
## 三、各个持久化存储的迁移
### （一）mysql 数据库之间的迁移
#### 1、用法示例
`./datamover mysql -h`
具体的flags解释如下；

| 标志全称 | 标志简称 | 标志类型 | 默认值 | 解释说明 |
| :---:  | :---:   | :---:  | :---: | :--- |
| --from | 无 | string | "root:root@tcp(localhost：3306)" | source database 连接串 |
| --target | 无 | string | "root:root@tcp(localhost：3306)" | target database 连接串 |
| --thread | -T | bool | false | 是否开启多线程模式| 
| --databases | -d | string | nil | mysql 数据库名称 | 
| --tables | -t | string | nil | 选择的数据表 tables |
| --output | -o | string | "" | 要输出的文件或目录, 可以省略 |
| --input | -i | string | "" | 数据库恢复所需要的输入的文件或目录 |
| --all-databases | -a | bool | false | mysql 全部的数据库，系统除外 |
| --without-create-database | -w | bool | false | 是否删除生成的创建数据库的sql文件（xxx-schema-create.sql）| 

#### 2、从源数据库导出sql文件, --output or -o 指明输出的文件目录，当不指明时，系统会自动生成相应的文件目录
###### 1）、单线程导出 sql文件，只能支持一个database的导出，用 flag --databases or -d 来指明具体的数据库名，用法如下： 
`./datamover mysql dump --from "root:root@tcp(localhost:3306)" -d gep (-o gep.sql)`
###### 2） 单线程导出全部数据库，除了系统数据库不导出，系统的数据库包括mysql, sys, performance_schema, information_schema
`./datamover mysql dump --from "user:password@tcp(host:port)" -a （-o all-databases.sql）`

上述命令中，没有用 -o 指定的输出文件，系统会默认保存在 all-databases.sql 文件中  
###### 3）、多线程成导出 sql 文件到目录，可以支持多个数据库的导出，多个每次都用 -d 指明, 多线程模式下一定要加上 -T flag， 用法如下：
`./datamover mysql dump --from "root:root@tcp(localhost:3306)" -d gep -d exer -d safe (-o gep_exer_safe)` -T
###### 4）、多线程成导出全部数据库，用法如下；
`./datamover mysql dump --from "user:password@tcp(host:port)" -a -T (-o all-databases)`

###### 5)、多线程模式下，还支持导出某个数据库的tables，可以支持一个或者多个，用 --tables 或者 -t 指明， 此种情况下只支持一个数据库

`./datamover mysql dump --from "user:password@tcp(host:port)" -d db_name -t tb_name1 -t tb_name2 -T `

###### 6)、在多线程模式下，如果只选择一个数据库的情况下，可以用 --without-create-database or -w 把生成的 create database sql 的文件删除，便于没有创建库表权限的用户迁移

`./datamover mysql dump --from "user:password@tcp(host:port)" -d db_name -t tb_name1 -t tb_name2 -T --without-create-database (-w)`

上述情况，如果不指定 --tables or -t， 则会输出 db_name 数据库下的所有 tables


#### 3、通过sql文件或目录迁入到指定数据库
###### 1）、单线程导入sql文件，用法如下：
`./datamover mysql restore --target "user:password@tcp(host:port)" -i gep.sql`

注意：该用法可以修改数据库，把要 update or drop 数据的时候，可以写成 xxx.sql，然后用上面的命令执行即可，输入的文件改成该 xxx.sql
###### 2）、多线程导入 sql 文件所在目录，用法如下，一定要加上多线程标志 --thread or -T ：
`./datamover mysql restore --target "user:password@tcp(host:port)" -i gep_exer_safe ` -T（--thread）

#### 4、在线迁移，默认就是多线程模式，不需要用 --thread or -T 来表示，支持用 -d 表示多个和 -a 所有的数据库
`./datamover mysql online --from "user1:password1@tcp(host1:port1)" --target "user2:password2@tcp(host2:port2)" -d exer -d safe ...`

`./datamover mysql online -from "user1:password1@tcp(host1:port1)" -target "user2:password2@tcp(host2:port2)" -a` 

`./datamover mysql online -from "user1:password1@tcp(host1:port1)" -target "user2:password2@tcp(host2:port2)" -d db_name -t tb_name1 -t tb_name2` 

`./datamover mysql online -from "user1:password1@tcp(host1:port1)" -target "user2:password2@tcp(host2:port2)" -d db_name -t tb_name1 -t tb_name2 -w` 

### (二) etcd 之间的迁移
+ etcd 之间的迁移，先通过命令行工具从源 etcd 导出 xxx.db 文件，然后再用命令行工具将 xxx.db 文件导入到另外一个 etcd 集群
+ etcd 的子命令包含了 save 和 restore，跟 etcdctl 保留一致
+ etcd 之间的数据迁移实现是通过 etcdctl 命令来实现的，如果大家对 etcdctl更熟悉，那就用 etcdctl 来进行迁移会更好。

`./datamover etcd -h`
#### 1、从源 etcd 中导出 xxx.db 文件，<db_file_name> 如果不填，则默认输出的是 etcd-snapshot-YYYY-MM-DD HH:mm:ss.db 文件
`./datamover etcd save <db_file_name> --endpoints http://host:port`

`./datamover etcd save etcd-node1.db --endpoinsts http://127.0.0.1:2379`

###### 如果是用到了 tls, 则命令行中还要明确 --cacert, --cert, --key 指明 tls 所需要的文件路径，例如：
`./datamover etcd save etcd-node2.db --cacert=/opt/etcd/ssl/ca.pem --cert=/opt/etcd/ssl/server.pem --key=/opt/etcd/ssl/server-key.pem --endpoints="https://192.168.1.61:2379`

#### 2、将 xxx.db 文件导入到新的 etcd 集群，此命令中要用 --data-dir 指明要导入的新etcd集群的数据目录，而且该数据目录必须为空
`./datamover etcd restore [db_file_name] --data-dir new-etcd`

`./datamover etcd restore etcd-node1.db --data-dir new-etcd-node1`

###### 当 etcd restore 命令行中要出现 --name 的时候，必须同时指明 --initial-cluster 和 initial-advertise-peer-urls 这两个标志位
`./datamover etcd restore etcd-node2.db --data-dir new-etcd-node2 --name node2 --initial-cluster node2=http://127.0.0.1:2380 --initial-advertise-peer-urls http://127.0.0.1:2380`

###### etcd restore 命令行中同时也可以带上 --endpoints，如下所示：
`./datamover etcd restore etcd-node1.db --data-dir new-etcd-node1 --endpoints http:127.0.0.1:2378`

`./datamover etcd restore etcd-node2.db --data-dir new-etcd-node2 --name node2 --initial-cluster node2=http://127.0.0.1:2380 --initial-advertise-peer-urls http://127.0.0.1:2380 --endpoints http://127.0.0.1:2378`

###### etcd restore 命令行中也支持 tls , 例如：
`./datamover etcd restore etcd-node2.db  --data-dir new-etcd-node2 --cacert=/opt/etcd/ssl/ca.pem --cert=/opt/etcd/ssl/server.pem --key=/opt/etcd/ssl/server-key.pem --endpoints="https://192.168.1.61:2379`

### (三)、redis 之间的迁移
+ 支持在线迁移
#### 1、在线迁移用法
`./datamover redis online --from [host1:port1] --target [host2:port2] --from-password <pwd1> --from-db <db1> --target-password <pwd2> --target-db <db2>`

以上命令中的flag，--from 和 --target 是必须的，其他的可以省略，默认为空或者是0 

`./datamover redis online --from 127.0.0.1:6379 --target 192.168.3.137:6379`

注意： 以上命令行中的 url 不要写成 http://127.0.0.1:6379，不需要带 http://, 不然会报错

#### 2、dump出rdb文件，通过重启redis加载该rdb文件进行迁移
`./datamover redis save [rdb_file_name] -u <redis-url> `
##### 例如
`./datamover redis save dump.db -u redis://localhost:6379`
##### 运行以上命令，会在当前目录中生成 dump.db 文件，然后将该 dump.db文件放在目标 redis 服务的数据目录中，然后重新启动目标 redis 服务即可完成数据的迁移。
##### 用法
`./datamover redis save --help`

`redis generates rdb snapshot files and outputs them to the specified directory`
 
 `Usage:`
 
   `datamover redis save [flags]`
 
 `Flags:`
 
   `-h, --help         help for save`
   
   `-u, --url string   redis server url (default "redis://127.0.0.1:6379")`
   
   
### (四)、zookeeper 之间的迁移 

#### 1、用法
`./datamover zookeeper online -h`

`move zookeeper data from source cluster to the target cluster`

`Usage:`

  `datamover zookeeper online [flags]`

`Flags:`

`  -f, --from string     source zookeeper cluster url`

`  -h, --help            help for online`

`  -t, --target string   target zookeeper cluster url`

#### 2、示例
`./datamover zookeeper online --from 127.0.0.1:2181 --target 192.168.34.165:2181`

`./datamover zookeeper online -f 127.0.0.1:2181 -t 192.168.34.165:2181`

### (五)、kafka 之间的迁移

#### 1、用法
`$ ./datamover kafka online -h` 

` move kafka data from source cluster target the target cluster`
 
 `Usage:`
 
 `  datamover kafka online [flags]`
 
 `Flags:`
 
 `  -f, --from string     source kafka cluster url`
 
 `  -h, --help            help for online `
 
 `  -t, --target string   target kafka cluster url`
 
 #### 2、示例
 `./datamover kafka online --from 127.0.0.1:9092 --target 192.168.34.165:9092`
 
 `./datamover kafka online -f 127.0.0.1:9092 -t 192.168.34.165:9092`
 