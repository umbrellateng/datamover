package main

import (
	"core.bank/datamover/cmd"
	"core.bank/datamover/log"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			log.Logger.Error("something wrong, received from panic: %v", r)
		}
	}()

	cmd.Execute()

	//flags.InitFlags()
	//from, to := flags.From, flags.To
	//if utils.OnlineMode(from, to) {
	//	fromUser, fromPwd, fromHost, fromPort, err := utils.ParseDBStringWithoutDB(from)
	//	if err != nil {
	//		log.Logger.Error("parse source db string error: %s", err.Error())
	//		return
	//	}
	//	toUser, toPwd, toHost, toPort, err := utils.ParseDBStringWithoutDB(to)
	//	if err != nil {
	//		log.Logger.Error("parse target db string error: %s", err.Error())
	//		return
	//	}
	//
	//	fromMysql := mover.NewMySql(fromUser, fromPwd, fromHost, fromPort, flags.All, flags.Databases)
	//	info := mover.NewBaseInfo(toUser, toPwd, toHost, toPort)
	//	err = fromMysql.MoveOnline([]mover.BaseInfo{info})
	//	if err != nil {
	//		log.Logger.Error("move database online error: " + err.Error())
	//	}
	//	return
	//}
	//
	//targetMysql := mover.NewMySql(flags.User, flags.Password, flags.Host, flags.Port, flags.All, flags.Databases)
	//
	//if flags.Thread {
	//	if flags.Restore {
	//		err := targetMysql.RestoreFromDirectory(flags.Input)
	//		if err != nil {
	//			log.Logger.Error("Restore DB from Directory " + flags.Input + " error: " + err.Error())
	//			return
	//		}
	//
	//	} else {
	//		err := targetMysql.DumpToDirectory(flags.Output)
	//		if err != nil {
	//			log.Logger.Error("Dump DB to Directory " + flags.Output + " error: " + err.Error())
	//			return
	//		}
	//	}
	//
	//} else {
	//	if flags.Restore {
	//		err := targetMysql.RestoreFromFile(flags.Input)
	//		if err != nil {
	//			log.Logger.Error("Restore DB from SqlFile error: " + err.Error())
	//			return
	//		}
	//
	//	} else {
	//		err := targetMysql.DumpToFile(flags.Output)
	//		if err != nil {
	//			log.Logger.Error("Dump DB to SqlFile error: " + err.Error())
	//			return
	//		}
	//	}
	//}
	//
	//log.Logger.Info("Success!")
}
