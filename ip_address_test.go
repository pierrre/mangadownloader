package mangadownloader

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestIpAddress(t *testing.T) {
	response, err := http.Get("http://bot.whatismyipaddress.com")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	ipAddress := string(data)
	t.Logf("Ip address: %s", ipAddress)
}
