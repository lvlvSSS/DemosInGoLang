package test

import (
	"regexp"
	"testing"
)

func TestRegexDemo(t *testing.T) {
	const (
		cityListReg = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`
	)
	contents := `<a href="http://www.zhenai.com/zhenghun/banan" class="">巴南</a>`
	compile := regexp.MustCompile(cityListReg)

	submatch := compile.FindAllSubmatch([]byte(contents), -1)

	for _, m := range submatch {
		t.Log("url:", string(m[1]), "city:", string(m[2]))
	}
}
