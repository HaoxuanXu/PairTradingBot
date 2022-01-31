package db

const (
	longGLDShortIAU           = "./db/pairtrading/gold/price_ratio/long_gld_short_iau.json"
	shortGLDLongIAU           = "./db/pairtrading/gold/price_ratio/short_gld_long_iau.json"
	longGLDShortIAURepeatNums = "./db/pairtrading/gold/repeat_num/long_gld_short_iau_num_repeat.json"
	shortGLDLongIAURepeatNums = "./db/pairtrading/gold/repeat_num/short_gld_long_iau_num_repeat.json"
	goldLogPath               = "./db/pairtrading/gold/log/"
	monitorLogPath            = "./"
)

func MapRecordPath(strat string) (string, string, string, string) {
	if strat == "gold" {
		return shortGLDLongIAU, longGLDShortIAU, longGLDShortIAURepeatNums, shortGLDLongIAURepeatNums
	}
	return "", "", "", ""
}

func MapLogPath(strat string) string {
	if strat == "gold" {
		return goldLogPath
	} else if strat == "monitor" {
		return monitorLogPath
	}
	return ""
}
