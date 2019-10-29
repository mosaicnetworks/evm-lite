package state

import (
	"github.com/ethereum/go-ethereum/common"
	ethState "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/trie"
)

func (s *State) ExportAllAccounts() (AccountRangeResult, error) {
	trie, err := s.main.stateDB.Database().OpenTrie(common.Hash{0})
	if err != nil {
		s.logger.Errorf("Error in ExportAllAccounts: %v", err)
		return AccountRangeResult{}, err
	}
	return accountRange(trie, nil, 999999999999)

}

// AccountRangeResult returns a mapping from the hash of an account addresses
// to its preimage. It will return the JSON null if no preimage is found.
// Since a query can return a limited amount of results, a "next" field is
// also present for paging.
type AccountRangeResult struct {
	Accounts map[common.Hash]*common.Address `json:"accounts"`
	Next     common.Hash                     `json:"next"`
}

func accountRange(st ethState.Trie, start *common.Hash, maxResults int) (AccountRangeResult, error) {

	//TODO remove this limit, it is debug code
	AccountRangeMaxResults := 24

	if start == nil {
		start = &common.Hash{0}
	}
	it := trie.NewIterator(st.NodeIterator(start.Bytes()))
	result := AccountRangeResult{Accounts: make(map[common.Hash]*common.Address), Next: common.Hash{}}

	if maxResults > AccountRangeMaxResults {
		maxResults = AccountRangeMaxResults
	}

	for i := 0; i < maxResults && it.Next(); i++ {
		if preimage := st.GetKey(it.Key); preimage != nil {
			addr := &common.Address{}
			addr.SetBytes(preimage)
			result.Accounts[common.BytesToHash(it.Key)] = addr
		} else {
			result.Accounts[common.BytesToHash(it.Key)] = nil
		}
	}

	if it.Next() {
		result.Next = common.BytesToHash(it.Key)
	}

	return result, nil
}
