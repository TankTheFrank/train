package train

import (
	"github.com/shaoshing/gotest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestHelpers(t *testing.T) {
	assert.Test = t
	Config.BundleAssets = false
	now := time.Now()
	stamp := strconv.FormatInt(now.Unix(), 10)
	updateAssetTimes(now)

	assert.Equal(`<script src="/assets/js/normal.js?`+stamp+`"></script>`, string(JavascriptTag("normal")))
	assert.Equal(`<script src="/assets/js/normal.js?`+stamp+`"></script>
<script src="/assets/js/normal1.js?`+stamp+`"></script>
<script src="/assets/js/sub/normal.js?`+stamp+`"></script>
<script src="/assets/js/sub/normal1.js?`+stamp+`"></script>
<script src="/assets/js/sub/require.js?`+stamp+`"></script>
<script src="/assets/js/require.js?`+stamp+`"></script>`, string(JavascriptTag("require")))

	assert.Equal(`<link type="text/css" rel="stylesheet" href="/assets/css/normal.css?`+stamp+`">`, string(StylesheetTag("normal")))
	assert.Equal(`<link type="text/css" rel="stylesheet" href="/assets/css/normal.css?`+stamp+`">
<link type="text/css" rel="stylesheet" href="/assets/css/sub/normal.css?`+stamp+`">
<link type="text/css" rel="stylesheet" href="/assets/css/sub/require.css?`+stamp+`">
<link type="text/css" rel="stylesheet" href="/assets/css/require.css?`+stamp+`">`, string(StylesheetTag("require")))

	assert.Equal(`<link type="text/css" rel="stylesheet" href="/assets/css/normal.css?`+stamp+`" media="print">`, string(StylesheetTagWithParam("normal", `media="print"`)))

	assert.Equal(`<script src="/assets/js/app.js?`+stamp+`"></script>`, string(JavascriptTag("app")))
	assert.Equal(`<link type="text/css" rel="stylesheet" href="/assets/css/app.css?`+stamp+`">`, string(StylesheetTag("app")))

	Config.Mode = PRODUCTION_MODE
	defer func() {
		Config.Mode = DEVELOPMENT_MODE
	}()
	ManifestInfo = FpAssets{
		"/assets/js/require.js":  "/assets/js/require-fingerprintinghash.js",
		"/assets/css/require.css": "/assets/css/require-fingerprintinghash.css",
	}

	assert.Equal(`<script src="/assets/js/require-fingerprintinghash.js"></script>`, string(JavascriptTag("require")))
	assert.Equal(`<link type="text/css" rel="stylesheet" href="/assets/css/require-fingerprintinghash.css">`, string(StylesheetTag("require")))

	ManifestInfo = FpAssets{}
}

func updateAssetTimes(t time.Time) {
	filepath.Walk(Config.AssetsPath, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			os.Chtimes(filePath, t, t)
		}
		return nil
	})
}
