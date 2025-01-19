package aksharamukha

// TODO add STDERR warning for these partially supported languages
var indicSubset = []string{"LaoTham", "LueTham", "KhuenTham", "PhagsPa", "TaiLaing", "Mon", "Ahom", "KhamtiShan", "Khmer", "Burmese", "Lao", "Thai", "Balinese", "Javanese", "Tibetan", "LaoPali", "TaiTham", "Cham", "Lepcha", "Ahom", "ZanabazarSquare"}

var Lang2Scripts = map[string][]string{
	"hin": {"Devanagari"},                     // Hindi (~320 million speakers)
	"ara": {"Arab"},                           // Arabic (~300 million speakers)
	"ben": {"Bengali"},                        // Bengali (~230 million speakers)
	"rus": {"RussianCyrillic"},                // Russian (~150 million speakers)
	"jpn": {"Hiragana", "Katakana"},           // Japanese (~125 million speakers)
	"pan": {"Gurmukhi", "Shahmukhi"},          // Punjabi (~90 million speakers)
	"mar": {"Devanagari"},                     // Marathi (~83 million speakers)
	"tel": {"Telugu"},                         // Telugu (~82 million speakers)
	"tam": {"Tamil", "TamilExtended"},         // Tamil (~70 million speakers)
	"fas": {"Arab-Fa", "Arab"},                // Persian (~70 million speakers)
	"urd": {"Urdu", "Arab", "Shahmukhi"},      // Urdu (~68 million speakers)
	"guj": {"Gujarati"},                       // Gujarati (~56 million speakers)
	"pus": {"Arab"},                           // Pashto (~40 million speakers)
	"mal": {"Malayalam"},                      // Malayalam (~37 million speakers)
	"mai": {"Devanagari", "Tirhuta", "Kaithi"},// Maithili (~34 million speakers)
	"mya": {"Burmese"},                        // Burmese (~33 million speakers)
	"ukr": {"RussianCyrillic"},                // Ukrainian (~30 million speakers)
	"uzb": {"RussianCyrillic"},                // Uzbek (~27 million speakers)
	"orm": {"Ethi"},                           // Oromo (~25 million speakers)
	"asm": {"Assamese"},                       // Assamese (~23 million speakers)
	"kur": {"Arab"},                           // Kurdish (~20 million speakers)
	"tha": {"Thai"},                           // Thai (~20 million speakers)
	"nep": {"Devanagari"},                     // Nepali (~16 million speakers)
	"khm": {"Khmer"},                          // Khmer (~16 million speakers)
	"sin": {"Sinhala"},                        // Sinhala (~16 million speakers)
	"kaz": {"RussianCyrillic"},                // Kazakh (~13 million speakers)
	"bul": {"RussianCyrillic"},                // Bulgarian (~7 million speakers)
	"bel": {"RussianCyrillic"},                // Belarusian (~7 million speakers)
	"sat": {"Santali"},                        // Santali (~7 million speakers)
	"srp": {"RussianCyrillic"},                // Serbian (~6 million speakers)
	"tir": {"Ethi"},                           // Tigrinya (~6 million speakers)
	"kas": {"Devanagari", "Sharada"},          // Kashmiri (~5 million speakers)
	"kir": {"RussianCyrillic"},                // Kyrgyz (~5 million speakers)
	"tgk": {"RussianCyrillic"},                // Tajik (~4 million speakers)

	// Languages with fewer speakers
	"mni": {"MeeteiMayek"},                   // Manipuri (~1.5 million speakers)
	"bod": {"Tibetan"},                       // Tibetan (~1.2 million speakers)
	"new": {"Newa"},                          // Nepal Bhasa (~800,000 speakers)
	"khb": {"KhuenTham", "LueTham"},          // Tai Khuen, Tai Lue
	"nod": {"TaiTham"},                       // Northern Thai
	"lep": {"Lepcha"},                        // Lepcha
	"lif": {"Limbu"},                         // Limbu
	"ccp": {"Chakma"},                        // Chakma
	"gon": {"GunjalaGondi", "MasaramGondi"},  // Gondi
	"hoc": {"WarangCiti"},                    // Ho
	"rhg": {"HanifiRohingya"},                // Rohingya
	"kht": {"KhamtiShan"},                    // Khamti Shan
	"kaw": {"Kawi"},                          // Kawi
	"jav": {"Javanese"},                      // Javanese
	"bug": {"Buginese"},                      // Buginese
	"syl": {"SylotiNagri"},                   // Sylheti
	"bho": {"Devanagari"},                    // additional: "Kaithi" // Bhojpuri
	"awa": {"Devanagari"},                    // Awadhi
	"kok": {"Devanagari"},                    // Konkani
	"dgo": {"Devanagari"},                    // Dogri
	"bra": {"Devanagari"},                    // Braj Bhasha
	"tjl": {"TaiLaing"},                      // Tai Laing
	"srb": {"SoraSompeng"},                   // Sora
	"rej": {"Rejang"},                        // Rejang
	"ban": {"Balinese"},                      // Balinese
	"saz": {"Saurashtra"},                    // Saurashtra
	"mak": {"Makasar"},                       // Makasar
	"div": {"Thaana"},                        // Dhivehi

	// Ancient or historic languages
	"san": {"Devanagari"},                   // additional: "Grantha", "GranthaPandya", "Kharoshthi", "Nandinagari", "Ranjana", "Sharada", "Siddham", "Brahmi" // Sanskrit
	"ave": {"Avestan"},                      // Avestan
	"pal": {"Phli"},                         // additional: "Phlp" // Pahlavi
	"xpr": {"Prti"},                         // Parthian
	"xna": {"Narb"},                         // Old North Arabian
	"xsa": {"Sarb"},                         // Old South Arabian
	"peo": {"OldPersian"},                   // Old Persian
	"sog": {"Sogd"},                         // additional: "Sogo" // Sogdian
	"arc": {"Armi"},                         // Imperial Aramaic
	"phn": {"Phnx"},                         // Phoenician
	"smp": {"Samr"},                         // Samaritan
	"uga": {"Ugar"},                         // Ugaritic
	"syr": {"Syre"},                         // additional: "Syrn", "Syrj" // Syriac

	// Languages with limited script usage data
	"aha": {"Ahom"},                         // Ahom
	"btx": {"BatakKaro"},                    // Batak Karo
	"btm": {"BatakManda"},                   // Batak Mandailing
	"btd": {"BatakPakpak"},                  // Batak Pakpak
	"bts": {"BatakSima"},                    // Batak Simalungun
	"bbc": {"BatakToba"},                    // Batak Toba
	"bku": {"Buhid"},                        // Buhid
	"hnn": {"Hanunoo"},                      // Hanunoo
	"mro": {"Mro"},                          // Mro
	"nnp": {"Wancho"},                       // Wancho
}



