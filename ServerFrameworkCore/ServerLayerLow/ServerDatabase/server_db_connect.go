package ServerDatabase

import (
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerLog"
	"database/sql"
	"fmt"
	_ "go-sql-driver"
)

// MariaDBConnect 连接 MariaDB 数据库.
func MariaDBConnect(User, Pass, Host, Port, DBName string) (*sql.DB, error) {
	ServerLog.StarLogFmt("INFO", "Connect DataBase: '"+DBName+"'", "Server", "Core", "DB", "Connect")
	// 构建数据库连接.
	ConnectDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", User, Pass, Host, Port, DBName)

	// 连接 SQL 数据库.
	DataBase, err := sql.Open("mysql", ConnectDSN)
	if err != nil {
		ServerLog.StarLogFmt("ERROR", "Connect MariaDB: "+err.Error(), "Server", "Core", "DB", "Connect")
		return nil, err
	}

	// 数据库测试连接 Ping.
	err = DataBase.Ping()
	if err != nil {
		ServerLog.StarLogFmt("ERROR", "Connect Test: "+err.Error(), "Server", "Core", "DB", "Connect")
		return nil, err
	} else {
		ServerLog.StarLogFmt("INFO", "Connect Success(Ping).", "Server", "Core", "DB", "Connect")
	}
	return DataBase, nil
}
