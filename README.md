百度uid-generator的go版本，目前只有cache-generator可以使用


func Test11(t *testing.T) {
	t.Run("uid test", func(t *testing.T) {
		var wg sync.WaitGroup
		c := &CacheGenerator{}
		w := &worker.WorkerDBConfig{
			UserName:  "root",
			Password:  "123456",
			IpAddrees: "127.0.0.1",
			Port:      "3306",
			DbName:    "test",
			Charset:   "utf8",
		}
		c.InitDefaultWithWorkDB(w)
		for i := 0; i < 10000000; i++ {
			wg.Add(1)
			go getUid(c, &wg)
		}
		wg.Wait()
	})
}

func getUid(c *CacheGenerator, wg *sync.WaitGroup) {
	defer wg.Done()
	//uid, err := c.GetUID()
	c.GetUID()
	//fmt.Println("uid=", uid, err)
	//fmt.Println("parseUID=",c.ParseUID(uid))
}
