package module_test

import (
	"errors"
	"fmt"

	cosmosMath "cosmossdk.io/math"
	"github.com/evoluteai-network/evoluteai-chain/app/params"
	state "github.com/evoluteai-network/evoluteai-chain/x/emissions"
	"github.com/evoluteai-network/evoluteai-chain/x/emissions/module"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	reputer1StartAmount = 1337
	reputer2StartAmount = 6969
	worker1StartAmount  = 4242
	worker2StartAmount  = 1111
)

func getConstWeights() [2][4]uint64 {
	return [2][4]uint64{{10, 20, 60, 100}, {30, 40, 50, 70}}
}

func getConstWeightsOnlyWorkers() [2][4]uint64 {
	return [2][4]uint64{{10, 20, 0, 0}, {30, 40, 0, 0}}
}

func getConstZeroWeights() [2][4]uint64 {
	return [2][4]uint64{{0, 0, 0, 0}, {0, 0, 0, 0}}
}

func (s *ModuleTestSuite) TestRegistrationTotalStakeAmountSet() {
	topicIds, err := mockCreateTopics(s, 1)
	s.NoError(err, "Error creating topic")
	topicId := topicIds[0]
	_, err = mockCreateTopics(s, 1)
	s.NoError(err, "Error creating topic 2")
	s.Equal(uint64(1), topicId, "Topic ID should start at 1")
	_, err = mockSomeReputers(s, topicId)
	s.NoError(err, "Error creating reputers")
	_, err = mockSomeWorkers(s, topicId)
	s.NoError(err, "Error creating workers")
	totalStakeExpected := cosmosMath.NewUint(
		reputer1StartAmount + reputer2StartAmount + worker1StartAmount + worker2StartAmount)
	totalStake, err := s.emissionsKeeper.GetTotalStake(s.ctx)
	s.NoError(err, "Error getting total stake")
	s.Equal(totalStakeExpected, totalStake, "Total stake should be the sum of all participants stakes")
}

func (s *ModuleTestSuite) TestRegistrationTopicStakeAmountSet() {
	topicIds, err := mockCreateTopics(s, 1)
	s.NoError(err, "Error creating topic")
	topicId := topicIds[0]
	_, err = mockCreateTopics(s, 1)
	s.NoError(err, "Error creating topic 2")
	s.Equal(uint64(1), topicId, "Topic ID should start at 1")
	_, err = mockSomeReputers(s, topicId)
	s.NoError(err, "Error creating reputers")
	_, err = mockSomeWorkers(s, topicId)
	s.NoError(err, "Error creating workers")
	topicStakeExpected := cosmosMath.NewUint(
		reputer1StartAmount + reputer2StartAmount + worker1StartAmount + worker2StartAmount)
	topicStake, err := s.emissionsKeeper.GetTopicStake(s.ctx, topicId)
	s.NoError(err, "Error getting topic stake")
	s.Equal(topicStakeExpected, topicStake, "Topic stake should be the sum of all participants stakes")
}

func (s *ModuleTestSuite) TestGetParticipantEmissionsForTopicSimple() {
	s.UtilSetParams()
	topicIds, err := mockCreateTopics(s, 1)
	s.NoError(err, "Error creating topic")
	topicId := topicIds[0]
	_, err = mockCreateTopics(s, 1)
	s.NoError(err, "Error creating topic 2")
	s.Equal(uint64(1), topicId, "Topic ID should start at 1")
	reputers, err := mockSomeReputers(s, topicId)
	s.NoError(err, "Error creating reputers")
	workers, err := mockSomeWorkers(s, topicId)
	s.NoError(err, "Error creating workers")
	err = mockSetWeights(s, topicId, reputers, workers, getConstWeights())
	s.NoError(err, "Error setting weights")
	topicStake, err := s.emissionsKeeper.GetTopicStake(s.ctx, topicId)
	s.NoError(err, "Error getting topic stake")
	totalStake, err := s.emissionsKeeper.GetTotalStake(s.ctx)
	s.NoError(err, "Error getting total stake")
	cumulativeEmissions := cosmosMath.NewUint(5000)
	metDemand := cosmosMath.NewUint(0)
	rewards, err := module.GetParticipantEmissionsForTopic(
		s.ctx,
		s.appModule,
		topicId,
		&topicStake,
		&cumulativeEmissions,
		&metDemand,
		&totalStake,
	)
	s.NoError(err, "Cumulative emissions for zero blocks should just be 0")
	s.Require().Equal(4, len(rewards), "every worker should get some emissions")
	for _, emission := range rewards {
		s.Require().Equal(
			cosmosMath.ZeroUint().LT(*emission),
			true,
			"Rewards emissions for every actor should not be 0")
	}
}

