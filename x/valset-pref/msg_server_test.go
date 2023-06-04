package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	appParams "github.com/osmosis-labs/osmosis/v16/app/params"
	valPref "github.com/osmosis-labs/osmosis/v16/x/valset-pref"
	"github.com/osmosis-labs/osmosis/v16/x/valset-pref/types"
)

func (s *KeeperTestSuite) TestSetValidatorSetPreference() {
	s.SetupTest()

	// setup 6 validators
	valAddrs := s.SetupMultipleValidators(6)

	tests := []struct {
		name                   string
		delegator              sdk.AccAddress
		preferences            []types.ValidatorPreference
		amountToDelegate       sdk.Coin // amount to delegate
		setExistingDelegations bool
		expectPass             bool
	}{
		{
			name:      "creation of new validator set, user does not have existing delegation",
			delegator: sdk.AccAddress([]byte("addr1---------------")),
			preferences: []types.ValidatorPreference{
				{
					ValOperAddress: valAddrs[0],
					Weight:         sdk.NewDecWithPrec(5, 1),
				},
				{
					ValOperAddress: valAddrs[1],
					Weight:         sdk.NewDecWithPrec(3, 1),
				},
				{
					ValOperAddress: valAddrs[2],
					Weight:         sdk.NewDecWithPrec(2, 1),
				},
			},
			expectPass: true,
		},
		{
			name:      "update 2 validator weights but leave the 3rd one as is",
			delegator: sdk.AccAddress([]byte("addr1---------------")),
			preferences: []types.ValidatorPreference{
				{
					ValOperAddress: valAddrs[0],
					Weight:         sdk.NewDecWithPrec(5, 1),
				},
				{
					ValOperAddress: valAddrs[1],
					Weight:         sdk.NewDecWithPrec(4, 1),
				},
				{
					ValOperAddress: valAddrs[2],
					Weight:         sdk.NewDecWithPrec(1, 1),
				},
			},
			expectPass: true,
		},
		{
			name:      "update existing validator with same valAddr and weights",
			delegator: sdk.AccAddress([]byte("addr1---------------")),
			preferences: []types.ValidatorPreference{
				{
					ValOperAddress: valAddrs[0],
					Weight:         sdk.NewDecWithPrec(5, 1),
				},
				{
					ValOperAddress: valAddrs[1],
					Weight:         sdk.NewDecWithPrec(4, 1),
				},
				{
					ValOperAddress: valAddrs[2],
					Weight:         sdk.NewDecWithPrec(1, 1),
				},
			},
			expectPass: false,
		},
		{
			name:      "update existing validator with same valAddr but different weights",
			delegator: sdk.AccAddress([]byte("addr1---------------")),
			preferences: []types.ValidatorPreference{
				{
					ValOperAddress: valAddrs[0],
					Weight:         sdk.NewDecWithPrec(1, 1),
				},
				{
					ValOperAddress: valAddrs[1],
					Weight:         sdk.NewDecWithPrec(2, 1),
				},
				{
					ValOperAddress: valAddrs[2],
					Weight:         sdk.NewDecWithPrec(7, 1),
				},
			},
			expectPass: true,
		},
		{
			name:      "create validator set with unknown validator address",
			delegator: sdk.AccAddress([]byte("addr1---------------")),
			preferences: []types.ValidatorPreference{
				{
					ValOperAddress: "addr1---------------",
					Weight:         sdk.NewDec(1),
				},
			},
			expectPass: false,
		},
		{
			name:      "user has valset, but does not have existing delegation",
			delegator: sdk.AccAddress([]byte("addr1---------------")),
			preferences: []types.ValidatorPreference{
				{
					ValOperAddress: valAddrs[3],
					Weight:         sdk.NewDecWithPrec(3, 1),
				},
				{
					ValOperAddress: valAddrs[4],
					Weight:         sdk.NewDecWithPrec(7, 1),
				},
			},
			expectPass: true, // SetValidatorSetPreference modifies the existing delegations
		},
		{
			name:      "user has existing valset and existing delegation",
			delegator: sdk.AccAddress([]byte("addr1---------------")),
			preferences: []types.ValidatorPreference{
				{
					ValOperAddress: valAddrs[3],
					Weight:         sdk.NewDecWithPrec(3, 1),
				},
				{
					ValOperAddress: valAddrs[4],
					Weight:         sdk.NewDecWithPrec(3, 1),
				},
				{
					ValOperAddress: valAddrs[5],
					Weight:         sdk.NewDecWithPrec(4, 1),
				},
			},
			amountToDelegate:       sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10_000_000)),
			setExistingDelegations: true,
			expectPass:             true,
		}, // SetValidatorSetPreference ignores the existing delegation and modifies the existing valset
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// setup message server
			msgServer := valPref.NewMsgServerImpl(s.App.ValidatorSetPreferenceKeeper)
			c := sdk.WrapSDKContext(s.Ctx)

			if test.setExistingDelegations {
				amountToFund := sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 100_000_000)} // 100 osmo
				s.FundAcc(test.delegator, amountToFund)

				err := s.PrepareExistingDelegations(s.Ctx, valAddrs, test.delegator, test.amountToDelegate.Amount)
				s.Require().NoError(err)
			}

			// call the sets new validator set preference
			_, err := msgServer.SetValidatorSetPreference(c, types.NewMsgSetValidatorSetPreference(test.delegator, test.preferences))
			if test.expectPass {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestDelegateToValidatorSet() {
	s.SetupTest()

	// valset test setup
	valAddrs, preferences, amountToFund := s.SetupValidatorsAndDelegations()

	tests := []struct {
		name                   string
		delegator              sdk.AccAddress
		amountToDelegate       sdk.Coin  // amount to delegate
		expectedShares         []sdk.Dec // expected shares after delegation
		setExistingDelegations bool      // if true, create new delegation (non-valset) with {delegator, valAddrs}
		setValSet              bool      // if true, create a new valset {delegator, preferences}
		expectPass             bool
	}{
		{
			name:             "Delegate to valid validators",
			delegator:        sdk.AccAddress([]byte("addr1---------------")),
			amountToDelegate: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10_000_000)),
			setValSet:        true,
			expectedShares:   []sdk.Dec{sdk.NewDec(2_000_000), sdk.NewDec(3_300_000), sdk.NewDec(1_200_000), sdk.NewDec(3_500_000)},
			expectPass:       true,
		},
		{
			name:             "Delegate more tokens to existing validator-set",
			delegator:        sdk.AccAddress([]byte("addr1---------------")),
			amountToDelegate: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10_000_000)),
			expectedShares:   []sdk.Dec{sdk.NewDec(4_000_000), sdk.NewDec(6_600_000), sdk.NewDec(2_400_000), sdk.NewDec(7_000_000)},
			expectPass:       true,
		},
		{
			name:             "User does not have enough tokens to stake",
			delegator:        sdk.AccAddress([]byte("addr2---------------")),
			amountToDelegate: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(200_000_000)),
			setValSet:        true,
			expectPass:       false,
		},
		{
			name:                   "Delegate to existing staking position (non valSet)",
			delegator:              sdk.AccAddress([]byte("addr3---------------")),
			amountToDelegate:       sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10_000_000)),
			expectedShares:         []sdk.Dec{sdk.NewDec(13_333_333), sdk.NewDec(13_333_333), sdk.NewDec(13_333_334)},
			setExistingDelegations: true,
			expectPass:             true,
		},
		{
			name:             "Delegate very small amount to existing valSet",
			delegator:        sdk.AccAddress([]byte("addr4---------------")),
			amountToDelegate: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0o10_013)), // small case
			expectedShares:   []sdk.Dec{sdk.NewDec(821), sdk.NewDec(1355), sdk.NewDec(492), sdk.NewDec(1439)},
			setValSet:        true,
			expectPass:       true,
		},
		{
			name:             "Delegate large amount to existing valSet",
			delegator:        sdk.AccAddress([]byte("addr5---------------")),
			amountToDelegate: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(96_386_414)),
			expectedShares:   []sdk.Dec{sdk.NewDec(19_277_282), sdk.NewDec(31_807_516), sdk.NewDec(11_566_369), sdk.NewDec(33_735_247)},
			setValSet:        true,
			expectPass:       true,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// setup message server
			msgServer := valPref.NewMsgServerImpl(s.App.ValidatorSetPreferenceKeeper)
			c := sdk.WrapSDKContext(s.Ctx)

			s.FundAcc(test.delegator, amountToFund)

			if test.setValSet {
				_, err := msgServer.SetValidatorSetPreference(c, types.NewMsgSetValidatorSetPreference(test.delegator, preferences))
				s.Require().NoError(err)
			}

			if test.setExistingDelegations {
				// we perform this operation len(valAddrs) times
				err := s.PrepareExistingDelegations(s.Ctx, valAddrs, test.delegator, test.amountToDelegate.Amount)
				s.Require().NoError(err)
			}

			_, err := msgServer.DelegateToValidatorSet(c, types.NewMsgDelegateToValidatorSet(test.delegator, test.amountToDelegate))
			if test.expectPass {
				s.Require().NoError(err)

				// check if the user balance decreased
				balance := s.App.BankKeeper.GetBalance(s.Ctx, test.delegator, sdk.DefaultBondDenom)
				expectedBalance := amountToFund[0].Amount.Sub(test.amountToDelegate.Amount)
				// valset has not been set so use the (expectedBalance = account balance)
				if !test.setValSet {
					expectedBalance = balance.Amount
				}

				s.Require().Equal(expectedBalance, balance.Amount)

				if test.setValSet {
					// check if the expectedShares matches after delegation
					for i, val := range preferences {
						valAddr, err := sdk.ValAddressFromBech32(val.ValOperAddress)
						s.Require().NoError(err)

						// guarantees that the delegator exists because we check it in DelegateToValidatorSet
						del, _ := s.App.StakingKeeper.GetDelegation(s.Ctx, test.delegator, valAddr)
						s.Require().Equal(test.expectedShares[i], del.Shares)
					}
				}

				if test.setExistingDelegations {
					delSharesAmt := sdk.NewDec(0)
					expectedAmt := sdk.NewDec(0)

					for i, val := range valAddrs {
						valAddr, err := sdk.ValAddressFromBech32(val)
						s.Require().NoError(err)

						// guarantees that the delegator exists because we check it in DelegateToValidatorSet
						del, _ := s.App.StakingKeeper.GetDelegation(s.Ctx, test.delegator, valAddr)
						delSharesAmt = delSharesAmt.Add(del.Shares)
						expectedAmt = expectedAmt.Add(test.expectedShares[i])
					}

					s.Require().Equal(expectedAmt, delSharesAmt)
				}
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestUnDelegateFromValidatorSet() {
	s.SetupTest()

	// valset test setup
	valAddrs, preferences, amountToFund := s.SetupValidatorsAndDelegations()

	tests := []struct {
		name                   string
		delegator              sdk.AccAddress
		coinToStake            sdk.Coin
		coinToUnStake          sdk.Coin
		expectedShares         []sdk.Dec // expected shares after undelegation
		setValSet              bool
		setExistingDelegations bool
		expectPass             bool
	}{
		{
			name:           "Unstake half from the ValSet",
			delegator:      sdk.AccAddress([]byte("addr1---------------")),
			coinToStake:    sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(20_000_000)), // delegate 20osmo
			coinToUnStake:  sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10_000_000)), // undelegate 10osmo
			expectedShares: []sdk.Dec{sdk.NewDec(2_000_000), sdk.NewDec(3_300_000), sdk.NewDec(1_200_000), sdk.NewDec(3_500_000)},
			setValSet:      true,
			expectPass:     true,
		},
		{
			name:           "Unstake x amount from ValSet",
			delegator:      sdk.AccAddress([]byte("addr2---------------")),
			coinToStake:    sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(20_000_000)),                                           // delegate 20osmo
			coinToUnStake:  sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(15_000_000)),                                           // undelegate 15osmo
			expectedShares: []sdk.Dec{sdk.NewDec(1_000_000), sdk.NewDec(1_650_000), sdk.NewDec(600_000), sdk.NewDec(1_750_000)}, // validatorDelegatedShares - (weight * coinToUnstake)
			setValSet:      true,
			expectPass:     true,
		},
		{
			name:          "Unstake everything",
			delegator:     sdk.AccAddress([]byte("addr3---------------")),
			coinToStake:   sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(20_000_000)),
			coinToUnStake: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(20_000_000)),
			setValSet:     true,
			expectPass:    true,
		},
		{
			name:          "Unstake more amount than the staked amount",
			delegator:     sdk.AccAddress([]byte("addr4---------------")),
			coinToStake:   sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(20_000_000)),
			coinToUnStake: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(40_000_000)),
			setValSet:     true,
			expectPass:    false,
		},
		{
			name:                   "UnDelegate from existing staking position (non valSet) ",
			delegator:              sdk.AccAddress([]byte("addr5---------------")),
			coinToStake:            sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10_000_000)),
			coinToUnStake:          sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10_000_000)),
			expectedShares:         []sdk.Dec{sdk.NewDec(1_000_000), sdk.NewDec(1_660_000), sdk.NewDec(600_000), sdk.NewDec(1_740_000)}, // validatorDelegatedShares - (weight * coinToUnstake)
			setExistingDelegations: true,
			expectPass:             true,
		},
		{
			name:           "Undelegate extreme amounts to check truncation, large amount",
			delegator:      sdk.AccAddress([]byte("addr6---------------")),
			coinToStake:    sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000)),
			coinToUnStake:  sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(87_461_351)),
			expectedShares: []sdk.Dec{sdk.NewDec(2_507_730), sdk.NewDec(4_137_755), sdk.NewDec(1_504_638), sdk.NewDec(4_388_526)}, // validatorDelegatedShares - (weight * coinToUnstake), for ex: 20_000_000 - (0.2 * 87_461_351)
			setValSet:      true,
			expectPass:     true,
		},
		{
			name:           "Undelegate extreme amounts to check truncation, small amount",
			delegator:      sdk.AccAddress([]byte("addr7---------------")),
			coinToStake:    sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10_000_000)),
			coinToUnStake:  sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1234)),
			expectedShares: []sdk.Dec{sdk.NewDec(1_999_754), sdk.NewDec(3_299_593), sdk.NewDec(1_199_852), sdk.NewDec(3_499_567)}, // validatorDelegatedShares - (weight * coinToUnstake),
			setValSet:      true,
			expectPass:     true,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.FundAcc(test.delegator, amountToFund) // 100 osmo

			// setup message server
			msgServer := valPref.NewMsgServerImpl(s.App.ValidatorSetPreferenceKeeper)
			c := sdk.WrapSDKContext(s.Ctx)

			if test.setValSet {
				// SetValidatorSetPreference sets a new list of val-set
				_, err := msgServer.SetValidatorSetPreference(c, types.NewMsgSetValidatorSetPreference(test.delegator, preferences))
				s.Require().NoError(err)

				// DelegateToValidatorSet delegate to existing val-set
				_, err = msgServer.DelegateToValidatorSet(c, types.NewMsgDelegateToValidatorSet(test.delegator, test.coinToStake))
				s.Require().NoError(err)
			}

			if test.setExistingDelegations {
				err := s.PrepareExistingDelegations(s.Ctx, valAddrs, test.delegator, test.coinToStake.Amount)
				s.Require().NoError(err)
			}

			_, err := msgServer.UndelegateFromValidatorSet(c, types.NewMsgUndelegateFromValidatorSet(test.delegator, test.coinToUnStake))
			if test.expectPass {
				s.Require().NoError(err)

				// check if the expectedShares matches after undelegation
				for i, val := range preferences {
					valAddr, err := sdk.ValAddressFromBech32(val.ValOperAddress)
					s.Require().NoError(err)

					// guarantees that the delegator exists because we check it in UnDelegateToValidatorSet
					del, found := s.App.StakingKeeper.GetDelegation(s.Ctx, test.delegator, valAddr)
					if found {
						s.Require().Equal(test.expectedShares[i], del.GetShares())
					}
				}
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestRedelegateToValidatorSet() {
	s.SetupTest()

	// prepare validators to delegate to
	preferences := s.PrepareDelegateToValidatorSet()

	valAddrs := s.SetupMultipleValidators(6)

	tests := []struct {
		name                        string
		delegator                   sdk.AccAddress
		newPreferences              []types.ValidatorPreference
		amountToDelegate            sdk.Coin  // amount to delegate
		expectedShares              []sdk.Dec // expected shares after delegation
		setExistingDelegation       bool      // ensures that there is existing delegations (non valset)
		setExistingValSetDelegation bool      // ensures that there is existing valset delegation
		expectPass                  bool
	}{
		{
			name:      "redelegate to a new set of validators",
			delegator: sdk.AccAddress([]byte("addr1---------------")),
			newPreferences: []types.ValidatorPreference{
				{
					ValOperAddress: valAddrs[0],
					Weight:         sdk.NewDecWithPrec(2, 1),
				},
				{
					ValOperAddress: valAddrs[1],
					Weight:         sdk.NewDecWithPrec(2, 1),
				},
				{
					ValOperAddress: valAddrs[2],
					Weight:         sdk.NewDecWithPrec(6, 1),
				},
			},
			amountToDelegate:            sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(20_000_000)),
			expectedShares:              []sdk.Dec{sdk.NewDec(4_000_000), sdk.NewDec(4_000_000), sdk.NewDec(12_000_000)},
			setExistingValSetDelegation: true,
			expectPass:                  true, // addr1 successfully redelegates to (valAddr0, valAddr1, valAddr2)
		},
		{
			name:      "redelegate to same set of validators",
			delegator: sdk.AccAddress([]byte("addr1---------------")),
			newPreferences: []types.ValidatorPreference{
				{
					ValOperAddress: valAddrs[0],
					Weight:         sdk.NewDecWithPrec(3, 1),
				},
				{
					ValOperAddress: valAddrs[1],
					Weight:         sdk.NewDecWithPrec(2, 1),
				},
				{
					ValOperAddress: valAddrs[2],
					Weight:         sdk.NewDecWithPrec(5, 1),
				},
			},
			amountToDelegate: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(20_000_000)),
			expectPass:       false, // first redelegation already in progress so must end that first
		},
		{
			name:      "redelegate to new set, but one validator from old set",
			delegator: sdk.AccAddress([]byte("addr1---------------")),
			newPreferences: []types.ValidatorPreference{
				{
					ValOperAddress: valAddrs[4],
					Weight:         sdk.NewDecWithPrec(5, 1),
				},
				{
					ValOperAddress: valAddrs[1],
					Weight:         sdk.NewDecWithPrec(3, 1),
				},
				{
					ValOperAddress: valAddrs[3],
					Weight:         sdk.NewDecWithPrec(2, 1),
				},
			},
			amountToDelegate: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(20_000_000)),
			expectedShares:   []sdk.Dec{sdk.NewDec(10_000_000), sdk.NewDec(6_000_000), sdk.NewDec(4_000_000)},
			expectPass:       false, // this fails because valAddrs[1] is being redelegated to in first test
		},
		{
			name:      "Redelegate to new valset with one existing delegation validator",
			delegator: sdk.AccAddress([]byte("addr2---------------")),
			newPreferences: []types.ValidatorPreference{
				{
					ValOperAddress: valAddrs[0], // validator that has existing delegation
					Weight:         sdk.NewDecWithPrec(5, 1),
				},
				{
					ValOperAddress: valAddrs[1],
					Weight:         sdk.NewDecWithPrec(3, 1),
				},
				{
					ValOperAddress: valAddrs[2],
					Weight:         sdk.NewDecWithPrec(2, 1),
				},
			},
			amountToDelegate:      sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10_000_000)),
			expectedShares:        []sdk.Dec{sdk.NewDec(5_000_000), sdk.NewDec(3_000_000), sdk.NewDec(2_000_000)},
			setExistingDelegation: true,
			expectPass:            true,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// setup message server
			msgServer := valPref.NewMsgServerImpl(s.App.ValidatorSetPreferenceKeeper)
			c := sdk.WrapSDKContext(s.Ctx)

			// fund the account that is trying to delegate
			s.FundAcc(test.delegator, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 100_000_000)})

			if test.setExistingDelegation {
				err := s.PrepareExistingDelegations(s.Ctx, []string{valAddrs[0]}, test.delegator, test.amountToDelegate.Amount)
				s.Require().NoError(err)
			}

			if test.setExistingValSetDelegation {
				// SetValidatorSetPreference sets a new list of val-set
				_, err := msgServer.SetValidatorSetPreference(c, types.NewMsgSetValidatorSetPreference(test.delegator, preferences))
				s.Require().NoError(err)

				// DelegateToValidatorSet delegate to existing val-set
				_, err = msgServer.DelegateToValidatorSet(c, types.NewMsgDelegateToValidatorSet(test.delegator, test.amountToDelegate))
				s.Require().NoError(err)
			}

			// RedelegateValidatorSet redelegates from an existing set to a new one
			_, err := msgServer.RedelegateValidatorSet(c, types.NewMsgRedelegateValidatorSet(test.delegator, test.newPreferences))
			if test.expectPass {
				s.Require().NoError(err)

				// check if the validator have received the correct amount of tokens
				for i, val := range test.newPreferences {
					valAddr, err := sdk.ValAddressFromBech32(val.ValOperAddress)
					s.Require().NoError(err)

					// guarantees that the delegator exists because we check it in DelegateToValidatorSet
					del, _ := s.App.StakingKeeper.GetDelegation(s.Ctx, test.delegator, valAddr)
					s.Require().Equal(test.expectedShares[i], del.Shares)
				}
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestWithdrawDelegationRewards() {
	s.SetupTest()

	// valset test setup
	valAddrs, preferences, amountToFund := s.SetupValidatorsAndDelegations()

	tests := []struct {
		name                  string
		delegator             sdk.AccAddress
		coinsToDelegate       sdk.Coin
		setValSetDelegation   bool
		setExistingDelegation bool
		expectPass            bool
	}{
		{
			name:                "Withdraw all rewards from existing valset delegations",
			delegator:           sdk.AccAddress([]byte("addr1---------------")),
			coinsToDelegate:     sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(20_000_000)), // delegate 20osmo
			setValSetDelegation: true,
			expectPass:          true,
		},
		{
			name:                  "Withdraw all rewards from existing staking delegations (no val-set)",
			delegator:             sdk.AccAddress([]byte("addr2---------------")),
			coinsToDelegate:       sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(20_000_000)),
			setExistingDelegation: true,
			expectPass:            true,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.FundAcc(test.delegator, amountToFund) // 100 osmo

			// setup message server
			msgServer := valPref.NewMsgServerImpl(s.App.ValidatorSetPreferenceKeeper)
			c := sdk.WrapSDKContext(s.Ctx)

			ctx := s.Ctx
			// setup test for only valset delegation
			if test.setValSetDelegation {
				// delegators have to set val-set before delegating tokens
				_, err := msgServer.SetValidatorSetPreference(c, types.NewMsgSetValidatorSetPreference(test.delegator, preferences))
				s.Require().NoError(err)

				// call the delegate to validator set preference message
				_, err = msgServer.DelegateToValidatorSet(c, types.NewMsgDelegateToValidatorSet(test.delegator, test.coinsToDelegate))
				s.Require().NoError(err)

				s.SetupDelegationReward(test.delegator, preferences, "", test.setValSetDelegation, test.setExistingDelegation)
			}

			// setup test for only existing staking position
			if test.setExistingDelegation {
				err := s.PrepareExistingDelegations(s.Ctx, valAddrs, test.delegator, test.coinsToDelegate.Amount)
				s.Require().NoError(err)

				s.SetupDelegationReward(test.delegator, nil, valAddrs[0], test.setValSetDelegation, test.setExistingDelegation)
			}

			_, err := msgServer.WithdrawDelegationRewards(c, types.NewMsgWithdrawDelegationRewards(test.delegator))
			if test.expectPass {
				s.Require().NoError(err)

				// the rewards for valset and existing delegations should be nil
				if test.setValSetDelegation {
					for _, val := range preferences {
						rewardAfterWithdrawValSet, _ := s.GetDelegationRewards(ctx, val.ValOperAddress, test.delegator)
						s.Require().Nil(rewardAfterWithdrawValSet)
					}
				}

				if test.setExistingDelegation {
					rewardAfterWithdrawExistingSet, _ := s.GetDelegationRewards(ctx, valAddrs[0], test.delegator)
					s.Require().Nil(rewardAfterWithdrawExistingSet)
				}
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestDelegateBondedTokens() {
	s.SetupTest()

	testLock := s.SetupLocks(sdk.AccAddress([]byte("addr1---------------")))

	tests := []struct {
		name                 string
		delegator            sdk.AccAddress
		lockId               uint64
		expectedUnlockedOsmo sdk.Coin
		expectedDelegations  []sdk.Dec
		setValSet            bool
		expectPass           bool
	}{
		{
			name:                 "DelegateBondedTokens with existing osmo denom lockId, bonded and <= 2 weeks bond duration",
			delegator:            sdk.AccAddress([]byte("addr1---------------")),
			lockId:               testLock[0].ID,
			expectedUnlockedOsmo: sdk.NewCoin(appParams.BaseCoinUnit, sdk.NewInt(60_000_000)), // delegator has 100osmo and creates 5 locks 10osmo each, forceUnlock only 1 lock
			expectedDelegations:  []sdk.Dec{sdk.NewDec(2_000_000), sdk.NewDec(3_300_000), sdk.NewDec(1_200_000), sdk.NewDec(3_500_000)},
			setValSet:            true,
			expectPass:           true,
		},
		{
			name:       "DelegateBondedTokens with existing stake denom lockId, bonded and <= 2 weeks bond duration",
			delegator:  sdk.AccAddress([]byte("addr1---------------")),
			lockId:     testLock[1].ID,
			expectPass: false,
		},
		{
			name:       "DelegateBondedTokens with non existing lockId",
			delegator:  sdk.AccAddress([]byte("addr1---------------")),
			lockId:     10,
			expectPass: false,
		},
		{
			name:       "DelegateBondedTokens with lockOwner != delegatorOwner",
			delegator:  sdk.AccAddress([]byte("addr1---------------")),
			lockId:     testLock[2].ID,
			expectPass: false,
		},
		{
			name:       "DelegateBondedTokens with lock duration > 2 weeks",
			delegator:  sdk.AccAddress([]byte("addr1---------------")),
			lockId:     testLock[3].ID,
			expectPass: false,
		},
		{
			name:       "DelegateBondedTokens with non bonded lockId",
			delegator:  sdk.AccAddress([]byte("addr1---------------")),
			lockId:     testLock[4].ID,
			expectPass: false,
		},
		{
			name:       "DelegateBondedTokens with synthetic locks",
			delegator:  sdk.AccAddress([]byte("addr1---------------")),
			lockId:     testLock[5].ID,
			expectPass: false,
		},
		{
			name:       "DelegateBondedTokens with multiple asset lock",
			delegator:  sdk.AccAddress([]byte("addr1---------------")),
			lockId:     testLock[6].ID,
			expectPass: false,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// setup message server
			msgServer := valPref.NewMsgServerImpl(s.App.ValidatorSetPreferenceKeeper)
			c := sdk.WrapSDKContext(s.Ctx)

			// creates a validator preference list to delegate to
			preferences := s.PrepareDelegateToValidatorSet()

			if test.setValSet {
				// SetValidatorSetPreference sets a new list of val-set
				_, err := msgServer.SetValidatorSetPreference(c, types.NewMsgSetValidatorSetPreference(test.delegator, preferences))
				s.Require().NoError(err)
			}

			_, err := msgServer.DelegateBondedTokens(c, types.NewMsgDelegateBondedTokens(test.delegator, test.lockId))
			if test.expectPass {
				s.Require().NoError(err)

				// check that the lock has been successfully unlocked
				// existingLocks should not contain the current lock
				existingLocks, err := s.App.LockupKeeper.GetPeriodLocks(s.Ctx)

				s.Require().NoError(err)
				s.Require().Equal(len(existingLocks), len(testLock)-1)

				balance := s.App.BankKeeper.GetBalance(s.Ctx, test.delegator, appParams.BaseCoinUnit)
				s.Require().Equal(test.expectedUnlockedOsmo, balance)

				// check if delegation has been done by checking if expectedDelegations matches after delegation
				for i, val := range preferences {
					valAddr, err := sdk.ValAddressFromBech32(val.ValOperAddress)
					s.Require().NoError(err)

					// guarantees that the delegator exists because we check it in DelegateToValidatorSet
					del, _ := s.App.StakingKeeper.GetDelegation(s.Ctx, test.delegator, valAddr)
					s.Require().Equal(test.expectedDelegations[i], del.Shares)
				}
			} else {
				s.Require().Error(err)
			}
		})
	}
}
