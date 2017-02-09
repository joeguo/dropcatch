package resell

import (
	"testing"
	"log"
)

var resell = &ResellAPI{UserId:0, Key:"apikey"}

func aTestResellAvailable(t *testing.T) {


	result, err := resell.Available([]string{"hostye", "freelancemedia"}, []string{"com", "net"})
	log.Println(err, result)
}

func TestResellRegister(t *testing.T) {
	customerId := 0
	contactId := 0

	result, err := resell.Register("ipadder.com", 1, customerId , contactId)
	log.Println(err, result)
}

func aTestResellDelete(t *testing.T){
	orderId:=0
	result, err :=resell.Delete(orderId)
	log.Println(err, result)
}
