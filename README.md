# DataMover
#### 数据迁移，目前只支持同构 mysql 之间的数据迁移。数据迁移过程中支持多线程模式，默认的多线程数是16，用flag --thread or -t 来开启多线程

#### 一、前置条件
+ 1、go 1.17.x version 以上（包括1.17） 
+ 2、系统中已经安装好了mysql, mysqldump 存在于环境变量 $PATH 中
#### 二、编译生成 datamover
`go build`

#### 三、命令行flag解释

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

#### 四、数据库 dump
数据库 dump 目前支持 mysql 中单个 database、多个database、及全量数据库dump（不包括infomation_schema、sys、mysql、performance_schema）
,其中，一定要保证所连机器的 mysql 数据库中有 exer 这个database，不然会报错
###### **1、单个 database 的 dump 示例**
`./datamover --user root --password root --host 127.0.0.1 --port 3306 --database gep` 
###### 输出如下：
` 2023/06/30 11:33:27.652477 main.go:229:         [INFO]         dump the certain database gep into file gep.sql ...
 
  2023/06/30 11:33:27.860979 main.go:251:         [INFO]         Success!
`
###### 如果没有指定 --output, 则默认输出的文件以 database name 命令，本命令行如果不指定output，则输出的是 gep.sql 文件， 命令行还可以简写如下： 
`./datamover -u root -p root -h 127.0.0.1 -P 3306 -d gep -o output.sql`
###### 输出如下：
` 2023/06/30 11:32:27.748643 main.go:229:         [INFO]         dump the certain database gep into file output.sql ...
 
  2023/06/30 11:32:28.105120 main.go:251:         [INFO]         Success!
`
###### 因为是dump同样的数据库 gep，上两个命令行输出的gep.sql 和 output.sql 内容完全一样
###### user 默认 root， password 默认 root, host 默认 127.0.0.1，port 默认 3306， 如果执行环境中一致，那么命令行中可省略这些 flag:
`./datamover -d gep`
###### 输出如下： 
` 2023/06/30 11:41:11.568631 main.go:229:         [INFO]         dump the certain database gep into file gep.sql ...
 
  2023/06/30 11:41:11.845842 main.go:251:         [INFO]         Success!`
###### 如果用多线程模式，只需要在单线程模式中的命令行中加 flag --thread 或者 -t 即可，但是默认输出的是目录，而非sql文件
`./datamover -d gep -t`
###### 输出如下：
` 2023/06/30 11:44:42.370203 main.go:172:         [INFO]         dump database exer into exer directory... 
 
  2023/06/30 11:44:42.389626 dumper.go:37:        [INFO]         dumping.database[exer].schema...
  
  2023/06/30 11:44:42.393374 dumper.go:47:        [INFO]         dumping.table[exer.t_person].schema...
  
  2023/06/30 11:44:42.393458 dumper.go:239:       [INFO]         dumping.table[exer.t_person].datas.thread[1]...
  
  2023/06/30 11:44:42.394382 dumper.go:47:        [INFO]         dumping.table[exer.t_role].schema...
  
  2023/06/30 11:44:42.394535 dumper.go:239:       [INFO]         dumping.table[exer.t_role].datas.thread[2]...
  
  2023/06/30 11:44:42.395472 dumper.go:151:       [INFO]         dumping.table[exer.t_role].done.allrows[1].allbytes[0MB].thread[2]...
  
  2023/06/30 11:44:42.395498 dumper.go:151:       [INFO]         dumping.table[exer.t_person].done.allrows[7].allbytes[0MB].thread[1]...
  
  2023/06/30 11:44:42.395490 dumper.go:241:       [INFO]         dumping.table[exer.t_role].datas.thread[2].done...
  
  2023/06/30 11:44:42.395513 dumper.go:241:       [INFO]         dumping.table[exer.t_person].datas.thread[1].done...
  
  2023/06/30 11:44:42.396229 dumper.go:260:       [INFO]         dumping.all.done.cost[0.01sec].allrows[8].allbytes[151].rate[0.00MB/s]
  
  2023/06/30 11:44:42.396884 main.go:251:         [INFO]         Success!
` 

