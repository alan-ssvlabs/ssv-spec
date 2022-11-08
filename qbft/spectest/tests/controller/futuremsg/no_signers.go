package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoSigners tests future msg with no signers
func NoSigners() *ControllerSyncSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	inputData := &qbft.Data{
		Root:   testingutils.TestAttesterConsensusDataRoot,
		Source: testingutils.TestAttesterConsensusDataByts,
	}
	msg := testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
		Height: 10,
		Round:  3,
	}, &qbft.Data{Root: inputData.Root})
	msg.Signers = []types.OperatorID{}
	msgEncoded, _ := msg.Encode()

	return &ControllerSyncSpecTest{
		Name: "future msgs no signer",
		InputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(identifier, types.ConsensusPrepareMsgType),
				Data: msgEncoded,
			},
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "5a1536414abb7928a962cc82e7307b48e3d6c17da15c3f09948c20bd89d41301",
		ExpectedError:        "invalid future msg: invalid future msg: message signers is empty",
	}
}