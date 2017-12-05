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

type StatYear struct { //StatYear[YY].Collect[MM].Collect[DD][".jpg"]
	GrandTotal uint64
	Years      map[string]StatMonth //[YY]=
}
type StatMonth struct {
	GrandTotal uint64
	Months     map[string]StatDay //[MM]=
}
type StatDay struct {
	GrandTotal uint64
	//Collect    map[string]Stats //[DD][".jpg"]=35
	Days map[string]uint64 //[DD][".jpg"]=35
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
	NotFoundBucket    = []byte("notfound")
	//errBucketNotFound = errors.New("Bucket not found")
	//errActIDExists    = errors.New("ActID already associated with Task")
)
