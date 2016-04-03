/*Package vcardenc is an opinionated (which is to say, incorrect) encoder of vcards.
 */
package vcardenc

import "strings"

// Vcard contains data that can be encoded.
// A Vcard begins with a "Begin:vCard" datum and ends similarly,
// and some standards require an immediate version datum. In this
// implementation, if those entries are present as data they'll be
// ignored, and the vCard 4 version will be used even if the results
// are completely inconsistent and terrible. Sorry.
// Because vCard is awful and IDGAF this makes no guarantees of
// output validity, and data beyond the most basic behaviour will
// probably fail to encode properly. Don't blame me, blame vCard.
type Vcard struct {
	Data []VcardDatum
}

// Encode returns something that might parse as a vCard in client software.
func (v Vcard) Encode(specialRules map[string]DatumEncoder) (string, error) {
	var output = "BEGIN:VCARD\nVERSION:4.0\n"
	for _, d := range v.Data {
		field := strings.ToUpper(d.FieldName)
		if field == "BEGIN" || field == "VERSION" || field == "END" {
			continue
		}
		dout, err := d.Output(specialRules)
		if err != nil {
			return "", err
		}
		output += dout
	}
	return output + "END:VCARD", nil
}
