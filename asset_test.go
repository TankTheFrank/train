package train

import (
	"testing"

	"github.com/shaoshing/gotest"
)

func init() {
	Config.Verbose = true
}

func TestReadingNormalAssets(t *testing.T) {
	assert.Test = t
	var content string
	var err error

	content, _ = ReadAsset("/assets/js/normal.js")
	assert.Equal("@normal.js\n", content)

	content, _ = ReadAsset("/assets/js/sub/normal.js")
	assert.Equal("@sub/normal.js\n", content)

	_, err = ReadAsset("/assets/not/exists/normal.js")
	assert.Equal("Asset Not Found: /assets/not/exists/normal.js", err.Error())

	content, _ = ReadAsset("/assets/css/normal.css")
	assert.Equal("normal.css\n", content)

	content, _ = ReadAsset("/assets/css/sub/normal.css")
	assert.Equal("sub/normal.css\n", content)

	_, err = ReadAsset("/assets/not/exists/normal.css")
	assert.Equal("Asset Not Found: /assets/not/exists/normal.css", err.Error())

	_, err = ReadAsset("/assets/static.txt")
	assert.Equal("Unsupported Asset: /assets/static.txt", err.Error())
}

func TestReadingSass(t *testing.T) {
	assert.Test = t
	content, err := ReadAsset("/assets/css/app.css")
	if err != nil {
		panic(err)
	}
	assert.Contain("h1", content)
	assert.Contain("h2", content)

	content, err = ReadAsset("/assets/css/app2.css")
	if err != nil {
		panic(err)
	}
	assert.Contain("h2", content)
	assert.Contain("h3", content)
}

func TestReadingCoffee(t *testing.T) {
	assert.Test = t
	content, err := ReadAsset("/assets/js/app.js")
	if err != nil {
		panic(err)
	}
	assert.Contain("square", content)
}

func TestRequireDirective(t *testing.T) {
	assert.Test = t
	Config.BundleAssets = true
	defer func() {
		Config.BundleAssets = false
	}()

	content, err := ReadAsset("/assets/js/require2.js")
	if err != nil {
		panic(err)
	}
	assert.Equal(`@normal.js

@sub/normal.js

`, content)

	content, err = ReadAsset("/assets/css/require2.css")
	if err != nil {
		panic(err)
	}
	assert.Equal(`normal.css

sub/normal.css

`, content)
}

func TestReadingAssetsWithRequire(t *testing.T) {
	assert.Test = t
	Config.BundleAssets = true
	var content string
	var err error

	content, _ = ReadAsset("/assets/js/require.js")
	assert.Contain("@sub/normal.js", content)
	assert.Contain("@sub/normal1.coffee", content)
	assert.Contain("@sub/require.js", content)
	assert.Contain("@normal.js", content)
	assert.Contain("@normal1.coffee", content)
	assert.Contain("@require.js", content)

	content, _ = ReadAsset("/assets/js/require.coffee")
	assert.Contain("@normal.js", content)
	assert.Contain("@require.coffee", content)

	content, _ = ReadAsset("/assets/css/require.css")
	assert.Equal(`normal.css

sub/normal.css

sub/require.css

require.css
`, content)

	content, _ = ReadAsset("/assets/css/require.scss")
	assert.Equal(`.normal1-scss {
  color: blue; }
sub/normal.css

sub/require.css

.normal1-scss {
  color: blue; }

.foo .bar {
  color: red; }`, content)

	// In current can not remove require comment lines
	content, _ = ReadAsset("/assets/css/require.sass")
	assert.Equal(`.normal1-scss {
  color: blue; }
sub/normal.css

sub/require.css

.foo {
  color: red; }`, content)

	_, err = ReadAsset("/assets/js/error.js")
	assert.Equal(`Asset Not Found: not/found.js
--- required by /assets/js/error.js`, err.Error())

	_, err = ReadAsset("/assets/js/errors.js")
	assert.Equal(`Asset Not Found: not/found.js
--- required by js/error.js
--- required by /assets/js/errors.js`, err.Error())

	Config.BundleAssets = false
	content, _ = ReadAsset("/assets/css/require.css")
	assert.Equal(`/*
 *= require css/normal
 *= require css/sub/require
 */
require.css
`, content)
}
