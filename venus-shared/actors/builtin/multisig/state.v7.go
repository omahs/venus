// FETCHED FROM LOTUS: builtin/multisig/state.go.template

package multisig

import (
	"bytes"
	"encoding/binary"

	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/venus/venus-shared/actors/adt"

	builtin7 "github.com/filecoin-project/specs-actors/v7/actors/builtin"

	msig7 "github.com/filecoin-project/specs-actors/v7/actors/builtin/multisig"
	adt7 "github.com/filecoin-project/specs-actors/v7/actors/util/adt"
)

var _ State = (*state7)(nil)

func load7(store adt.Store, root cid.Cid) (State, error) {
	out := state7{store: store}
	err := store.Get(store.Context(), root, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func make7(store adt.Store, signers []address.Address, threshold uint64, startEpoch abi.ChainEpoch, unlockDuration abi.ChainEpoch, initialBalance abi.TokenAmount) (State, error) {
	out := state7{store: store}
	out.State = msig7.State{}
	out.State.Signers = signers
	out.State.NumApprovalsThreshold = threshold
	out.State.StartEpoch = startEpoch
	out.State.UnlockDuration = unlockDuration
	out.State.InitialBalance = initialBalance

	em, err := adt7.StoreEmptyMap(store, builtin7.DefaultHamtBitwidth)
	if err != nil {
		return nil, err
	}

	out.State.PendingTxns = em

	return &out, nil
}

type state7 struct {
	msig7.State
	store adt.Store
}

func (s *state7) LockedBalance(currEpoch abi.ChainEpoch) (abi.TokenAmount, error) {
	return s.State.AmountLocked(currEpoch - s.State.StartEpoch), nil
}

func (s *state7) StartEpoch() (abi.ChainEpoch, error) {
	return s.State.StartEpoch, nil
}

func (s *state7) UnlockDuration() (abi.ChainEpoch, error) {
	return s.State.UnlockDuration, nil
}

func (s *state7) InitialBalance() (abi.TokenAmount, error) {
	return s.State.InitialBalance, nil
}

func (s *state7) Threshold() (uint64, error) {
	return s.State.NumApprovalsThreshold, nil
}

func (s *state7) Signers() ([]address.Address, error) {
	return s.State.Signers, nil
}

func (s *state7) ForEachPendingTxn(cb func(id int64, txn Transaction) error) error {
	arr, err := adt7.AsMap(s.store, s.State.PendingTxns, builtin7.DefaultHamtBitwidth)
	if err != nil {
		return err
	}
	var out msig7.Transaction
	return arr.ForEach(&out, func(key string) error {
		txid, n := binary.Varint([]byte(key))
		if n <= 0 {
			return fmt.Errorf("invalid pending transaction key: %v", key)
		}
		return cb(txid, (Transaction)(out)) //nolint:unconvert
	})
}

func (s *state7) PendingTxnChanged(other State) (bool, error) {
	other7, ok := other.(*state7)
	if !ok {
		// treat an upgrade as a change, always
		return true, nil
	}
	return !s.State.PendingTxns.Equals(other7.PendingTxns), nil
}

func (s *state7) transactions() (adt.Map, error) {
	return adt7.AsMap(s.store, s.PendingTxns, builtin7.DefaultHamtBitwidth)
}

func (s *state7) decodeTransaction(val *cbg.Deferred) (Transaction, error) {
	var tx msig7.Transaction
	if err := tx.UnmarshalCBOR(bytes.NewReader(val.Raw)); err != nil {
		return Transaction{}, err
	}
	return Transaction(tx), nil
}

func (s *state7) GetState() interface{} {
	return &s.State
}
