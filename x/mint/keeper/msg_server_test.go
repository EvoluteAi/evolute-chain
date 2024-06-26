package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/evoluteai-network/evoluteai-chain/x/mint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *IntegrationTestSuite) TestUpdateParams() {
	testCases := []struct {
		name      string
		request   *types.MsgUpdateParams
		expectErr bool
	}{
		{
			name: "set invalid authority (not an address)",
			request: &types.MsgUpdateParams{
				Authority: "foo",
			},
			expectErr: true,
		},
		{
			name: "set invalid authority (not defined authority)",
			request: &types.MsgUpdateParams{
				Authority: "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
			},
			expectErr: true,
		},
		{
			name: "set invalid params",
			request: &types.MsgUpdateParams{
				Authority: s.mintKeeper.GetAuthority(),
				Params: types.Params{
					MintDenom:             sdk.DefaultBondDenom,
					InflationRateChange:   sdkmath.LegacyNewDecWithPrec(-13, 2),
					InflationMax:          sdkmath.LegacyNewDecWithPrec(20, 2),
					InflationMin:          sdkmath.LegacyNewDecWithPrec(7, 2),
					GoalBonded:            sdkmath.LegacyNewDecWithPrec(67, 2),
					BlocksPerYear:         uint64(60 * 60 * 8766 / 5),
					MaxSupply:             sdkmath.NewUintFromString("1000000000000000000000000000"),
					HalvingInterval:       uint64(25246080),
					CurrentBlockProvision: sdkmath.NewUintFromString("2831000000000000000000"),
				},
			},
			expectErr: true,
		},
		{
			name: "set full valid params",
			request: &types.MsgUpdateParams{
				Authority: s.mintKeeper.GetAuthority(),
				Params: types.Params{
					MintDenom:             sdk.DefaultBondDenom,
					InflationRateChange:   sdkmath.LegacyNewDecWithPrec(8, 2),
					InflationMax:          sdkmath.LegacyNewDecWithPrec(20, 2),
					InflationMin:          sdkmath.LegacyNewDecWithPrec(2, 2),
					GoalBonded:            sdkmath.LegacyNewDecWithPrec(37, 2),
					BlocksPerYear:         uint64(60 * 60 * 8766 / 5),
					MaxSupply:             sdkmath.NewUintFromString("1000000000000000000000000000"),
					HalvingInterval:       uint64(25246080),
					CurrentBlockProvision: sdkmath.NewUintFromString("2831000000000000000000"),
				},
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			_, err := s.msgServer.UpdateParams(s.ctx, tc.request)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
