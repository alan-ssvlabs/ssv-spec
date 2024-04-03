package messages

import (
	"crypto/rsa"

	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SignedMsgDuplicateSigners tests SignedMessage with duplicate signers
func SignedMsgDuplicateSigners() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingCommitMultiSignerMessage(
		[]*rsa.PrivateKey{ks.NetworkKeys[1], ks.NetworkKeys[1], ks.NetworkKeys[2]},
		[]types.OperatorID{1, 2, 3},
	)
	msg.OperatorID = []types.OperatorID{1, 1, 2}

	return &tests.MsgSpecTest{
		Name: "duplicate signers",
		Messages: []*types.SignedSSVMessage{
			msg,
		},
		ExpectedError: "non unique signer",
	}
}
