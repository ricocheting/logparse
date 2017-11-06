package internal

type Stat struct {
	Name  string
	Value uint64
}

type Stats map[string]uint64

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
