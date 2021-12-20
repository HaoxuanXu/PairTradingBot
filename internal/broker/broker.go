package broker

import (
	"log"
	"sync"
	"time"

	"github.com/HaoxuanXu/TradingBot/configs"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/shopspring/decimal"
)

type AlpacaBroker struct {
	client              alpaca.Client
	account             *alpaca.Account
	Clock               alpaca.Clock
	PortfolioValue      float64
	TransactionNums     int
	MaxPortfolioPercent float64
	HasPosition         bool
	MinProfitThreshold  float64
}

// You can treat this as a constructor of the broker class
func GetBroker(accountType string, entryPercent float64) *AlpacaBroker {
	generatedBroker := &AlpacaBroker{}
	generatedBroker.initialize(accountType, entryPercent)
	return generatedBroker
}

func (broker *AlpacaBroker) initialize(accountType string, entryPercent float64) {
	cred := configs.GetCredentials(accountType)
	broker.client = alpaca.NewClient(
		alpaca.ClientOpts{
			ApiKey:    cred.API_KEY,
			ApiSecret: cred.API_SECRET,
			BaseURL:   cred.BASE_URL,
		})
	account, _ := broker.client.GetAccount()
	clock, _ := broker.client.GetClock()
	broker.account = account
	broker.Clock = *clock
	broker.PortfolioValue = broker.account.PortfolioValue.InexactFloat64()
	broker.TransactionNums = 0
	broker.MaxPortfolioPercent = entryPercent
	broker.HasPosition = false
	broker.MinProfitThreshold = broker.CalculateMinProfitThreshold(1.0)
}

func (broker *AlpacaBroker) CalculateMinProfitThreshold(baseNum float64) float64 {
	return baseNum * (broker.PortfolioValue * broker.MaxPortfolioPercent) / 120000
}

func (broker *AlpacaBroker) refreshOrderStatus(orderID string) (string, *alpaca.Order) {
	newOrder, _ := alpaca.GetOrder(orderID)
	orderStatus := newOrder.Status

	return orderStatus, newOrder
}

func (broker *AlpacaBroker) monitorOrder(order *alpaca.Order) (*alpaca.Order, bool) {
	var success bool
	orderID := order.ID
	status, updatedOrder := broker.refreshOrderStatus(orderID)
	for success {
		switch status {
		case "new", "accepted", "partially_filled":
			time.Sleep(time.Second)
			status, updatedOrder = broker.refreshOrderStatus(orderID)
		case "filled":
			success = true
		case "done_for_day", "canceled", "expired", "replaced":
			success = false
		default:
			time.Sleep(time.Second)
			status, updatedOrder = broker.refreshOrderStatus(orderID)
		}
	}
	log.Printf("The final state of the order is %s\n", status)
	return updatedOrder, success
}

func (broker *AlpacaBroker) SubmitOrderAsync(qty float64, symbol, side, orderType, timeInForce string, channel chan *alpaca.Order, wg *sync.WaitGroup) {
	defer wg.Done()
	quantity := decimal.NewFromFloat(qty)
	order, _ := broker.client.PlaceOrder(
		alpaca.PlaceOrderRequest{
			AccountID:   broker.account.ID,
			AssetKey:    &symbol,
			Qty:         &quantity,
			Side:        alpaca.Side(side),
			Type:        alpaca.OrderType(orderType),
			TimeInForce: alpaca.TimeInForce(timeInForce),
		},
	)
	finalOrder, _ := broker.monitorOrder(order)
	channel <- finalOrder
}

func (broker *AlpacaBroker) CloseAllPositions() {
	broker.client.CloseAllPositions()
}

func (broker *AlpacaBroker) GetDailyProfit() float64 {
	return broker.account.PortfolioValue.InexactFloat64() - broker.PortfolioValue
}
