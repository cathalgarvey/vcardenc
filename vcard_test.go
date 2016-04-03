package vcardenc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	wikipediaCard         = "BEGIN:VCARD\nVERSION:4.0\nN:Gump;Forrest;;;\nFN:Forrest Gump\nORG:Bubba Gump Shrimp Co.\nTITLE:Shrimp Man\nPHOTO;MEDIATYPE=image/gif:http://www.example.com/dir_photos/my_photo.gif\nTEL;VALUE=uri;TYPE=work,voice:tel:+11115551212\nTEL;VALUE=uri;TYPE=home,voice:tel:+14045551212\nADR;TYPE=work;LABEL=\"100 Waters Edge\\nBaytown, LA 30314\\nUnited States of A\n merica\":;;100 Waters Edge;Baytown;LA;30314;United States of America\nADR;TYPE=home;LABEL=\"42 Plantation St.\\nBaytown, LA 30314\\nUnited States of\n  America\":;;42 Plantation St.;Baytown;LA;30314;United States of America\nEMAIL:forrestgump@example.com\nREV:20080424T195243Z\nEND:VCARD"
	wikipediaCardTestCase = Vcard{
		Data: []VcardDatum{
			SemicolonStructuredDatum("N", nil, "Gump", "Forrest", "", "", ""),
			StringDatum("FN", nil, "Forrest Gump"),
			StringDatum("ORG", nil, "Bubba Gump Shrimp Co."),
			StringDatum("TITLE", nil, "Shrimp Man"),
			StringDatum("PHOTO", map[string][]string{"MEDIATYPE": []string{"image/gif"}}, "http://www.example.com/dir_photos/my_photo.gif"),
			StringDatum("TEL", map[string][]string{"TYPE": []string{"work", "voice"}, "VALUE": []string{"uri"}}, "tel:+11115551212"),
			StringDatum("TEL", map[string][]string{"TYPE": []string{"home", "voice"}, "VALUE": []string{"uri"}}, "tel:+14045551212"),
			SemicolonStructuredDatum("ADR", map[string][]string{"TYPE": []string{"work"}, "LABEL": []string{"\"100 Waters Edge\nBaytown, LA 30314\nUnited States of America\""}}, "", "", "100 Waters Edge", "Baytown", "LA", "30314", "United States of America"),
			SemicolonStructuredDatum("ADR", map[string][]string{"TYPE": []string{"home"}, "LABEL": []string{"\"42 Plantation St.\nBaytown, LA 30314\nUnited States of America\""}}, "", "", "42 Plantation St.", "Baytown", "LA", "30314", "United States of America"),
			StringDatum("EMAIL", nil, "forrestgump@example.com"),
			StringDatum("REV", nil, "20080424T195243Z"),
		},
	}
)

func TestSimpleEncode(t *testing.T) {
	enctest, err := wikipediaCardTestCase.Encode(nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, wikipediaCard, enctest)
}
