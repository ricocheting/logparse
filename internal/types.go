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

////////////////////////////////////
type StatTotal struct { //StatTotal[YY].Months[MM].Days[DD][".jpg"]
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
func (st *StatTotal) AddTotal(y, m, d []byte, n []byte) *StatTotal {
	i := Btoi64(n)

	st.Total += i
	sy := st.Get(y)
	sy.Total += i
	sm := sy.Get(m)
	sm.Total += i
	sm.AddDayNum(d, n)

	if sm.Date.IsZero() {
		sm.Date, _ = time.Parse("200601", string(y)+string(m))
	}

	return st
}

////////////////////////////////////
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

////////////////////////////////////
type StatMonth struct {
	Date  time.Time             // timestamp for this month. so we can get that "201701" means save info to path /2017/01-January.html
	Total uint64                //this months's total
	Days  map[uint8]interface{} //.Days[DD]=35
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

////////////////////////////////////
var StatusCodeNames = map[string]string{
	"200": "OK",
	"206": "Partial Content", //resume
	"301": "Moved Permanently",
	"302": "Found",        //redirect
	"304": "Not Modified", //use cache
	"400": "Bad Request",
	"401": "Unauthorized",
	"403": "Forbidden",
	"404": "Not Found",
	"405": "Method Not Allowed",
	"408": "Request Timeout",
	"416": "Range Not Satisfiable", //asked to resume on part of the file that doesn't exist
	"499": "Client Closed Request",
	"500": "Internal Server Error",
}

/*func (sc *StatCollection) Add(dateKey, name string, val uint64) {
	st := sc.Collect[dateKey]
	if st == nil {
		st = Stats{}
		sc.Collect[dateKey] = st
	}
	st[name] = val
}*/

var (
	//namesBucket       = []byte("buckets")
	YearBucket        = []byte("year")
	HitsBucket        = []byte("hits")
	ExtensionsBucket  = []byte("extensions")
	StatusCodesBucket = []byte("statuscodes")
	IPSBucket         = []byte("ips")
	NotFoundBucket    = []byte("notfound")
	//errBucketNotFound = errors.New("Bucket not found")
	//errActIDExists    = errors.New("ActID already associated with Task")
)
