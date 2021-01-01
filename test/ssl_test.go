package test

import (
	"fmt"
	"go_tools/ssl_creator"
	"testing"
)

func TestCreateSimpleSSLCert(t *testing.T) {
	cert, key, _ := ssl_creator.CreateSimpleSSLCert("2021-Feb-03", "127.0.0.1", []string{"Example org."}, []string{"bookd"}, "Web app")

	fmt.Println(cert.String())
	fmt.Println("----------")
	fmt.Println(key.String())

}
