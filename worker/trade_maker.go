package worker

import (
	"github.com/gitslagga/gitbitex-spot/matching"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/service"
	"github.com/prometheus/common/log"
	"time"
)

type TradeMaker struct {
	tradeCh   chan *models.Trade
	logReader matching.LogReader
	logOffset int64
	logSeq    int64
}

func NewTradeMaker(logReader matching.LogReader) *TradeMaker {
	t := &TradeMaker{
		tradeCh:   make(chan *models.Trade, 1000),
		logReader: logReader,
	}

	lastTrade, err := mysql.SharedStore().GetLastTradeByProductId(logReader.GetProductId())
	if err != nil {
		panic(err)
	}
	if lastTrade != nil {
		t.logOffset = lastTrade.LogOffset
		t.logSeq = lastTrade.LogSeq
	}

	t.logReader.RegisterObserver(t)
	return t
}

func (t *TradeMaker) Start() {
	if t.logOffset > 0 {
		t.logOffset++
	}
	go t.logReader.Run(t.logSeq, t.logOffset)
	go t.runFlusher()
}

func (t *TradeMaker) OnOpenLog(log *matching.OpenLog, offset int64) {
	// do nothing
}

func (t *TradeMaker) OnDoneLog(log *matching.DoneLog, offset int64) {
	// do nothing
}

func (t *TradeMaker) OnMatchLog(log *matching.MatchLog, offset int64) {
	t.tradeCh <- &models.Trade{
		Id:           log.TradeId,
		ProductId:    log.ProductId,
		TakerOrderId: log.TakerOrderId,
		MakerOrderId: log.MakerOrderId,
		Price:        log.Price,
		Size:         log.Size,
		Side:         log.Side,
		Time:         log.Time,
		LogOffset:    offset,
		LogSeq:       log.Sequence,
	}
}

func (t *TradeMaker) runFlusher() {
	var trades []*models.Trade

	for {
		select {
		case trade := <-t.tradeCh:
			trades = append(trades, trade)

			if len(t.tradeCh) > 0 && len(trades) < 1000 {
				continue
			}

			// 确保入库成功
			for {
				err := service.AddTrades(trades)
				if err != nil {
					log.Error(err)
					time.Sleep(time.Second)
					continue
				}
				trades = nil
				break
			}
		}
	}
}
