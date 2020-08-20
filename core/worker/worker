package worker

import (
	"baidu-uid-go/core/util"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"time"
)

type WorkerDBConfig struct {
	UserName  string
	Password  string
	IpAddrees string
	Port      string
	DbName    string
	Charset   string
	db        *sqlx.DB
}

func NewWorkerDB(w *WorkerDBConfig) (int64, error) {
	w.connectMysql()
	return w.addRecord()
}

func (w *WorkerDBConfig) connectMysql() {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v", w.UserName, w.Password, w.IpAddrees, w.Port, w.DbName, w.Charset)
	Db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	w.db = Db
}

func (w *WorkerDBConfig) addRecord() (int64, error) {
	hostname, environment := util.GetHostAndE()
	sql := "INSERT INTO WORKER_NODE(HOST_NAME,PORT,TYPE,LAUNCH_DATE,MODIFIED,CREATED)VALUES (?,?,?,?,NOW(),NOW())"
	rand.Seed(time.Now().Unix())
	port := fmt.Sprintf("%v-%v", time.Now().UnixNano()/1e6, rand.Intn(100000))
	//fmt.Println(hostname, port, environment)
	result, err := w.db.Exec(sql, hostname, port, environment, time.Now())
	if err != nil {
		panic(err)
	}
	id, _ := result.LastInsertId()
	//fmt.Printf("insert success, last id:[%d]\n", id)
	return id, nil
}
