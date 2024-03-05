/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/19 3:25 下午
 */
package mysql

import (
	"core.bank/datamover/log"
	"core.bank/datamover/utils"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xelabs/go-mydumper/common"
	"github.com/xelabs/go-mydumper/config"
	querypb "github.com/xelabs/go-mysqlstack/sqlparser/depends/query"
	sqlcommon "github.com/xelabs/go-mysqlstack/sqlparser/depends/common"
	"github.com/xelabs/go-mysqlstack/xlog"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	output string
	separator string
	databases []string
	tables []string

	all bool
	withoutCreateDatabase bool
)

var dumpCmd = &cobra.Command{
	Use: "dump",
	Short: "dump data from mysql database",
	Run: dumpCommandFunc,
	Args: cobra.NoArgs,
}

func init() {
	dumpCmd.Flags().StringArrayVarP(&databases, "databases", "d", nil, "the dump databases of mysql")
	dumpCmd.Flags().StringArrayVarP(&tables, "tables", "t", nil, "the table name of some database")
	dumpCmd.Flags().StringVarP(&output, "output", "o", "", "the location that save the dump file or directory ")
	dumpCmd.Flags().StringVarP(&separator, "separator", "s", "|@|", "the separator of the database fields values")

	dumpCmd.Flags().BoolVarP(&all, "all-databases", "a", false, "all mysql databases except(mysql|sys|performance_schema|information_schema)")
	dumpCmd.Flags().BoolVarP(&withoutCreateDatabase, "without-create-database", "w", false, "if true the create-database.sql will be removed from the output directory")
}

func dumpCommandFunc(cmd *cobra.Command, args []string) {
	username, password, host, port, err := utils.ParseDBStringWithoutDB(from)
	if err != nil {
		log.Logger.Error("parse mysql connection string error: " + err.Error())
		return
	}

	if thread {
		outputDir, err := dumpToDirectory(username, password, host, port, output, separator)
		if err != nil {
			log.Logger.Error("dump mysql databases in multi-threaded mode error: " + err.Error())
			return
		}

		if withoutCreateDatabase && len(databases) == 1{
			fileName := filepath.Join(outputDir, fmt.Sprintf("%s-schema-create.sql", databases[0]))
			err := os.Remove(fileName)
			log.Logger.Info("remove the sql file: %s", fileName)
			if err != nil {
				log.Logger.Warning("remove the %s file error: %s", fileName, err.Error())
			}
		}

	} else {
		err := dumpToSqlFile(username, password, host, port, output)
		if err != nil {
			log.Logger.Error("dump mysql database in single-threaded mode error: " + err.Error())
			return
		}
	}
	fmt.Println()
	log.Logger.Info("dump database on success!")
}

func dumpToDirectory(username, password, host, port, outputDir, separator string) (string, error) {
	dumperArgs := defaultConfig()
	dumperArgs.User = username
	dumperArgs.Password = password
	dumperArgs.Address = fmt.Sprintf("%s:%s", host, port)
	if all {
		if len(outputDir) == 0 {
			outputDir = "all-databases"
		}
		log.Logger.Info("dump all databases into file " + outputDir + " at multi-threaded mode...")
		fmt.Println()
	} else {
		dumperArgs.DatabaseRegexp = ""
		if len(databases) == 0 {
			return "", fmt.Errorf(" please provide at least one database name with flag --databases or -d.")
		}
		databasesStr := strings.Join(databases, ",")
		dumperArgs.Database = databasesStr
		if len(tables) != 0 {
			if len(databases) != 1 {
				return "", fmt.Errorf("if you specify the tables, the database can only specify one")
			}
			dumperArgs.Table = strings.Join(tables, ",")
		}
		if len(outputDir) == 0 {
			if len(databases) == 1 {
				outputDir = databases[0]
			} else {
				outputDir = strings.Join(databases, "_")
			}
		}
		log.Logger.Info("dump database " + databasesStr + " into " + outputDir +  " directory... " )
		fmt.Println()
	}

	dumperArgs.Outdir = outputDir
	if _, err := os.Stat(dumperArgs.Outdir); os.IsNotExist(err) {
		x := os.MkdirAll(dumperArgs.Outdir, 0o777)
		common.AssertNil(x)
	}
	Dumper(log.Logger, dumperArgs, separator)

	return outputDir, nil
}

