package base

var (
	CAN_PUT_FLAG    uint32 = 0
	CAN_PUT_FLAGV2  uint64 = 0
	CAN_TAKE_FLAG   uint32 = 1
	CAN_TAKE_FLAGV2 uint64 = 1
	START_POINT     int64  = -1
	NOT_RUNNING     uint64 = 0
	RUNNING         uint64 = 1
)

type Rb interface {
	Put(uid uint64) bool
	Take() (uint64, error)
	AsyncPadding()
}
