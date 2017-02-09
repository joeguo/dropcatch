package resellerclub

import (
	"fmt"
	"net/http"
	"log"
	"encoding/json"
	"strings"
	"time"
	"sync"
	"errors"
)

const (
	BaseUrl = "https://httpapi.com/api/domains/%s.json?auth-userid=%d&api-key=%s&%s"
)

type ResellAPI struct {
	UserId  int
	Key     string
	Debug   bool
	mutex   sync.Mutex
	t1      time.Time
	t2      time.Time
}

func (resell *ResellAPI) request(u string, post bool) (*http.Response, error) {
	resell.mutex.Lock()
	defer func() {
		resell.t1 = resell.t2
		resell.t2 = time.Now()
		resell.mutex.Unlock()
	}()

	dur := time.Since(resell.t1)
	if dur.Nanoseconds() < int64(time.Second) {
		time.Sleep(time.Duration(int64(time.Second) - dur.Nanoseconds()))
	}
	if (post) {
		return http.PostForm(u, nil)
	}
	return http.Get(u)


}


func (resell *ResellAPI) Register(domain string, duration int, customerId int, contactId int) (bool, error) {

	command := "register"
	o := fmt.Sprintf("domain-name=%s&years=%d&ns=ns1.linode.com&ns=ns2.linode.com&customer-id=%d&reg-contact-id=%d&admin-contact-id=%d&tech-contact-id=%d&billing-contact-id=%d&invoice-option=KeepInvoice&protect-privacy=False",
		domain, duration, customerId, contactId, contactId, contactId, contactId)
	//https://httpapi.com/api/domains/register.json?auth-userid=0&api-key=key&domain-name=domain.com&years=1&ns=ns1.domain.com&ns=ns2.domain.com&customer-id=0&reg-contact-id=0&admin-contact-id=0&tech-contact-id=0&billing-contact-id=0&invoice-option=KeepInvoice&protect-privacy=False


	u := fmt.Sprintf(BaseUrl, command, resell.UserId, resell.Key, o)

	//log.Printf(u)
	res, err := resell.request(u, true)
	if err != nil {
		log.Println(err)
		return false,err
	}
	dec := json.NewDecoder(res.Body)
	var s map[string]string
	dec.Decode(&s)
	//log.Printf("%+v", s)
	if s["status"] == "error" {
		log.Printf(s["error"])
		return false, errors.New(s["error"])
	}
	if s["status"] == "Success" {
		return true, nil
	}
	//map[actiontype:AddNewDomain actionstatusdesc:Domain registration completed Successfully actiontypedesc:Registration of aureole-tech.com for 1 year unutilisedsellingamount:-11.160 sellingamount:-11.160 entityid:53356784 actionstatus:Success status:Success eaqid:236729894 customerid:10771231 description:aureole-tech.com invoiceid:41651548 sellingcurrencysymbol:USD]
	//return silo.register(command, domain, duration, coupon)
	return true, nil
}

func (resell *ResellAPI) Delete(orderId int) (bool, error) {
	command := "delete"

	//https://httpapi.com/api/domains/register.json?auth-userid=0&api-key=key&domain-name=domain.com&years=1&ns=ns1.domain.com&ns=ns2.domain.com&customer-id=0&reg-contact-id=0&admin-contact-id=0&tech-contact-id=0&billing-contact-id=0&invoice-option=KeepInvoice&protect-privacy=False


	u := fmt.Sprintf(BaseUrl, command, resell.UserId, resell.Key, fmt.Sprintf("order-id=%d", orderId))

	//log.Printf(u)
	res, err := resell.request(u, true)
	if err != nil {
		log.Println(err)
	}
	dec := json.NewDecoder(res.Body)
	var s map[string]string
	dec.Decode(&s)
	log.Printf("%+v", s)
	if s["status"] == "error" {
		log.Printf(s["error"])
		return false, errors.New(s["error"])
	}
	//return silo.register(command, domain, duration, coupon)
	return true, nil
}

func (resell *ResellAPI) Available(domains []string, tlds []string) (map[string]bool, error) {
	//https://httpapi.com/api/domains/available.json?auth-userid=0&api-key=key&domain-name=domain1&domain-name=domain2&tlds=com&tlds=net
	command := "available"
	ds := make([]string, len(domains))
	for i, d := range (domains) {
		ds[i] = fmt.Sprintf("domain-name=%s", d)
	}
	ts := make([]string, len(tlds))
	for i, t := range (tlds) {
		ts[i] = fmt.Sprintf("tlds=%s", t)
	}
	ds = append(ds, ts...)

	u := fmt.Sprintf(BaseUrl, command, resell.UserId, resell.Key, strings.Join(ds, "&"))
	log.Println(u)
	res, err := resell.request(u, false)
	if err != nil {
		log.Println(err)
	}
	dec := json.NewDecoder(res.Body)
	var s map[string]map[string]string
	dec.Decode(&s)
	//log.Printf("%+v", s)
	result := make(map[string]bool)
	for k, v := range (s) {
		result[k] = v["status"] == "available"
	}
	return result, nil

}
