package godaddy

var templateFile = `
{{define "Common"}}
<sCLTRID>{{.ActionID}}</sCLTRID>
<credential>
	<Account>{{.Account}}</Account>
	<Password>{{.Password}}</Password>
</credential>
{{end}}

{{define "CheckAvailability"}}
	<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	  <soap:Body>
		<CheckAvailability xmlns="http://wildwestdomains.com/webservices/">
	      {{template "Common" .}}
		  {{if .Domains}}
		  <sDomainArray>
			{{range .Domains}}
			<string>{{.}}</string>
			{{end}}
		  </sDomainArray>
		  {{end}}
		  {{if .Hosts}}
		  <sHostArray>
			{{range .Hosts}}
			<string>{{.}}</string>
			{{end}}
		  </sHostArray>
		  {{end}}
		  {{if .NS}}
		  <sNSArray>
			{{range .NS}}
			<string>string</string>
			{{end}}
		  </sNSArray>
		  {{end}}
		</CheckAvailability>
	   </soap:Body>
	</soap:Envelope>
{{end}}

{{define "ProcessRequest"}}
	<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	  <soap:Body>
		<ProcessRequest xmlns="http://wildwestdomains.com/webservices/">
		{{template "Common" .}}
		  <sRequestXML>
		   {{.XMLCommand|safe}}
		  </sRequestXML>
		</ProcessRequest>
	  </soap:Body>
	</soap:Envelope>
{{end}}

{{define "Describe"}}
  <soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <Describe xmlns="http://wildwestdomains.com/webservices/">
     {{template "Common" .}}
    </Describe>
  </soap:Body>
</soap:Envelope>
{{end}}

{{define "Cancel"}}
  <soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <Cancel xmlns="http://wildwestdomains.com/webservices/">
     {{template "Common" .}}
      <sType>{{.Type}}</sType>
      <sIDArray>
          {{range .IDs}}
           <string>{{.}}</string>
          {{end}}
      </sIDArray>
    </Cancel>
  </soap:Body>
</soap:Envelope>
{{end}}

{{define "GetAvailableBalance" }}
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetAvailableBalance xmlns="http://wildwestdomains.com/webservices/">
      {{template "Common" .}}
    </GetAvailableBalance>
  </soap:Body>
</soap:Envelope>
{{end}}

{{define "ContactInfo"}}
	 <fname>{{.FirstName}}</fname>
	 <lname>{{.LastName}}</lname>
	 {{if .Org}}<org>{{.Org}}</org> {{end}}
	 <email>{{.Email}}</email>
	 <sa1>{{.Street1}}</sa1>
	 {{if .Street2}}<sa2>{{.Street2}}</sa2> {{end}}
	 <city>{{.City}}</city>
	 <sp>{{.State}}</sp>
	 <pc>{{.Postcode}}</pc>
	 <cc>{{.Country}}</cc>
	 <phone>{{.Phone|safe}}</phone>
	 {{if .Fax}}<fax>{{.Fax|safe}}</fax>{{end}}
{{end}}

{{define "OrderDomains"}}
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <OrderDomains xmlns="http://wildwestdomains.com/webservices/">
      {{template "Common" .}}
      <shopper>
        <acceptOrderTOS>{{.Shopper.AcceptOrderTOS}}</acceptOrderTOS>
        <user>{{.Shopper.User}}</user>
        <pwd>{{.Shopper.Pwd}}</pwd>
        {{if .Shopper.PwdHint}}<pwdhint>{{.Shopper.PwdHint}}</pwdhint>  {{end}}
        <email>{{.Shopper.Email}}</email>
        <firstname>{{.Shopper.FirstName}}</firstname>
        <lastname>{{.Shopper.LastName}}</lastname>
        <phone>{{.Shopper.Phone|safe}}</phone>
        {{if .Shopper.Pin}} <pin>{{.Shopper.Pin}}</pin>{{end}}
      </shopper>
      <items>
      {{range .Domains}}
        <DomainRegistration>
          <order>
            <productid>{{.Order.ProductId}}</productid>
            <quantity>{{.Order.Quantity}}</quantity>
            <riid>{{.Order.RIID}}</riid>
            <duration>{{.Order.Duration}}</duration>
          </order>
          <sld>{{.Sld}}</sld>
          <tld>{{.Tld}}</tld>
          <period>{{.Period}}</period>
          {{if .Nexus}}
          <nexus>
            <category>{{.Nexus.Category}}</category>
            <use>{{.Nexus.Use}}</use>
            {{if .Nexus.Country}}
            <country>{{.Nexus.Country}}</country>
            {{end}}
          </nexus>
          {{end}}
          <nsArray>{{range .NS}}<NS><name>{{.}}</name></NS>{{end}}</nsArray>
        <registrant>
              {{template "ContactInfo" .Registrant}}
        </registrant>
        <admin>
            {{template "ContactInfo" .Admin}}
        </admin>
        <billing>
            {{template "ContactInfo" .Billing}}
        </billing>
        <tech>
            {{template "ContactInfo" .Tech}}
        </tech>
          <autorenewflag>{{.AutoRenew}}</autorenewflag>
        </DomainRegistration>
        {{end}}
      </items>
    </OrderDomains>
  </soap:Body>
</soap:Envelope>
{{end}}

{{define "Poll"}}
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <Poll xmlns="http://wildwestdomains.com/webservices/">
      {{template "Common" .}}
    </Poll>
  </soap:Body>
</soap:Envelope>
{{end}}
`
