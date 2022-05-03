package db

const (
	repeatNumPath  = "./db/pairtrading/repeat_num/"
	priceRatioPath = "./db/pairtrading/price_ratio/"
	volatilityPath = "./db/pairtrading/volatility/"

	// Gold
	longGLDShortIAU           = priceRatioPath + "long_gld_short_iau.json"
	shortGLDLongIAU           = priceRatioPath + "short_gld_long_iau.json"
	longGLDShortIAURepeatNums = repeatNumPath + "long_gld_short_iau_num_repeat.json"
	shortGLDLongIAURepeatNums = repeatNumPath + "short_gld_long_iau_num_repeat.json"
	goldVolatility            = volatilityPath + "gold.json"

	// Bond
	longAGGShortBND           = priceRatioPath + "long_agg_short_bnd.json"
	shortAGGLongBND           = priceRatioPath + "short_agg_long_bnd.json"
	longAGGShortBNDRepeatNums = repeatNumPath + "long_agg_short_bnd_num_repeat.json"
	shortAGGLongBNDRepeatNums = repeatNumPath + "short_agg_long_bnd_num_repeat.json"

	// S&P 500 Value
	longMDYShortIJH           = priceRatioPath + "long_mdy_short_ijh.json"
	shortMDYLongIJH           = priceRatioPath + "short_mdy_long_ijh.json"
	longMDYShortIJHRepeatNums = repeatNumPath + "long_mdy_short_ijh_num_repeat.json"
	shortMDYLongIJHRepeatNums = repeatNumPath + "short_mdy_long_ijh_num_repeat.json"

	// Utilities Sector
	longVPUShortXLU           = priceRatioPath + "long_vpu_short_xlu.json"
	shortVPULongXLU           = priceRatioPath + "short_vpu_long_xlu.json"
	longVPUShortXLURepeatNums = repeatNumPath + "long_vpu_short_xlu_num_repeat.json"
	shortVPULongXLURepeatNums = repeatNumPath + "short_vpu_long_xlu_num_repeat.json"

	// Russell 2000
	longIWMShortVTWO           = priceRatioPath + "long_iwm_short_vtwo.json"
	shortIWMLongVTWO           = priceRatioPath + "short_iwm_long_vtwo.json"
	longIWMShortVTWORepeatNums = repeatNumPath + "long_iwm_short_vtwo_num_repeat.json"
	shortIWMLongVTWORepeatNums = repeatNumPath + "short_iwm_long_vtwo_num_repeat.json"
	logPath                    = "./db/pairtrading/log/"

	monitorLogPath = "./"
)

type AssetParamConfig struct {
	AssetType                             string
	ShortExensiveLongCheapPriceRatioPath  string
	LongExpensiveShortCheapPriceRatioPath string
	ShortExpensiveLongCheapRepeatNumPath  string
	LongExpensiveShortCheapRepeatNumPath  string
	VolatilityPath                        string
}

func getAssetParamConfig(strat string) *AssetParamConfig {
	assetParamConfig := &AssetParamConfig{}

	if strat == "gold" {
		assetParamConfig.AssetType = "gold"
		assetParamConfig.ShortExensiveLongCheapPriceRatioPath = shortGLDLongIAU
		assetParamConfig.LongExpensiveShortCheapPriceRatioPath = longGLDShortIAU
		assetParamConfig.ShortExpensiveLongCheapRepeatNumPath = shortGLDLongIAURepeatNums
		assetParamConfig.LongExpensiveShortCheapRepeatNumPath = longGLDShortIAURepeatNums
		assetParamConfig.VolatilityPath = goldVolatility
	} else if strat == "bond" {
		assetParamConfig.AssetType = "bond"
		assetParamConfig.ShortExensiveLongCheapPriceRatioPath = shortAGGLongBND
		assetParamConfig.LongExpensiveShortCheapPriceRatioPath = longAGGShortBND
		assetParamConfig.ShortExpensiveLongCheapRepeatNumPath = shortAGGLongBNDRepeatNums
		assetParamConfig.LongExpensiveShortCheapRepeatNumPath = longAGGShortBNDRepeatNums
	} else if strat == "spvalue" {
		assetParamConfig.AssetType = "spvalue"
		assetParamConfig.ShortExensiveLongCheapPriceRatioPath = shortMDYLongIJH
		assetParamConfig.LongExpensiveShortCheapPriceRatioPath = longMDYShortIJH
		assetParamConfig.ShortExpensiveLongCheapRepeatNumPath = shortMDYLongIJHRepeatNums
		assetParamConfig.LongExpensiveShortCheapRepeatNumPath = longMDYShortIJHRepeatNums
	} else if strat == "utilities" {
		assetParamConfig.AssetType = "utilities"
		assetParamConfig.ShortExensiveLongCheapPriceRatioPath = shortVPULongXLU
		assetParamConfig.LongExpensiveShortCheapPriceRatioPath = longVPUShortXLU
		assetParamConfig.ShortExpensiveLongCheapRepeatNumPath = shortVPULongXLURepeatNums
		assetParamConfig.LongExpensiveShortCheapRepeatNumPath = longVPUShortXLURepeatNums
	} else if strat == "russell2000" {
		assetParamConfig.AssetType = "russell2000"
		assetParamConfig.ShortExensiveLongCheapPriceRatioPath = shortIWMLongVTWO
		assetParamConfig.LongExpensiveShortCheapPriceRatioPath = longIWMShortVTWO
		assetParamConfig.ShortExpensiveLongCheapRepeatNumPath = shortIWMLongVTWORepeatNums
		assetParamConfig.LongExpensiveShortCheapRepeatNumPath = longIWMShortVTWORepeatNums
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
