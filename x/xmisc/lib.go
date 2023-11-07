package xmisc

import (
	"regexp"
	"strings"

	"github.com/fatih/color"
)

var (
	styleRegex  = regexp.MustCompile(`<style[\s\S]*?</style>`)
	bRegex      = regexp.MustCompile("<b>(.*?)</b>")
	strongRegex = regexp.MustCompile("<strong>(.*?)</strong>")
	iRegex      = regexp.MustCompile("<i>(.*?)</i>")
	hrRegex     = regexp.MustCompile("<hr(.*?)>")
	divRegex    = regexp.MustCompile("<div>(.*?)</div>")
	brRegex     = regexp.MustCompile("<br(.*?)>")
	pRegex      = regexp.MustCompile("<p>(.*?)</p>")
	ulRegex     = regexp.MustCompile("<ul>(.*?)</ul>")
	liRegex     = regexp.MustCompile("<li>(.*?)</li>")
	olRegex     = regexp.MustCompile("<ol>(.*?)</ol>")
	imgRegex    = regexp.MustCompile("<img(.*?)>")
	aRegex      = regexp.MustCompile(`<a href="([^"]*)">(.*?)</a>`)
	codeRegex   = regexp.MustCompile("<code>(.*?)</code>")
	tagRegex    = regexp.MustCompile("<(.*?)>")
)

func PurgeStyle(html string) string {
	return styleRegex.ReplaceAllString(html, "")
}

func TtyColor(text string) string {
	// b or strong
	text = bRegex.ReplaceAllString(text, color.YellowString("$1"))
	text = strongRegex.ReplaceAllString(text, color.YellowString("$1"))
	// i or em
	text = iRegex.ReplaceAllString(text, color.GreenString("$1"))
	// hr
	text = hrRegex.ReplaceAllString(text, "---")
	// div
	text = divRegex.ReplaceAllString(text, "$1")
	// br or p
	text = brRegex.ReplaceAllString(text, "\n")
	text = pRegex.ReplaceAllString(text, "$1")
	// ul
	text = ulRegex.ReplaceAllString(text, "$1")
	text = liRegex.ReplaceAllString(text, " - $1\n")
	// ol
	text = olRegex.ReplaceAllString(text, "$1")
	text = liRegex.ReplaceAllString(text, " - $1\n")
	// img
	text = imgRegex.ReplaceAllString(text, "(image)")

	// a
	text = aRegex.ReplaceAllString(text, color.BlueString("$2")+color.WhiteString(" ($1)"))
	// code
	text = codeRegex.ReplaceAllString(text, color.HiBlueString("$1"))
	// space
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	// fullwidth space
	text = strings.ReplaceAll(text, "&ensp;", "  ")

	// any other tags
	text = tagRegex.ReplaceAllString(text, "")
	return text
}
