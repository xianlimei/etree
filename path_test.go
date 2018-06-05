// Copyright 2015 Brett Vickers.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package etree

import "testing"

var testXML = `
<?xml version="1.0" encoding="UTF-8"?>
<bookstore xmlns:p="urn:books-com:prices">

	<!Directive>

	<book category="COOKING">
		<title lang="en">Everyday Italian</title>
		<author>Giada De Laurentiis</author>
		<year>2005</year>
		<p:price>30.00</p:price>
		<editor>Clarkson Potter</editor>
	</book>

	<book category="CHILDREN">
		<title lang="en" sku="150">Harry Potter</title>
		<author>J K. Rowling</author>
		<year>2005</year>
		<p:price>29.99</p:price>
		<editor></editor>
		<editor/>
	</book>

	<book category="WEB">
		<title lang="en">XQuery Kick Start</title>
		<author>James McGovern</author>
		<author>Per Bothner</author>
		<author>Kurt Cagle</author>
		<author>James Linn</author>
		<author>Vaidyanathan Nagarajan</author>
		<year>2003</year>
		<p:price>49.99</p:price>
		<editor>
		</editor>
	</book>

	<!-- Final book -->
	<book category="WEB" path="/books/xml">
		<title lang="en">Learning XML</title>
		<author>Erik T. Ray</author>
		<year>2003</year>
		<p:price>39.95</p:price>
	</book>

</bookstore>
`

type test struct {
	path   string
	result interface{}
}

type errorResult string

var tests = []test{

	// basic queries
	{"./bookstore/book/title", []string{"Everyday Italian", "Harry Potter", "XQuery Kick Start", "Learning XML"}},
	{"./bookstore/book/author", []string{"Giada De Laurentiis", "J K. Rowling", "James McGovern", "Per Bothner", "Kurt Cagle", "James Linn", "Vaidyanathan Nagarajan", "Erik T. Ray"}},
	{"./bookstore/book/year", []string{"2005", "2005", "2003", "2003"}},
	{"./bookstore/book/p:price", []string{"30.00", "29.99", "49.99", "39.95"}},
	{"./bookstore/book/isbn", nil},

	// descendant queries
	{"//title", []string{"Everyday Italian", "Harry Potter", "XQuery Kick Start", "Learning XML"}},
	{"//book/title", []string{"Everyday Italian", "Harry Potter", "XQuery Kick Start", "Learning XML"}},
	{".//title", []string{"Everyday Italian", "Harry Potter", "XQuery Kick Start", "Learning XML"}},
	{".//bookstore//title", []string{"Everyday Italian", "Harry Potter", "XQuery Kick Start", "Learning XML"}},
	{".//book/title", []string{"Everyday Italian", "Harry Potter", "XQuery Kick Start", "Learning XML"}},
	{".//p:price/.", []string{"30.00", "29.99", "49.99", "39.95"}},
	{".//price", []string{"30.00", "29.99", "49.99", "39.95"}},

	// positional queries
	{"./bookstore/book[1]/title", "Everyday Italian"},
	{"./bookstore/book[4]/title", "Learning XML"},
	{"./bookstore/book[5]/title", nil},
	{"./bookstore/book[3]/author[0]", "James McGovern"},
	{"./bookstore/book[3]/author[1]", "James McGovern"},
	{"./bookstore/book[3]/author[3]/./.", "Kurt Cagle"},
	{"./bookstore/book[3]/author[6]", nil},
	{"./bookstore/book[-1]/title", "Learning XML"},
	{"./bookstore/book[-4]/title", "Everyday Italian"},
	{"./bookstore/book[-5]/title", nil},

	// text queries
	{"./bookstore/book[author='James McGovern']/title", "XQuery Kick Start"},
	{"./bookstore/book[author='Per Bothner']/title", "XQuery Kick Start"},
	{"./bookstore/book[author='Kurt Cagle']/title", "XQuery Kick Start"},
	{"./bookstore/book[author='James Linn']/title", "XQuery Kick Start"},
	{"./bookstore/book[author='Vaidyanathan Nagarajan']/title", "XQuery Kick Start"},
	{"//book[p:price='29.99']/title", "Harry Potter"},
	{"//book[price='29.99']/title", "Harry Potter"},
	{"//book/price[text()='29.99']", "29.99"},
	{"//book/author[text()='Kurt Cagle']", "Kurt Cagle"},
	{"//book/editor[text()]", []string{"Clarkson Potter", "\n\t\t"}},

	// attribute queries
	{"./bookstore/book[@category='WEB']/title", []string{"XQuery Kick Start", "Learning XML"}},
	{"./bookstore/book[@path='/books/xml']/title", []string{"Learning XML"}},
	{"./bookstore/book[@category='COOKING']/title[@lang='en']", "Everyday Italian"},
	{"./bookstore/book/title[@lang='en'][@sku='150']", "Harry Potter"},
	{"./bookstore/book/title[@lang='fr']", nil},

	// parent queries
	{"./bookstore/book[@category='COOKING']/title/../../book[4]/title", "Learning XML"},

	// root queries
	{"/bookstore/book[1]/title", "Everyday Italian"},
	{"/bookstore/book[4]/title", "Learning XML"},
	{"/bookstore/book[5]/title", nil},
	{"/bookstore/book[3]/author[0]", "James McGovern"},
	{"/bookstore/book[3]/author[1]", "James McGovern"},
	{"/bookstore/book[3]/author[3]/./.", "Kurt Cagle"},
	{"/bookstore/book[3]/author[6]", nil},
	{"/bookstore/book[-1]/title", "Learning XML"},
	{"/bookstore/book[-4]/title", "Everyday Italian"},
	{"/bookstore/book[-5]/title", nil},

	// extra-whitespace queries
	{" / bookstore  / book [ 1 ] / title  ", "Everyday Italian"},
	{" / bookstore / book[ -5 ] / title ", nil},

	// quotes
	{".//book[@category='COOKING']/title[@lang='en']", "Everyday Italian"},
	{`.//book[@category="COOKING"]/title[@lang="en"]`, "Everyday Italian"},

	// union queries
	{"./bookstore/book[1]|book[4]/title", []string{"Everyday Italian", "Learning XML"}},
	{"./bookstore/(book[1]|book[4])/title", []string{"Everyday Italian", "Learning XML"}},
	{"./bookstore/book[author='Kurt Cagle'|author='James Linn']/title", "XQuery Kick Start"},
	{"./bookstore/book[(author='Kurt Cagle')|(author='James Linn')]/title", "XQuery Kick Start"},
	{"./bookstore/book[((author='Kurt Cagle')|(author='James Linn'))]/title", "XQuery Kick Start"},
	{"./bookstore/book[((author='Kurt Cagle')|author='James Linn')]/title", "XQuery Kick Start"},

	// bad paths
	{"./bookstore/book[]", errorResult("etree: path contains an empty filter expression.")},
	{"./bookstore/book[@category='WEB'", errorResult("etree: path has invalid filter [brackets].")},
	{"./bookstore/book[@category='WEB]", errorResult("etree: path has mismatched filter quotes.")},
	{"./bookstore/book[author]a", errorResult("etree: path has invalid filter [brackets].")},
}

