package name

import (
	"fmt"
	"net/http"
	"log"
	"encoding/json"
	"bytes"
	"errors"
	"io/ioutil"
)

const (
	DevBaseUrl        = "http://api.dev.name.com%s"
	ProductionBaseUrl = "http://api.name.com%s"
	BaseUrl           = ProductionBaseUrl
)

type NameAPI struct {
	Account      string
	Token        string
	sessionToken string
}

type Status struct {
	Code    int `json:"code"`
	Message string `json:"message"`
}

func (api *NameAPI) Login() (bool, error) {
	u := "/api/login"
	form := map[string]string{"username": api.Account, "api_token": api.Token}
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.Encode(form)
	r, err := http.Post(fmt.Sprintf(BaseUrl, u), "application/json", &buffer)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer r.Body.Close()

	type Result struct {
		Status       Status   `json:"result"`
		SessionToken string `json:"session_token"`
	}
	decoder := json.NewDecoder(r.Body)
	result := Result{}
	err = decoder.Decode(&result)

	if err != nil {
		log.Println(err)
		return false, err
	}
	if result.Status.Code == 100 {
		api.sessionToken = result.SessionToken
		return true, nil
	}
	return false, errors.New(result.Status.Message)
}

func (api *NameAPI) Hello() {
	url := "/api/hello"

	req, err := http.NewRequest("GET", fmt.Sprintf(BaseUrl, url), nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Api-Session-Token", api.sessionToken)

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(r)
	}
	defer r.Body.Close()
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(r)
	}
	log.Println(string(content))

}

func (api *NameAPI) Available(keyword,tld string)(bool,error){
   url:="/api/domain/check"
   contents:=`
      {"keyword":"%s",
       "tlds":["%s"],
       "services":["availability"]
      }
      `
      var buffer bytes.Buffer
	contents=fmt.Sprintf(contents, keyword,tld)
	//log.Println(contents)
	buffer.WriteString(contents)
	req, err := http.NewRequest("POST", fmt.Sprintf(BaseUrl, url),  &buffer)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Api-Session-Token", api.sessionToken)
	req.Header.Add("content-type",  "application/json")
	r, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println(err)
		return false, err
	}
	defer r.Body.Close()

      

	type Result struct {
		Status       Status   `json:"result"`
                Domains      map[string]map[string]interface{}  `json:"domains"`
	}
	decoder := json.NewDecoder(r.Body)
	result := Result{}
	err = decoder.Decode(&result)

	if err != nil {
		log.Println(err)
		return false, err
	}
	//
	if result.Status.Code == 100 {
                if a, ok := result.Domains[fmt.Sprintf("%s.%s",keyword, tld)]["avail"].(bool); ok {
			    return a,nil
		}  
		log.Println(fmt.Sprintf("%s.%s",keyword, tld),result) 
                return false,nil         
		
	}
	return false, errors.New(result.Status.Message)
}


func (api *NameAPI) Register(domain string) (bool, error) {
	url := "/api/domain/create"
	contents := `
	{"domain_name":"%s",
	  "period":1,
	  "nameservers":["ns1,name.com","ns2.name.com"],
	  "contacts":[
	  {"type": ["registrant","administrative","billing","technical"],
	   "first_name":"FirstName",
	   "last_name":"LastName",
	    "organization":"organization",
	    "address_1":"address_1",
	    "address_2":"",
	    "city":"city",
	    "state":"state",
	    "zip":"",
	    "country":"",
	    "phone":"",
	    "fax":"",
	    "email":"aa@gmail.com"
	   }
	  ]
	}
	`
	var buffer bytes.Buffer
	contents=fmt.Sprintf(contents, domain)
	log.Println(contents)
	buffer.WriteString(contents)
	req, err := http.NewRequest("POST", fmt.Sprintf(BaseUrl, url),  &buffer)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Api-Session-Token", api.sessionToken)
	req.Header.Add("content-type",  "application/json")
	r, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println(err)
		return false, err
	}
	defer r.Body.Close()

	type Result struct {
		Status       Status   `json:"result"`
	}
	decoder := json.NewDecoder(r.Body)
	result := Result{}
	err = decoder.Decode(&result)

	if err != nil {
		log.Println(err)
		return false, err
	}
	log.Println(result)
	if result.Status.Code == 100 {
		return true, nil
	}
	return false, errors.New(result.Status.Message)
}

