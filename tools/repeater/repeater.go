package repeater

import "github.com/HaoxuanXu/TradingBot/tools/util"

// This function calculates the hypothetical profit given that we choose to enter a position after repeatNum number of repeats
func calculateHypotheticalProfit(repeatMapper map[int]int, repeatNum int) float64 {
	var profit float64
	profit = 0.0
	for key, val := range repeatMapper {
		if key < repeatNum {
			continue
		} else if key == repeatNum {
			profit -= 6.0 * float64(val)
		} else {
			profit += float64(val)
		}
	}
	return profit
}

func CalculateOptimalRepeatNum(repeatArray []int) int {
	repeatMapper := make(map[int]int)
	highestProfit := 0.0
	optimalNum := 0
	for _, val := range repeatArray {
		if _, ok := repeatMapper[val]; !ok {
			repeatMapper[val] = 0
		}
		repeatMapper[val]++
	}

	for repeatNum := 1; repeatNum < util.GetMaxInt(repeatArray); repeatNum++ {
		currentProfit := calculateHypotheticalProfit(repeatMapper, repeatNum)
		if currentProfit > highestProfit {
			highestProfit = currentProfit
			optimalNum = repeatNum
		}
	}
	return optimalNum
}
