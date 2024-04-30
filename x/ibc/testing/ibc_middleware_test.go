package testing

import (
	"encoding/json"

	"github.com/evoluteai-network/evoluteai-chain/x/ibc/gmp"
)

func (s *IBCTestSuite) TestGMPMessageFrom_Success() {
	s.IBCTransferProviderToevoluteai(s.providerAddr, s.evoluteaiAddr, nativeDenom, ibcTransferAmount, "")

	generalMsg := gmp.Message{
		SourceChain:   "axelar",
		SourceAddress: "evoluteai",
		Payload:       []byte("Hello evoluteai, I am Axelar"),
		Type:          gmp.TypeGeneralMessage,
	}
	generalMsgJson, _ := json.Marshal(generalMsg)
	s.IBCTransferProviderToevoluteai(s.providerAddr, s.evoluteaiAddr, nativeDenom, ibcTransferAmount, string(generalMsgJson))

	generalMsgWithToken := gmp.Message{
		SourceChain:   "axelar",
		SourceAddress: "evoluteai",
		Payload:       []byte("Hello evoluteai, I am Axelar"),
		Type:          gmp.TypeGeneralMessageWithToken,
	}
	generalMsgWithTokenJson, _ := json.Marshal(generalMsgWithToken)
	s.IBCTransferProviderToevoluteai(s.providerAddr, s.evoluteaiAddr, nativeDenom, ibcTransferAmount, string(generalMsgWithTokenJson))
}

func (s *IBCTestSuite) TestGMPMessageTo_Success() {
	generalMsg := gmp.Message{
		SourceChain:   "evoluteai",
		SourceAddress: "axelar",
		Payload:       []byte("Hello Axelar, I am evoluteai"),
		Type:          gmp.TypeGeneralMessage,
	}
	generalMsgJson, _ := json.Marshal(generalMsg)
	s.IBCTransferevoluteaiToProvider(s.evoluteaiAddr, s.providerAddr, nativeDenom, ibcTransferAmount, string(generalMsgJson))
}
