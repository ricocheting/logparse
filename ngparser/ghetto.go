package ngparser

import (
	"bufio"
	"io"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ricocheting/logparse/storage"
)

type StatType uint8

const ( // StatTypes
	IPs StatType = iota
	StatusCodes
	Pages
	Hits
	UserAgents
	Extensions

	maxType
)
const (
	timeFmt = `2/Jan/2006:15:04:05 -0700`
)

var re = regexp.MustCompile(`(.+?)\s[^[]+\[([^\]]+)\]\s"(\w+) (.+?)\sHTTP/(\d\.\d)"\s+(\d+)\s+(\d+)\s+"([^"]+)"\s+"([^"]+)"`)
var store *storage.Store

// $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" "$http_x_forwarded_for";
type Record struct {
	IP        string
	TS        time.Time
	Method    string
	Filename  string
	Status    string
	Referer   string
	UserAgent string
}

type Stat struct {
	Name  string
	Value uint64
}

type stats map[string]uint64

func (s stats) ToSlice(min uint64) []Stat {
	out := make([]Stat, 0, len(s))
	for k, v := range s {
		if min > 0 && v < min {
			continue
		}
		out = append(out, Stat{k, v})
	}
	return out[:len(out):len(out)] // trim the slice to release the unused memory
}

type Parser struct {
	mux   sync.RWMutex
	data  [maxType]stats
	count uint64
	ipv6  uint64
}

func New() *Parser {
	var p Parser
	for i := range p.data {
		p.data[i] = stats{}
	}
	return &p
}

func (p *Parser) Parse(r io.Reader, fn func(r *Record)) {
	// p.mux.Lock()
	// defer p.mux.Unlock()
	var (
		sc  = bufio.NewScanner(r)
		in  = make(chan string, 1) //runtime.NumCPU()
		out = make(chan *Record, 1)
		wg  sync.WaitGroup
	)

	for i := 0; i < cap(in); i++ {
		wg.Add(1)
		go p.parseLine(&wg, in, out)
	}

	go func() {
		for sc.Scan() {
			in <- sc.Text()
		}
		close(in)
		wg.Wait()
		close(out)
	}()

	var startDate time.Time

	for r := range out {
		if fn != nil {
			fn(r)
		}

		cleanPath := r.Filename
		if idx := strings.IndexByte(cleanPath, '?'); idx != -1 {
			cleanPath = cleanPath[:idx]
		}

		if startDate.IsZero() {
			startDate = r.TS
		} else if isNewerDay(startDate, r.TS) {
			// insert p.data into buckets
			// clear p.data
			//fmt.Printf("isNewerDay: %+v  OR  %+v\n", r.TS, startDate)
			p.saveData([]byte(r.TS.Format("20060102")))
			startDate = r.TS
		}

		p.mux.Lock()

		// how many IPv6 addresses
		ipCnt := p.data[IPs][r.IP]
		if strings.IndexByte(r.IP, ':') > -1 && ipCnt == 0 {
			p.ipv6++
		}

		// log the rest of the data
		p.count++
		p.data[IPs][r.IP] = ipCnt + 1
		p.data[StatusCodes][r.Status]++

		// if the page exists, save its info
		if r.Status != "404" {
			p.data[Pages][cleanPath]++
			//p.data[Hits][r.Filename]++
			//p.data[UserAgents][r.UserAgent]++ // probably should parse the agent and store something like Chrome-XX, IE11, Edge, etc.
			p.data[Extensions][strings.ToLower(filepath.Ext(cleanPath))]++
		}

		p.mux.Unlock()
	}
}

// get a specific type of stat in the log
// Example:  fmt.Printf("%+v\n", p.Stats(ngparser.Pages, 1000))
func (p *Parser) Stats(t StatType, filterMin uint64) (out []Stat) {
	p.mux.RLock()
	out = p.data[t].ToSlice(filterMin)
	p.mux.RUnlock()

	// sorting outside the lock
	sort.Slice(out, func(i, j int) bool { return out[i].Value > out[j].Value })

	return
}

// get the total number of unique entries
// fmt.Printf("Unique Files: %+v\n", p.StatsCount(ngparser.Pages))
func (p *Parser) StatsCount(t StatType) (total uint64) {
	p.mux.RLock()
	total = uint64(len(p.data[t]))
	p.mux.RUnlock()
	return
}

// get the total number of hits in the log
func (p *Parser) Count() (l uint64) {
	p.mux.RLock()
	l = p.count
	p.mux.RUnlock()
	return
}

// get the number of IPv4 and IPv6 in the log
func (p *Parser) IPsCount() (v4, v6 uint64) {
	p.mux.RLock()
	total := uint64(len(p.data[IPs]))
	ipv6 := p.ipv6
	p.mux.RUnlock()
	return total - ipv6, ipv6
}

// convert the log line into a record
func (p *Parser) parseLine(wg *sync.WaitGroup, in chan string, out chan *Record) {
	cp := re.Copy()
	for l := range in {
		var line []string
		if parsed := cp.FindAllStringSubmatch(l, -1); len(parsed) == 1 {
			line = parsed[0]
		} else {
			continue
		}

		r := &Record{
			IP:        line[1],
			Method:    line[3],
			Filename:  line[4],
			Status:    line[6],
			Referer:   line[7],
			UserAgent: line[8],
		}

		r.TS, _ = time.Parse(timeFmt, line[2])
		out <- r
	}
	wg.Done()
}

func isNewerDay(startDate, compDate time.Time) bool {
	/*y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2*/

	d1 := startDate.Truncate(24 * time.Hour)
	d2 := compDate.Truncate(24 * time.Hour)

	return !(d1.Equal(d2) || d2.Before(d1))
}

func (p *Parser) saveData(dateKey []byte) {
	store = storage.NewStore(filepath.Join("data", "db"))
	if err := store.Open(); err != nil {
		panic("Error opening storage (db possibly still open by another process): " + err.Error())
	}

	err := store.SaveHits(dateKey, p.Count())

	if err != nil {
		panic("Error saveData(): " + err.Error())
	}

}
