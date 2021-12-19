package db

const (
	longGLDShortIAU = "./db/pairtrading/gold/price_ratio/long_gld_short_iau.json"
	shortGLDLongIAU = "./db/pairtrading/gold/price_ratio/short_gld_long_iau.json"
	goldRepeatNums  = "./db/pairtrading/gold/repeat_num/price_ratio_num_repeat.json"
)

func MapRecordPath(strat string) (string, string, string) {
	if strat == "gold" {
		return shortGLDLongIAU, longGLDShortIAU, goldRepeatNums
	}
	return "", "", ""
}