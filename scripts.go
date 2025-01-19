package aksharamukha

import "slices"
// IsValidScript checks if a script is valid
func IsValidScript(s Script) bool {
	return slices.Contains(validScripts, s)
}


type Script string

const (
	Ahom            Script = "Ahom"
	Arab            Script = "Arab"
	Ariyaka         Script = "Ariyaka"
	Assamese        Script = "Assamese"
	Avestan         Script = "Avestan"
	Balinese        Script = "Balinese"
	BatakKaro       Script = "BatakKaro"
	BatakManda      Script = "BatakManda"
	BatakPakpak     Script = "BatakPakpak"
	BatakSima       Script = "BatakSima"
	BatakToba       Script = "BatakToba"
	Bengali         Script = "Bengali"
	Bhaiksuki       Script = "Bhaiksuki"
	Brahmi          Script = "Brahmi"
	Buginese        Script = "Buginese"
	Buhid           Script = "Buhid"
	Burmese         Script = "Burmese"
	Chakma          Script = "Chakma"
	Cham            Script = "Cham"
	RussianCyrillic Script = "RussianCyrillic"
	Devanagari      Script = "Devanagari"
	Dogra           Script = "Dogra"
	Elym            Script = "Elym"
	Ethi            Script = "Ethi"
	GunjalaGondi    Script = "GunjalaGondi"
	MasaramGondi    Script = "MasaramGondi"
	Grantha         Script = "Grantha"
	GranthaPandya   Script = "GranthaPandya"
	Gujarati        Script = "Gujarati"
	Hanunoo         Script = "Hanunoo"
	Hatr            Script = "Hatr"
	Hebrew          Script = "Hebrew"
	HebrAr          Script = "Hebr-Ar"
	Armi            Script = "Armi"
	Phli            Script = "Phli"
	Prti            Script = "Prti"
	Hiragana        Script = "Hiragana"
	Katakana        Script = "Katakana"
	Javanese        Script = "Javanese"
	Kaithi          Script = "Kaithi"
	Kannada         Script = "Kannada"
	Kawi            Script = "Kawi"
	KhamtiShan      Script = "KhamtiShan"
	Kharoshthi      Script = "Kharoshthi"
	Khmer           Script = "Khmer"
	Khojki          Script = "Khojki"
	KhomThai        Script = "KhomThai"
	Khudawadi       Script = "Khudawadi"
	Lao             Script = "Lao"
	LaoPali         Script = "LaoPali"
	Lepcha          Script = "Lepcha"
	Limbu           Script = "Limbu"
	Mahajani        Script = "Mahajani"
	Makasar         Script = "Makasar"
	Malayalam       Script = "Malayalam"
	Mani            Script = "Mani"
	Marchen         Script = "Marchen"
	MeeteiMayek     Script = "MeeteiMayek"
	Modi            Script = "Modi"
	Mon             Script = "Mon"
	Mongolian       Script = "Mongolian"
	Mro             Script = "Mro"
	Multani         Script = "Multani"
	Nbat            Script = "Nbat"
	Nandinagari     Script = "Nandinagari"
	Newa            Script = "Newa"
	Narb            Script = "Narb"
	OldPersian      Script = "OldPersian"
	Sogo            Script = "Sogo"
	Sarb            Script = "Sarb"
	Oriya           Script = "Oriya"
	Pallava         Script = "Pallava"
	Palm            Script = "Palm"
	ArabFa          Script = "Arab-Fa"
	PhagsPa         Script = "PhagsPa"
	Phnx            Script = "Phnx"
	Phlp            Script = "Phlp"
	Gurmukhi        Script = "Gurmukhi"
	Ranjana         Script = "Ranjana"
	Rejang          Script = "Rejang"
	HanifiRohingya  Script = "HanifiRohingya"
	BarahaNorth     Script = "BarahaNorth"
	BarahaSouth     Script = "BarahaSouth"
	RomanColloquial Script = "RomanColloquial"
	PersianDMG      Script = "PersianDMG"
	HK              Script = "HK"
	IAST            Script = "IAST"
	IASTPali        Script = "IASTPali"
	IPA             Script = "IPA"
	ISO             Script = "ISO"
	ISOPali         Script = "ISOPali"
	ISO233          Script = "ISO233"
	ISO259          Script = "ISO259"
	Itrans          Script = "Itrans"
	IASTLOC         Script = "IASTLOC"
	RomanReadable   Script = "RomanReadable"
	HebrewSBL       Script = "HebrewSBL"
	SLP1            Script = "SLP1"
	Type            Script = "Type"
	Latn            Script = "Latn"
	Titus           Script = "Titus"
	Velthuis        Script = "Velthuis"
	WX              Script = "WX"
	Samr            Script = "Samr"
	Santali         Script = "Santali"
	Saurashtra      Script = "Saurashtra"
	Shahmukhi       Script = "Shahmukhi"
	Shan            Script = "Shan"
	Sharada         Script = "Sharada"
	Siddham         Script = "Siddham"
	Sinhala         Script = "Sinhala"
	Sogd            Script = "Sogd"
	SoraSompeng     Script = "SoraSompeng"
	Soyombo         Script = "Soyombo"
	Sundanese       Script = "Sundanese"
	SylotiNagri     Script = "SylotiNagri"
	Syrn            Script = "Syrn"
	Syre            Script = "Syre"
	Syrj            Script = "Syrj"
	Tagalog         Script = "Tagalog"
	Tagbanwa        Script = "Tagbanwa"
	TaiLaing        Script = "TaiLaing"
	Takri           Script = "Takri"
	Tamil           Script = "Tamil"
	TamilExtended   Script = "TamilExtended"
	TamilBrahmi     Script = "TamilBrahmi"
	Telugu          Script = "Telugu"
	Thaana          Script = "Thaana"
	Thai            Script = "Thai"
	TaiTham         Script = "TaiTham"
	LaoTham         Script = "LaoTham"
	KhuenTham       Script = "KhuenTham"
	LueTham         Script = "LueTham"
	Tibetan         Script = "Tibetan"
	Tirhuta         Script = "Tirhuta"
	Ugar            Script = "Ugar"
	Urdu            Script = "Urdu"
	Vatteluttu      Script = "Vatteluttu"
	Wancho          Script = "Wancho"
	WarangCiti      Script = "WarangCiti"
	ZanabazarSquare Script = "ZanabazarSquare"
)