func TestPath(t *testing.T) {
	doc := NewDocument()
	err := doc.ReadFromString(testXML)
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		path, err := CompilePath(test.path)
		if err != nil {
			if r, ok := test.result.(errorResult); !ok || err.Error() != string(r) {
				fail(t, test)
			}
			continue
		}

		// Test both FindElementsPath and FindElementPath
		element := doc.FindElementPath(path)
		elements := doc.FindElementsPath(path)

		switch s := test.result.(type) {
		case errorResult:
			fail(t, test)
		case nil:
			if element != nil || len(elements) != 0 {
				fail(t, test)
			}
		case string:
			if element == nil || element.Text() != s ||
				len(elements) != 1 || elements[0].Text() != s {
				fail(t, test)
			}
		case []string:
			if element == nil || element.Text() != s[0] || len(elements) != len(s) {
				fail(t, test)
				continue
			}
			for i := 0; i < len(elements); i++ {
				if elements[i].Text() != s[i] {
					fail(t, test)
					break
				}
			}
		}

	}
}

func fail(t *testing.T, test test) {
	t.Errorf("etree: failed test '%s'\n", test.path)
}

func TestAbsolutePath(t *testing.T) {
	doc := NewDocument()
	err := doc.ReadFromString(testXML)
	if err != nil {
		t.Error(err)
	}

	elements := doc.FindElements("//book/author")
	for _, e := range elements {
		title := e.FindElement("/bookstore/book[1]/title")
		if title == nil || title.Text() != "Everyday Italian" {
			t.Errorf("etree: absolute path test failed")
		}

		title = e.FindElement("//book[p:price='29.99']/title")
		if title == nil || title.Text() != "Harry Potter" {
			t.Errorf("etree: absolute path test failed")
		}
	}
}

func TestTokenizer(t *testing.T) {
	pathStr := `./ ((apple[ ((banana) | cat="dog") ] [-1])) | b /p:price`

	var c compiler2
	toks, err := c.tokenizePath(pathStr)
	if err != nil {
		t.Errorf("ERR: %v\n", err)
	}

	for i, tok := range toks {
		if tok.value != "" {
			t.Logf("%2d: tok=%s v='%s'\n", i, tokName[tok.id], tok.value)
		} else {
			t.Logf("%2d: tok=%s\n", i, tokName[tok.id])
		}
	}

	p, err := CompilePath2(pathStr)
	if err != nil {
		t.Errorf("ERR: %v\n", err)
	}

	_ = p
}

func TestPath2(t *testing.T) {
	doc := NewDocument()
	err := doc.ReadFromString(testXML)
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		path, err := CompilePath2(test.path)
		if err != nil {
			if _, ok := test.result.(errorResult); !ok {
				fail(t, test)
			}
			continue
		}

		// Test both FindElementsPath and FindElementPath
		element := doc.FindElementPath2(path)
		elements := doc.FindElementsPath2(path)

		switch s := test.result.(type) {
		case errorResult:
			fail(t, test)
		case nil:
			if element != nil || len(elements) != 0 {
				fail(t, test)
			}
		case string:
			if element == nil || element.Text() != s ||
				len(elements) != 1 || elements[0].Text() != s {
				fail(t, test)
			}
		case []string:
			if element == nil || element.Text() != s[0] || len(elements) != len(s) {
				fail(t, test)
				continue
			}
			for i := 0; i < len(elements); i++ {
				if elements[i].Text() != s[i] {
					fail(t, test)
					break
				}
			}
		}

	}
}

var tokName = []string{
	"nil",
	"'/'",
	"'//'",
	"'['",
	"']'",
	"'('",
	"')'",
	"'|'",
	"'='",
	"':'",
	"'@'",
	"'.'",
	"'..'",
	"'*'",
	"string",
	"ident",
	"num",
	"EOL",
}
