package consensus

import (
	"crypto/rsa"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	ssvcomparable "github.com/ssvlabs/ssv-spec/ssv/spectest/comparable"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// futureDecidedSyncCommitteeContributionSC returns runner with 2 stored instances.
// One undecided and one decided for height 13.
// This is because we are processing messages for height 13 while height 12 is still running.
// There are also pre consensus messages that start the new instance.
func futureDecidedSyncCommitteeContributionSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestSyncCommitteeContributionConsensusData

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.SyncCommitteeContributionRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SignedSSVMessage{
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3))),
					},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SignedSSVMessage{},
				),
				StartingDuty: &cd.Duty,
				Finished:     false,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					CommitteeMember:   testingutils.TestingCommitteeMember(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            testingutils.TestingDutySlot,
					LastPreparedRound: qbft.NoRound,
					Decided:           false,
				},
				StartValue: comparable.NoErrorEncoding(comparable.FixIssue178(cd, spec.DataVersionBellatrix)),
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SignedSSVMessage{})

			decidedInstance := &qbft.Instance{
				State: &qbft.State{
					CommitteeMember:   testingutils.TestingCommitteeMember(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            testingutils.TestingDutySlot + 1,
					LastPreparedRound: qbft.NoRound,
					Decided:           true,
					DecidedValue:      testingutils.TestingQBFTFullData,
				},
			}
			comparable.SetMessages(
				decidedInstance,
				[]*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
						[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
						[]types.OperatorID{1, 2, 3},
						qbft.Height(testingutils.TestingDutySlot+1),
						ret.GetBaseRunner().QBFTController.Identifier,
					),
				},
			)

			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, decidedInstance)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			ret.GetBaseRunner().QBFTController.Height = testingutils.TestingDutySlot + 1

			return ret
		}(),
	}
}

// futureDecidedProposerSC returns runner with 2 stored instances. One undecided and one decided for height 13.
// This is because we are processing messages for height 13 while height 12 is still running.
// There are also pre consensus messages that start the new instance.
func futureDecidedProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestProposerConsensusDataV(version)
	cdBytes := testingutils.TestProposerConsensusDataBytsV(version)

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.ProposerRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SignedSSVMessage{
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[2], 2, version))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[3], 3, version))),
					},
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SignedSSVMessage{},
				),
				StartingDuty: &cd.Duty,
				Finished:     false,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					CommitteeMember:   testingutils.TestingCommitteeMember(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            qbft.Height(testingutils.TestingDutySlotV(version)),
					LastPreparedRound: qbft.NoRound,
					Decided:           false,
				},
				StartValue: cdBytes,
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SignedSSVMessage{})

			decidedInstance := &qbft.Instance{
				State: &qbft.State{
					CommitteeMember:   testingutils.TestingCommitteeMember(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            qbft.Height(testingutils.TestingDutySlotV(version) + 1),
					LastPreparedRound: qbft.NoRound,
					Decided:           true,
					DecidedValue:      testingutils.TestingQBFTFullData,
				},
			}
			comparable.SetMessages(
				decidedInstance,
				[]*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
						[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
						[]types.OperatorID{1, 2, 3},
						qbft.Height(testingutils.TestingDutySlotV(version)+1),
						ret.GetBaseRunner().QBFTController.Identifier,
					),
				},
			)

			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, decidedInstance)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			ret.GetBaseRunner().QBFTController.Height = qbft.Height(testingutils.TestingDutySlotV(version) + 1)

			return ret
		}(),
	}
}

// futureDecidedBlindedProposerSC returns runner with 2 stored instances. One undecided and one decided for height 13.
// This is because we are processing messages for height 13 while height 12 is still running.
// There are also pre consensus messages that start the new instance.
func futureDecidedBlindedProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestProposerBlindedBlockConsensusDataV(version)
	cdBytes := testingutils.TestProposerBlindedBlockConsensusDataBytsV(version)

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.ProposerBlindedBlockRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SignedSSVMessage{
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[2], 2, version))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgV(ks.Shares[3], 3, version))),
					}),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SignedSSVMessage{},
				),
				StartingDuty: &cd.Duty,
				Finished:     false,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					CommitteeMember:   testingutils.TestingCommitteeMember(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            qbft.Height(testingutils.TestingDutySlotV(version)),
					LastPreparedRound: qbft.NoRound,
					Decided:           false,
				},
				StartValue: cdBytes,
			}
			comparable.SetMessages(ret.GetBaseRunner().State.RunningInstance, []*types.SignedSSVMessage{})

			decidedInstance := &qbft.Instance{
				State: &qbft.State{
					CommitteeMember:   testingutils.TestingCommitteeMember(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            qbft.Height(testingutils.TestingDutySlotV(version) + 1),
					LastPreparedRound: qbft.NoRound,
					Decided:           true,
					DecidedValue:      testingutils.TestingQBFTFullData,
				},
			}
			comparable.SetMessages(
				decidedInstance,
				[]*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeightAndIdentifier(
						[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
						[]types.OperatorID{1, 2, 3},
						qbft.Height(testingutils.TestingDutySlotV(version)+1),
						ret.GetBaseRunner().QBFTController.Identifier,
					),
				},
			)

			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, decidedInstance)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			ret.GetBaseRunner().QBFTController.Height = qbft.Height(testingutils.TestingDutySlotV(version) + 1)

			return ret
		}(),
	}
}
