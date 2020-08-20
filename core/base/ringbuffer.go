package base

import (
	"baidu-uid-go/core/util"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var (
	CAN_PUT_FLAG  uint32 = 0
	CAN_TAKE_FLAG uint32 = 1
	START_POINT   int64  = -1
	NOT_RUNNING   uint64 = 0
	RUNNING       uint64 = 1
)

type RingBuffer struct {
	tail   int64
	cursor int64

	bufferSize uint64
	indexMask  int64
	slots      []uint64
	flags      []uint32
	//slots            []atomic.Value
	//flags            []atomic.Value
	paddingThreshold int64
	mu               sync.Mutex

	//padding
	running uint64

	UidProvider func(uint64) []uint64

	lastSecond uint64
}

//初始化ringbuffer
func InitRingBuffer(bufferSize uint64, paddingFactor uint64) *RingBuffer {
	r := &RingBuffer{
		bufferSize: bufferSize,
		indexMask:  (int64)(bufferSize - 1),
		slots:      make([]uint64, bufferSize),
		flags:      make([]uint32, bufferSize),
		//slots:            make([]atomic.Value, bufferSize),
		//flags:            make([]atomic.Value, bufferSize),
		paddingThreshold: (int64)(bufferSize * paddingFactor / 100),
	}
	var i uint64
	for i = 0; i < bufferSize; i++ {
		//r.flags[i].Store(CAN_PUT_FLAG)
		r.flags[i] = CAN_PUT_FLAG
	}
	r.tail = START_POINT
	r.cursor = START_POINT
	r.lastSecond = uint64(time.Now().Unix())
	return r
}

func (r *RingBuffer) Put(uid uint64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	currentTail := atomic.LoadInt64(&r.tail)
	currentCursor := atomic.LoadInt64(&r.cursor)
	distance := currentTail - currentCursor
	if distance == r.indexMask {
		//到达最大buffersize，拒绝再放入uid
		fmt.Println(fmt.Sprintf("Rejected putting buffer for uid:%v",uid))
		return false
	}
	nextTailIndex := r.calSlotIndex(currentTail + 1)
	if atomic.LoadUint32(&r.flags[nextTailIndex]) != CAN_PUT_FLAG {
		//标志不是可put，拒绝
		fmt.Println("Curosr not in can put status")
		return false
	}
	//nextTailIndexflag := r.flags[nextTailIndex].Load().(uint8)
	//if nextTailIndexflag != CAN_PUT_FLAG {
	//	//标志不是可put，拒绝
	//	fmt.Println("标志不是可put，拒绝")
	//	return false
	//}
	atomic.StoreUint64(&r.slots[nextTailIndex], uid)
	atomic.StoreUint32(&r.flags[nextTailIndex], CAN_TAKE_FLAG)
	//r.slots[nextTailIndex].Store(uid)
	//r.flags[nextTailIndex].Store(CAN_TAKE_FLAG)
	atomic.AddInt64(&r.tail, 1)
	return true
}

func (r *RingBuffer) Take() (uint64, error) {
	currentCursor := atomic.LoadInt64(&r.cursor)
	nextCursor := util.Uint64UpdateAndGet(&r.cursor, func(o int64) int64 {
		if o == atomic.LoadInt64(&r.tail) {
			return o
		} else {
			return o + 1
		}
	})
	if nextCursor < currentCursor {
		panic("Curosr can't move back")
	}
	currentTail := atomic.LoadInt64(&r.tail)
	if currentTail-nextCursor < r.paddingThreshold {
		//异步填充
		go r.AsyncPadding()
	}
	if currentTail == currentCursor {
		//拒绝
		return 0, errors.New("Rejected take uid")
	}
	nextCursorIndex := r.calSlotIndex(nextCursor)
	if atomic.LoadUint32(&r.flags[nextCursorIndex]) != CAN_TAKE_FLAG {
		return 0, errors.New("Curosr not in can take status")
	}
	//nextCursorIndexflag := r.flags[nextCursorIndex].Load().(uint8)
	//if nextCursorIndexflag != CAN_TAKE_FLAG {
	//	return 0, errors.New("Curosr not in can take status")
	//}
	uid := atomic.LoadUint64(&r.slots[nextCursorIndex])
	atomic.StoreUint32(&r.flags[nextCursorIndex], CAN_PUT_FLAG)
	//uid := r.slots[nextCursorIndex].Load().(uint64)
	//r.flags[nextCursorIndex].Store(CAN_PUT_FLAG)
	return uid, nil
}

func (r *RingBuffer) AsyncPadding() {
	// is still running
	if !atomic.CompareAndSwapUint64(&r.running, NOT_RUNNING, RUNNING) {
		return
	}

	isFullRingBuffer := true
	// fill the rest slots until to catch the cursor
	for isFullRingBuffer {
		uidList := r.UidProvider(atomic.AddUint64(&r.lastSecond, 1))
		for _, uid := range uidList {
			isFullRingBuffer = !r.Put(uid)
			if isFullRingBuffer {
				break
			}
		}
	}

	// not running now
	atomic.CompareAndSwapUint64(&r.running, RUNNING, NOT_RUNNING)
}

func (r *RingBuffer) calSlotIndex(sequence int64) int64 {
	//return sequence % r.indexMask
	return sequence & r.indexMask
}
