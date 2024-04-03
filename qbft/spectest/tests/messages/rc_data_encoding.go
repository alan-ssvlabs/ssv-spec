package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// RoundChangeDataEncoding tests encoding RoundChangeData
func RoundChangeDataEncoding() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingRoundChangeMessageWithParams(
		ks.NetworkKeys[1], types.OperatorID(1), qbft.FirstRound, qbft.FirstHeight, testingutils.TestingQBFTRootData, 2,
		testingutils.MarshalJustifications([]*types.SignedSSVMessage{
			testingutils.TestingPrepareMessageWithRound(ks.NetworkKeys[1], types.OperatorID(1), 2),
			testingutils.TestingPrepareMessageWithRound(ks.NetworkKeys[2], types.OperatorID(2), 2),
			testingutils.TestingPrepareMessageWithRound(ks.NetworkKeys[3], types.OperatorID(3), 2),
		}))

	r, _ := msg.GetRoot()
	b, _ := msg.Encode()

	return &tests.MsgSpecTest{
		Name: "round change data encoding",
		Messages: []*types.SignedSSVMessage{
			msg,
		},
		EncodedMessages: [][]byte{
			b,
		},
		ExpectedRoots: [][32]byte{
			r,
		},
	}
}
