package types

import (
	"fmt"
	"strings"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeUploadCosmWasmPoolCodeAndWhiteList = "UploadCosmWasmPoolCodeAndWhiteListProposal"
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeUploadCosmWasmPoolCodeAndWhiteList)
	govtypes.RegisterProposalTypeCodec(&UploadCosmWasmPoolCodeAndWhiteListProposal{}, "osmosis/UploadCosmWasmPoolCodeAndWhiteListProposal")
}

var (
	_ govtypes.Content = &UploadCosmWasmPoolCodeAndWhiteListProposal{}
)

// NewUploadCosmWasmPoolCodeAndWhiteListProposal returns a new instance of an upload cosmwasm pool code and whitelist proposal struct.
func NewUploadCosmWasmPoolCodeAndWhiteListProposal(title, description string, wasmByteCode []byte) govtypes.Content {
	return &UploadCosmWasmPoolCodeAndWhiteListProposal{
		Title:        title,
		Description:  description,
		WASMByteCode: wasmByteCode,
	}
}

func (p *UploadCosmWasmPoolCodeAndWhiteListProposal) GetTitle() string { return p.Title }

// GetDescription gets the description of the proposal
func (p *UploadCosmWasmPoolCodeAndWhiteListProposal) GetDescription() string { return p.Description }

// ProposalRoute returns the router key for the proposal
func (p *UploadCosmWasmPoolCodeAndWhiteListProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of the proposal
func (p *UploadCosmWasmPoolCodeAndWhiteListProposal) ProposalType() string {
	return ProposalTypeUploadCosmWasmPoolCodeAndWhiteList
}

// ValidateBasic validates a governance proposal's abstract and basic contents
func (p *UploadCosmWasmPoolCodeAndWhiteListProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if len(p.WASMByteCode) == 0 {
		return fmt.Errorf("wasm byte code cannot be nil")
	}

	return nil
}

// String returns a string containing the pool incentives proposal.
func (p UploadCosmWasmPoolCodeAndWhiteListProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Upload CosmWasm Pool Code and WhiteList Proposal:
Title:       %s
Description: %s
`, p.Title, p.Description))
	return b.String()
}
