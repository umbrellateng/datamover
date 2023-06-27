# DataMover
#### 数据迁移，目前只支持同构 mysql 之间的数据迁移

#### 编译
`go build -o mover`

#### 默认设置

###### user: root  
###### host: 127.0.0.1  
###### port: 3306

#### 命令行缩写
###### --user: -u
###### --password: -p
###### --host: -h
###### --port: -P（大写）
###### --database: -d
###### --output: -o
###### --all-databases: -a
###### --file： -f
###### --restore： -r


#### 1、全库dump并输出到sql文件
###### 用法：
`./mover -u [username] -p [password] -h [mysql host] --port [mysql port] -a --output [输出的 sql 文件路径]`

###### 如果是在本机执行，可以使用默认设置，示例如下：
`./mover -p root -a -o all.sql` 

#### 2、指定某个库的dump并输出到sql文件
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
      
