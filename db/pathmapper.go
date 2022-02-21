package db

const (
	longGLDShortIAU           = "./db/pairtrading/price_ratio/long_gld_short_iau.json"
	shortGLDLongIAU           = "./db/pairtrading/price_ratio/short_gld_long_iau.json"
	longGLDShortIAURepeatNums = "./db/pairtrading/repeat_num/long_gld_short_iau_num_repeat.json"
	shortGLDLongIAURepeatNums = "./db/pairtrading/repeat_num/short_gld_long_iau_num_repeat.json"
	goldLogPath               = "./db/pairtrading/log/"

	monitorLogPath = "./"
)

type AssetParamConfig struct {
	ShortExensiveLongCheapPriceRatioPath  string
	LongExpensiveShortCheapPriceRatioPath string
	ShortExpensiveLongCheapRepeatNumPath  string
	LongExpensiveShortCheapRepeatNumPath  string
}

func getAssetParamConfig(strat string) *AssetParamConfig {
	assetParamConfig := &AssetParamConfig{}

	if strat == "gold" {
		assetParamConfig.ShortExensiveLongCheapPriceRatioPath = shortGLDLongIAU
		assetParamConfig.LongExpensiveShortCheapPriceRatioPath = longGLDShortIAU
		assetParamConfig.ShortExpensiveLongCheapRepeatNumPath = shortGLDLongIAURepeatNums
		assetParamConfig.LongExpensiveShortCheapRepeatNumPath = longGLDShortIAURepeatNums
	}
	return assetParamConfig
}

func MapRecordPath(strat string) *AssetParamConfig {
	return getAssetParamConfig(strat)
}

func MapLogPath(strat string) string {
	if strat == "gold" {
		return goldLogPath
	} else if strat == "monitor" {
		return monitorLogPath
	}
	return ""
}
