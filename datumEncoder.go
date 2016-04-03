package vcardenc

// DatumEncoder is a function that can handle a Datum.
// It is used to establish a map of special rules for particular
// fieldnames which should be used when encoding a fieldname, because
// vCard is absolute shit and requires crazy special-case rules.
type DatumEncoder func(VcardDatum) (string, error)
