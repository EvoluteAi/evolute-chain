package module_test

import (
	"time"

	cosmosMath "cosmossdk.io/math"
	state "github.com/evoluteai-network/evoluteai-chain/x/emissions"
)

// this test must live in the module_test repo to actually do real non-mocked funds transfer tests
func (s *ModuleTestSuite) TestRequestInferenceInvalidCustomerNotEnoughFunds() {
	timeNow := uint64(time.Now().UTC().Unix())
	var initialStake int64 = 100
	_, err := mockCreateTopics(s, 1)
	s.Require().NoError(err)
	r := state.MsgRequestInference{
		Sender: s.addrsStr[0],
		Requests: []*state.RequestInferenceListItem{
			{
				Nonce:                0,
				TopicId:              1,
				Cadence:              0,
				MaxPricePerInference: cosmosMath.NewUint(uint64(initialStake)),
				BidAmount:            cosmosMath.NewUint(uint64(initialStake)),
				TimestampValidUntil:  timeNow + 100,
				ExtraData:            []byte("Test"),
			},
		},
	}
	_, err = s.msgServer.RequestInference(s.ctx, &r)
	s.Require().Error(err, "spendable balance 0uallo is smaller than 100uallo: insufficient funds")
}
