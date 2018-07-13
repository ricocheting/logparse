package internal

import (
	"time"
)

type Stat struct { // still used in records.go
	Name  string
	Value uint64
}

type Stats map[string]uint64 //[".jpg"]=35

/*type StatCollection struct {
	GrandTotal uint64
	Collect    map[string]Stats //[YYYYMMDD][".jpg"]=35
}*/

//#####################################
type StatTotal struct { //StatTotal[YY].Months[MM].Days[DD] is [".jpg"]=1500 or directly =1500
	Total uint64               // total for all in database
	Years map[uint16]*StatYear //.Years[YYYY]=
}

func (st *StatTotal) Get(y []byte) (sm *StatYear) {
	if st.Years == nil {
		st.Years = map[uint16]*StatYear{}
	}

	n := Btoi16(y)

	if sm = st.Years[n]; sm == nil {
		sm = &StatYear{}
		st.Years[n] = sm
	}
	return
}
func (st *StatTotal) AddTotal(y, m, d []byte, ni []byte) *StatTotal {
	i := Btoi64(ni)

	st.Total += i
	sy := st.Get(y)
	sy.Total += i
	sm := sy.Get(m)
	sm.Total += i
	sm.AddDayNum(d, ni)

	sm.Avg = (sm.Total / uint64(len(sm.Days)))

	if sm.Min == 0 || sm.Min > i { // if it's today, ignore
		sm.Min = i
	}
	if sm.Max < i {
		sm.Max = i
	}

	if sm.Date.IsZero() {
		sm.Date, _ = time.Parse("200601", string(y)+string(m))
	}

	return st
}
func (st *StatTotal) AddStat(y, m, d []byte, name string, i uint64) *StatTotal {

	st.Total += i
	sy := st.Get(y)
	sy.Total += i
	sm := sy.Get(m)
	sm.Total += i
	sm.AddDayStats(d, name, i)

	if sm.Date.IsZero() {
		sm.Date, _ = time.Parse("200601", string(y)+string(m))
	}

	return st
}

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

/*func (sm *StatMonth) AddDay(d []byte, n []byte) *StatMonth {
	if sm.Days == nil {
		sm.Days = map[uint8]uint64{}
	}
	sm.Days[Btoi8(d)] += Btoi64(n)
	return sm
}*/

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

/*func (sc *StatCollection) Add(dateKey, name string, val uint64) {
	st := sc.Collect[dateKey]
	if st == nil {
		st = Stats{}
		sc.Collect[dateKey] = st
	}
	st[name] = val
}*/

//#####################################
type StatErrors struct {
	Page map[string]StatErrorPage //.Page["https://www.example.com/errors.html"].Missing["/404/missing.jog"] = 36
}

type StatErrorPage struct {
	Page    string            //used when sorting into []StatErrorPage instead of StatErrors.Page map
	Total   uint64            //.Page["https://www.example.com/errors.html"].Total = 329687296
	Missing map[string]uint64 //.Page["https://www.example.com/errors.html"].Missing["/404/missing.jog"] = 36
}

func (se *StatErrors) Increment(page, missing string) *StatErrors {
	if se.Page == nil {
		se.Page = map[string]StatErrorPage{}
	}

	if sep, ok := se.Page[page]; ok {
		sep.Missing[missing]++
		sep.Total++
		se.Page[page] = sep
	} else {
		//initialize
		sep := StatErrorPage{}
		sep.Missing = map[string]uint64{}

		sep.Missing[missing] = 1
		sep.Total++
		se.Page[page] = sep
	}

	return se
}

func (se *StatErrors) SetVal(page, missing string, val uint64) *StatErrors {
	if se.Page == nil {
		se.Page = map[string]StatErrorPage{}
	}

	if sep, ok := se.Page[page]; ok {
		sep.Missing[missing] = val
		sep.Total += val
		se.Page[page] = sep
	} else {
		//initialize
		sep := StatErrorPage{}
		sep.Missing = map[string]uint64{}

		sep.Missing[missing] = val
		sep.Total += val
		se.Page[page] = sep
	}

	return se
}

func (se StatErrors) ToSlice(min uint64) []StatErrorPage { //still used on ghetto.go
	out := make([]StatErrorPage, 0, len(se.Page))
	for k, v := range se.Page {
		if min > 0 && v.Total < min {
			continue
		}
		out = append(out, StatErrorPage{Page: k, Total: v.Total, Missing: v.Missing})
	}
	return out[:len(out):len(out)] // trim the slice to release the unused memory
}
