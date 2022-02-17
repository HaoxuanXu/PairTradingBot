package broker

import (
	"log"
	"sync"
	"time"

	"github.com/HaoxuanXu/TradingBot/configs"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/shopspring/decimal"
)

var lock = &sync.Mutex{}

type AlpacaBroker struct {
	client              alpaca.Client
	account             *alpaca.Account
	Clock               alpaca.Clock
	PortfolioValue      float64
	TransactionNums     int
	MaxPortfolioPercent float64
	HasPosition         bool
	LastTradeTime       time.Time
}

var (
	generatedBroker *AlpacaBroker
)

// You can treat this as a constructor of the broker class
func GetBroker(accountType string, entryPercent float64) *AlpacaBroker {

	lock.Lock()
	defer lock.Unlock()

	if generatedBroker == nil {
		generatedBroker = &AlpacaBroker{}
		generatedBroker.initialize(accountType, entryPercent)
	}
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
	broker.PortfolioValue = broker.account.Equity.InexactFloat64()
	broker.TransactionNums = 0
	broker.MaxPortfolioPercent = entryPercent
	broker.HasPosition = false
	broker.LastTradeTime = time.Now()
}

func (broker *AlpacaBroker) refreshOrderStatus(orderID string) (string, *alpaca.Order) {
	newOrder, _ := broker.client.GetOrder(orderID)
	orderStatus := newOrder.Status

	return orderStatus, newOrder
}

func (broker *AlpacaBroker) UpdateLastTradeTime() {
	broker.LastTradeTime = time.Now()
}

func (broker *AlpacaBroker) MonitorOrder(order *alpaca.Order) (*alpaca.Order, bool) {
	success := false
	orderID := order.ID
	status, updatedOrder := broker.refreshOrderStatus(orderID)
	for !success {
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
	return updatedOrder, success
}

func (broker *AlpacaBroker) SubmitOrderAsync(qty float64, symbol, side, orderType, timeInForce string, channel chan *alpaca.Order) {
	quantity := decimal.NewFromFloat(qty)
	order, _ := broker.client.PlaceOrder(
		alpaca.PlaceOrderRequest{
			AssetKey:    &symbol,
			AccountID:   broker.account.ID,
			Qty:         &quantity,
			Side:        alpaca.Side(side),
			Type:        alpaca.OrderType(orderType),
			TimeInForce: alpaca.TimeInForce(timeInForce),
		},
	)
	finalOrder, _ := broker.MonitorOrder(order)
	channel <- finalOrder
}

func (broker *AlpacaBroker) ListPositions() []alpaca.Position {
	positions, err := broker.client.ListPositions()
	if err != nil {
		log.Panic(err)
	}
	return positions
}

func (broker *AlpacaBroker) GetPosition(symbol string) *alpaca.Position {
	position, err := broker.client.GetPosition(symbol)
	if err != nil {
		log.Println(err)
	}
	return position
}

func (broker *AlpacaBroker) CloseAllPositions() {
	broker.client.CloseAllPositions()
}

func (broker *AlpacaBroker) GetDailyProfit() float64 {
	newAccount, _ := broker.client.GetAccount()
	return newAccount.Equity.InexactFloat64() - broker.PortfolioValue
}
