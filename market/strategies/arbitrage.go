package strategies

import (
    // "fmt"

    "../../market/core"
    "../../market/constants"
)

type BestPrice struct {
    Exchange string
    Price float64
}

type BestPrices struct {
    Symbol string
    Ask BestPrice   // Самое дорогое предложение на покупку
    Bid BestPrice  // Самое дешевое предложение преддожение на продажу
}

func EmptyBestRpices(symbol string) *BestPrices {
    return &BestPrices { Symbol: symbol,
                         Ask: BestPrice { Exchange: "", Price: 0 },
                         Bid: BestPrice { Exchange: "", Price: 0 } }
}

type Quotations struct {
    Depth *core.Depth
    Signal *core.Signal
    IndexDepth int
    Limit float64
}

type Arbitrage struct {
    BaseStrategy
}

const (
    PROPERTY_QUOTATIONS = "quotations"
    PROPERTY_BESTPRICES = "best_prices"
    PROPERTY_STRATEFYACTION = "strategy_action"
    PROPERTY_INDEXDEPTH = "index_depth"
    PROPERTY_LIMIT = "limit"

)

func NewArbitrage(name string) *Arbitrage {
    arbitrage := new(Arbitrage)
    arbitrage.init(name)
    arbitrage.SetProperty(PROPERTY_QUOTATIONS, make(map[string]map[string]*Quotations))
    // Значения торговых условий поумолчанию
    arbitrage.SetPropertySymbol(PROPERTY_DEF_SYMBOL, PROPERTY_INDEXDEPTH, 0)
    arbitrage.SetPropertySymbol(PROPERTY_DEF_SYMBOL, PROPERTY_LIMIT, 0.001)
    // Инициализируем расчет торговой стратегии
    arbitrage.calculate = func (signal *core.Signal) bool {
        // Данная стратегия работает только по сигналам из стакана
        if signal.Entity != constants.ENTITY_DEPTH {
            return false
        }
        // Собираем все котировки в одно хранлище, группируя их по торговым парам и биржам
        arbitrage.GetProperty(PROPERTY_QUOTATIONS).(map[string]map[string]*Quotations)[signal.Symbol] = arbitrage.readQuotations(signal)
        // Выбираем лучшие цена на покупку и продажу
        arbitrage.SetProperty(PROPERTY_BESTPRICES, arbitrage.chooseBestPrices())
        // Создаем торговое действие на основе лучших цен
        strategyAction := arbitrage.createStrategyAction()
        // Отправляем торговое действие на исполнение
        if strategyAction != nil {
            arbitrage.SetProperty(PROPERTY_STRATEFYACTION, strategyAction)
            arbitrage.chAction<- arbitrage.GetProperty(PROPERTY_STRATEFYACTION).(*core.StrategyAction)
        }
        return true
    }
    return arbitrage
}
// Задача данной функции оставить предыдущие значения в стакане в случае когда придет пустой стакан из сигнала
func (arbitrage *Arbitrage) readQuotations(signal *core.Signal) map[string]*Quotations {
    // Делаем выборку из хранилища уже имеющиеся данные по торговой паре
    symbols := arbitrage.GetProperty(PROPERTY_QUOTATIONS).(map[string]map[string]*Quotations)[signal.Symbol]
    // Если еще ни разу по данной паре небыло информации, то инициализируем коллекцию
    if symbols == nil {
        symbols = make(map[string]*Quotations)
    }
    // Делаем выборку по бирже
    quotations := symbols[signal.Exchange]
    // Если по данной бирже еще небыло информации, то инициализируем коллекцию
    if quotations == nil {
        quotations = new(Quotations)
        quotations.Depth = signal.Depth()
    }
    if quotations.Depth == nil {
        quotations.Depth = signal.Depth()
    }
    // Если пришедшие данные по стакану не пустые то записываем их
    if len(signal.Depth().Asks) > 0 {
        quotations.Depth.Asks = signal.Depth().Asks
        // for price, ask := range signal.Depth().Asks {
        //     if ask.Amount > 0 {
        //         quotations.Depth.Asks[price] = ask
        //     }
        // }
    }
    if len(signal.Depth().Bids) > 0 {
        quotations.Depth.Bids = signal.Depth().Bids
        // for price, bid := range signal.Depth().Bids {
        //     if bid.Amount > 0 {
        //         quotations.Depth.Bids[price] = bid
        //     }
        // }
    }
    // Сохраняем некоторую вспомогательную информацию
    quotations.Signal = signal
    if arbitrage.GetPropertySymbol(signal.Symbol, PROPERTY_INDEXDEPTH) != nil {
        quotations.IndexDepth = int(arbitrage.GetPropertySymbol(signal.Symbol, PROPERTY_INDEXDEPTH).(float64))
    } else {
        quotations.IndexDepth = int(arbitrage.GetPropertySymbol(PROPERTY_DEF_SYMBOL, PROPERTY_INDEXDEPTH).(float64))
    }
    if arbitrage.GetPropertySymbol(signal.Symbol, PROPERTY_LIMIT) != nil {
        quotations.Limit = arbitrage.GetPropertySymbol(signal.Symbol, PROPERTY_LIMIT).(float64)
    } else {
        quotations.Limit = arbitrage.GetPropertySymbol(PROPERTY_DEF_SYMBOL, PROPERTY_LIMIT).(float64)
    }
    symbols[signal.Exchange] = quotations
    return symbols
}
// Функция выбора лучших цен на покупку и рподажу сгрупированных по торговым парам
func (arbitrage *Arbitrage) chooseBestPrices() map[string]*BestPrices {
    quotations := arbitrage.GetProperty(PROPERTY_QUOTATIONS).(map[string]map[string]*Quotations)
    bestBySymbol := make(map[string]*BestPrices)
    // Находим самое дорогое предложение на покупку
    for symbol, exchanges := range quotations {
        best := new(BestPrices)
        best.Symbol = symbol
        var limit float64
        for exchange, quotation := range exchanges {
            limit = quotation.Limit
            if quotation.Depth == nil {
                continue
            }
            if quotation.Depth.Asks != nil {
                if len(quotation.Depth.Asks) > 0 {
                    ask := core.GetOrderByIndex(quotation.Depth.Asks, quotation.IndexDepth)
                    if ask != nil {
                        if best.Ask.Price == 0 {
                            best.Ask.Price = ask.Price
                            best.Ask.Exchange = quotation.Signal.Exchange
                        } else {
                            if ask.Price > best.Ask.Price {
                                best.Ask.Price = ask.Price
                                best.Ask.Exchange = exchange
                            }
                        }
                    }
                }
            }
            // Находим самую дешевое предложение на продажу
            if quotation.Depth.Bids != nil {
                if len(quotation.Depth.Bids) >= 0 {
                    bid := core.GetOrderByIndex(quotation.Depth.Bids, quotation.IndexDepth)
                    if bid != nil {
                        if best.Bid.Price == 0 {
                            best.Bid.Price = bid.Price
                            best.Bid.Exchange = quotation.Signal.Exchange
                        } else {
                            if bid.Price < best.Bid.Price {
                                best.Bid.Price = bid.Price
                                best.Bid.Exchange = exchange
                            }
                        }
                    }
                }
            }
        }
        if (best.Ask.Price - best.Bid.Price) >= limit {
            bestBySymbol[symbol] = best
        } else {
            bestBySymbol[symbol] = EmptyBestRpices(symbol)
        }
    }
    return bestBySymbol
}
//
func (arbitrage *Arbitrage) createStrategyAction() *core.StrategyAction {
    return nil
}
