package devd

import (
	"testing"

	"github.com/GeertJohan/go.rice"
	"github.com/cortesi/devd/ricetemp"
	"github.com/cortesi/devd/termlog"
)

var formatURLTests = []struct {
	tls    bool
	addr   string
	port   int
	output string
}{
	{true, "127.0.0.1", 8000, "https://devd.io:8000"},
	{false, "127.0.0.1", 8000, "http://devd.io:8000"},
	{false, "127.0.0.1", 80, "http://devd.io"},
	{true, "127.0.0.1", 443, "https://devd.io"},
	{false, "127.0.0.1", 443, "http://devd.io:443"},
}

func TestFormatURL(t *testing.T) {
	for i, tt := range formatURLTests {
		url := formatURL(tt.tls, tt.addr, tt.port)
		if url != tt.output {
			t.Errorf("Test %d, expected \"%s\" got \"%s\"", i, tt.output, url)
		}
	}
}

func TestPickPort(t *testing.T) {
	_, err := pickPort("127.0.0.1", 8000, 10000, true)
	if err != nil {
		t.Errorf("Could not bind to any port: %s", err)
	}
	_, err = pickPort("127.0.0.1", 8000, 8000, true)
	if err == nil {
		t.Errorf("Expected not to be able to bind to any port")
	}

}

func fsEndpoint(s string) *filesystemEndpoint {
	e, _ := newFilesystemEndpoint(s)
	return e
}

func TestDevdRouteHandler(t *testing.T) {
	logger := termlog.NewLog()
	logger.Quiet()
	r := Route{"", "/", fsEndpoint("./testdata")}
	templates := ricetemp.MustMakeTemplates(rice.MustFindBox("templates"))

	devd := Devd{LivereloadRoutes: true}
	h := devd.RouteHandler(logger, r, templates)
	ht := handlerTester{t, h}

	AssertCode(t, ht.Request("GET", "/", nil), 200)
}

func TestDevdHandler(t *testing.T) {
	logger := termlog.NewLog()
	logger.Quiet()
	templates := ricetemp.MustMakeTemplates(rice.MustFindBox("templates"))

	devd := Devd{LivereloadRoutes: true, WatchPaths: []string{"./"}}
	devd.AddRoutes([]string{"./"})
	h, err := devd.Handler(logger, templates)
	if err != nil {
		t.Error(err)
	}
	ht := handlerTester{t, h}

	AssertCode(t, ht.Request("GET", "/", nil), 200)
	AssertCode(t, ht.Request("GET", "/nonexistent", nil), 404)
}

func TestGetTLSConfig(t *testing.T) {
	_, err := getTLSConfig("nonexistent")
	if err == nil {
		t.Error("Expected failure, found success.")
	}
	_, err = getTLSConfig("./testdata/certbundle.pem")
	if err != nil {
		t.Errorf("Could not get TLS config: %s", err)
	}
}
