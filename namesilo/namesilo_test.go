package namesilo

import (
	"testing"
	"log"
)
var(silo=&NamesiloAPI{Key:"YournameSiloAPIKey",Debug:true})


//func aTestNamesiloRegister(t *testing.T){
//	result,err:=silo.Register("dailyhiker.com",1,"")
//	log.Println(err, result)
//}
//
func TestNamesiloAvailable(t *testing.T) {
	result, err := silo.Available([]string{"hostye.com", "dailyhiker.com"})
	log.Println(err, result)
}


