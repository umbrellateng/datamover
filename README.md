# DataMover
#### 数据迁移，目前只支持同构 mysql 之间的数据迁移

#### 编译
`go build -o mover`

#### 命令行flag解释

| 标志全称 | 标志简称 | 标志类型 | 默认值 | 解释说明 |
| :---:  | :---:   | :---:  | :---: | :--- |
| --user | -u | string | root | mysql 用户名称 |
| --password | -p | string | "" | mysql 用户密码 |
| --host | -h | string | "127.0.0.1" | mysql ip 地址 |
| --port | -P（大写）| string | "3306" | mysql port 端口号 |
| --database | -d | string | "" | mysql 数据库名称 | 
| --output | -o | string | default.sql | 要输出的文件或目录（多线程模式下输出的是目录）|
| --file | -f | string | "" |  要导入的sql文件名称，用于单线程情况下的某个数据库恢复 | 
| --all-databases | -a | bool | false | mysql 全部的数据库，infomation_schema、sys、mysql、performance_schema 除外 |
| --restore | -r | bool | false | 数据库恢复标志，命令行中不出现，那就意味着是 dump | 
| --thread | -t | bool | false | 是否开启多线程模式，命令行中出现，则开启多线程模式，默认开启的线程数是 16| 

#### 1、全库dump并输出到sql文件
###### 用法：
`./mover -u [username] -p [password] -h [mysql host] --port [mysql port] -a --output [输出的 sql 文件路径]`

###### 如果是在本机执行，可以使用默认设置，示例如下：
`./mover -p root -a -o all.sql` 

#### 2、指定某个库的dump并输出到sql文件，目前只支持单个库的操作
###### 用法：
`./mover -u [username] -p [password] -h [mysql host] --port [mysql port] --database [数据库名称] --output [输出的 sql 文件路径]`
###### 如果是在本机执行，可以使用默认设置，示例如下：
`./mover -p root -d exer -o exer.sql`
###### 其中，一定要保证所连机器的 mysql 数据库中有 exer 这个database，不然会报错

#### 3、从某个具体的sql文件进行恢复
###### 用法：
`./mover -u [username] -p [password] -h [mysql host] --port [mysql port] -r --file [sql 文件路径]`
###### 如果是在本机执行，可以使用默认设置，示例如下：
`./mover -p root -r -f exer.sql`        
      
