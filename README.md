# DataMover
#### 数据迁移，目前只支持同构 mysql 之间的数据迁移。数据迁移过程中支持多线程模式，默认的多线程数是16，用flag --thread or -t 来开启多线程

#### 一、前置条件
+ 1、go 1.16 以上（包括1.16） 
+ 2、系统中已经安装好了mysql, mysqldump 存在于环境变量 $PATH 中

#### 二、编译生成 datamover
`go build`

#### 三、命令行flag解释
| 标志全称 | 标志简称 | 标志类型 | 默认值 | 适用子命令域 | 解释说明 |
| :---:  | :---:   | :---:  | :---: | :--- | :--- |
| --from | -f| string | "root:root@tcp(localhost：3306)" | datamover mysql (dump or online) | source database 连接串 |
| --target | -t | string | "root:root@tcp(localhost：3306)" | datamover mysql (restore or online )| target database 连接串 |
| --thread | -T | bool | false | datamover mysql (dump or restore) | 是否开启多线程模式| 
| --databases | -d | string | "" | datamover mysql (全部子命令) | mysql 数据库名称 | 
| --output | -o | string | "" | datamover mysql dump | 要输出的文件或目录, 可以省略 |
| --input | -i | string | "" | datamover mysql restore | 数据库恢复所需要的输入的文件或目录 |
| --all-databases | -a | bool | false | datamover mysql (dump or online) | mysql 全部的数据库，系统除外 |

#### 四、mysql 数据迁移命令
##### 1、从源数据库导出sql文件, --output or -o 指明输出的文件目录，当不指明时，系统会自动生成相应的文件目录
###### 1）、单线程导出 sql文件，只能支持一个database的导出，用 flag --databases or -d 来指明具体的数据库名，用法如下： 
`./datamover mysql dump --from "root:root@tcp(localhost:3306)" -d gep (-o gep.sql)`
###### 2） 单线程导出全部数据库，除了系统数据库不导出，系统的数据库包括mysql, sys, performance_schema, information_schema
`./datamover mysql dump --from "user:password@tcp(host:port)" -a （-o all-databases.sql）`

上述命令中，没有用 -o 指定的输出文件，系统会默认保存在 all-databases.sql 文件中  
###### 3）、多线程成导出 sql 文件到目录，可以支持多个数据库的导出，多个每次都用 -d 指明, 多线程模式下一定要加上 -T flag， 用法如下：
`./datamover mysql dump --from "root:root@tcp(localhost:3306)" -d gep -d exer -d safe (-o gep_exer_safe)` -T
###### 4）、多线程成导出全部数据库，用法如下；
`./datamover mysql dump --from "user:password@tcp(host:port)" -a -T (-o all-databases)`
###### 2、通过sql文件或目录迁入到指定数据库
###### 1）、单线程导入sql文件，用法如下：
`./datamover mysql restore --to "user:password@tcp(host:port)" -i gep.sql`

注意：该用法可以修改数据库，把要 update or drop 数据的时候，可以写成 xxx.sql，然后用上面的命令执行即可，输入的文件改成该 xxx.sql
###### 2）、多线程导入 sql 文件所在目录，用法如下，一定要加上多线程标志 --thread or -T ：
`./datamover mysql restore --to "user:password@tcp(host:port)" -i gep_exer_safe ` -T（--thread）
##### 3、在线迁移，默认就是多线程模式，不需要用 --thread or -T 来表示，支持用 -d 表示多个和 -a 所有的数据库
`./datamover mysql online --from "user1:password1@tcp(host1:port1)" --target "user2:password2@tcp(host2:port2)" -d exer -d safe ...`

`./datamover mysql online -f "user1:password1@tcp(host1:port1)" -t "user2:password2@tcp(host2:port2)" -a  ` 

