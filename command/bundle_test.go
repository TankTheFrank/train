package trainCommand

import (
	"io/ioutil"
	"testing"

	"github.com/shaoshing/gotest"
	"github.com/tankthefrank/train"
)

func init() {
	train.Config.Verbose = true
}

func assertEqual(path, content string) {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	assert.Equal(content, string(c))
}

func TestCommand(t *testing.T) {
	assert.Test = t

	assert.TrueM(prepareEnv(), "Unable to prepare env for cmd tests")

	removeAssets()

	copyAssets()

	assertEqual("static/assets/js/normal.js", "normal.js\n")
	assertEqual("static/assets/js/require.js", `//= require js/normal
//= require js/sub/require
require.js
`)
	assertEqual("static/assets/css/require.css", `/*
 *= require css/normal
 *= require css/sub/require
 */
require.css
`)

	bundleAssets()
	assertEqual("static/assets/js/normal.js", "normal.js\n")
	assertEqual("static/assets/js/require.js", `normal.js

sub/normal.js

sub/require.js

require.js
`)
	assertEqual("static/assets/css/require.css", `normal.css

sub/normal.css

sub/require.css

require.css
`)

	assertEqual("static/assets/css/font.css", `h1 {
  color: green; }`)
	assertEqual("static/assets/css/app.css", `h1 {
  color: green; }

h2 {
  color: green; }`)

	assertEqual("static/assets/css/scss.css", `h2 {
  color: green; }`)

	assertEqual("static/assets/js/app.js", `(function() {
  var a;

  a = 12;

}).call(this);
`)

	compressAssets()
	assertEqual("static/assets/js/require.js", `normal.js;sub/normal.js;sub/require.js;require.js;`)
	assertEqual("static/assets/js/require-min.js", `Please
Do
Not
Compresee
Me
`)

	fingerPrintAssets()
	assertEqual("static/assets/css/font.css", `h1{color:green}`) // should keep original assets
	train.LoadManifestInfo()
	assertEqual("static"+train.ManifestInfo["/assets/css/font.css"], `h1{color:green}`)

	removeAssets()
}
