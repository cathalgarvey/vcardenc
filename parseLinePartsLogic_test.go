package vcardenc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type expectedParsedValues struct {
	ExpectedValue, Remaining string
}

var (
	StringParsingTestCases = map[string]expectedParsedValues{
		"This is the \\\"parsed\\\" bit\" and this remains": expectedParsedValues{
			"This is the \"parsed\" bit", " and this remains",
		},
		"Here's \\\"another\\\" test case with no overhang\"": expectedParsedValues{
			"Here's \"another\" test case with no overhang", "",
		},
	}

	SCParsingTestCases = map[string]expectedParsedValues{
		"18 something grove;town land;country\\;of\\;semicolons;foo continent": expectedParsedValues{
			"18 something grove", ";town land;country\\;of\\;semicolons;foo continent",
		},
		"No semicolon at all": expectedParsedValues{
			"No semicolon at all", "",
		},
	}

	CommaParsingTestCases = map[string]expectedParsedValues{
		"18 something grove,town land,country\\,of\\,semicolons,foo continent": expectedParsedValues{
			"18 something grove", ",town land,country\\,of\\,semicolons,foo continent",
		},
	}
)

func TestQuotedStringParsing(t *testing.T) {
	for tc, tcvals := range StringParsingTestCases {
		parsed, remaining, err := parseQuotedValue(tc, []rune{'"'}, true)
		assert.Nil(t, err)
		assert.Equal(t, tcvals.ExpectedValue, parsed)
		assert.Equal(t, tcvals.Remaining, remaining)
	}
}

func TestSemicolonDelimiterParsing(t *testing.T) {
	for tc, tcvals := range SCParsingTestCases {
		parsed, remaining, err := parseQuotedValue(tc, []rune{';'}, false)
		assert.Nil(t, err)
		assert.Equal(t, tcvals.ExpectedValue, parsed)
		assert.Equal(t, tcvals.Remaining, remaining)
	}
}

func TestCommaDelimiterParsing(t *testing.T) {
	for tc, tcvals := range CommaParsingTestCases {
		parsed, remaining, err := parseQuotedValue(tc, []rune{','}, false)
		assert.Nil(t, err)
		assert.Equal(t, tcvals.ExpectedValue, parsed)
		assert.Equal(t, tcvals.Remaining, remaining)
	}
}

var (
	semicolonStructuredExamples = map[string][]string{
		"foo;bar;qux\\;baz;fooobarrr": []string{"foo", "bar", "qux;baz", "fooobarrr"},
	}

	commaStructuredExamples = map[string][]string{
		"foo,bar,qux\\,baz,fooobarrr": []string{"foo", "bar", "qux,baz", "fooobarrr"},
	}
)

func TestSemicolonStructuredValueParsing(t *testing.T) {
	for tc, tcval := range semicolonStructuredExamples {
		tcparsed, err := parseStructuredValue(tc, ';')
		assert.Nil(t, err)
		assert.EqualValues(t, tcval, tcparsed)
	}
}

func TestCommaStructuredValueParsing(t *testing.T) {
	for tc, tcval := range commaStructuredExamples {
		tcparsed, err := parseStructuredValue(tc, ',')
		assert.Nil(t, err)
		assert.EqualValues(t, tcval, tcparsed)
	}
}

type expectedMetaParseValues struct {
	Attrs     AttrMap
	Remaining string
}

var (
	metaParseTCs = map[string]expectedMetaParseValues{
		";foo=bar;baz=qux,qum:some value we don't care about": expectedMetaParseValues{
			AttrMap{"foo": []string{"bar"}, "baz": []string{"qux", "qum"}}, "some value we don't care about",
		},
	}
)

func TestMetaDataParse(t *testing.T) {
	for metaEGstring, parsedTC := range metaParseTCs {
		attrParsed, remainderParsed, err := parseAttrs(metaEGstring)
		assert.Nil(t, err)
		assert.EqualValues(t, parsedTC.Attrs, attrParsed)
		assert.EqualValues(t, parsedTC.Remaining, remainderParsed)
	}
}
