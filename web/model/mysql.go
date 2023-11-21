package model

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var db *sql.DB

func InitDB() bool { //连接到MySQL
	var err error
	// 初始化数据库连接
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.ip"),
		viper.GetString("db.port"),
		viper.GetString("db.dbname"),
	))
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return false
	}
	err = db.Ping() //尝试连接数据库
	if err != nil {
		fmt.Println("connect mysql failed,", err)
		return false
	}
	//设置数据库连接池的最大连接数
	db.SetMaxIdleConns(viper.GetInt("db.maxconnect"))
	return true
}

func query(SQL string, args ...interface{}) ([]map[string]string, bool) { //通用查询
	rows, err := db.Query(SQL, args...) //执行SQL语句，比如select * from user
	if err != nil {
		panic(err)
	}
	columns, _ := rows.Columns()            //获取列的信息
	count := len(columns)                   //列的数量
	var values = make([]interface{}, count) //创建一个与列的数量相当的空接口
	for i := range values {
		var ii interface{} //为空接口分配内存
		values[i] = &ii    //取得这些内存的指针，因后继的Scan函数只接受指针
	}
	ret := make([]map[string]string, 0) //创建返回值：不定长的map类型切片
	for rows.Next() {
		err := rows.Scan(values...)  //开始读行，Scan函数只接受指针变量
		m := make(map[string]string) //用于存放1列的 [键/值] 对
		if err != nil {
			panic(err)
		}
		for i, colName := range columns {
			var raw_value = *(values[i].(*interface{})) //读出raw数据，类型为byte
			var v string
			switch t := raw_value.(type) {
			case int64:
				v = fmt.Sprintf("%d", t) //将raw数据转换成字符串
			default:
				b, _ := raw_value.([]byte)
				v = string(b) //将raw数据转换成字符串
			}
			m[colName] = v //colName是键，v是值
		}
		ret = append(ret, m) //将单行所有列的键值对附加在总的返回值上（以行为单位）
	}

	defer rows.Close()

	if len(ret) != 0 {
		return ret, true
	}
	return nil, false
}

func SelectMySql(tableName string, conditions map[string]interface{}) ([]map[string]string, bool) {
	var args []interface{}
	whereParts := make([]string, 0, len(conditions))
	for k, v := range conditions {
		whereParts = append(whereParts, fmt.Sprintf("%s = ?", k))
		args = append(args, v)
	}
	SQL := "SELECT * FROM " + tableName
	if len(whereParts) > 0 {
		SQL += " WHERE " + strings.Join(whereParts, " AND ")
	}
	return query(SQL, args...)
}

func InsertMySql(tableName string, values map[string]interface{}) (sql.Result, bool) { //sql.Result为了获得id
	var args []interface{}
	columns := make([]string, 0, len(values))
	placeholders := make([]string, 0, len(values))
	for k, v := range values {
		columns = append(columns, k)
		placeholders = append(placeholders, "?")
		args = append(args, v)
	}
	SQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))
	ret, err := db.Exec(SQL, args...)
	if err != nil {
		fmt.Printf("insert failed, err: %v\n", err)
		return nil, false
	}
	return ret, true
}

func UpdateMySQL(tableName string, values map[string]interface{}, conditions map[string]interface{}) bool {
	var args []interface{}
	setParts := make([]string, 0, len(values))
	for k, v := range values {
		setParts = append(setParts, fmt.Sprintf("%s = ?", k))
		args = append(args, v)
	}
	whereParts := make([]string, 0, len(conditions))
	for k, v := range conditions {
		whereParts = append(whereParts, fmt.Sprintf("%s = ?", k))
		args = append(args, v)
	}
	sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, strings.Join(setParts, ", "), strings.Join(whereParts, " AND "))
	_, err := db.Exec(sqlStr, args...)
	if err != nil {
		fmt.Printf("update failed, err: %v\n", err)
		return false
	}
	return true
}

func DeleteMySQL(tableName string, conditions map[string]interface{}) (sql.Result, bool) {
	var args []interface{}
	whereParts := make([]string, 0, len(conditions))
	for k, v := range conditions {
		whereParts = append(whereParts, fmt.Sprintf("%s = ?", k))
		args = append(args, v)
	}
	SQL := "DELETE FROM " + tableName
	if len(whereParts) > 0 {
		SQL += " WHERE " + strings.Join(whereParts, " AND ")
	}
	ret, err := db.Exec(SQL, args...)
	if err != nil {
		fmt.Printf("delete failed, err: %v\n", err)
		return nil, false
	}
	return ret, true
}
