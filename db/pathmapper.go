package db

const (
	repeatNumPath  = "./db/pairtrading/repeat_num/"
	priceRatioPath = "./db/pairtrading/price_ratio/"

	longGLDShortIAU           = priceRatioPath + "long_gld_short_iau.json"
	shortGLDLongIAU           = priceRatioPath + "short_gld_long_iau.json"
	longGLDShortIAURepeatNums = repeatNumPath + "long_gld_short_iau_num_repeat.json"
	shortGLDLongIAURepeatNums = repeatNumPath + "short_gld_long_iau_num_repeat.json"

	longAGGShortBND           = priceRatioPath + "long_agg_short_bnd.json"
	shortAGGLongBND           = priceRatioPath + "short_agg_long_bnd.json"
	longAGGShortBNDRepeatNums = repeatNumPath + "long_agg_short_bnd_num_repeat.json"
	shortAGGLongBNDRepeatNums = repeatNumPath + "short_agg_long_bnd_num_repeat.json"
	logPath                   = "./db/pairtrading/log/"

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
	} else if strat == "bond" {
		assetParamConfig.ShortExensiveLongCheapPriceRatioPath = shortAGGLongBND
		assetParamConfig.LongExpensiveShortCheapPriceRatioPath = longAGGShortBND
		assetParamConfig.ShortExpensiveLongCheapRepeatNumPath = shortAGGLongBNDRepeatNums
		assetParamConfig.LongExpensiveShortCheapRepeatNumPath = longAGGShortBNDRepeatNums
	}
	return assetParamConfig
}

func MapRecordPath(strat string) *AssetParamConfig {
	return getAssetParamConfig(strat)
}

func MapLogPath(strat string) string {
	if strat == "monitor" {
		return monitorLogPath
	}
	return logPath
}
