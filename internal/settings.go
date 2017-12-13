package internal

//#####################################
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

var (
	//namesBucket       = []byte("buckets")
	YearBucket        = []byte("year")
	HitsBucket        = []byte("hits")
	ExtensionsBucket  = []byte("extensions")
	StatusCodesBucket = []byte("statuscodes")
	IPSBucket         = []byte("ips")
	ErrorsBucket      = []byte("errors")
	//errBucketNotFound = errors.New("Bucket not found")
	//errActIDExists    = errors.New("ActID already associated with Task")
)