func dumpToSqlFile(username, password, host, port, outputFile string) error {
	var execCmd *exec.Cmd
	if all {
		if len(outputFile) == 0 {
			outputFile = "all-databases.sql"
		}
		log.Logger.Info("dump all databases into file " + outputFile + " ...")
		fmt.Println()
		execCmd = exec.Command("mysqldump", "-u", username, "-p"+password, "--host", host, "--port", port, "-A")
	} else {
		// 检查数据库名是否为空
		if  len(databases) == 0{
			return fmt.Errorf("please provide the mysql database name with flag --databases or -d")
		}

		if len(databases) > 1 {
			return fmt.Errorf("only one mysql database dump is supported in single-threaded mode")
		}

		if len(outputFile) == 0 {
			outputFile = databases[0] + ".sql"
		}
		log.Logger.Info("dump the certain database " + databases[0] +  " into file " + outputFile + " ...")
		fmt.Println()
		execCmd = exec.Command("mysqldump", "-u", username, "-p"+password, "--host", host, "--port", port,"--databases", databases[0])
	}

	f, err := os.Create(outputFile)
	if err != nil {
		log.Logger.Error("create " +  outputFile + " file error: " + err.Error())
	}
	defer f.Close()

	execCmd.Stdout = f
	err = execCmd.Run()
	if err != nil {
		log.Logger.Info(execCmd.String())
		return fmt.Errorf("mysqldump command run error: " + err.Error())
	}

	return nil
}

