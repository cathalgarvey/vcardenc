package vcardenc

import (
	"encoding/base64"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/cathalgarvey/wraprunes"
)

type valueType string

const (
	// StringValueType is the valueType for strings.
	StringValueType valueType = "string"
	// CommaStructuredValueType is the valueType for comma-separated values
	CommaStructuredValueType valueType = "commaStructured"
	// SemicolonStructuredValueType is the valueType for semicolon-separated values
	SemicolonStructuredValueType valueType = "semicolonStructured"
	// BinaryValueType is the valueType for binary data, base64 encoded and line-wrapped.
	BinaryValueType valueType = "binary"
)

var (
	validTypes = map[valueType]interface{}{"binary": nil, "string": nil, "commaStructured": nil, "semicolonStructured": nil}

	// ErrBadDatumType is returned on Output if the ValueType string is not part
	// of the expected set.
	ErrBadDatumType = errors.New("Unknown or malformed value type.")
)

func isValidType(t valueType) bool {
	_, ok := validTypes[t]
	return ok
}

// AttrMap is a map of attributes. When represented in card form, this is sorted,
// but this is only to assist in testing and isn't strictly part of the spec.
type AttrMap map[string][]string

// VcardDatum is a single or multiple line key:value entry in a vcard.
type VcardDatum struct {
	// FieldName is the name of the field. For "Begin:vcard"
	// for example, the fieldname is "begin" and the value is "vcard".
	FieldName string `json:"fieldName"`

	// Attrs are the semicolon-delimited key:value pairs following
	// the field name but prior to the colon and the value.
	Attrs AttrMap `json:"attrs,omitempty"`

	// ValueType is either "structured" or "string".
	ValueType valueType `json:"valueType"`

	// StructuredValue is the value if ValueType is "structured",
	// and is represented in vcard as a semicolon-delimited list of strings.
	StructuredValue []string `json:"structuredValue,omitempty"`

	// StringValue is the value if ValueType is "string",
	// and is represented in vcard as a free (possibly wrapped)string
	// following the colon.
	StringValue string `json:"stringValue,omitempty"`

	// BinaryValue is the value if ValueType is "binary"
	BinaryValue []byte `json:"binaryValue,omitempty"`
}

// StringDatum is a shortcut for making a string-type datum.
func StringDatum(fieldName string, attrs AttrMap, fieldValue string) VcardDatum {
	return VcardDatum{
		FieldName:   fieldName,
		ValueType:   StringValueType,
		StringValue: fieldValue,
		Attrs:       attrs,
	}
}

// DateDatum is a representable date, simply encoded as a string as YYYYMMDD.
func DateDatum(fieldName string, attrs AttrMap, t time.Time) VcardDatum {
	return StringDatum(fieldName, attrs, t.Format("20060102"))
}

// CommaStructuredDatum makes construcing a CommaStructuredValueType easy.
func CommaStructuredDatum(fieldName string, attrs AttrMap, fieldValues ...string) VcardDatum {
	return VcardDatum{
		FieldName:       fieldName,
		ValueType:       CommaStructuredValueType,
		StructuredValue: fieldValues,
		Attrs:           attrs,
	}
}

// SemicolonStructuredDatum makes construcing a SemicolonStructuredValueType easy.
func SemicolonStructuredDatum(fieldName string, attrs AttrMap, fieldValues ...string) VcardDatum {
	return VcardDatum{
		FieldName:       fieldName,
		ValueType:       SemicolonStructuredValueType,
		StructuredValue: fieldValues,
		Attrs:           attrs,
	}

}

type orderableKV struct {
	Key   string
	Value string
}

type orderableKVs []orderableKV

// Implement sort.Interface for orderableKVs

// Len is the number of elements in the collection.
func (okvs orderableKVs) Len() int {
	return len(okvs)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (okvs orderableKVs) Less(i, j int) bool {
	return okvs[i].Key > okvs[j].Key
}

// Swap swaps the elements with indexes i and j.
func (okvs orderableKVs) Swap(i, j int) {
	iv := okvs[i]
	okvs[i] = okvs[j]
	okvs[j] = iv
}

// Output converts a datum into a string for printing to a vCard buffer.
// specialRules, if provided, is a map of FieldNames to encoding functions
// that override this default behaviour, because vCard is the shittiest
// encoding format ever.
func (datum VcardDatum) Output(specialRules map[string]DatumEncoder) (string, error) {
	if specialFunc, ok := specialRules[datum.FieldName]; ok {
		return specialFunc(datum)
	}
	if !isValidType(datum.ValueType) {
		return "", ErrBadDatumType
	}
	var buf string
	buf += strings.ToUpper(datum.FieldName)
	kvs := make(orderableKVs, 0, len(datum.Attrs))
	for key, values := range datum.Attrs {
		safekey := escape(key)
		safeval := escapedJoin(values, ",")
		kvs = append(kvs, orderableKV{safekey, safeval})
	}
	sort.Sort(kvs)
	for _, kv := range kvs {
		buf += ";" + kv.Key + "=" + kv.Value
	}
	switch datum.ValueType {
	case StringValueType:
		{
			// Escape and rune-wrap the output data.
			buf += ":" + escape(datum.StringValue)
		}
	case SemicolonStructuredValueType:
		{
			buf += ":" + escapedJoin(datum.StructuredValue, ";")
		}
	case CommaStructuredValueType:
		{
			buf += ":" + escapedJoin(datum.StructuredValue, ",")
		}
	case BinaryValueType:
		{
			b64datum := base64.StdEncoding.EncodeToString(datum.BinaryValue)
			buf += ":" + b64datum
		}
	}
	buf = strings.Join(wraprunes.Wrap(buf, 75), "\n ")
	buf += "\n"
	return buf, nil
}
