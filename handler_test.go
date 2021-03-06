package train

import (
	"github.com/shaoshing/gotest"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"
)

var httpClient = http.Client{}
var httpServer *httptest.Server

func initServer() {
	http.DefaultServeMux = http.NewServeMux()
	ConfigureHttpHandler(nil)
	httpServer = httptest.NewServer(nil)
}

func TestDeliverUnbundledAssets(t *testing.T) {
	initServer()

	assert.Test = t
	Config.BundleAssets = true

	assertAsset("/assets/static.txt", "static.txt\n", "text/plain")
	assertAsset("/assets/images/dummy.png", "dummy\n", "image/png")
	assert404("/assets/not/found.js")

	assertAsset("/assets/js/normal.js", "@normal.js\n", "application/javascript")
	assertAsset("/assets/css/normal.css", "normal.css\n", "text/css")

	bundledJs := assertAsset("/assets/js/require.js", "", "application/javascript")
	assert.Contain("@sub/normal.js", bundledJs)
	assert.Contain("@sub/normal1.coffee", bundledJs)
	assert.Contain("@sub/require.js", bundledJs)
	assert.Contain("@normal.js", bundledJs)
	assert.Contain("@normal1.coffee", bundledJs)
	assert.Contain("@require.js", bundledJs)

	assertAsset("/assets/css/require.css", `normal.css

sub/normal.css

sub/require.css

require.css
`, "text/css")

	assertAsset("/assets/css/app.css", `h1 {
  color: green; }

h2 {
  color: red; }`, "text/css")

	assertAsset("/assets/css/app2.css", `h2 {
  color: green; }

h3 {
  color: green; }`, "text/css")

	body, _, status := get("/assets/css/app.err.css")
	assert.Contain("Could not compile sass", body)
	assert.Equal(500, status)

	body, _, _ = get("/assets/js/app.js")
	assert.Contain("square = function(x)", body)

	body, _, status = get("/assets/js/app.err.js")
	assert.Contain("Could not compile coffee", body)
	assert.Equal(500, status)
}

func TestDeliverBundledAssets(t *testing.T) {
	Config.Mode = PRODUCTION_MODE
	initServer()

	assert.Test = t
	exec.Command("cp", "-rf", "assets/static", "./").Run()
	defer func() {
		exec.Command("rm", "-rf", "static").Run()
		Config.Mode = DEVELOPMENT_MODE
	}()

	assertAsset("/assets/app.js", "app.js\n", "application/javascript")
	assert404("/assets/normal.js")
}

func get(url string) (body, contentType string, status int) {
	response, err := httpClient.Get(httpServer.URL + url)
	if err != nil {
		panic(err)
	}
	b_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	body = string(b_body)
	contentType = response.Header.Get("Content-Type")
	status = response.StatusCode

	return
}

func assertAsset(url, expectedBody, expectedContentType string) string {
	body, contentType, _ := get(url)
	if len(expectedBody) != 0 {
		assert.Equal(expectedBody, body)
	}
	assert.Equal(true, strings.Index(contentType, expectedContentType) != -1)
	return body
}

func assert404(url string) {
	_, _, status := get(url)
	assert.Equal(404, status)
}
