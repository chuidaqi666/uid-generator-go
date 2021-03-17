package base

import (
	"errors"
	"fmt"
	"github.com/chuidaqi666/uid-generator-go/core/util"
	"sync"
	"sync/atomic"
	"time"
)

type PaddedTailStruct struct {
	tail int64
	_    CacheLinePad
}

type PaddedCursorStruct struct {
	cursor int64
	_      CacheLinePad
}

type PaddedRunStruct struct {
	running uint64
	_       CacheLinePad
}

type PaddedSlotStruct struct {
	slot uint64
	_    CacheLinePad
}

type PaddedflagStruct struct {
	flag uint64
	_    CacheLinePad
}

type CacheLinePad struct {
	_ [CacheLinePadSizeV2]byte
}

const CacheLinePadSizeV2 = 32

type RingBufferV2 struct {
	bufferSize       uint64
	indexMask        int64
	paddingThreshold int64
	lastSecond       uint64
	running          PaddedRunStruct
	tail             PaddedTailStruct
	cursor           PaddedCursorStruct
	slots            []PaddedSlotStruct
	flags            []PaddedflagStruct
	mu               sync.Mutex
	UidProvider      func(uint64) []uint64
}

//初始化ringbuffer
func InitRingBufferV2(bufferSize uint64, paddingFactor uint64) *RingBufferV2 {
	r := &RingBufferV2{
		bufferSize: bufferSize,
		indexMask:  (int64)(bufferSize - 1),
		slots:      make([]PaddedSlotStruct, bufferSize),
		flags:      make([]PaddedflagStruct, bufferSize),
		//slots:            make([]atomic.Value, bufferSize),
		//flags:            make([]atomic.Value, bufferSize),
		paddingThreshold: (int64)(bufferSize * paddingFactor / 100),
	}
	var i uint64
	for i = 0; i < bufferSize; i++ {
		//r.flags[i].Store(CAN_PUT_FLAG)
		r.flags[i].flag = CAN_PUT_FLAGV2
	}
	r.tail = PaddedTailStruct{
		tail: -1,
	}
	r.cursor = PaddedCursorStruct{
		cursor: -1,
	}
	r.lastSecond = uint64(time.Now().Unix())
	return r
}

func (r *RingBufferV2) Put(uid uint64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	currentTail := atomic.LoadInt64(&r.tail.tail)
	currentCursor := atomic.LoadInt64(&r.cursor.cursor)
	distance := currentTail - currentCursor
	if distance == r.indexMask {
		//到达最大buffersize，拒绝再放入uid
		fmt.Printf("Rejected putting buffer for uid:%v,tail:%v,cursor:%v\n", uid, currentTail, currentCursor)
		return false
	}
	nextTailIndex := r.calSlotIndexV2(currentTail + 1)
	if atomic.LoadUint64(&r.flags[nextTailIndex].flag) != CAN_PUT_FLAGV2 {
		//标志不是可put，拒绝
		fmt.Printf("Curosr not in can put status,rejected uid:%v,tail:%v,cursor:%v\n", uid, currentTail, currentCursor)
		return false
	}
	//nextTailIndexflag := r.flags[nextTailIndex].Load().(uint8)
	//if nextTailIndexflag != CAN_PUT_FLAG {
	//	//标志不是可put，拒绝
	//	fmt.Println("标志不是可put，拒绝")
	//	return false
	//}
	atomic.StoreUint64(&r.slots[nextTailIndex].slot, uid)
	atomic.StoreUint64(&r.flags[nextTailIndex].flag, CAN_TAKE_FLAGV2)
	//r.slots[nextTailIndex].Store(uid)
	//r.flags[nextTailIndex].Store(CAN_TAKE_FLAG)
	atomic.AddInt64(&r.tail.tail, 1)
	return true
}

func (r *RingBufferV2) Take() (uint64, error) {
	currentCursor := atomic.LoadInt64(&r.cursor.cursor)
	nextCursor := util.Uint64UpdateAndGet(&r.cursor.cursor, func(o int64) int64 {
		if o == atomic.LoadInt64(&r.tail.tail) {
			return o
		} else {
			return o + 1
		}
	})
	if nextCursor < currentCursor {
		panic("Curosr can't move back")
	}
	currentTail := atomic.LoadInt64(&r.tail.tail)
	if currentTail-nextCursor < r.paddingThreshold {
		//异步填充
		go r.AsyncPadding()
	}
	if currentTail == currentCursor {
		//拒绝
		return 0, errors.New("Rejected take uid")
	}
	nextCursorIndex := r.calSlotIndexV2(nextCursor)
	if atomic.LoadUint64(&r.flags[nextCursorIndex].flag) != CAN_TAKE_FLAGV2 {
		return 0, errors.New("Curosr not in can take status")
	}
	//nextCursorIndexflag := r.flags[nextCursorIndex].Load().(uint8)
	//if nextCursorIndexflag != CAN_TAKE_FLAG {
	//	return 0, errors.New("Curosr not in can take status")
	//}
	uid := atomic.LoadUint64(&r.slots[nextCursorIndex].slot)
	atomic.StoreUint64(&r.flags[nextCursorIndex].flag, CAN_PUT_FLAGV2)
	//uid := r.slots[nextCursorIndex].Load().(uint64)
	//r.flags[nextCursorIndex].Store(CAN_PUT_FLAG)
	return uid, nil
}

func (r *RingBufferV2) AsyncPadding() {
	// is still running
	if !atomic.CompareAndSwapUint64(&r.running.running, NOT_RUNNING, RUNNING) {
		return
	}

	isFullRingBuffer := false
	// fill the rest slots until to catch the cursor
	for !isFullRingBuffer {
		uidList := r.UidProvider(atomic.AddUint64(&r.lastSecond, 1))
		for _, uid := range uidList {
			isFullRingBuffer = !r.Put(uid)
			if isFullRingBuffer {
				break
			}
		}
	}

	// not running now
	atomic.CompareAndSwapUint64(&r.running.running, RUNNING, NOT_RUNNING)
}

func (r *RingBufferV2) calSlotIndexV2(sequence int64) int64 {
	//return sequence % r.indexMask
	return sequence & r.indexMask
}
