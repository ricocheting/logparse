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
	"504": "Gateway Timeout Error",
}

var (
	YearBucket        = []byte("year")
	HitsBucket        = []byte("hits")
	ExtensionsBucket  = []byte("extensions")
	DirectoriesBucket = []byte("directories")
	StatusCodesBucket = []byte("statuscodes")
	IPSBucket         = []byte("ips")
	ErrorsBucket      = []byte("errors")

	// this is used for parsing refers and for printing headers to the template if no -domain="" command line flag is passed in (both parser and printer)
	DefaultDomain = "sitename.com"
)
