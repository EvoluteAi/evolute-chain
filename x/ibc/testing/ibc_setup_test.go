package testing

import (
	"cosmossdk.io/log"
	"encoding/json"
	"fmt"
	"testing"

	"cosmossdk.io/math"
	app2 "github.com/evoluteai-network/evoluteai-chain/app"
	"github.com/evoluteai-network/evoluteai-chain/app/params"
	dbm "github.com/cosmos/cosmos-db"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"
	"github.com/stretchr/testify/suite"
)

var (
	nativeDenom            = params.DefaultBondDenom
	ibcTransferAmount      = math.NewInt(100_000)
	genesisWalletAmount, _ = math.NewIntFromString("10000000000000000000")
)

type IBCTestSuite struct {
	suite.Suite

	coordinator   *ibctesting.Coordinator
	evoluteaiChain   *ibctesting.TestChain // aka chainA
	providerChain *ibctesting.TestChain

	path *ibctesting.Path

	providerAddr          sdk.AccAddress
	evoluteaiAddr            sdk.AccAddress
	providerToevoluteaiDenom string
}

func TestIBCTestSuite(t *testing.T) {
	suite.Run(t, new(IBCTestSuite))
}

func evoluteaiAppInitializer() (ibctesting.TestingApp, map[string]json.RawMessage) {
	app, err := app2.NewevoluteaiApp(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		simtestutil.EmptyAppOptions{},
	)
	if err != nil {
		fmt.Printf("Initializing app error: %v\n", err)
		return nil, nil
	}

	return app, app.DefaultGenesis()
}

func (s *IBCTestSuite) SetupTest() {
	// we need to redefine this variable to make tests work cause we use untrn as default bond denom in evoluteai
	sdk.DefaultBondDenom = nativeDenom
	ibctesting.DefaultTestingAppInit = evoluteaiAppInitializer
	//params.InitSDKConfig()

	s.coordinator = ibctesting.NewCoordinator(s.T(), 2)
	s.evoluteaiChain = s.coordinator.GetChain(ibctesting.GetChainID(1))
	s.providerChain = s.coordinator.GetChain(ibctesting.GetChainID(2))

	s.path = ibctesting.NewPath(s.evoluteaiChain, s.providerChain)
	s.path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	s.path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort
	s.path.EndpointA.ChannelConfig.Version = transfertypes.Version
	s.path.EndpointB.ChannelConfig.Version = transfertypes.Version

	s.coordinator.Setup(s.path)

	// Store ibc transfer denom for providerChain=>evoluteai for test convenience
	fullTransferDenomPath := transfertypes.GetPrefixedDenom(
		transfertypes.PortID,
		s.path.EndpointB.ChannelID,
		nativeDenom,
	)
	transferDenom := transfertypes.ParseDenomTrace(fullTransferDenomPath).IBCDenom()
	s.providerToevoluteaiDenom = transferDenom

	// Store default addresses from evoluteai and provider chain for test convenience
	s.providerAddr = s.providerChain.SenderAccount.GetAddress()
	s.evoluteaiAddr = s.evoluteaiChain.SenderAccount.GetAddress()

	// ensure genesis balances are as expected
	s.assertevoluteaiBalance(s.evoluteaiAddr, nativeDenom, genesisWalletAmount)
	s.assertProviderBalance(s.providerAddr, nativeDenom, genesisWalletAmount)
}

func (s *IBCTestSuite) IBCTransfer(
	path *ibctesting.Path,
	sourceEndpoint *ibctesting.Endpoint,
	fromAddr sdk.AccAddress,
	toAddr sdk.AccAddress,
	transferDenom string,
	transferAmount math.Int,
	memo string,
) {
	timeoutHeight := clienttypes.NewHeight(1, 110)

	// Create Transfer Msg
	transferMsg := transfertypes.NewMsgTransfer(sourceEndpoint.ChannelConfig.PortID,
		sourceEndpoint.ChannelID,
		sdk.NewCoin(transferDenom, transferAmount),
		fromAddr.String(),
		toAddr.String(),
		timeoutHeight,
		0,
		memo,
	)

	// Send message from provider chain
	res, err := sourceEndpoint.Chain.SendMsgs(transferMsg)
	s.Assert().NoError(err)

	// Relay transfer msg to evoluteai chain
	packet, err := ibctesting.ParsePacketFromEvents(res.GetEvents())
	s.Require().NoError(err)

	//nolint:errcheck // this will return an error for multi-hop routes; that's expected
	path.RelayPacket(packet)
}

func (s *IBCTestSuite) IBCTransferProviderToevoluteai(
	providerAddr sdk.AccAddress,
	evoluteaiAddr sdk.AccAddress,
	transferDenom string,
	transferAmount math.Int,
	memo string,
) {
	s.IBCTransfer(
		s.path,
		s.path.EndpointB,
		providerAddr,
		evoluteaiAddr,
		transferDenom,
		transferAmount,
		memo,
	)
}

func (s *IBCTestSuite) IBCTransferevoluteaiToProvider(
	evoluteaiAddr sdk.AccAddress,
	providerAddr sdk.AccAddress,
	transferDenom string,
	transferAmount math.Int,
	memo string,
) {
	s.IBCTransfer(
		s.path,
		s.path.EndpointA,
		evoluteaiAddr,
		providerAddr,
		transferDenom,
		transferAmount,
		memo,
	)
}

func (s *IBCTestSuite) getBalance(
	bk bankkeeper.Keeper,
	chain *ibctesting.TestChain,
	addr sdk.AccAddress,
	denom string,
) sdk.Coin {
	ctx := chain.GetContext()
	return bk.GetBalance(ctx, addr, denom)
}

func (s *IBCTestSuite) assertBalance(
	bk bankkeeper.Keeper,
	chain *ibctesting.TestChain,
	addr sdk.AccAddress,
	denom string,
	expectedAmt math.Int,
) {
	actualAmt := s.getBalance(bk, chain, addr, denom).Amount
	s.Assert().
		Equal(expectedAmt, actualAmt, "Expected amount of %s: %s; Got: %s", denom, expectedAmt, actualAmt)
}

func (s *IBCTestSuite) assertevoluteaiBalance(
	addr sdk.AccAddress,
	denom string,
	expectedAmt math.Int,
) {
	app, _ := s.evoluteaiChain.App.(*app2.evoluteaiApp)
	s.assertBalance(app.BankKeeper, s.evoluteaiChain, addr, denom, expectedAmt)
}

func (s *IBCTestSuite) assertProviderBalance(
	addr sdk.AccAddress,
	denom string,
	expectedAmt math.Int,
) {
	app, _ := s.providerChain.App.(*app2.evoluteaiApp)
	s.assertBalance(app.BankKeeper, s.providerChain, addr, denom, expectedAmt)
}
