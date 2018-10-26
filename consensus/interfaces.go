// Copyright (C) 2018 go-dappley authors
//
// This file is part of the go-dappley library.
//
// the go-dappley library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// the go-dappley library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-dappley library.  If not, see <http://www.gnu.org/licenses/>.
//

package consensus

import "github.com/dappley/go-dappley/core"

type MinedBlock struct {
	block   *core.Block
	isValid bool
}

type BlockProducer interface {
	// Setup tells the producer to give rewards to beneficiaryAddr and return the new block through retChan
	Setup(bc *core.Blockchain, beneficiaryAddr string, retChan chan *MinedBlock)

	SetPrivKey(key string)

	GetPrivKey() string

	// Beneficiary returns the address to receive rewards
	Beneficiary() string

	Start()

	Stop()

	// Validate returns true if blk is valid according to consensus
	Validate(blk *core.Block) bool
}
