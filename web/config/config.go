package config

import (
	"database/sql"
	"io"
	"os"

	"fmt"

	"github.com/gin-gonic/gin"
)

var Db *sql.DB

// InitConfig initializes the configuration of the application
func InitConfig() {
	initDB()

	//打印输出信息到文件
	f, _ := os.OpenFile("./api/gin.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func initDB() bool { //连接到MySQL
	var err error
	// 初始化数据库连接
	Db, err = sql.Open("mysql", "root:lq3525926@tcp(127.0.0.1:3306)/mytest")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return false
	}
	err = Db.Ping() //尝试连接数据库
	if err != nil {
		fmt.Println("connect mysql failed,", err)
		return false
	}
	//设置数据库连接池的最大连接数
	Db.SetMaxIdleConns(10)
	return true
}

func Query(SQL string) ([]map[string]string, bool) { //通用查询
	rows, err := Db.Query(SQL) //执行SQL语句，比如select * from users
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
			b, _ := raw_value.([]byte)
			v := string(b) //将raw数据转换成字符串
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
