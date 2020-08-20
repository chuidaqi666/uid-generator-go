package base

import (
	"os"
	"strconv"
	util "uid-generator-go/core/util"
)

const (
	TimestampBits = "TIME_STAMP_BITS"
	WorkerIdBits  = "WORKER_ID_BITS"
	SequenceBits  = "SEQUENCE_BITS"
	EpochStr      = "EPOCH_STR"
	//EpochSeconds  = "EPOCH_SECONDS"
	BoostPower    = "BOOST_POWER"
	PaddingFactor = "PADDING_FACTOR"
)

type UidglobalCfg struct {
	TimestampBits int
	WorkerIdBits  int
	SequenceBits  int
	//workId        int
	EpochStr      string
	BoostPower    int
	PaddingFactor uint64
}

func (u *UidglobalCfg) Init(t int, w int, s int, e string) {
	u.InitDefault()
	u.TimestampBits = t
	u.WorkerIdBits = w
	u.SequenceBits = s
	u.EpochStr = e
	//u.BoostPower = 3
	//u.PaddingFactor = 50
	//err := configFromSystemEnv(u)
	//if err != nil {
	//	panic(err)
	//}
}

func (u *UidglobalCfg) InitDefault() {
	u.TimestampBits = 29
	u.WorkerIdBits = 22
	u.SequenceBits = 13
	u.EpochStr = "2020-08-20"
	u.BoostPower = 3
	u.PaddingFactor = 50
	err := configFromSystemEnv(u)
	if err != nil {
		panic(err)
	}
}

func configFromSystemEnv(uc *UidglobalCfg) (err error) {
	if timestampBits := os.Getenv(TimestampBits); !util.IsBlank(timestampBits) {
		uc.TimestampBits, err = strconv.Atoi(timestampBits)
	}
	if workerIdBits := os.Getenv(WorkerIdBits); !util.IsBlank(workerIdBits) {
		uc.WorkerIdBits, err = strconv.Atoi(workerIdBits)
	}
	if sequenceBits := os.Getenv(SequenceBits); !util.IsBlank(sequenceBits) {
		uc.SequenceBits, err = strconv.Atoi(sequenceBits)
	}
	if epochStr := os.Getenv(EpochStr); !util.IsBlank(epochStr) {
		uc.EpochStr = epochStr
	}
	if boostPower := os.Getenv(BoostPower); !util.IsBlank(boostPower) {
		uc.BoostPower, err = strconv.Atoi(boostPower)
	}
	if paddingFactor := os.Getenv(PaddingFactor); !util.IsBlank(paddingFactor) {
		uc.PaddingFactor, err = strconv.ParseUint(PaddingFactor, 10, 64)
	}
	return
}
