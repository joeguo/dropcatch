package godaddy

import (
	"fmt"
	"net/http"
	"log"
	"bytes"
	"html/template"
	"encoding/xml"
	"strings"
)

const (
	ProductionBaseUrl = "https://api.wildwestdomains.com/wswwdapi/wapi.asmx"
	DevBaseUrl        = "https://api.ote.wildwestdomains.com/wswwdapi/wapi.asmx"
	BaseUrl           = DevBaseUrl
	Register          = "Register"
	Renew             = "renew"
	Transfer          = "transfer"
)

type GodaddyAPI struct {
	Account  string
	Password string
	plt   *template.Template
}

func New(account, password string) *GodaddyAPI {
	funcMap := template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
	plt := template.Must(template.New("Template").Funcs(funcMap).Parse(templateFile))
	return &GodaddyAPI{account, password, plt}
}

func (api *GodaddyAPI) send(action string, templ string, data map[string]interface {}) (string, error) {
	data["Account"] = api.Account
	data["Password"] = api.Password
	t := api.plt.Lookup(templ)
	var buffer bytes.Buffer
	err := t.Execute(&buffer, data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	l := buffer.Len()
	log.Println(buffer.String())
	req, err := http.NewRequest("POST", BaseUrl, &buffer)
	if err != nil {
		log.Println(err)
		return "", err
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", l))
	req.Header.Set("SOAPAction", fmt.Sprintf("http://wildwestdomains.com/webservices/%s", action))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	decoder := xml.NewDecoder(res.Body)
	defer res.Body.Close()

	var result string
	for {
		token, _ := decoder.Token()
		//fmt.Println(token)
		if token == nil {
			break
		}
		switch startElement := token.(type) {
		case xml.StartElement:
			if startElement.Name.Local == fmt.Sprintf("%sResult", action) {
				decoder.DecodeElement(&result, &startElement)
			}
		}
	}
	return result, nil
}

func (api *GodaddyAPI) Cancel(actionID string, ctype string, ids[]string) {
	data := make(map[string]interface {})
	data["ActionID"] = actionID
	data["Type"] = ctype
	data["IDs"] = ids
	result, err := api.send("Cancel", "Cancel", data)
	log.Println(result, err)
}
func (api *GodaddyAPI) CheckAvailability(actionID string, domains, hosts, ns []string) (map[string]bool, error) {
	data := make(map[string]interface {})
	data["ActionID"] = actionID
	data["Domains"] = domains
	data["Hosts"] = hosts
	data["NS"] = ns
	result, err := api.send("CheckAvailability", "CheckAvailability", data)
	if err != nil {
		return nil, err
	}
	type Domain struct {
		Name          string `xml:"name,attr"`
		Available     int  `xml:"avail,attr"`
		Backorderable int   `xml:"canBackorder,attr"`
	}
	type Check struct {
		Domain []Domain   `xml:"domain"`
	}
	available := make(map[string]bool)
	check := Check{}
	err = xml.Unmarshal([]byte(result), &check)
	if err != nil {
		log.Println(err)
	}
	for _, d := range check.Domain {
		available[d.Name] = d.Available > 0
	}
	return available, nil
}
func (api *GodaddyAPI) CheckUser() {

}
func (api *GodaddyAPI) Describe(actionID string) (string, error) {
	data := make(map[string]interface {})
	data["ActionID"] = actionID
	return api.send("Describe", "Describe", data)
}
func (api *GodaddyAPI) DomainForwarding() {

}

func (api *GodaddyAPI) GetDomainAlertCredits() {

}

func (api *GodaddyAPI) GetExpiringNameList() {

}
func (api *GodaddyAPI) GetMonitoredDomainList() {

}
func (api *GodaddyAPI) Info() {

}
func (api *GodaddyAPI) ManageTransfer() {

}
func (api *GodaddyAPI) NameGen() {

}

func (api *GodaddyAPI) NameGenDB() {

}

func (api *GodaddyAPI) NameGenDBWithTimeLimit() {

}

func (api *GodaddyAPI) OrderCredits() {

}
func (api *GodaddyAPI) OrderDomainBackOrders(actionID string, shopper *Shopper, domains []DomainRegistration) {
	data := make(map[string]interface {})
	data["ActionID"] = actionID
	data["Shopper"] = shopper
	data["Domains"] = domains
	log.Println(api.send("OrderDomainBackOrders", "OrderDomainBackOrders", data))
}
func (api *GodaddyAPI) OrderDomainPrivacy() {

}
func (api *GodaddyAPI) OrderDomainRenewals() {

}
func (api *GodaddyAPI) OrderDomains(actionID string, shopper *Shopper, domains []DomainRegistration) {
	data := make(map[string]interface {})
	data["ActionID"] = actionID
	data["Shopper"] = shopper
	data["Domains"] = domains
	log.Println(api.send("OrderDomains", "OrderDomains", data))
}
func (api *GodaddyAPI) OrderPrivateDomainRenewals() {

}
func (api *GodaddyAPI) OrderDomainTransfers() {

}
func (api *GodaddyAPI) OrderResourceRenewals() {

}
func (api *GodaddyAPI) Poll(actionID string)(string,error) {
	data := make(map[string]interface {})
	data["ActionID"] = actionID
	return api.send("Poll", "Poll", data)
}
func (api *GodaddyAPI) ProcessRequest(actionID string) bool {
	data := make(map[string]interface {})
	var buffer bytes.Buffer
	command := fmt.Sprintf(`<wapi clTRID="%s" account="%s" pwd="%s"><manage><script cmd="reset"/></manage></wapi>`, actionID, api.Account, api.Password)
	xml.EscapeText(&buffer, []byte(command))
	data["ActionID"] = actionID
	data["XMLCommand"] = buffer.String()
	result, err := api.send("ProcessRequest", "ProcessRequest", data)
	if err != nil {
		return false
	}
	if strings.Contains(result, "scripting status reset") {
		return true
	}
	return false

}
func (api *GodaddyAPI) RemoveDomainAlert() {

}
func (api *GodaddyAPI) ResetPassword() {

}
func (api *GodaddyAPI) SetDomainLocking() {

}
func (api *GodaddyAPI) SetShopperInfo() {

}
func (api *GodaddyAPI) SetupDomainAlert() {

}
func (api *GodaddyAPI) UpdateDomainAlert() {

}

func (api *GodaddyAPI) UpdateDomainContact() {

}
func (api *GodaddyAPI) UpdateDomainForwarding() {

}
func (api *GodaddyAPI) UpdateDomainMasking() {

}
func (api *GodaddyAPI) UpdateDomainOwnership() {

}
func (api *GodaddyAPI) UpdateNameServer() {

}
func (api *GodaddyAPI) CreateNewShopper() {

}
func (api *GodaddyAPI) ValidateRegistration() {

}
func (api *GodaddyAPI) GetAvailableBalance(actionID string) {
	data := make(map[string]interface {})
	data["ActionID"] = actionID
	result, err := api.send("GetAvailableBalance", "GetAvailableBalance", data)
	log.Println(result, err)
}
