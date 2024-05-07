package ssv

import (
	"fmt"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

type Committee struct {
	Runners                 map[spec.Slot]*CommitteeRunner
	Operator                types.Operator
	SignatureVerifier       types.SignatureVerifier
	CreateRunnerFn          func() *CommitteeRunner
	HighestAttestingSlotMap map[types.ValidatorPK]spec.Slot
}

// NewCommittee creates a new cluster
func NewCommittee(
	operator types.Operator,
	verifier types.SignatureVerifier,
	createRunnerFn func() *CommitteeRunner,
) *Committee {
	return &Committee{
		Runners:                 make(map[spec.Slot]*CommitteeRunner),
		Operator:                operator,
		SignatureVerifier:       verifier,
		CreateRunnerFn:          createRunnerFn,
		HighestAttestingSlotMap: make(map[types.ValidatorPK]spec.Slot),
	}

}

// StartDuty starts a new duty for the given slot
func (c *Committee) StartDuty(duty *types.CommitteeDuty) error {
	if _, exists := c.Runners[duty.Slot]; exists {
		return errors.New(fmt.Sprintf("CommitteeRunner for slot %d already exists", duty.Slot))
	}
	c.Runners[duty.Slot] = c.CreateRunnerFn()
	var validatorToStopMap map[spec.Slot]types.ValidatorPK
	// Filter old duties based on highest attesting slot
	duty, validatorToStopMap, c.HighestAttestingSlotMap = FilterCommitteeDuty(duty, c.HighestAttestingSlotMap)
	// Stop validators with old duties
	c.stopDuties(validatorToStopMap)
	c.updateAttestingSlotMap(duty)
	return c.Runners[duty.Slot].StartNewDuty(duty)
}

func (c *Committee) stopDuties(validatorToStopMap map[spec.Slot]types.ValidatorPK) {
	for slot, validator := range validatorToStopMap {
		runner, exists := c.Runners[slot]
		if exists {
			runner.StopDuty(validator)
		}
	}
}

// FilterCommitteeDuty filters the committee duties by the slots given per validator.
// It returns the filtered duties, the validators to stop and updated slot map.
func FilterCommitteeDuty(duty *types.CommitteeDuty, slotMap map[types.ValidatorPK]spec.Slot) (
	*types.CommitteeDuty,
	map[spec.Slot]types.ValidatorPK,
	map[types.ValidatorPK]spec.Slot,
) {
	validatorsToStop := make(map[spec.Slot]types.ValidatorPK)

	for i, beaconDuty := range duty.BeaconDuties {
		validatorPK := types.ValidatorPK(beaconDuty.PubKey)
		slot, exists := slotMap[validatorPK]
		if exists {
			if slot < beaconDuty.Slot {
				validatorsToStop[beaconDuty.Slot] = validatorPK
				slotMap[validatorPK] = beaconDuty.Slot
			} else { // else don't run duty with old slot
				duty.BeaconDuties[i] = nil
			}
		}
	}
	return duty, validatorsToStop, slotMap
}

// ProcessMessage processes Network Message of all types
func (c *Committee) ProcessMessage(signedSSVMessage *types.SignedSSVMessage) error {
	// Validate message
	if err := signedSSVMessage.Validate(); err != nil {
		return errors.Wrap(err, "invalid SignedSSVMessage")
	}

	// Verify SignedSSVMessage's signature
	if err := c.SignatureVerifier.Verify(signedSSVMessage, c.Operator.Committee); err != nil {
		return errors.Wrap(err, "SignedSSVMessage has an invalid signature")
	}

	msg := signedSSVMessage.SSVMessage

	switch msg.GetType() {
	case types.SSVConsensusMsgType:
		qbftMsg := &qbft.Message{}
		if err := qbftMsg.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get consensus Message from network Message")
		}

		if err := qbftMsg.Validate(); err != nil {
			return errors.Wrap(err, "invalid qbft Message")
		}

		runner, exists := c.Runners[spec.Slot(qbftMsg.Height)]
		if !exists {
			return errors.New("no runner found for message's slot")
		}
		return runner.ProcessConsensus(signedSSVMessage)
	case types.SSVPartialSignatureMsgType:
		pSigMessages := &types.PartialSignatureMessages{}
		if err := pSigMessages.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get post consensus Message from network Message")
		}

		// Validate
		if len(signedSSVMessage.OperatorIDs) != 1 {
			return errors.New("PartialSignatureMessage has more than 1 signer")
		}

		if err := pSigMessages.ValidateForSigner(signedSSVMessage.OperatorIDs[0]); err != nil {
			return errors.Wrap(err, "invalid PartialSignatureMessages")
		}

		if pSigMessages.Type == types.PostConsensusPartialSig {
			runner, exists := c.Runners[pSigMessages.Slot]
			if !exists {
				return errors.New("no runner found for message's slot")
			}
			return runner.ProcessPostConsensus(pSigMessages)
		}
	default:
		return errors.New("unknown msg")
	}
	return nil

}

// updateAttestingSlotMap updates the highest attesting slot map from beacon duties
func (c *Committee) updateAttestingSlotMap(duty *types.CommitteeDuty) {
	for _, beaconDuty := range duty.BeaconDuties {
		if beaconDuty.Type == types.BNRoleAttester {
			validatorPK := types.ValidatorPK(beaconDuty.PubKey)
			if _, ok := c.HighestAttestingSlotMap[validatorPK]; !ok {
				c.HighestAttestingSlotMap[validatorPK] = beaconDuty.Slot
			}
			if c.HighestAttestingSlotMap[validatorPK] < beaconDuty.Slot {
				c.HighestAttestingSlotMap[validatorPK] = beaconDuty.Slot
			}
		}
	}
}
