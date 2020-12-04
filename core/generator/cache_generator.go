package generator

import (
	"fmt"
	"github.com/chuidaqi666/uid-generator-go/core/base"
	"github.com/chuidaqi666/uid-generator-go/core/util"
	"github.com/chuidaqi666/uid-generator-go/core/worker"
	"sync"
)

type CacheGenerator struct {
	ringbuffer *base.RingBuffer
	bita       *base.BitsAllocator
	/** Customer epoch, unit as second. For example 2016-05-20 (ms: 1463673600000)*/
	//epochStr     string
	epochSeconds uint64
	workId       uint64
	boostPower   int
	config       *base.UidglobalCfg
	initOnce     sync.Once
}

func (c *CacheGenerator) InitDefaultWithWorkDB(w *worker.WorkerDBConfig) {
	c.initOnce.Do(func() {
		workId, _ := worker.NewWorkerDB(w)
		u := &base.UidglobalCfg{}
		//初始化数据库
		u.InitDefault()
		c.config = u
		c.bita = &base.BitsAllocator{}
		//初始化BitsAllocator
		c.bita.Init(c.config.TimestampBits, c.config.WorkerIdBits, c.config.SequenceBits)
		c.workId = uint64(workId)
		if c.workId > c.bita.MaxWorkerId {
			panic(fmt.Sprintf("Worker id %v exceeds the max %v", c.workId, c.bita.MaxWorkerId))
		}
		c.epochSeconds = util.DateToSecond(c.config.EpochStr)
		c.boostPower = c.config.BoostPower
		bufferSize := (c.bita.MaxSequence + 1) << c.boostPower
		c.ringbuffer = base.InitRingBuffer(bufferSize, c.config.PaddingFactor)
		c.ringbuffer.UidProvider = c.nextIdsForOneSecond
		c.ringbuffer.AsyncPadding()
	})
}

func (c *CacheGenerator) InitUidglobalCfgAndWorkerDB(u *base.UidglobalCfg, w *worker.WorkerDBConfig) {
	c.initOnce.Do(func() {
		workId, _ := worker.NewWorkerDB(w)
		c.config = u
		c.bita = &base.BitsAllocator{}
		c.bita.Init(c.config.TimestampBits, c.config.WorkerIdBits, c.config.SequenceBits)
		c.workId = uint64(workId)
		if c.workId > c.bita.MaxWorkerId {
			panic(fmt.Sprintf("Worker id %v exceeds the max %v", c.workId, c.bita.MaxWorkerId))
		}
		c.epochSeconds = util.DateToSecond(c.config.EpochStr)
		c.boostPower = c.config.BoostPower
		bufferSize := (c.bita.MaxSequence + 1) << c.boostPower
		c.ringbuffer = base.InitRingBuffer(bufferSize, c.config.PaddingFactor)
		c.ringbuffer.UidProvider = c.nextIdsForOneSecond
		c.ringbuffer.AsyncPadding()
	})
}

func (c *CacheGenerator) GetUID() (uid uint64, err error) {
	uid, err = c.ringbuffer.Take()
	return
}

// parse UID
func (c *CacheGenerator) ParseUID(uid uint64) string {
	sequence := (uid << (c.bita.TotalBits - c.bita.SequenceBits)) >> (c.bita.TotalBits - c.bita.SequenceBits)
	workerId := (uid << (c.bita.TimestampBits)) >> (c.bita.TotalBits - c.bita.WorkerIdBits)
	deltaSeconds := uid >> (c.bita.WorkerIdBits + c.bita.SequenceBits)
	return fmt.Sprintf("{\"UID\":\"%v\",\"deltaSeconds\":\"%v\",\"workerId\":\"%v\",\"sequence\":\"%v\"}",
		uid, c.epochSeconds+deltaSeconds, workerId, sequence)
}

func (c *CacheGenerator) nextIdsForOneSecond(currentSecond uint64) (uidList []uint64) {
	// Initialize result list size of (max sequence + 1)
	listSize := c.bita.MaxSequence + 1
	// Allocate the first sequence of the second, the others can be calculated with the offset
	firstSeqUid := c.bita.Allocate(currentSecond-c.epochSeconds, c.workId, 0)
	var offset uint64
	for offset = 0; offset < listSize; offset++ {
		uidList = append(uidList, firstSeqUid+offset)
	}
	return uidList
}
