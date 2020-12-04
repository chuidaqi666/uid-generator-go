package generator

import (
	"fmt"
	"github.com/chuidaqi666/uid-generator-go/core/base"
	"github.com/chuidaqi666/uid-generator-go/core/worker"
	"sync"
	"testing"
)

func Test11(t *testing.T) {
	t.Run("uid test", func(t *testing.T) {
		c := &CacheGenerator{}
		w := &worker.WorkerDBConfig{
			UserName:  "root",
			Password:  "123456",
			IpAddrees: "localhost",
			Port:      "3306",
			DbName:    "db1",
			Charset:   "utf8",
		}
		c.InitDefaultWithWorkDB(w)
		fmt.Println(c.ParseUID(10611077482164756))
		fmt.Println(c.ParseUID(10611249280856541))
		fmt.Println(c.ParseUID(10611421079548299))
		fmt.Println(c.ParseUID(10611592878239417))
		fmt.Println(c.ParseUID(10611764676932158))
		fmt.Println(c.ParseUID(10611936475623861))
		fmt.Println(c.ParseUID(10612108274316327))
		fmt.Println("----------------------")
		fmt.Println(c.ParseUID(10616059644229570))
		fmt.Println(c.ParseUID(10616231442921690))
		fmt.Println(c.ParseUID(10616403241612322))

		//Rejected putting buffer for uid:10616059644229570
		//Curosr not in can put status,rejected uid:%v 10616231442921690
		//Rejected putting buffer for uid:10616403241612322

		//{"UID":"10600735200894975","deltaSeconds":"1598161322","workerId":"14","sequence":"8191"}
		//{"UID":"10600494682718208","deltaSeconds":"1598161315","workerId":"14","sequence":"0"}

		//{"UID":"10610665165299712","deltaSeconds":"1598161611","workerId":"17","sequence":"0"}
		//Rejected putting buffer for uid:10611077482164756
		//Rejected putting buffer for uid:10611249280856541
		//Rejected putting buffer for uid:10611421079548299
		//Rejected putting buffer for uid:10611592878239417
		//Rejected putting buffer for uid:10611764676932158
		//Rejected putting buffer for uid:10611936475623861
		//Rejected putting buffer for uid:10612108274316327
	})

	t.Run("uid test", func(t *testing.T) {
		var wg sync.WaitGroup
		c := &CacheGenerator{}
		w := &worker.WorkerDBConfig{
			UserName:  "root",
			Password:  "123456",
			IpAddrees: "localhost",
			Port:      "3306",
			DbName:    "db1",
			Charset:   "utf8",
		}
		c.InitDefaultWithWorkDB(w)
		//uid, _ := c.GetUID()
		//fmt.Println(c.ParseUID(uid))
		for i := 0; i < 300; i++ {
			wg.Add(1)
			go getUid(c, &wg)
		}
		wg.Wait()
		//uid, _ = c.GetUID()
		//fmt.Println(c.ParseUID(uid))
	})
}

func getUid(c base.UidGenerator, wg *sync.WaitGroup) {
	defer wg.Done()
	//var once sync.Once
	//uid, _ := c.GetUID()
	//fmt.Println(1)
	for i := 0; i < 200000; i++ {
		c.GetUID()
		//once.Do(func() {
		//	uid, err := c.GetUID()
		//	if err!=nil{
		//		fmt.Println("___________________parseUID=",c.ParseUID(uid),uid,err)
		//	}
		//})
		//uid, _ := c.GetUID()
		//fmt.Println("parseUID=",c.ParseUID(uid))
	}
	//fmt.Println("uid=", uid, err)
	//fmt.Println("parseUID=",c.ParseUID(uid))
}
