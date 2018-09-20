package internal

import (
	"time"
)

type Stat struct { // still used in records.go
	Name  string
	Value uint64
}

type Stats map[string]uint64 //[".jpg"]=35

//#####################################
type StatYear struct {
	Total  uint64               //this year's total
	Months map[uint8]*StatMonth //.Months[MM]=
}

func (sm *StatYear) Get(m []byte) (sd *StatMonth) {
	if sm.Months == nil {
		sm.Months = map[uint8]*StatMonth{}
	}

	n := Btoi8(m)

	if sd = sm.Months[n]; sd == nil {
		sd = &StatMonth{}
		sm.Months[n] = sd
	}
	return
}

//#####################################
type StatMonth struct {
	Date  time.Time             // timestamp for this month. so we can get that "201701" means save info to path /2017/01-January.html
	Total uint64                //this months's total
	Days  map[uint8]interface{} //.Days[DD]=35
	Min   uint64
	Max   uint64
	Avg   uint64
}

func (sm *StatMonth) Get(d []byte) interface{} {
	return sm.Days[Btoi8(d)] // doesn't check for nil because you can read a nil map
}

func (sm *StatMonth) AddDayStats(d []byte, name string, val uint64) {
	if sm.Days == nil {
		sm.Days = map[uint8]interface{}{}
	}

	if st, ok := sm.Days[Btoi8(d)].(Stats); ok {
		st[name] = val
		sm.Days[Btoi8(d)] = st
	} else {
		//initialize
		st := Stats{}
		st[name] = val
		sm.Days[Btoi8(d)] = st
	}

}

func (sm *StatMonth) AddDayNum(d []byte, n []byte) *StatMonth {
	if sm.Days == nil {
		sm.Days = map[uint8]interface{}{}
	}

	if val, ok := sm.Days[Btoi8(d)].(uint64); ok {
		sm.Days[Btoi8(d)] = val + Btoi64(n)
	} else {
		//initialize
		sm.Days[Btoi8(d)] = Btoi64(n)
	}

	return sm
}

/*
// SortDaysData converts sm.Days[uint8] = Stats{} to sorted sm.Days[uint8] = []Stat
func (sm *StatMonth) SortDaysData() {
	for i, e := range sm.Days {

		// i is the index, e the element
	}

}*/