var validScripts = []Script{"Ahom", "Arab", "Ariyaka", "Assamese", "Avestan", "Balinese", "BatakKaro", "BatakManda", "BatakPakpak", "BatakSima", "BatakToba", "Bengali", "Bhaiksuki", "Brahmi", "Buginese", "Buhid", "Burmese", "Chakma", "Cham", "RussianCyrillic", "Devanagari", "Dogra", "Elym", "Ethi", "GunjalaGondi", "MasaramGondi", "Grantha", "GranthaPandya", "Gujarati", "Hanunoo", "Hatr", "Hebrew", "Hebr-Ar", "Armi", "Phli", "Prti", "Hiragana", "Katakana", "Javanese", "Kaithi", "Kannada", "Kawi", "KhamtiShan", "Kharoshthi", "Khmer", "Khojki", "KhomThai", "Khudawadi", "Lao", "LaoPali", "Lepcha", "Limbu", "Mahajani", "Makasar", "Malayalam", "Mani", "Marchen", "MeeteiMayek", "Modi", "Mon", "Mongolian", "Mro", "Multani", "Nbat", "Nandinagari", "Newa", "Narb", "OldPersian", "Sogo", "Sarb", "Oriya", "Pallava", "Palm", "Arab-Fa", "PhagsPa", "Phnx", "Phlp", "Gurmukhi", "Ranjana", "Rejang", "HanifiRohingya", "BarahaNorth", "BarahaSouth", "RomanColloquial", "PersianDMG", "HK", "IAST", "IASTPali", "IPA", "ISO", "ISOPali", "ISO233", "ISO259", "Itrans", "IASTLOC", "RomanReadable", "HebrewSBL", "SLP1", "Type", "Latn", "Titus", "Velthuis", "WX", "Samr", "Santali", "Saurashtra", "Shahmukhi", "Shan", "Sharada", "Siddham", "Sinhala", "Sogd", "SoraSompeng", "Soyombo", "Sundanese", "SylotiNagri", "Syrn", "Syre", "Syrj", "Tagalog", "Tagbanwa", "TaiLaing", "Takri", "Tamil", "TamilExtended", "TamilBrahmi", "Telugu", "Thaana", "Thai", "TaiTham", "LaoTham", "KhuenTham", "LueTham", "Tibetan", "Tirhuta", "Ugar", "Urdu", "Vatteluttu", "Wancho", "WarangCiti", "ZanabazarSquare"}

