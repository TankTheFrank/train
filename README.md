# Train

[![Build Status](https://travis-ci.org/shaoshing/train.png?branch=master)](https://travis-ci.org/shaoshing/train)

Assets Management Package for web app in Go. The main purpose of it is to introduce some good practices already existed in Ruby on Rails' [Assets Pipeline](http://guides.rubyonrails.org/asset_pipeline.html).

## Main features

* Organize assets with the [Include Directive](#include-directive).
* [Pipeline](#pipeline) for [SASS](https://github.com/sass/sassc) and [CoffeeScript](coffeescript.org) in the runtime.
* [Bundling and Fingerprinting Assets](#bundling-and-fingerprinting-assets) for production.

## Installation

Get the package:

```bash
$ go get github.com/tankthefrank/train
```

Install the command-line tool:

```bash
$ go build -o $GOPATH/bin/train github.com/tankthefrank/train/cmd
```

Install [node-sass](https://github.com/sass/node-sass), [CoffeeCcript](http://coffeescript.org/)

```bash
$ npm install -g node-sass
$ npm install -g coffee-script
```

### Prepare for the Pipeline feature

If planning to use SASS or CoffeeScript, you should run the `diagnose` command to see whether your environment is fit for the feature. Otherwise, skip to [the next section](#quick-example).

```bash
# Diagnose and follow the instructions to get your environment prepared
$ train diagnose

# If you experience `command not found` error, you should add $GOPATH/bin to $PATH
# or run the command as follow:
$ $GOPATH/bin/train
```

### Quick Example

```bash
$ cd $GOPATH/src/github.com/tankthefrank/train
$ go run example/main.go
# Visit localhost:8000 and play with the `include` directive and the SASS and CoffeeScript Pipeline.
```

In the example page, you can toggle the Include Directive feature, or try out the Pipeline feature by requesting a sass or coffee file.

### Use it in a beego project

```go

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/tankthefrank/train"
)

func init() {
	// tags
	beego.AddFuncMap("javascript_tag", train.JavascriptTag)
	beego.AddFuncMap("stylesheet_tag", train.StylesheetTag)
	beego.AddFuncMap("stylesheet_tag_with_param", train.StylesheetTagWithParam)

	// add the assets handler
	train.SetFileServer()
	beego.Handler(train.Config.AssetsUrl+"*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		train.ServeRequest(w, r)
	}))
}

```

### Use it in a generic project

First, allow train to handle assets requests by adding handler to the http.ServeMux:

```go
// Adding handler to the http.DefaultServeMux
train.ConfigureHttpHandler(nil)
http.ListenAndServe(":8000", nil)
```

For custom ServeMux that overwrites the DefaultServeMux, you will need to pass the mux to `train.ConfigureHttpHandler`:

```go
mux := http.NewServeMux()

...

train.ConfigureHttpHandler(mux)
http.ListenAndServe(":8000", mux)
```

Next, add the helper functions to templates so that Train can generate assets links for you:

```go
import "github.com/tankthefrank/train"
import "html/template"

...

tpl := template.New("home")
// Adding helpers
tpl.Funcs(template.FuncMap{
  "javascript_tag":            train.JavascriptTag,
  "stylesheet_tag":            train.StylesheetTag,
  "stylesheet_tag_with_param": train.StylesheetTagWithParam,
})
tpl.ParseFiles("home.html")
tpl.Execute(wr, nil)
```

Now in your template file, you can use the above helpers to include your assets:

```html
(example: home.html)

<html>
  <head>
  {{stylesheet_tag "main"}}
  {{stylesheet_tag "home"}}

  {{javascript_tag "main"}}
  {{javascript_tag "home"}}
  ...
  </head>
…
</html>
```

Train enforce the following assets hierarchy and generate asset paths accordingly:

```
Project Root
├── assets
│   ├── js // put js and coffee scripts here
│   │   ├── main.js
│   │   ├── home.coffee
│   └── css // put css and sass here
│       ├── main.sass
│       ├── home.css
```

## Include Directive

Train allows you specify dependency inside asset file by using the `include` directive, and when you include the file using Train's helper, Train will check the dependency and expand the file into related files.

Say you have the following files:

```
├── assets
│   ├── js
│   │   ├── base.js
│   │   ├── app.js   // depends on base.js
```

The regular way of insuring the dependency would be including both javascript files in the html file, something like this:

```html
<script src="/assets/js/basic.js"></script>
<script src="/assets/js/app.js"></script>
```

In the Train way, you can do it by specifying the dependency in app.js:

```js
//= require js/base
...
```

And then use the helper to include app.js:

```html
{{javascript_tag "app"}}
```

When request for the html, the content will become:

```html
<script src="/assets/js/basic.js?3392212"></script>
<script src="/assets/js/app.js?3392212"></script>
```

To use the include directive in css is similar to js:

```css
/*
 *= require css/base
 */
...
```

### SASS and CoffeeScript

The Include Directive is only available for js and css. However, SASS already has the @import directive, which is doing the same thing. For CoffeeScript, you will have to manage the dependencies in a regular way.

## Pipeline

When handling js or css request, Train will first look for the asset file with the same extension in the assets folder. If the file cannot be found, it will keep searching for a alternative extension, which is .sass/.scss for .css and .coffee for .js . When found, Train will convert the file into the desired extension.

Take a look at an simple example:

```
├── assets
│   ├── css
│   │   ├── app.sass
```

In the html, you include the sass file as if it is a css file:

```html
{{stylesheet_tag "app"}}
```

### Configuration

There are several configuration options related to the Pipeline feature:

```go
// From SASS's doc:
// When set to true, causes the line number and file where a selector is defined to be
// emitted into the compiled CSS in a format that can be understood by the browser. Useful in
// conjunction with [the FireSass Firebug extension](https://addons.mozilla.org/en-US/firefox/addon/103988)
// for displaying the Sass filename and line number.
train.Config.SASS.DebugInfo = true // false by default


// From SASS's doc:
// When set to true, causes the line number and file where a selector is defined to be emitted
// into the compiled CSS as a comment. Useful for debugging, especially when using imports and mixins.
train.Config.SASS.LineNumbers = true // false by default


// Show SASS and CoffeeScript errors.
train.Config.Verbose = true // false by default
```

## Bundling and Fingerprinting Assets

You probably want to merge or convert the assets in production site for performance concern. This is done by running Train's command-line tool `train` without any option:

```bash
$ cd project/root
$ train
-> clean bundled assets
-> copy assets from assets
-> bundle and compile assets
-> compress assets
-> Fingerprinting Assets
-> cleaup assets directory
```

The following example is what were generated after running the `train` command:

```
Project Root
├── assets
│   ├── js
│   │   ├── main.js
│   │   ├── app.js
│   │   ├── home.coffee
│   └── css
│       ├── main.sass
│       ├── home.css
├── static
│   ├── assets // generated by train
│   │   ├── manifest.txt
│   │   ├── js
│   │   │   ├── main.js
│   │   │   ├── main-223e2f3f9ca508630ead4db28042cc42.js
│   │   │   ├── app.js
│   │   │   ├── app-c5d14af50112f85c0aee9181b14f02e4.js
│   │   │   ├── home.js
│   │   │   ├── home-c471ecdacdaf77f591100c4cffd51f41.js
│   │   └── css
│   │       ├── main.css
│   │       ├── main-d208d2ef0e80f9a7d372f0bd681f8ade.css
│   │       ├── home.css
│   │       ├── home-924c344bccc46742a90835cc104dbe20.css
```

When Train detect the static/assets folder, it will disable the Include Directive and Pipeline features and serve from these static files directly. Template helpers will also stop expanding assets and generate with fingerprinted paths:

```go
{{javascript_tag "app"}}
{{stylesheet_tag "home"}}

// to

<script src="/static/assets/js/app-c5d14af50112f85c0aee9181b14f02e4.js"></script>
<link rel="stylesheet" href="/static/assets/css/home-924c344bccc46742a90835cc104dbe20.css">
```

### Why Fingerprinting?

From Rails' Assets Pipeline Document:

```
Fingerprinting is a technique that makes the name of a file dependent on the contents of the file.
When the file contents change, the filename is also changed. For content that is static or infrequently
changed, this provides an easy way to tell whether two versions of a file are identical, even across
different servers or deployment dates.
```

Checkout [its document](http://guides.rubyonrails.org/asset_pipeline.html#what-is-fingerprinting-and-why-should-i-care) for more details about this technique.

### Deploy to Production Server

There are two ways to deploy the Bundled and Fingerprinted assets to your server:

1. Run the `train` command in the production server after each deployment. By doing this you can make sure to update static/assets to the latest. This is the simples way, but it requires your server have NodeJS and required npm if you are using the Pipeline feature.

2. Run the `train` command in your local machine and upload the assets to the production server. With this way, the production server doesn't need to have NodeJS and required npm for the command.

Here is bash snippet to deploy assets using the second way:

```bash
SERVER="replace to your server's ssh address"
SERVER_PUBLIC="replace to your server's static path"

echo "Bundle assets"
$GOPATH/bin/train

if [[ $? != 0 ]] ; then
  echo "== fail to bundle assets"
  exit 1
fi

echo "Copy assets to $SERVER"
cd static
tar zcf assets.zip assets
scp assets.zip "$SERVER":assets.zip
ssh $SERVER "tar mxf assets.zip && sudo cp -r assets/* $SERVER_PUBLIC/assets/ && rm assets.zip"
rm -f assets.zip assets
cd -
```

## Status

Train is production ready, and has been used in our production site [Qortex](https://qortex.net). You are very welcome to report usage in your project.

Tested language / lib versions:

* Go: go1.2.1 darwin/amd64
* node-sass: 2.0.1
* CoffeeScript: 2.2.0

## Contribution

* Fork & Clone
* Make awesome changes (as well as tests)
* Run the tests
* Pull Request

### Run the tests

* Install Go (1.2.1)
* Run all the tests ``./test_all.sh``

## License

Train is released under the [MIT License](http://www.opensource.org/licenses/MIT).
