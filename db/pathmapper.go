package db

const (
	longGLDShortIAU = "./db/pairtrading/gold/price_ratio/long_gld_short_iau.json"
	shortGLDLongIAU = "./db/pairtrading/gold/price_ratio/short_gld_long_iau.json"
	goldRepeatNums  = "./db/pairtrading/gold/repeat_num/price_ratio_num_repeat.json"
	goldLogPath     = "./db/pairtrading/gold/log/"
	monitorLogPath  = "./"
)

func MapRecordPath(strat string) (string, string, string) {
	if strat == "gold" {
		return shortGLDLongIAU, longGLDShortIAU, goldRepeatNums
	}
	return "", "", ""
}

func MapLogPath(strat string) string {
	if strat == "gold" {
		return goldLogPath
	} else if strat == "monitor" {
		return monitorLogPath
	}
	return ""
}
