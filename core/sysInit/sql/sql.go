package sql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go-admin/core/sysInit"
	"go-admin/core/utils/log"
	"reflect"
	"strings"
)

var (
	db *sql.DB
	isSqlInit bool
	// struceName ->{
	//		field1: {
	//			columeName: colume1 / nil
	//			null: null / nil
	//			unique: nil / unique
	//			pk: pk / nil
	//			type: int64 / string / int32 / time
	//		}
	//}
	entities map[string]map[string]map[string]string
)

func init(){
	dataSource := string(sysInit.GetConfigValue("mysql.user") + ":" + sysInit.GetConfigValue("mysql.password") + "@tcp(" + sysInit.GetConfigValue("mysql.host") + ":" +
		sysInit.GetConfigValue("mysql.port") + ")/" + sysInit.GetConfigValue("mysql.dbname") + "?charset=utf8");
	tdb, err := sql.Open("mysql", dataSource)

	if err != nil {
		log.Error.Println(err)
	}else {
		db = tdb
		isSqlInit = true
	}
	entities = make(map[string]map[string]map[string]string)
	log.Debug.Println("init")
}

func Insert(instance interface{}) interface{}{
	instType := reflect.TypeOf(instance)
	instValue := reflect.ValueOf(instance)

	fieldMap := make([]string, instType.NumField())
	valueMap := make([]interface{}, instType.NumField())

	for i := 0; i < len(fieldMap); i++{
		fieldMap[i] = instType.Field(i).Name
		valueMap[i] = instValue.Field(i)
	}

	sqlString := fmt.Sprint("INSERT INTO `", "", "` (")

	for index, key := range fieldMap{
		if index == len(fieldMap) - 1 {
			sqlString += string("`" + key + "`) VALUES(?, ?, ?, ?)")
		}else {
			sqlString += string("`" + key + "`, ")
		}

	}
	log.Debug.Println(sqlString, " - ")

	return nil
}

func RegisterOrm(instance interface{}) error{
	instType := reflect.TypeOf(instance)
	instValue := reflect.ValueOf(instance)

	structName := instValue.Type().Name()
	entities[structName] = make(map[string]map[string]string)
	fieldMap := make([]string, instType.NumField())

	if len(fieldMap) == 0 {
		return fmt.Errorf("struct has not element")
	}

	for i := 0; i < len(fieldMap); i++{
		if orm := instType.Field(i).Tag.Get("orm"); orm != "" {
			entities[structName][instType.Field(i).Name] = make(map[string]string)
			var ormField = strings.Split(orm, ";")
			var fieldParam = entities[structName][instType.Field(i).Name]

			for j := 0; j < len(ormField); j++ {
				if strings.Contains(ormField[j], "column") {
					var column = strings.Split(strings.Split(ormField[j], "(")[1], ")")[0]
					fieldParam["columnName"] = column
				}else {
					fieldParam[ormField[j]] = ormField[j]
				}
			}

			fieldParam["type"] = reflect.ValueOf(instance).Field(i).Type().Name()

			if fieldParam["columnName"] == "" {
				fieldParam["columnName"] = instType.Field(i).Name
			}
		}else {
			log.WARN.Println(instType.String(), instType.Field(i).Name, " does have orm note")
		}
	}

	RunCreateTable()
	return nil
}

func RunCreateTable() error{
	dbx := sqlx.NewDb(db, "mysql")
	defer dbx.Close()

	if err := dbx.Ping(); err != nil {
		return err
	}

	for v, m := range entities {
		var create = fmt.Sprint(" \n\tCREATE TABLE IF NOT EXISTS `", strings.ToLower(v), "` {\n")

		var index = 0
		for v2, m2 := range m{
			var p string
			create += fmt.Sprint("\t\t`", v2, "` ")

			if p = m2["type"]; p == "int64" {
				create += fmt.Sprint("bigint ")
			}else if p == "string"{
				create += fmt.Sprint("varchar(255) ")
			}

			if p = m2["null"]; p == "" {
				create += fmt.Sprint("NOT NULL ")

				if p = m2["default"]; p != "" && m2["type"] != "int64" {
					create += fmt.Sprint("DEFAULT '", m2["default"], "' ")
				}else if m2["type"] != "int64" {
					create += fmt.Sprint("DEFAULT '' ")
				}
			}

			if p = m2["pk"]; p != "" {
				create += fmt.Sprint("PRIMARY KEY ")
			}else if p = m2["unique"]; p != "" {
				create += fmt.Sprint("UNIQUE ")
			}

			if index != len(m) - 1{
				create += ",\n"
			}else {
				create += "\n\t} ENGINE=InnoDB;"
			}
			index++
		}

		log.Debug.Println(create)
	}

	

	return nil
}

