package latemsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// LateProposal tests process late proposal msg for an instance which just decided
func LateProposal() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()

	inputData := &qbft.Data{
		Root:   testingutils.TestAttesterConsensusDataRoot,
		Source: testingutils.TestAttesterConsensusDataByts,
	}
	msgs := testingutils.DecidingMsgsForHeight(inputData, identifier, qbft.FirstHeight, ks)

	signedMsgEncoded, _ := testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, inputData).Encode()
	lateMsg := &types.Message{
		ID:   types.PopulateMsgType(identifier, types.ConsensusProposeMsgType),
		Data: signedMsgEncoded,
	}
	multiSignMsg := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
		[]types.OperatorID{1, 2, 3},
		&qbft.Message{
			Height: qbft.FirstHeight,
			Round:  qbft.FirstRound,
		}, inputData)

	msgs = append(msgs, lateMsg)

	return &tests.ControllerSpecTest{
		Name: "late proposal",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         inputData,
				InputMessages:      msgs,
				DecidedVal:         inputData.Source,
				DecidedCnt:         1,
				SavedDecided:       multiSignMsg,
				BroadcastedDecided: multiSignMsg,
				ControllerPostRoot: "df9b2787df60e1e15b0c840410592d27803d44cd5fb086cfa8fc23181cea6293",
			},
		},
		ExpectedError: "could not process msg: proposal invalid: proposal is not valid with current state",
	}
}