func (s *ModuleTestSuite) TestGetParticipantEmissionsForTopicNoReputerEmissions() {
	s.UtilSetParams()
	topicIds, err := mockCreateTopics(s, 1)
	s.NoError(err, "Error creating topic")
	topicId := topicIds[0]
	_, err = mockCreateTopics(s, 1)
	s.NoError(err, "Error creating topic 2")
	s.Equal(uint64(1), topicId, "Topic ID should start at 1")
	reputers, err := mockSomeReputers(s, topicId)
	s.NoError(err, "Error creating reputers")
	workers, err := mockSomeWorkers(s, topicId)
	s.NoError(err, "Error creating workers")
	err = mockSetWeights(s, topicId, reputers, workers, getConstWeightsOnlyWorkers())
	s.NoError(err, "Error setting weights")
	topicStake, err := s.emissionsKeeper.GetTopicStake(s.ctx, topicId)
	s.NoError(err, "Error getting topic stake")
	totalStake, err := s.emissionsKeeper.GetTotalStake(s.ctx)
	s.NoError(err, "Error getting total stake")
	cumulativeEmissions := cosmosMath.NewUint(5000)
	metDemand := cosmosMath.NewUint(0)
	rewards, err := module.GetParticipantEmissionsForTopic(
		s.ctx,
		s.appModule,
		topicId,
		&topicStake,
		&cumulativeEmissions,
		&metDemand,
		&totalStake,
	)
	s.NoError(err, "Cumulative emissions for zero blocks should just be 0")
	s.Require().Equal(2, len(rewards), "every worker should get some emissions")
	for _, emission := range rewards {
		s.Require().Equal(
			cosmosMath.ZeroUint().LT(*emission),
			true,
			"Worker emissions should not be 0")
	}
}

func (s *ModuleTestSuite) TestGetParticipantEmissionsForTopicNoWeights() {
	topicIds, err := mockCreateTopics(s, 1)
	s.NoError(err, "Error creating topic")
	topicId := topicIds[0]
	s.Equal(uint64(1), topicId, "Topic ID should start at 1")
	reputers, err := mockSomeReputers(s, topicId)
	s.NoError(err, "Error creating reputers")
	workers, err := mockSomeWorkers(s, topicId)
	s.NoError(err, "Error creating workers")
	err = mockSetWeights(s, topicId, reputers, workers, getConstZeroWeights())
	s.NoError(err, "Error setting weights")
	topicStake, err := s.emissionsKeeper.GetTopicStake(s.ctx, topicId)
	s.NoError(err, "Error getting topic stake")
	totalStake, err := s.emissionsKeeper.GetTotalStake(s.ctx)
	s.NoError(err, "Error getting total stake")
	cumulativeEmissions := cosmosMath.NewUint(5000)
	metDemand := cosmosMath.NewUint(0)
	rewards, err := module.GetParticipantEmissionsForTopic(
		s.ctx,
		s.appModule,
		topicId,
		&topicStake,
		&cumulativeEmissions,
		&metDemand,
		&totalStake,
	)
	s.NoError(err, "Cumulative emissions for zero blocks should just be 0")
	s.Require().Equal(2, len(rewards), "every worker should get some emissions")
	reputersAsString := make(map[string]struct{})
	for _, reputer := range reputers {
		reputersAsString[reputer.String()] = struct{}{}
	}
	for participant, reward := range rewards {
		s.Require().Equal(
			cosmosMath.ZeroUint().LT(*reward),
			true,
			"Rewards emissions for every actor should not be 0")
		_, isReputer := reputersAsString[participant]
		s.Require().True(isReputer, "Only reputers get paid when all weights are zero")
	}
}

/*************************************************
 *               HELPER FUNCTIONS				 *
 *												 *
 *************************************************/

// mock mint coins to participants
func mockMintRewardCoins(s *ModuleTestSuite, amount []cosmosMath.Int, target []sdk.AccAddress) error {
	if len(amount) != len(target) {
		return fmt.Errorf("amount and target must be the same length")
	}
	for i, addr := range target {
		coins := sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, amount[i]))
		s.bankKeeper.MintCoins(s.ctx, state.evoluteaiStakingAccountName, coins)
		s.bankKeeper.SendCoinsFromModuleToAccount(s.ctx, state.evoluteaiStakingAccountName, addr, coins)
	}
	return nil
}