###### **2、多个数据库的dump示例**
###### 多个数据库的 dump 只能在多线程模式下进行，单线程模式下不支持，所以命令行中必须加  --thread or -t
###### 比如 mysql 数据库中有 gep、safe、exer， 多个数据库用 -d 进行标识， 则 dump 这些数据库命令行如下：
`./datamover -u root -p root -h 127.0.0.1 -P 3306 -d gep -d safe -d exer -t`
###### 输出如下
`2023/06/30 15:00:54.560199 main.go:172:         [INFO]         dump database gep,safe,exer into gep_safe_exer directory... 
 
  2023/06/30 15:00:54.639412 dumper.go:37:        [INFO]         dumping.database[gep].schema...
  
  2023/06/30 15:00:54.641738 dumper.go:37:        [INFO]         dumping.database[safe].schema...
  
  2023/06/30 15:00:54.643623 dumper.go:37:        [INFO]         dumping.database[exer].schema...
  
  2023/06/30 15:00:54.669894 dumper.go:47:        [INFO]         dumping.table[gep.c2c_finance_coin].schema...
  
  2023/06/30 15:00:54.669984 dumper.go:239:       [INFO]         dumping.table[gep.c2c_finance_coin].datas.thread[1]...
  
  2023/06/30 15:00:54.673816 dumper.go:47:        [INFO]         dumping.table[gep.center].schema...
  
  2023/06/30 15:00:54.673918 dumper.go:239:       [INFO]         dumping.table[gep.center].datas.thread[2]...
  
  ......
`

###### **3、所有数据库的dump示例**
###### 注意： 所有数据库不包括infomation_schema、sys、mysql、performance_schema 
###### 单线程模式和多线程模式下都支持所有数据库的dump, 单线程默认输出到 all-databases.sql 文件中，多线程模式默认输出到 all-databases 目录中
###### 用 --all-databases or -a 来标识所有数据库
###### 单线程模式下 dump 所有数据库命令行如下：
`./datamover -u root -p root -h 127.0.0.1 -P 3306 -a`
###### 输出如下：
`2023/06/30 15:11:04.992774 main.go:217:         [INFO]         dump all databases into file all-databases.sql ...
 
  2023/06/30 15:11:05.519289 main.go:251:         [INFO]         Success!
`
###### 多线程模式下 dump 所有数据库只需要在单线程命令行后面加  --thread or -t 标志：
`./datamover -u root -p root -h 127.0.0.1 -P 3306 -a -t`
###### 输出如下：
`2023/06/30 15:13:19.745027 main.go:152:         [INFO]         dump all databases into file  at multi-threaded mode...
 
  2023/06/30 15:13:19.784618 dumper.go:37:        [INFO]         dumping.database[db1].schema...
  
  2023/06/30 15:13:19.785555 dumper.go:37:        [INFO]         dumping.database[exer].schema...
  
  2023/06/30 15:13:19.785992 dumper.go:37:        [INFO]         dumping.database[fileserver].schema...
  
  2023/06/30 15:13:19.786334 dumper.go:37:        [INFO]         dumping.database[gep].schema...
  
  2023/06/30 15:13:19.786994 dumper.go:37:        [INFO]         dumping.database[safe].schema...
  
  2023/06/30 15:13:19.787439 dumper.go:37:        [INFO]         dumping.database[seckill].schema...
  
  2023/06/30 15:13:19.796049 dumper.go:47:        [INFO]         dumping.table[db1.user_infos].schema...
  
  2023/06/30 15:13:19.796119 dumper.go:239:       [INFO]         dumping.table[db1.user_infos].datas.thread[1]...
  
  ......
`
      
#### 五、数据库 restore
###### 数据库恢复命令行中要加上标志 --restore or -r，并且需要 --input or -i 指定输入文件或目录，单线程输入指定的是sql文件，多线程模式下输入指定的是目录
###### 1、单线程数据库恢复命令行示例：
`./datamover -u root -p root -h 127.0.0.1 -P 3306 -r -i gep.sql`
###### 输出如下
` 2023/06/30 15:19:00.029118 main.go:91:          [ERROR]        the input gep.sql is not a directory

  2023/06/30 15:19:00.031260 main.go:197:         [INFO]         restore the database from the certain sql file: gep.sql ...
 
  2023/06/30 15:19:01.244988 main.go:251:         [INFO]         Success!`
  
###### 2、多线程数据库恢复命令行示例：
`./datamover -u root -p root -h 127.0.0.1 -P 3306 -r -i gep -t`      
###### 上述命令行中的 gep 是个目录，上述第四部分内容多线程命令行生成出来的，输出如下：
`2023/06/30 15:23:29.844539 main.go:141:         [INFO]         restore databases from the directory: gep ...
 
  2023/06/30 15:23:29.874821 loader.go:77:        [INFO]         restoring.database[gep]

  2023/06/30 15:23:29.874852 loader.go:90:        [INFO]         working.table[gep.c2c_finance_coin]

  2023/06/30 15:23:29.908439 loader.go:114:       [INFO]         restoring.schema[gep.c2c_finance_coin]

  2023/06/30 15:23:29.908466 loader.go:90:        [INFO]         working.table[gep.center]

  2023/06/30 15:23:29.943051 loader.go:114:       [INFO]         restoring.schema[gep.center]

  2023/06/30 15:23:29.943073 loader.go:90:        [INFO]         working.table[gep.deposit]
  
  ......
`