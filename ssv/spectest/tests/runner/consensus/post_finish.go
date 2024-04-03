package consensus

import (
	"crypto/rsa"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostFinish tests a valid commit msg after runner finished
func PostFinish() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid post finish",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestSyncCommitteeContributionConsensusData, ks, types.BNRoleSyncCommitteeContribution),
					// post consensus
					testingutils.SSVMsgSyncCommitteeContribution(1, ks.NetworkKeys[1], nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgSyncCommitteeContribution(2, ks.NetworkKeys[2], nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgSyncCommitteeContribution(3, ks.NetworkKeys[3], nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[3], 3)),
					// commit msg
					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
						[]*rsa.PrivateKey{ks.NetworkKeys[4]},
						[]types.OperatorID{4},
						qbft.Height(testingutils.TestingDutySlot),
						testingutils.SyncCommitteeContributionMsgID,
						testingutils.TestSyncCommitteeContributionConsensusDataByts,
					),
				),
				PostDutyRunnerStateRoot: postFinishSyncCommitteeContributionSC().Root(),
				PostDutyRunnerState:     postFinishSyncCommitteeContributionSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[0], testingutils.TestingContributionProofsSigned[0], ks)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[1], testingutils.TestingContributionProofsSigned[1], ks)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[2], testingutils.TestingContributionProofsSigned[2], ks)),
				},
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestSyncCommitteeConsensusData, ks, types.BNRoleSyncCommittee),
					// post consensus
					testingutils.SSVMsgSyncCommittee(1, ks.NetworkKeys[1], nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgSyncCommittee(2, ks.NetworkKeys[2], nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgSyncCommittee(3, ks.NetworkKeys[3], nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[3], 3)),
					// commit msg
					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
						[]*rsa.PrivateKey{ks.NetworkKeys[4]},
						[]types.OperatorID{4},
						qbft.Height(testingutils.TestingDutySlot),
						testingutils.SyncCommitteeMsgID,
						testingutils.TestSyncCommitteeConsensusDataByts,
					),
				),
				PostDutyRunnerStateRoot: postFinishSyncCommitteeSC().Root(),
				PostDutyRunnerState:     postFinishSyncCommitteeSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeBlockRoot(ks)),
				},
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestAggregatorConsensusData, ks, types.BNRoleAggregator),
					// post consensus
					testingutils.SSVMsgAggregator(1, ks.NetworkKeys[1], nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgAggregator(2, ks.NetworkKeys[2], nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgAggregator(3, ks.NetworkKeys[3], nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[3], 3)),
					// commit msg
					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
						[]*rsa.PrivateKey{ks.NetworkKeys[4]},
						[]types.OperatorID{4},
						qbft.Height(testingutils.TestingDutySlot),
						testingutils.AggregatorMsgID,
						testingutils.TestAggregatorConsensusDataByts,
					),
				),
				PostDutyRunnerStateRoot: postFinishAggregatorSC().Root(),
				PostDutyRunnerState:     postFinishAggregatorSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], 1),
					testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedAggregateAndProof(ks)),
				},
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: append(
					testingutils.SSVDecidingMsgsV(testingutils.TestAttesterConsensusData, ks, types.BNRoleAttester),
					// post consensus
					testingutils.SSVMsgAttester(1, ks.NetworkKeys[1], nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, testingutils.TestingDutySlot)),
					testingutils.SSVMsgAttester(2, ks.NetworkKeys[2], nil, testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, testingutils.TestingDutySlot)),
					testingutils.SSVMsgAttester(3, ks.NetworkKeys[3], nil, testingutils.PostConsensusAttestationMsg(ks.Shares[3], 3, testingutils.TestingDutySlot)),
					// commit msg
					testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
						[]*rsa.PrivateKey{ks.NetworkKeys[4]},
						[]types.OperatorID{4},
						qbft.Height(testingutils.TestingDutySlot),
						testingutils.AttesterMsgID,
						testingutils.TestAttesterConsensusDataByts,
					),
				),
				PostDutyRunnerStateRoot: postFinishAttesterSC().Root(),
				PostDutyRunnerState:     postFinishAttesterSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, testingutils.TestingDutySlot),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedAttestation(ks)),
				},
			},
		},
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("proposer (%s)", version.String()),
			Runner: testingutils.ProposerRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: append(
				testingutils.SSVDecidingMsgsV(
					testingutils.TestProposerConsensusDataV(version),
					ks,
					types.BNRoleProposer,
				), // consensus
				// post consensus
				testingutils.SSVMsgProposer(1, ks.NetworkKeys[1], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version)),
				testingutils.SSVMsgProposer(2, ks.NetworkKeys[2], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version)),
				testingutils.SSVMsgProposer(3, ks.NetworkKeys[3], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version)),
				// commit msg
				testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
					[]*rsa.PrivateKey{ks.NetworkKeys[4]},
					[]types.OperatorID{4},
					qbft.Height(testingutils.TestingDutySlotV(version)),
					testingutils.ProposerMsgID,
					testingutils.TestProposerConsensusDataBytsV(version),
				),
			),
			PostDutyRunnerStateRoot: postFinishProposerSC(version).Root(),
			PostDutyRunnerState:     postFinishProposerSC(version).ExpectedState,
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
				testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version),
			},
			BeaconBroadcastedRoots: []string{
				testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, version)),
			},
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("proposer blinded block (%s)", version.String()),
			Runner: testingutils.ProposerBlindedBlockRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: append(
				testingutils.SSVDecidingMsgsV(testingutils.TestProposerBlindedBlockConsensusDataV(version), ks, types.BNRoleProposer), // consensus
				// post consensus
				testingutils.SSVMsgProposer(1, ks.NetworkKeys[1], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version)),
				testingutils.SSVMsgProposer(2, ks.NetworkKeys[2], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version)),
				testingutils.SSVMsgProposer(3, ks.NetworkKeys[3], nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version)),
				// commit msg
				testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
					[]*rsa.PrivateKey{ks.NetworkKeys[4]},
					[]types.OperatorID{4},
					qbft.Height(testingutils.TestingDutySlotV(version)),
					testingutils.ProposerMsgID,
					testingutils.TestProposerBlindedBlockConsensusDataBytsV(version),
				),
			),
			PostDutyRunnerStateRoot: postFinishBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     postFinishBlindedProposerSC(version).ExpectedState,
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
				testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version),
			},
			BeaconBroadcastedRoots: []string{
				testingutils.GetSSZRootNoError(testingutils.TestingSignedBlindedBeaconBlockV(ks, version)),
			},
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
