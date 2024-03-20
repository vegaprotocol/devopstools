package generate

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"time"
)

var (
	adjectives = []string{
		"bright", "colossal", "condescending", "confused", "cooperative", "costly", "courageous",
		"oppulent", "magnificent",
		"ashy", "blue", "gray", "green", "icy", "lemon", "mango", "orange", "purple", "red", "salmon",
		"important", "inexpensive", "fluky", "odd", "promising", "powerful", "rich", "shy", "tender",
		"opportune", "timely", "vast", "wrong",
		"crashing", "echoing", "faint", "harsh", "howling", "loud", "melodic", "noisy",
		"purring", "quiet", "rhythmic", "thundering", "wailing", "whining", "whispering",
		"adorable", "adventurous", "agitated", "alert", "bored", "brave",
		"determined", "distressed", "disturbed", "dizzy",
		"excited", "extensive", "exuberant", "frustrating", "funny", "fuzzy",
		"hungry", "icy", "ideal", "immense", "impressionable", "intrigued", "irate", "foolish", "frantic", "fresh",
		"clean",
		"graceful", "gritty", "happy", "hollow",
		"quaint", "scruffy", "shapely", "short", "skinny", "stocky", "ugly", "unkempt", "unsightly",
		"friendly", "frightened", "glorious", "happy", "harebrained", "healthy",
		"helpful", "helpless", "high", "hollow", "loose", "lovely", "lucky",
		"alive", "better", "careful", "clever", "dead", "easy", "famous", "gifted", "hallowed", "helpful",
		"mysterious", "narrow", "perfect", "perplexed", "quizzical", "tender", "tense",
		"terrible", "tricky", "wicked", "yummy", "zippy",
	}
	nouns = []string{
		"table", "chair", "vacuum", "man", "river", "cloud", "beach", "party", "ball", "bat",
		"book", "pen", "pencil", "cookie", "biscuit", "butterfly", "bee", "hill", "baby", "pool",
		"apple", "mango", "day", "week", "month", "flowers", "honey", "train", "bus", "car", "hat", "shirt", "sister",
		"brother", "milk", "coffee", "bed", "plant", "tree", "horse", "wall", "cat", "monkey",
		"dog", "roof", "chimney", "tile", "shoe", "park", "bird", "pond", "duck", "farmer", "sheep", "computer", "fan", "television",
		"stove", "spoon", "door", "window", "doctor", "teacher", "fisherman", "barber", "clouds", "wind", "road", "lake", "policeman",
		"swing", "phone", "school", "taxi", "burger", "pizza", "pickle", "movie", "music", "town", "building", "candle",
		"dress", "mouse", "river", "toy", "key", "country",
	}
	// see: https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
	countries = []string{
		"AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ",
		"BA", "BB", "BC", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BK", "BL", "BM", "BN", "BO", "BP", "BQ", "BR", "BS", "BT", "BU", "BV", "BW", "BX", "BY", "BZ",
		"CA", "CB", "CC", "CD", "CE", "CF", "CG", "CH", "CI", "CJ", "CK", "CL", "CM", "CN", "CO", "CP", "CQ", "CR", "CS", "CT", "CU", "CV", "CW", "CX", "CY", "CZ",
		"DA", "DB", "DC", "DD", "DE", "DF", "DG", "DH", "DI", "DJ", "DK", "DL", "DM", "DN", "DO", "DP", "DQ", "DR", "DS", "DT", "DU", "DV", "DW", "DX", "DY", "DZ",
		"EA", "EB", "EC", "ED", "EE", "EF", "EG", "EH", "EI", "EJ", "EK", "EL", "EM", "EN", "EO", "EP", "EQ", "ER", "ES", "ET", "EU", "EV", "EW", "EX", "EY", "EZ",
		"FA", "FB", "FC", "FD", "FE", "FF", "FG", "FH", "FI", "FJ", "FK", "FL", "FM", "FN", "FO", "FP", "FQ", "FR", "FS", "FT", "FU", "FV", "FW", "FX", "FY", "FZ",
		"GA", "GB", "GC", "GD", "GE", "GF", "GG", "GH", "GI", "GJ", "GK", "GL", "GM", "GN", "GO", "GP", "GQ", "GR", "GS", "GT", "GU", "GV", "GW", "GX", "GY", "GZ",
		"HA", "HB", "HC", "HD", "HE", "HF", "HG", "HH", "HI", "HJ", "HK", "HL", "HM", "HN", "HO", "HP", "HQ", "HR", "HS", "HT", "HU", "HV", "HW", "HX", "HY", "HZ",
		"IA", "IB", "IC", "ID", "IE", "IF", "IG", "IH", "II", "IJ", "IK", "IL", "IM", "IN", "IO", "IP", "IQ", "IR", "IS", "IT", "IU", "IV", "IW", "IX", "IY", "IZ",
		"JA", "JB", "JC", "JD", "JE", "JF", "JG", "JH", "JI", "JJ", "JK", "JL", "JM", "JN", "JO", "JP", "JQ", "JR", "JS", "JT", "JU", "JV", "JW", "JX", "JY", "JZ",
		"KA", "KB", "KC", "KD", "KE", "KF", "KG", "KH", "KI", "KJ", "KK", "KL", "KM", "KN", "KO", "KP", "KQ", "KR", "KS", "KT", "KU", "KV", "KW", "KX", "KY", "KZ",
		"LA", "LB", "LC", "LD", "LE", "LF", "LG", "LH", "LI", "LJ", "LK", "LL", "LM", "LN", "LO", "LP", "LQ", "LR", "LS", "LT", "LU", "LV", "LW", "LX", "LY", "LZ",
		"MA", "MB", "MC", "MD", "ME", "MF", "MG", "MH", "MI", "MJ", "MK", "ML", "MM", "MN", "MO", "MP", "MQ", "MR", "MS", "MT", "MU", "MV", "MW", "MX", "MY", "MZ",
		"NA", "NB", "NC", "ND", "NE", "NF", "NG", "NH", "NI", "NJ", "NK", "NL", "NM", "NN", "NO", "NP", "NQ", "NR", "NS", "NT", "NU", "NV", "NW", "NX", "NY", "NZ",
		"OA", "OB", "OC", "OD", "OE", "OF", "OG", "OH", "OI", "OJ", "OK", "OL", "OM", "ON", "OO", "OP", "OQ", "OR", "OS", "OT", "OU", "OV", "OW", "OX", "OY", "OZ",
		"PA", "PB", "PC", "PD", "PE", "PF", "PG", "PH", "PI", "PJ", "PK", "PL", "PM", "PN", "PO", "PP", "PQ", "PR", "PS", "PT", "PU", "PV", "PW", "PX", "PY", "PZ",
		"QA", "QB", "QC", "QD", "QE", "QF", "QG", "QH", "QI", "QJ", "QK", "QL", "QM", "QN", "QO", "QP", "QQ", "QR", "QS", "QT", "QU", "QV", "QW", "QX", "QY", "QZ",
		"RA", "RB", "RC", "RD", "RE", "RF", "RG", "RH", "RI", "RJ", "RK", "RL", "RM", "RN", "RO", "RP", "RQ", "RR", "RS", "RT", "RU", "RV", "RW", "RX", "RY", "RZ",
		"SA", "SB", "SC", "SD", "SE", "SF", "SG", "SH", "SI", "SJ", "SK", "SL", "SM", "SN", "SO", "SP", "SQ", "SR", "SS", "ST", "SU", "SV", "SW", "SX", "SY", "SZ",
		"TA", "TB", "TC", "TD", "TE", "TF", "TG", "TH", "TI", "TJ", "TK", "TL", "TM", "TN", "TO", "TP", "TQ", "TR", "TS", "TT", "TU", "TV", "TW", "TX", "TY", "TZ",
		"UA", "UB", "UC", "UD", "UE", "UF", "UG", "UH", "UI", "UJ", "UK", "UL", "UM", "UN", "UO", "UP", "UQ", "UR", "US", "UT", "UU", "UV", "UW", "UX", "UY", "UZ",
		"VA", "VB", "VC", "VD", "VE", "VF", "VG", "VH", "VI", "VJ", "VK", "VL", "VM", "VN", "VO", "VP", "VQ", "VR", "VS", "VT", "VU", "VV", "VW", "VX", "VY", "VZ",
		"WA", "WB", "WC", "WD", "WE", "WF", "WG", "WH", "WI", "WJ", "WK", "WL", "WM", "WN", "WO", "WP", "WQ", "WR", "WS", "WT", "WU", "WV", "WW", "WX", "WY", "WZ",
		"XA", "XB", "XC", "XD", "XE", "XF", "XG", "XH", "XI", "XJ", "XK", "XL", "XM", "XN", "XO", "XP", "XQ", "XR", "XS", "XT", "XU", "XV", "XW", "XX", "XY", "XZ",
		"YA", "YB", "YC", "YD", "YE", "YF", "YG", "YH", "YI", "YJ", "YK", "YL", "YM", "YN", "YO", "YP", "YQ", "YR", "YS", "YT", "YU", "YV", "YW", "YX", "YY", "YZ",
		"ZA", "ZB", "ZC", "ZD", "ZE", "ZF", "ZG", "ZH", "ZI", "ZJ", "ZK", "ZL", "ZM", "ZN", "ZO", "ZP", "ZQ", "ZR", "ZS", "ZT", "ZU", "ZV", "ZW", "ZX", "ZY", "ZZ",
	}
)

