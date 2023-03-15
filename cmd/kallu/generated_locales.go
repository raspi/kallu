// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

package main

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type dictionary struct {
	index []uint32
	data  string
}

func (d *dictionary) Lookup(key string) (data string, ok bool) {
	p, ok := messageKeyToIndex[key]
	if !ok {
		return "", false
	}
	start, end := d.index[p], d.index[p+1]
	if start == end {
		return "", false
	}
	return d.data[start:end], true
}

func init() {
	dict := map[string]catalog.Dictionary{
		"en":    &dictionary{index: enIndex, data: enData},
		"fi_FI": &dictionary{index: fi_FIIndex, data: fi_FIData},
	}
	fallback := language.MustParse("en")
	cat, err := catalog.NewFromMap(dict, catalog.Fallback(fallback))
	if err != nil {
		panic(err)
	}
	message.DefaultCatalog = cat
}

var messageKeyToIndex = map[string]int{
	"Disable color output":                          7,
	"Examples:":                                     11,
	"Full year:":                                    12,
	"How many months per line":                      2,
	"How many next months":                          0,
	"How many previous months":                      1,
	"Month 1-12 (defaults to current month)":        4,
	"One calendar at a time:":                       14,
	"Only one month, equivalent to -next 0 -prev 0": 6,
	"Only this month:":                              13,
	"Parameters:":                                   9,
	"Print full year":                               5,
	"Start day for week 0-6 (sun-sat)":              8,
	"Year (defaults to current year)":               3,
	"april":                                         28,
	"august":                                        32,
	"count must be > 0":                             17,
	"december":                                      36,
	"default: %q":                                   10,
	"february":                                      26,
	"invalid month: %d":                             15,
	"invalid starting day of week: %d":              16,
	"january":                                       25,
	"july":                                          31,
	"june":                                          30,
	"march":                                         27,
	"may":                                           29,
	"november":                                      35,
	"october":                                       34,
	"september":                                     33,
	"short.friday":                                  22,
	"short.monday":                                  18,
	"short.saturday":                                23,
	"short.sunday":                                  24,
	"short.thursday":                                21,
	"short.tuesday":                                 19,
	"short.wednesday":                               20,
	"week":                                          37,
}

var enIndex = []uint32{ // 39 elements
	// Entry 0 - 1F
	0x00000000, 0x00000015, 0x0000002e, 0x00000047,
	0x00000067, 0x0000008e, 0x0000009e, 0x000000cc,
	0x000000e1, 0x00000102, 0x0000010e, 0x0000011d,
	0x00000127, 0x00000132, 0x00000143, 0x0000015b,
	0x00000170, 0x00000194, 0x000001a6, 0x000001aa,
	0x000001ae, 0x000001b2, 0x000001b6, 0x000001ba,
	0x000001be, 0x000001c2, 0x000001ca, 0x000001d3,
	0x000001d9, 0x000001df, 0x000001e3, 0x000001e8,
	// Entry 20 - 3F
	0x000001ed, 0x000001f4, 0x000001fe, 0x00000206,
	0x0000020f, 0x00000218, 0x0000021d,
} // Size: 180 bytes

const enData string = "" + // Size: 541 bytes
	"\x02How many next months\x02How many previous months\x02How many months " +
	"per line\x02Year (defaults to current year)\x02Month 1-12 (defaults to c" +
	"urrent month)\x02Print full year\x02Only one month, equivalent to -next " +
	"0 -prev 0\x02Disable color output\x02Start day for week 0-6 (sun-sat)" +
	"\x02Parameters:\x02default: %[1]q\x02Examples:\x02Full year:\x02Only thi" +
	"s month:\x02One calendar at a time:\x02invalid month: %[1]d\x02invalid s" +
	"tarting day of week: %[1]d\x02count must be > 0\x02mon\x02tue\x02wed\x02" +
	"thu\x02fri\x02sat\x02sun\x02january\x02february\x02march\x02april\x02may" +
	"\x02june\x02july\x02august\x02september\x02october\x02november\x02decemb" +
	"er\x02week"

var fi_FIIndex = []uint32{ // 39 elements
	// Entry 0 - 1F
	0x00000000, 0x00000021, 0x00000043, 0x00000063,
	0x00000080, 0x000000a8, 0x000000bb, 0x000000f3,
	0x0000010a, 0x00000131, 0x0000013d, 0x0000014a,
	0x00000158, 0x00000164, 0x0000017a, 0x00000195,
	0x000001ae, 0x000001d1, 0x000001ea, 0x000001ed,
	0x000001f0, 0x000001f3, 0x000001f6, 0x000001f9,
	0x000001fc, 0x000001ff, 0x00000208, 0x00000211,
	0x0000021b, 0x00000224, 0x0000022d, 0x00000236,
	// Entry 20 - 3F
	0x00000240, 0x00000247, 0x0000024f, 0x00000257,
	0x00000261, 0x0000026a, 0x0000026e,
} // Size: 180 bytes

const fi_FIData string = "" + // Size: 622 bytes
	"\x02Kuinka monta seuraavaa kuukautta\x02Kuinka monta edellistä kuukautta" +
	"\x02Kuinka monta kuukautta per rivi\x02Vuosi (vakiona kuluva vuosi)\x02K" +
	"uukausi 1-12 (vakiona kuluva kuukausi)\x02Tulosta koko vuosi\x02Vain yks" +
	"i kuukausi, vastaa parametrejä -next 0 -prev 0\x02Älä tulosta värejä\x02" +
	"Ensimmäinen viikonpäivä 0-6 (su-la)\x02Parametrit:\x02vakio: %[1]q\x02Es" +
	"imerkkejä:\x02Koko vuosi:\x02Vain tämä kuukausi:\x02Yksi kalenteri kerra" +
	"llaan:\x02väärä kuukausi: %[1]d\x02tunnistamaton viikonpäivä: %[1]d\x02m" +
	"äärän tulee olla > 0\x02ma\x02ti\x02ke\x02to\x02pe\x02la\x02su\x02tammi" +
	"kuu\x02helmikuu\x02maaliskuu\x02huhtikuu\x02toukokuu\x02kesäkuu\x02heinä" +
	"kuu\x02elokuu\x02syyskuu\x02lokakuu\x02marraskuu\x02joulukuu\x02vko"

	// Total table size 1523 bytes (1KiB); checksum: 98EB07F9