var Script2RomanScheme = map[string]string{
	// Modern Indic scripts - using ISO 15919 for comprehensive sound representation
	"Devanagari":     "ISO",             // Modern languages like Hindi, Marathi
	"Bengali":        "ISO",
	"Gujarati":       "ISO",
	"Gurmukhi":       "ISO",
	"Kannada":        "ISO",
	"Malayalam":      "ISO",
	"Oriya":          "ISO",
	"Tamil":          "ISO",
	"Telugu":         "ISO",
	"Sinhala":        "ISO",
	"MeeteiMayek":    "ISO",
	"Tirhuta":        "ISO",
	"SylotiNagri":    "ISO",

	// Southeast Asian scripts
	"Thai":           "ISO",
	"Lao":            "ISO",
	"LaoPali":        "ISO",
	"Burmese":        "ISO",
	"Khmer":          "ISO",
	"Javanese":       "ISO",
	"Balinese":       "ISO",
	"Cham":           "ISO",
	"TaiTham":        "ISO",
	"LaoTham":        "ISO",
	"KhuenTham":      "ISO",
	"LueTham":        "ISO",
	"Chakma":         "ISO",
	"Lepcha":         "ISO",
	"Limbu":          "ISO",
	"Ahom":           "ISO",

	// Semitic scripts
	"Arab":           "ISO233",          // alt: PersianDMG
	"Arab-Fa":        "PersianDMG",
	"Hebrew":         "ISO259",          // alt: HebrewSBL
	"Hebr-Ar":        "ISO259",
	"Syrn":           "Latn",            // Eastern Syriac
	"Syrj":           "Latn",            // Western Syriac
	"Syre":           "Latn",            // Estrangela Syriac
	"Armi":           "Latn",            // Imperial Aramaic
	"Phnx":           "Latn",            // Phoenician
	"OldPersian":     "PersianDMG",

	// East Asian scripts
	"Hiragana":       "ISO",
	"Katakana":       "ISO",

	// Other scripts
	"RussianCyrillic": "ISO",
	"IPA":             "ISO",
	"Thaana":          "ISO",
	"Tibetan":         "ISO",

	// Classical/Ancient Indic scripts
	"Grantha":        "IAST",
	"GranthaPandya":  "IAST",
	"Brahmi":         "IAST",
	"Siddham":        "IAST",
	"Sharada":        "IAST",
	"Modi":           "IAST",
	"Nandinagari":    "IAST",
	"Kharoshthi":     "IAST",
	"Bhaiksuki":      "IAST",
	"TamilBrahmi":    "IAST",

	// Ancient scripts
	"Narb":           "Latn",            // Old North Arabian
	"Sarb":           "Latn",            // Old South Arabian
	"Phli":           "PersianDMG",      // Inscriptional Pahlavi
	"Phlp":           "PersianDMG",      // Psalter Pahlavi
	"Prti":           "PersianDMG",      // Inscriptional Parthian
	"Sogo":           "Latn",            // Old Sogdian
	"Sogd":           "Latn",            // Sogdian
	"Ugar":           "Latn",            // Ugaritic
	"Samr":           "Latn",            // Samaritan

	// Lesser used scripts
	"Vatteluttu":     "ISO",
	"Mahajani":       "ISO",
	"Multani":        "ISO",
	"Khudawadi":      "ISO",
	"Khojki":         "ISO",
	"Shan":           "ISO",
	"TaiLaing":       "ISO",
	"KhamtiShan":     "ISO",
	"Mongolian":      "ISO",
	"PhagsPa":        "ISO",
	"Marchen":        "ISO",
	"ZanabazarSquare":"ISO",
	"Soyombo":        "ISO",
	"Dogra":          "ISO",
	"GunjalaGondi":   "ISO",
	"MasaramGondi":   "ISO",
	"SoraSompeng":    "ISO",
	"WarangCiti":     "ISO",
}


