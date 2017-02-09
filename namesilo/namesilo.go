package namesilo

//rate limit: 2 requests/seconds
import (
	"fmt"
	"net/http"
	"log"
	"encoding/xml"
	"strings"
	"time"
	"sync"
)

const (
	NamesiloBaseUrl      = "https://www.namesilo.com/api/%s?version=1&type=xml&key=%s&%s"
	NamesiloBatchBaseUrl = "https://www.namesilo.com/apibatch/%s?version=1&type=xml&key=%s&%s"
)

type NamesiloAPI struct {
	Key   string
	Debug bool
	mutex sync.Mutex
	t1    time.Time
	t2    time.Time
}


type Request struct {
	XMLName       xml.Name `xml:"request"`
	Operation     string  `xml:"operation"`
	Ip            string  `xml:"ip"`

}

type Reply struct {
	XMLName        xml.Name `xml:"reply"`
	Code           int  `xml:"code"`
	Detail         string  `xml:"detail"`
	Message        string  `xml:"message"`
	Domain         string    `xml:"domain"`
	Amount         float64 `xml:"order_amount"`
	Available      []string   `xml:"available>domain"`
	Unavailable    []string   `xml:"unavailable>domain"`
	Invalid        []string   `xml:"invalid>domain"`
}

type Namesilo struct {
	XMLName       xml.Name `xml:"namesilo"`
	Request       Request
	Reply         Reply
}

func (silo *NamesiloAPI) request(u string) (*http.Response, error) {
	silo.mutex.Lock()
	defer func() {
		silo.t1 = silo.t2
		silo.t2 = time.Now()
		silo.mutex.Unlock()
	}()

	dur := time.Since(silo.t1)
	if dur.Nanoseconds() < int64(time.Second) {
		time.Sleep(time.Duration(int64(time.Second) - dur.Nanoseconds()))
	}

	return http.Get(u)

}

func (silo *NamesiloAPI) register(baseUrl string, command string, domain string, duration int, coupon string) (bool, error) {
	u := fmt.Sprintf(baseUrl, command, silo.Key, fmt.Sprintf("domain=%s&years=%d&private=1&auto_renew=0&coupon", domain, duration, coupon))
	r, err := silo.request(u)

	if err != nil {
		return false, err
	}
	defer r.Body.Close()
	result := Namesilo{}
	decoder := xml.NewDecoder(r.Body)
	decoder.Decode(&result)
	if silo.Debug {
		log.Println(result)
	}
	return result.Reply.Code == 300, nil
}

func (silo *NamesiloAPI) Register(domain string, duration int, coupon string) (bool, error) {
	command := "registerDomain"
	return silo.register(NamesiloBaseUrl, command, domain, duration, coupon)
}

func (silo *NamesiloAPI) RegisterDrop(domain string, duration int, coupon string) (bool, error) {
	command := "registerDomainDrop"
	return silo.register(NamesiloBatchBaseUrl, command, domain, duration, coupon)
}

func (silo *NamesiloAPI) Available(domains []string) (map[string]bool, error) {

	command := "checkRegisterAvailability"
	u := fmt.Sprintf(NamesiloBaseUrl, command, silo.Key, fmt.Sprintf("domains=%s", strings.Join(domains, ",")))
	r, err := silo.request(u)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	namesilo := Namesilo{}
	decoder := xml.NewDecoder(r.Body)
	decoder.Decode(&namesilo)
	if silo.Debug {
		log.Println(namesilo)
	}
	result := make(map[string]bool)
	for _, d := range (namesilo.Reply.Available) {
		result[d] = true
	}
	for _, d := range (namesilo.Reply.Unavailable) {
		result[d] = false
	}
	for _, d := range (namesilo.Reply.Invalid) {
		result[d] = false
	}
	return result, nil
}

