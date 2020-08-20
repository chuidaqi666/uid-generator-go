package base

type UidGenerator interface {
	GetUID() (uint64, error)
	ParseUID(uint64) string
}
