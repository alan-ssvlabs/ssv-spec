package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// MultiDecidedInstances tests deciding multiple instances
func MultiDecidedInstances() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	inputData := &qbft.Data{
		Root:   testingutils.TestAttesterConsensusDataRoot,
		Source: testingutils.TestAttesterConsensusDataByts,
	}
	instanceData := func(height qbft.Height, postRoot string) *tests.RunInstanceData {
		multiSignMsg := testingutils.MultiSignQBFTMsg(
			[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2], testingutils.Testing4SharesSet().Shares[3]},
			[]types.OperatorID{1, 2, 3},
			&qbft.Message{
				Height: height,
				Round:  qbft.FirstRound,
			}, inputData)
		multiSignMsgEncoded, _ := multiSignMsg.Encode()
		return &tests.RunInstanceData{
			InputValue: inputData,
			InputMessages: []*types.Message{
				{
					ID:   types.PopulateMsgType(identifier, types.DecidedMsgType),
					Data: multiSignMsgEncoded,
				},
			},
			SavedDecided:       multiSignMsg,
			DecidedVal:         inputData.Source,
			DecidedCnt:         1,
			ControllerPostRoot: postRoot,
		}
	}

	return &tests.ControllerSpecTest{
		Name: "multi decide instances",
		RunInstanceData: []*tests.RunInstanceData{
			instanceData(qbft.FirstHeight, "70c5218e3832249ba51590e5e8850d60863e3a2281e090669448075e68795a25"),
			instanceData(1, "7311ad2d2b9d480bb17e49a281192a4f298f21cfc03d05cae776f1e8cefc3fea"),
			instanceData(2, "78a6e414620e88f63ad81f6a3834087620345a7236d052352fb11a11bb00f29e"),
			instanceData(3, "4bb484d62ad78d5c310b4473583bcb1ef450198a19e0efa3e61e6d692bd1ec8e"),
			instanceData(4, "85fb859090ab4b932754654cf9ad17ad281aa2d62833bba31cd1120e2ee9a98e"),
			instanceData(5, "3d42cfabe35b3d229e114b5eb3e250cfb9214e97afe3034f7041f67230a57a00"),
			instanceData(8, "505dcff70f958af26c14c289c58dbe8678db8af6eb9c8fb59a2e19a63190c531"),
			instanceData(9, "18c3869993b7bb7af00ec0346a370fc8062b763748ee6d445e79d0781c2c0cef"),
			instanceData(10, "fe635fafdb3bcd91bd5df56c7aea5e7809d36967a422468bf2585eadaee813bd"),
		},
	}
}