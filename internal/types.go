package internal

type Stat struct {
	Name  string
	Value uint64
}

type Stats map[string]uint64

type StatCollection struct {
	GrandTotal uint64
	Collect    map[string]Stats //[YYYYMMDD][".jpg"]=35
}

type StatTotal struct { //StatTotal[YY].Months[MM].Days[DD][".jpg"]
	Total uint64              // total for all in database
	Years map[uint8]*StatYear //.Years[YYYY]=
}

func (st *StatTotal) Get(y []byte) (sm *StatYear) {
	if st.Years == nil {
		st.Years = map[uint8]*StatYear{}
	}

	n := Btoi8(y) // this will break badly if it's not yyyy

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
	sm.AddDay(d, n)

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
	Total uint64           //this months's total
	Days  map[uint8]uint64 //.Days[DD]=35
}

func (sm *StatMonth) Get(d []byte) uint64 {
	return sm.Days[Btoi8(d)] // doesn't check for nil because you can read a nil map
}

func (sm *StatMonth) AddDay(d []byte, n []byte) *StatMonth {
	if sm.Days == nil {
		sm.Days = map[uint8]uint64{}
	}
	sm.Days[Btoi8(d)] += Btoi64(n)
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

func (sc *StatCollection) Add(dateKey, name string, val uint64) {
	st := sc.Collect[dateKey]
	if st == nil {
		st = Stats{}
		sc.Collect[dateKey] = st
	}
	st[name] = val
}

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