// WriteFile used to write datas to file.
func WriteFile(file string, data string) error {
	flag := os.O_RDWR | os.O_APPEND
	if _, err := os.Stat(file); os.IsNotExist(err) {
		flag |= os.O_CREATE
	}
	f, err := os.OpenFile(file, flag, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := f.Write(sqlcommon.StringToBytes(data))
	if err != nil {
		return err
	}
	if n != len(data) {
		return io.ErrShortWrite
	}
	return nil
}

func dumpTable(log *xlog.Log, conn *common.Connection, args *config.Config, database string, table string, separator string) {
	var allBytes uint64
	var allRows uint64
	var where string
	var selfields []string

	fields := make([]string, 0, 16)
	{
		cursor, err := conn.StreamFetch(fmt.Sprintf("SELECT * FROM `%s`.`%s` LIMIT 1", database, table))
		common.AssertNil(err)

		flds := cursor.Fields()
		for _, fld := range flds {
			log.Debug("dump -- %#v, %s, %s", args.Filters, table, fld.Name)
			if _, ok := args.Filters[table][fld.Name]; ok {
				continue
			}

			fields = append(fields, fmt.Sprintf("`%s`", fld.Name))
			replacement, ok := args.Selects[table][fld.Name]
			if ok {
				selfields = append(selfields, fmt.Sprintf("%s AS `%s`", replacement, fld.Name))
			} else {
				selfields = append(selfields, fmt.Sprintf("`%s`", fld.Name))
			}
		}
		err = cursor.Close()
		common.AssertNil(err)
	}

	if v, ok := args.Wheres[table]; ok {
		where = fmt.Sprintf(" WHERE %v", v)
	}

	cursor, err := conn.StreamFetch(fmt.Sprintf("SELECT %s FROM `%s`.`%s` %s", strings.Join(selfields, ", "), database, table, where))
	common.AssertNil(err)

	fileNo := 1
	stmtsize := 0
	chunkbytes := 0
	rows := make([]string, 0, 256)
	inserts := make([]string, 0, 256)
	for cursor.Next() {
		row, err := cursor.RowValues()
		common.AssertNil(err)

		values := make([]string, 0, 16)
		for _, v := range row {
			if v.Raw() == nil {
				values = append(values, "NULL")
			} else {
				str := v.String()
				switch {
				case v.IsSigned(), v.IsUnsigned(), v.IsFloat(), v.IsIntegral(), v.Type() == querypb.Type_DECIMAL:
					values = append(values, str)
				default:
					values = append(values, fmt.Sprintf("\"%s\"", common.EscapeBytes(v.Raw())))
				}
			}
		}
		//r := "(" + strings.Join(values, ",") + ")"
		r := strings.Join(values, separator)
		rows = append(rows, r)

		allRows++
		stmtsize += len(r)
		chunkbytes += len(r)
		allBytes += uint64(len(r))
		atomic.AddUint64(&args.Allbytes, uint64(len(r)))
		atomic.AddUint64(&args.Allrows, 1)

		if stmtsize >= args.StmtSize {
			//insertone := fmt.Sprintf("INSERT INTO `%s`(%s) VALUES\n%s", table, strings.Join(fields, ","), strings.Join(rows, ",\n"))
			insertone := strings.Join(rows, "\n")
			inserts = append(inserts, insertone)
			rows = rows[:0]
			stmtsize = 0
		}

		if (chunkbytes / 1024 / 1024) >= args.ChunksizeInMB {
			query := strings.Join(inserts, "\n") + "\n"
			//file := fmt.Sprintf("%s/%s.%s.%09d.csv", args.Outdir, database, table, fileNo)
			file := fmt.Sprintf("%s/%s.%s.csv", args.Outdir, database, table)
			WriteFile(file, query)

			log.Info("dumping.table[%s.%s].rows[%v].bytes[%vMB].part[%v].thread[%d]", database, table, allRows, (allBytes / 1024 / 1024), fileNo, conn.ID)
			inserts = inserts[:0]
			chunkbytes = 0
			fileNo++
		}
	}
	if chunkbytes > 0 {
		if len(rows) > 0 {
			//insertone := fmt.Sprintf("INSERT INTO `%s`(%s) VALUES\n%s", table, strings.Join(fields, ","), strings.Join(rows, ",\n"))
			insertone := strings.Join(rows, "\n")
			inserts = append(inserts, insertone)
		}

		query := strings.Join(inserts, "\n") + "\n"
		//file := fmt.Sprintf("%s/%s.%s.%09d.csv", args.Outdir, database, table, fileNo)
		file := fmt.Sprintf("%s/%s.%s.csv", args.Outdir, database, table)
		WriteFile(file, query)
	}
	err = cursor.Close()
	common.AssertNil(err)

	log.Info("dumping.table[%s.%s].done.allrows[%v].allbytes[%vMB].thread[%d]...", database, table, allRows, (allBytes / 1024 / 1024), conn.ID)
}

func writeMetaData(args *config.Config) {
	file := fmt.Sprintf("%s/metadata", args.Outdir)
	WriteFile(file, "")
}

func allDatabases(log *xlog.Log, conn *common.Connection) []string {
	qr, err := conn.Fetch("SHOW DATABASES")
	common.AssertNil(err)

	databases := make([]string, 0, 128)
	for _, t := range qr.Rows {
		databases = append(databases, t[0].String())
	}
	return databases
}

func filterDatabases(log *xlog.Log, conn *common.Connection, filter *regexp.Regexp, invert bool) []string {
	qr, err := conn.Fetch("SHOW DATABASES")
	common.AssertNil(err)

	databases := make([]string, 0, 128)
	for _, t := range qr.Rows {
		if (!invert && filter.MatchString(t[0].String())) || (invert && !filter.MatchString(t[0].String())) {
			databases = append(databases, t[0].String())
		}
	}
	return databases
}

func dumpTableSchema(log *xlog.Log, conn *common.Connection, args *config.Config, database string, table string) {
	qr, err := conn.Fetch(fmt.Sprintf("SHOW CREATE TABLE `%s`.`%s`", database, table))
	common.AssertNil(err)
	schema := qr.Rows[0][1].String() + ";\n"

	file := fmt.Sprintf("%s/%s.%s-schema.sql", args.Outdir, database, table)
	common.WriteFile(file, schema)
	log.Info("dumping.table[%s.%s].schema...", database, table)
}

func allTables(log *xlog.Log, conn *common.Connection, database string) []string {
	qr, err := conn.Fetch(fmt.Sprintf("SHOW TABLES FROM `%s`", database))
	common.AssertNil(err)

	tables := make([]string, 0, 128)
	for _, t := range qr.Rows {
		tables = append(tables, t[0].String())
	}
	return tables
}

func dumpDatabaseSchema(log *xlog.Log, conn *common.Connection, args *config.Config, database string) {
	err := conn.Execute(fmt.Sprintf("USE `%s`", database))
	common.AssertNil(err)

	schema := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", database)
	file := fmt.Sprintf("%s/%s-schema-create.sql", args.Outdir, database)
	common.WriteFile(file, schema)
	log.Info("dumping.database[%s].schema...", database)
}

// Dumper used to start the dumper worker.
func Dumper(log *xlog.Log, args *config.Config, separator string) {
	pool, err := common.NewPool(log, args.Threads, args.Address, args.User, args.Password, args.SessionVars)
	common.AssertNil(err)
	defer pool.Close()

	// Meta data.
	writeMetaData(args)

	// database.
	var wg sync.WaitGroup
	conn := pool.Get()
	var databases []string
	t := time.Now()
	if args.DatabaseRegexp != "" {
		r := regexp.MustCompile(args.DatabaseRegexp)
		databases = filterDatabases(log, conn, r, args.DatabaseInvertRegexp)
	} else {
		if args.Database != "" {
			databases = strings.Split(args.Database, ",")
		} else {
			databases = allDatabases(log, conn)
		}
	}
	for _, database := range databases {
		dumpDatabaseSchema(log, conn, args, database)
	}

	// tables.
	tables := make([][]string, len(databases))
	for i, database := range databases {
		if args.Table != "" {
			tables[i] = strings.Split(args.Table, ",")
		} else {
			tables[i] = allTables(log, conn, database)
		}
	}
	pool.Put(conn)

	for i, database := range databases {
		for _, table := range tables[i] {
			conn := pool.Get()
			dumpTableSchema(log, conn, args, database, table)

			wg.Add(1)
			go func(conn *common.Connection, database string, table string) {
				defer func() {
					wg.Done()
					pool.Put(conn)
				}()
				log.Info("dumping.table[%s.%s].datas.thread[%d]...", database, table, conn.ID)
				dumpTable(log, conn, args, database, table, separator)
				log.Info("dumping.table[%s.%s].datas.thread[%d].done...", database, table, conn.ID)
			}(conn, database, table)
		}
	}

	tick := time.NewTicker(time.Millisecond * time.Duration(args.IntervalMs))
	defer tick.Stop()
	go func() {
		for range tick.C {
			diff := time.Since(t).Seconds()
			allbytesMB := float64(atomic.LoadUint64(&args.Allbytes) / 1024 / 1024)
			allrows := atomic.LoadUint64(&args.Allrows)
			rates := allbytesMB / diff
			log.Info("dumping.allbytes[%vMB].allrows[%v].time[%.2fsec].rates[%.2fMB/sec]...", allbytesMB, allrows, diff, rates)
		}
	}()

	wg.Wait()
	elapsed := time.Since(t).Seconds()
	log.Info("dumping.all.done.cost[%.2fsec].allrows[%v].allbytes[%v].rate[%.2fMB/s]", elapsed, args.Allrows, args.Allbytes, (float64(args.Allbytes/1024/1024) / elapsed))
}
