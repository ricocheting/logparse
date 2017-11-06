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

/*func (sc *StatCollection) Add(key0, key1 string, val uint64) {
	st := sc.Collect[key0]
	if st == nil {
		st = Stats{}
		sc.Collect[key0] = st
	}
	st[key1] += val
}*/
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
	//errBucketNotFound = errors.New("Bucket not found")
	//errActIDExists    = errors.New("ActID already associated with Task")
)
