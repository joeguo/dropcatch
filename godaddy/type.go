package godaddy

type OrderItem struct {
	ProductId int
	Quantity  int
	RIID      string
	Duration  int
}

type Contact struct {
	FirstName string
	LastName  string
	Org       string
	Email     string
	Street1   string
	Street2   string
	City      string
	State     string
	Postcode  string
	Country   string
	Phone     string
	Fax       string
}
type Shopper struct {
	AcceptOrderTOS string
	User           string
	Pwd            string
	PwdHint        string
	Email          string
	FirstName      string
	LastName       string
	Phone          string
	Pin            string
}

type Nexus struct {
	Category string
	Use      string
	Country  string
}


type DomainRegistration struct {
	Order     OrderItem
	Sld       string
	Tld       string
	Period    int
	Registrant *Contact
	Admin *Contact
	Billing *Contact
	Tech *Contact
	AutoRenew int
	NS       []string
	Nexus     *Nexus
}

