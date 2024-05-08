package attestationsynccommittee

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// StartMixedDuties starts a cluster runner with 30 attestation and 30 sync committee duties
func StartMixedDuties() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "start mixed duties",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:           "30 attestations 30 sync committees",
				Runner:         testingutils.CommitteeRunner(ks),
				Duty:           testingutils.TestingCommitteeDuty(testingutils.TestingDutySlot, validatorIndexList(30), validatorIndexList(30)),
				Messages:       []*types.SignedSSVMessage{},
				OutputMessages: []*types.PartialSignatureMessages{},
			},
		},
	}

	return multiSpecTest
}