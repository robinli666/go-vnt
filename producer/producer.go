// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package producer implements Hubble block creation and producing.
package producer

import (
	"fmt"
	"sync/atomic"

	"github.com/vntchain/go-vnt/accounts"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/consensus"
	"github.com/vntchain/go-vnt/core"
	"github.com/vntchain/go-vnt/core/state"
	"github.com/vntchain/go-vnt/core/types"
	"github.com/vntchain/go-vnt/event"
	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/go-vnt/params"
	"github.com/vntchain/go-vnt/vnt/downloader"
	"github.com/vntchain/go-vnt/vntdb"
)

// Backend wraps all methods required for producing.
type Backend interface {
	AccountManager() *accounts.Manager
	BlockChain() *core.BlockChain
	TxPool() *core.TxPool
	ChainDb() vntdb.Database
}

// Producer creates blocks and searches for proof-of-work values.
type Producer struct {
	mux *event.TypeMux

	worker    *worker
	coinbase  common.Address
	producing int32
	vnt       Backend
	engine    consensus.Engine

	canStart    int32 // can start indicates whether we can start the block producing operation
	shouldStart int32 // should start indicates whether we should start after sync
}

func New(vnt Backend, config *params.ChainConfig, mux *event.TypeMux, engine consensus.Engine) *Producer {
	producer := &Producer{
		vnt:      vnt,
		mux:      mux,
		engine:   engine,
		worker:   newWorker(config, engine, common.Address{}, vnt, mux),
		canStart: 1,
	}
	go producer.update()

	return producer
}

// update keeps track of the downloader events. Please be aware that this is a one shot type of update loop.
// It's entered once and as soon as `Done` or `Failed` has been broadcasted the events are unregistered and
// the loop is exited. This to prevent a major security vuln where external parties can DOS you with blocks
// and halt your producing operation for as long as the DOS continues.
func (self *Producer) update() {
	events := self.mux.Subscribe(downloader.StartEvent{}, downloader.DoneEvent{}, downloader.FailedEvent{})
out:
	for ev := range events.Chan() {
		switch ev.Data.(type) {
		case downloader.StartEvent:
			atomic.StoreInt32(&self.canStart, 0)
			if self.Producing() {
				self.Stop()
				atomic.StoreInt32(&self.shouldStart, 1)
				log.Info("Producing aborted due to sync")
			}
		case downloader.DoneEvent, downloader.FailedEvent:
			shouldStart := atomic.LoadInt32(&self.shouldStart) == 1

			atomic.StoreInt32(&self.canStart, 1)
			atomic.StoreInt32(&self.shouldStart, 0)
			if shouldStart {
				self.Start(self.coinbase)
			}
			// unsubscribe. we're only interested in this event once
			events.Unsubscribe()
			// stop immediately and ignore all further pending events
			break out
		}
	}
}

func (self *Producer) Start(coinbase common.Address) {
	atomic.StoreInt32(&self.shouldStart, 1)
	self.SetCoinbase(coinbase)

	if atomic.LoadInt32(&self.canStart) == 0 {
		log.Info("Network syncing, will start producer afterwards")
		return
	}
	atomic.StoreInt32(&self.producing, 1)

	log.Info("Starting block producing operation")
	self.worker.start()
	self.worker.commitNewWork()
}

func (self *Producer) Stop() {
	self.worker.stop()
	atomic.StoreInt32(&self.producing, 0)
	atomic.StoreInt32(&self.shouldStart, 0)
}

func (self *Producer) Producing() bool {
	return atomic.LoadInt32(&self.producing) > 0
}

func (self *Producer) SetExtra(extra []byte) error {
	if uint64(len(extra)) > params.MaximumExtraDataSize {
		return fmt.Errorf("Extra exceeds max length. %d > %v", len(extra), params.MaximumExtraDataSize)
	}
	self.worker.setExtra(extra)
	return nil
}

// Pending returns the currently pending block and associated state.
func (self *Producer) Pending() (*types.Block, *state.StateDB) {
	return self.worker.pending()
}

// PendingBlock returns the currently pending block.
//
// Note, to access both the pending block and the pending state
// simultaneously, please use Pending(), as the pending state can
// change between multiple method calls
func (self *Producer) PendingBlock() *types.Block {
	return self.worker.pendingBlock()
}

func (self *Producer) SetCoinbase(addr common.Address) {
	self.coinbase = addr
	self.worker.setCoinbase(addr)
}
