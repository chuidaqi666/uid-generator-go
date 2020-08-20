package base

import "sync"

type BitsAllocator struct {
	TotalBits      int
	timestampShift int
	workerIdShift  int

	TimestampBits int
	WorkerIdBits  int
	SequenceBits  int

	maxDeltaSeconds uint64
	MaxWorkerId     uint64
	MaxSequence     uint64
	initOnce        sync.Once
}

func (b *BitsAllocator) Init(timestampBits int, workerIdBits int, sequenceBits int) {
	b.initOnce.Do(func() {
		if timestampBits+workerIdBits+sequenceBits < 64 {
			panic("Less than 64 bits")
		}
		if timestampBits+workerIdBits+sequenceBits > 64 {
			panic("more than 64 bits")
		}
		//initialize bits
		b.TimestampBits = timestampBits
		b.WorkerIdBits = workerIdBits
		b.SequenceBits = sequenceBits
		//initialize max value
		var m int64
		m = -1
		b.maxDeltaSeconds = uint64(^(m << timestampBits))
		b.MaxWorkerId = uint64(^(m << workerIdBits))
		b.MaxSequence = uint64(^(m << sequenceBits))
		// initialize shift
		b.timestampShift = workerIdBits + sequenceBits
		b.workerIdShift = sequenceBits
		b.TotalBits = 1 << 6
	})
}

func (b *BitsAllocator) Allocate(deltaSeconds uint64, workerId uint64, sequence uint64) uint64 {
	return (deltaSeconds << b.timestampShift) | (workerId << b.workerIdShift) | sequence
}
