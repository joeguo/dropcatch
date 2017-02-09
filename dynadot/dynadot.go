package dynadot

import (
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
	"errors"
	"log"
	"sync"
)

//rate limit:request only after the previous request returns

const (
	BaseUrl = "https://api.dynadot.com/api2.html?key=%s&command=%s&%s"
)

type DynadotAPI struct {
	Key   string
	Debug bool
	mutex sync.Mutex
}

func (api *DynadotAPI) Available(domains []string) (map[string]bool, error) {
	ps := make([]string, len(domains))
	for i, d := range (domains) {
		ps[i] = fmt.Sprintf("domain%d=%s", i, d)
	}
	command := "search"
	u := fmt.Sprintf(BaseUrl, api.Key, command, strings.Join(ps, "&"))
	api.mutex.Lock()
	r, err := http.Get(u)
	api.mutex.Unlock()
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	result := make(map[string]bool)
	response := string(bs)
	if strings.HasPrefix(response, "ok,") {
		array := strings.Split(response, "\n")
		array = array[2:len(array) - 1]
		for _, line := range (array) {
			temp := strings.Split(line, ",")
			result[temp[1]] = temp[3] == "yes"
		}
		return result, nil

	}
	return nil, errors.New(response)
}

func (api *DynadotAPI) Register(domain string, duration int) (bool, error) {
	command := "register"
	u := fmt.Sprintf(BaseUrl, api.Key, command, fmt.Sprintf("domain=%s&duration=%d", domain, duration))
	api.mutex.Lock()
	r, err := http.Get(u)
	api.mutex.Unlock()
	if err != nil {
		return false, err
	}
	defer r.Body.Close()
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false, err
	}

	response := string(bs)
	if strings.HasPrefix(response, "ok,") {
		array := strings.Split(response, "\n")
		if strings.HasPrefix(array[2], "success") {
			return true, nil
		}
		if api.Debug {
			log.Println(domain, array[2])
		}
		return false, nil
	}
	return false, errors.New(response)
}
