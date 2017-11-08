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
	//errBucketNotFound = errors.New("Bucket not found")
	//errActIDExists    = errors.New("ActID already associated with Task")
)