// give some reputers coins, have them stake those coins
func mockSomeReputers(s *ModuleTestSuite, topicId uint64) ([]sdk.AccAddress, error) {
	reputerAddrs := []sdk.AccAddress{
		s.addrs[0],
		s.addrs[1],
	}
	reputerAmounts := []cosmosMath.Int{
		cosmosMath.NewInt(reputer1StartAmount),
		cosmosMath.NewInt(reputer2StartAmount),
	}
	err := mockMintRewardCoins(
		s,
		reputerAmounts,
		reputerAddrs,
	)
	if err != nil {
		return nil, err
	}
	_, err = s.msgServer.Register(s.ctx, &state.MsgRegister{
		Creator:      reputerAddrs[0].String(),
		LibP2PKey:    "libp2pkeyReputer1",
		MultiAddress: "multiaddressReputer1",
		TopicIds:     []uint64{topicId},
		InitialStake: cosmosMath.NewUintFromBigInt(reputerAmounts[0].BigInt()),
		IsReputer:    true,
	})
	if err != nil {
		return nil, err
	}
	_, err = s.msgServer.Register(s.ctx, &state.MsgRegister{
		Creator:      reputerAddrs[1].String(),
		LibP2PKey:    "libp2pkeyReputer2",
		MultiAddress: "multiaddressReputer2",
		TopicIds:     []uint64{topicId},
		InitialStake: cosmosMath.NewUintFromBigInt(reputerAmounts[1].BigInt()),
		IsReputer:    true,
	})
	if err != nil {
		return nil, err
	}
	return reputerAddrs, nil
}

// give some workers coins, have them stake those coins
func mockSomeWorkers(s *ModuleTestSuite, topicId uint64) ([]sdk.AccAddress, error) {
	workerAddrs := []sdk.AccAddress{
		s.addrs[2],
		s.addrs[3],
	}
	workerAmounts := []cosmosMath.Int{
		cosmosMath.NewInt(worker1StartAmount),
		cosmosMath.NewInt(worker2StartAmount),
	}
	err := mockMintRewardCoins(
		s,
		workerAmounts,
		workerAddrs,
	)
	if err != nil {
		return nil, err
	}
	_, err = s.msgServer.Register(s.ctx, &state.MsgRegister{
		Creator:      workerAddrs[0].String(),
		LibP2PKey:    "libp2pkeyWorker1",
		MultiAddress: "multiaddressWorker1",
		TopicIds:     []uint64{topicId},
		InitialStake: cosmosMath.NewUintFromBigInt(workerAmounts[0].BigInt()),
		Owner:        workerAddrs[0].String(),
	})
	if err != nil {
		return nil, err
	}
	_, err = s.msgServer.Register(s.ctx, &state.MsgRegister{
		Creator:      workerAddrs[1].String(),
		LibP2PKey:    "libp2pkeyWorker2",
		MultiAddress: "multiaddressWorker2",
		TopicIds:     []uint64{topicId},
		InitialStake: cosmosMath.NewUintFromBigInt(workerAmounts[1].BigInt()),
		Owner:        workerAddrs[1].String(),
	})
	if err != nil {
		return nil, err
	}
	return workerAddrs, nil
}

func mockSetWeights(
	s *ModuleTestSuite,
	topicId uint64,
	reputers []sdk.AccAddress,
	workers []sdk.AccAddress,
	weights [2][4]uint64) error {
	var target sdk.AccAddress
	for i := 0; i < 2; i++ {
		reputer := reputers[i]
		for j := 0; j < 4; j++ {
			weight := weights[i][j]
			if j > 1 {
				target = reputers[j-2]
			} else {
				target = workers[j]
			}
			err := s.emissionsKeeper.SetWeight(
				s.ctx,
				topicId,
				reputer,
				target,
				cosmosMath.NewUint(weight),
			)
			if err != nil {
				if !(errors.Is(err, state.ErrDoNotSetMapValueToZero)) {
					return err
				}
			}
		}
	}
	return nil
}

// create a topic
func mockCreateTopics(s *ModuleTestSuite, numToCreate uint64) ([]uint64, error) {
	ret := make([]uint64, 0)
	var i uint64
	for i = 0; i < numToCreate; i++ {
		topicMessage := state.MsgCreateNewTopic{
			Creator:          s.addrsStr[0],
			Metadata:         "metadata",
			WeightLogic:      "logic",
			WeightMethod:     "whatever",
			WeightCadence:    10800,
			InferenceLogic:   "morelogic",
			InferenceMethod:  "whatever2",
			InferenceCadence: 60,
		}
		response, err := s.msgServer.CreateNewTopic(s.ctx, &topicMessage)
		if err != nil {
			return nil, err
		}
		ret = append(ret, response.TopicId)
	}
	return ret, nil
}
