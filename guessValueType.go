package vcardenc

var (
	// Field names on which to guess value type, bearing in mind this is vCard 4.0
	// and that I'm largely working from wikipedia to avoid headaches.
	knownSemicolonStructuredFields = []string{"ADR", "N"}
	knownCommaStructuredFields     = []string{"NICKNAME"}
	knownBase64Fields              = []string{}
)

// guess and return a valueType. Default in case of stupid is StringValueType
// TODO: this should look for v4.0 style type hints in attrs, to disambiguate
// fieldNames that can have URI, data-URI, or raw base64 datatypes.
func guessValueType(fieldName string, attrs AttrMap, rawValue string) valueType {
	if stringSliceContains(knownSemicolonStructuredFields, fieldName) {
		return SemicolonStructuredValueType
	}
	if stringSliceContains(knownCommaStructuredFields, fieldName) {
		return CommaStructuredValueType
	}
	if stringSliceContains(knownBase64Fields, fieldName) {
		return BinaryValueType
	}
	// TODO: Unfinished
	return StringValueType
}
