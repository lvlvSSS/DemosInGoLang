package test

import (
	"go_tools/util/yaml"
	"os"
	"path/filepath"
	"testing"
)

type SiteConfig struct {
	HttpPort  int
	HttpsOn   bool
	Domain    string
	HttpsPort int
}

func TestResourcePath(t *testing.T) {
	target, _ := os.Getwd()
	t.Log(filepath.Join(target, "resources", "config", "log4go.yaml"))
	engine := yaml.New()

	siteConfig := &SiteConfig{}
	if value, err := engine.GetStruct("Site", siteConfig); err != nil {
		t.Logf("%v", err)
	} else {
		t.Logf("%v", value.(*SiteConfig).HttpsPort)
	}

}
