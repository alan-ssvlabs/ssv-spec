package partialsigcontainer

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func Duplicate() tests.SpecTest {

	// Create a test key set
	ks := testingutils.Testing4SharesSet()

	// Create PartialSignatureMessage for testing
	msg1 := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, spec.DataVersionPhase0)
	msg2 := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, spec.DataVersionPhase0)
	msg3 := testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, spec.DataVersionPhase0)
	msgs := []*types.PartialSignatureMessage{msg1.Messages[0], msg2.Messages[0], msg3.Messages[0]}

	return &PartialSigContainerTest{
		Name:            "duplicate",
		Quorum:          ks.Threshold,
		ValidatorPubKey: ks.ValidatorPK.Serialize(),
		SignatureMsgs:   msgs,
		ExpectedError:   "could not reconstruct a valid signature",
		ExpectedQuorum:  false,
	}
}
