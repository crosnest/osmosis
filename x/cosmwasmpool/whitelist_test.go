package cosmwasmpool_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/osmosis-labs/osmosis/v15/app/apptesting"
)

type WhitelistSuite struct {
	apptesting.KeeperTestHelper
}

func TestWhitelistSuite(t *testing.T) {
	suite.Run(t, new(WhitelistSuite))
}

func (s *WhitelistSuite) TestWhitelist() {
	s.Setup()

	const (
		defaultPoolId uint64 = 5
	)

	// Check that the pool is not whitelisted
	s.Require().False(s.App.CosmwasmPoolKeeper.IsWhitelisted(s.Ctx, defaultPoolId))

	// Add the pool to the whitelist and check that it is now whitelisted.
	s.App.CosmwasmPoolKeeper.AddToWhitelist(s.Ctx, defaultPoolId)
	s.Require().True(s.App.CosmwasmPoolKeeper.IsWhitelisted(s.Ctx, defaultPoolId))

	// Add another pool to the whitelist and check that it is now whitelisted.
	s.App.CosmwasmPoolKeeper.AddToWhitelist(s.Ctx, defaultPoolId+1)
	s.Require().True(s.App.CosmwasmPoolKeeper.IsWhitelisted(s.Ctx, defaultPoolId+1))

	// Remove the first pool from the whitelist and check that it is no longer whitelisted
	// while the other pool still is.
	s.App.CosmwasmPoolKeeper.RemoveFromWhitelist(s.Ctx, defaultPoolId)
	s.Require().False(s.App.CosmwasmPoolKeeper.IsWhitelisted(s.Ctx, defaultPoolId))
	s.Require().True(s.App.CosmwasmPoolKeeper.IsWhitelisted(s.Ctx, defaultPoolId+1))
}
