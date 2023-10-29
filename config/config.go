package config

import "database/sql"

var (
	DatabaseIP       = "127.0.0.1"
	DatabasePort     = "3306"
	DatabaseName     = "fusioncdn"
	DatabasePassword = "fusioncdn"
	DatabaseDbName   = "fusioncdn"
	Key              = "6iQJasdasdasdasde7ti"
	ApiKey           = ""
)

func GetDb() (*sql.DB, error) {
	db, err := sql.Open("mysql", GetDatabaseName()+":"+GetDatabasePassword()+"@tcp("+GetDatabaseIP()+":"+GetDatabasePort()+")/"+GetDatabaseDbName()+"?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		return &sql.DB{}, err
	}
	err = db.Ping()
	if err != nil {
		return &sql.DB{}, err
	}
	return db, nil
}

func GetDatabaseIP() string {
	return DatabaseIP
}

func GetDatabasePort() string {
	return DatabasePort
}

func GetDatabaseName() string {
	return DatabaseName
}

func GetDatabasePassword() string {
	return DatabasePassword
}

func GetDatabaseDbName() string {
	return DatabaseDbName
}
