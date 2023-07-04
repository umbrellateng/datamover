package main

import (
	"core.bank/datamover/core/mysql"
	"core.bank/datamover/flags"
	"core.bank/datamover/log"
	"core.bank/datamover/utils"
)


var (
	onlineTmpDir string
)

func main() {

	defer func() {

		if r := recover(); r != nil {
			if mysql.IsOnlineMode(flags.From, flags.To) {
				_ = utils.DeleteDirAndFiles(onlineTmpDir)
			}
			log.Logger.Error("something wrong, received from panic: %v", r)
		}
	}()

	flags.InitFlags()

	var err error
	dbInfo := mysql.DBInfo{
		Username: flags.User,
		Password: flags.Password,
		Host: flags.Host,
		Port: flags.Port,
	}

	if mysql.IsOnlineMode(flags.From, flags.To) {
		err = mysql.OnlineMove(dbInfo, flags.From, flags.To, flags.Databases, flags.All)
		if err != nil {
			log.Logger.Error("move database online error: " + err.Error())
		}
		return
	}

	if flags.Thread {
		if flags.Restore {
			err = mysql.RestoreDBFromDirectory(dbInfo, flags.Input)
			if err != nil {
				log.Logger.Error("Restore DB from Directory " + flags.Input + " error: " + err.Error())
				return
			}

		} else {
			err = mysql.DumpDBToDirectory(dbInfo, flags.Output, flags.Databases, flags.All)
			if err != nil {
				log.Logger.Error("Dump DB to Directory " + flags.Output + " error: " + err.Error())
				return
			}
		}

	} else {
		if flags.Restore {
			err = mysql.RestoreDBFromSqlFile(dbInfo, flags.Input)
			if err != nil {
				log.Logger.Error("Restore DB from SqlFile error: " + err.Error())
				return
			}

		} else {
			err = mysql.DumpDBToSqlFile(dbInfo, flags.Output, flags.Databases, flags.All)
			if err != nil {
				log.Logger.Error("Dump DB to SqlFile error: " + err.Error())
				return
			}
		}
	}

	log.Logger.Info("Success!")
}
