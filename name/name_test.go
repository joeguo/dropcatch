package name

import (
	"testing"
	"log"
)

func TestNameLogin(t *testing.T) {
	account := ""
	token := ""
	name := NameAPI{Account: account, Token: token}
	//ok, err := name.Login()
	//log.Println(ok, err)
	//name.Hello()
	ok, err := name.Available("depthsinme", "net")
	log.Println(ok, err)
	ok, err = name.Available("calart", "edu")
	log.Println(ok, err)
}
