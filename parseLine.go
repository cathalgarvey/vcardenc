package vcardenc

import (
	"encoding/base64"
	"errors"
	"strings"
)

// The following code assumes that lines have already been parsed and reconstructed
// from wrapped representations!

var (
	emptyDatum = VcardDatum{}

	// ErrDatumLineColonNotFound is returned when no FieldName[;:]
	// delimiter is found
	ErrDatumLineColonNotFound = errors.New("Could not find a colon while parsing a datum line")

	// ErrFailedToParseQuotedString is returned if the closing
	// quotation mark is not found in a quoted string
	ErrFailedToParseQuotedString = errors.New("Failed to parse the end of a quoted string")

	// ErrFailedToParseStructuredValue is returned if a next value
	// from a structured value list is not parsed successfully.
	ErrFailedToParseStructuredValue = errors.New("Failed to extract next value from structured value list")

	// ErrNoMetadataFound is returned if a line suggests metadata but fails to deliver
	ErrNoMetadataFound = errors.New("Metadata suggested by non-immediate colon, but leading semicolon not found")

	// ErrBadMetadata is returned if metadata is apparently malformed
	ErrBadMetadata = errors.New("Metadata appears malformed")
)

// ParseDatumLine accepts a pre-unwrapped line of data and parses it into
// three chunks; name, attr, value. These are then decoded to a VcardDatum.
func ParseDatumLine(line string) (parsed VcardDatum, err error) {
	fn, attrMap, val, err := splitDatumLine(line)
	if err != nil {
		return emptyDatum, err
	}
	vt := guessValueType(fn, attrMap, val)
	finishedDatum := VcardDatum{
		FieldName: fn,
		Attrs:     attrMap,
		ValueType: vt,
	}
	switch vt {
	case StringValueType:
		{
			finishedDatum.StringValue = val
		}
	case CommaStructuredValueType:
		{
			sval, err := parseStructuredValue(val, ',')
			if err != nil {
				return emptyDatum, err
			}
			finishedDatum.StructuredValue = sval
		}
	case SemicolonStructuredValueType:
		{
			sval, err := parseStructuredValue(val, ';')
			if err != nil {
				return emptyDatum, err
			}
			finishedDatum.StructuredValue = sval
		}
	case BinaryValueType:
		{
			dval, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				return emptyDatum, err
			}
			finishedDatum.BinaryValue = dval
		}
	default:
		panic("Unexpected valueType returned from guessValueType (?)")
	}
	return finishedDatum, nil
}

// Parse a comma-or-semicolon seperated value string of escaped shit into a slice.
func parseStructuredValue(valS string, delimiter rune) (sval []string, err error) {
	var (
		parsedVal string
	)
	delimString := string([]rune{delimiter})
	for {
		parsedVal, valS, err = parseQuotedValue(valS, []rune{delimiter}, false)
		if err != nil {
			return nil, err
		}
		if len(valS) == 0 || valS == "\n" {
			break
		}
		sval = append(sval, parsedVal)
		if valS[0:1] == delimString {
			valS = valS[1:]
		}
	}
	return sval, nil
}

// Parses an un-wrapped line to the three key portions of a vCard datum.
func splitDatumLine(line string) (fieldName string, attrs AttrMap, value string, err error) {
	fieldName, line, err = parseFieldName(line)
	if err != nil {
		return "", nil, "", err
	}
	attrs, value, err = parseAttrs(line)
	if err != nil {
		return "", nil, "", err
	}
	return fieldName, attrs, value, nil
}

func parseFieldName(line string) (fieldName, remainder string, err error) {
	colonIndex := strings.IndexRune(line, ':')
	if colonIndex == -1 {
		return "", "", ErrDatumLineColonNotFound
	}
	fieldNameDelimitIndex := colonIndex
	semicolonIndex := strings.IndexRune(line, ';')
	if semicolonIndex != -1 && semicolonIndex < fieldNameDelimitIndex {
		fieldNameDelimitIndex = semicolonIndex
	}
	fieldName = line[:fieldNameDelimitIndex]
	remainder = line[fieldNameDelimitIndex:]
	return fieldName, remainder, nil
}

// parses key=<value>;key=<value>:datumValue where <value> may be escaped
// or quoted.
func parseAttrs(line string) (attrs AttrMap, value string, err error) {
	// First deal with case where there's no attrs at all:
	if line[:1] == ":" {
		return nil, line[1:], nil
	}
	// If it's not a colon, first rune must be a semicolon.
	if line[:1] != ";" {
		return nil, "", ErrNoMetadataFound
	}
	line = line[1:]
	attrs = make(AttrMap)
	for {
		var rawMetaValue string
		nextDelimiter := strings.IndexRune(line, '=')
		if nextDelimiter == -1 && line[:1] != ":" {
			return nil, "", ErrBadMetadata
		}
		metaFieldName := line[:nextDelimiter]
		line = line[nextDelimiter+1:]
		if line[:1] == "\"" {
			rawMetaValue, line, err = parseQuotedValue(line[1:], []rune{'"'}, true)
		} else {
			rawMetaValue, line, err = parseQuotedValue(line, []rune{';', ':'}, false)
		}
		if err != nil {
			return nil, "", err
		}
		if len(rawMetaValue) == 0 {
			return nil, "", ErrBadMetadata
		}
		attrs[metaFieldName] = []string{} // []string{metaValue}
		for {
			var nextVal string
			if len(rawMetaValue) == 0 {
				break
			}
			nextVal, rawMetaValue, err = parseQuotedValue(rawMetaValue, []rune{',', ':'}, false)
			if err != nil {
				return nil, "", ErrBadMetadata
			}
			if len(nextVal) == 0 {
				break
			}
			if len(rawMetaValue) > 0 && rawMetaValue[:1] == "," {
				rawMetaValue = rawMetaValue[1:]
			}
			attrs[metaFieldName] = append(attrs[metaFieldName], nextVal)
		}
		if line[:1] == ":" {
			line = line[1:]
			break
		}
		if line[:1] == ";" {
			line = line[1:]
			continue
		}
	}
	return attrs, line, nil
}
