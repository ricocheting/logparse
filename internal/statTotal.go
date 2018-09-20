package internal

import "time"

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