func GenerateName() (string, error) {
	adjLen := big.NewInt(int64(len(adjectives)))
	nounLen := big.NewInt(int64(len(nouns)))

	randomAdjIndex, err := rand.Int(rand.Reader, adjLen)
	if err != nil {
		return "", fmt.Errorf("failed to generate a random name, %w", err)
	}
	randomAdj := adjectives[randomAdjIndex.Int64()]

	randomNounIndex, err := rand.Int(rand.Reader, nounLen)
	if err != nil {
		return "", fmt.Errorf("failed to generate a random name, %w", err)
	}
	randomNoun := nouns[randomNounIndex.Int64()]

	return fmt.Sprintf("%s-%s", randomAdj, randomNoun), nil
}

func GenerateAvatarURL() (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	lettersLen := big.NewInt(int64(len(letters)))
	ret := make([]byte, 32)
	for i := int64(0); i < 32; i++ {
		num, err := rand.Int(rand.Reader, lettersLen)
		if err != nil {
			return "", fmt.Errorf("failed to generate random avatar url, %w", err)
		}
		ret[i] = letters[num.Int64()]
	}
	return fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=identicon", string(ret)), nil
}

func GenerateRandomWikiURL() (string, error) {
	errMsg := "failed to generate random wiki url"
	req, err := http.NewRequest(http.MethodGet, "https://en.wikipedia.org/wiki/Special:Random", nil)
	if err != nil {
		return "", fmt.Errorf("%s, %w", errMsg, err)
	}

	httpClient := http.Client{
		Timeout: time.Second * 5,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	response, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s, %w", errMsg, err)
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusFound { // status code 302
		url, err := response.Location()
		if err != nil {
			return "", fmt.Errorf("%s, %w", errMsg, err)
		}
		return url.String(), nil
	}
	return "", fmt.Errorf("%s, wrong response from wiki backend, expected 302, but got %d, %v", errMsg, response.StatusCode, response)
}

func GenerateCountryCode() (string, error) {
	countryLen := big.NewInt(int64(len(countries)))

	randomIndex, err := rand.Int(rand.Reader, countryLen)
	if err != nil {
		return "", fmt.Errorf("failed to generate a contry code, %w", err)
	}
	randomContry := countries[randomIndex.Int64()]
	return randomContry, nil
}
