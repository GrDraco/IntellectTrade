package strategies

import (
    "../../market/core"
    "../../market/constants"
)

type Arbitrage struct {
    BaseStrategy
}

func NewArbitrage(name string) *Arbitrage {
    arbitrage := new(Arbitrage)
    arbitrage.init(name)
    return arbitrage
}

func (arbitrage *Arbitrage) CalculateAction(signal *core.Signal) bool {
    switch signal.Entity {
    case constants.ENTITY_TICK:
        arbitrage.SetProperty("test", signal.TickStr(true))
    case constants.ENTITY_DEPTH:
        arbitrage.SetProperty("test", signal.DepthStr(true))
    case constants.ENTITY_CANDLE:
        arbitrage.SetProperty("test", signal.CandlesStr(true))
    }
    return false
}
