package vcardenc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	DatumLineTCs = map[string]VcardDatum{
		"FN:Forrest Gump": StringDatum("FN", nil, "Forrest Gump"),
		"PHOTO;MEDIATYPE=image/gif:http://www.example.com/dir_photos/my_photo.gif": StringDatum("PHOTO", map[string][]string{"MEDIATYPE": []string{"image/gif"}}, "http://www.example.com/dir_photos/my_photo.gif"),
		"TEL;VALUE=uri;TYPE=home,voice:tel:+14045551212":                           StringDatum("TEL", map[string][]string{"TYPE": []string{"home", "voice"}, "VALUE": []string{"uri"}}, "tel:+14045551212"),
		// Currently fails because quoted values are still broken into lists on commas,
		// so metadata parsing needs an overhaul.
		"ADR;TYPE=work;LABEL=\"100 Waters Edge\\nBaytown, LA 30314\\nUnited States of America\":;;100 Waters Edge;Baytown;LA;30314;United States of America": SemicolonStructuredDatum("ADR", map[string][]string{"TYPE": []string{"work"}, "LABEL": []string{"\"100 Waters Edge\nBaytown, LA 30314\nUnited States of America\""}}, "", "", "100 Waters Edge", "Baytown", "LA", "30314", "United States of America"),
	}
)

func TestParseDatumLine(t *testing.T) {
	for line, expected := range DatumLineTCs {
		parsed, err := ParseDatumLine(line)
		assert.Nil(t, err)
		assert.EqualValues(t, expected, parsed)
	}
}
