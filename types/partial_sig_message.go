package types

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type PartialSigMsgType uint64

const (
	// PostConsensusPartialSig is a partial signature over a decided duty (attestation data, block, etc)
	PostConsensusPartialSig PartialSigMsgType = iota
	// RandaoPartialSig is a partial signature over randao reveal
	RandaoPartialSig
	// SelectionProofPartialSig is a partial signature for aggregator selection proof
	SelectionProofPartialSig
	// ContributionProofs is the partial selection proofs for sync committee contributions (it's an array of sigs)
	ContributionProofs
	// ValidatorRegistrationPartialSig is a partial signature over a ValidatorRegistration object
	ValidatorRegistrationPartialSig
	// VoluntaryExitPartialSig is a partial signature over a VoluntaryExit object
	VoluntaryExitPartialSig
)

type PartialSignatureMessages struct {
	Type     PartialSigMsgType
	Slot     phase0.Slot
	Messages []*PartialSignatureMessage `ssz-max:"13"`
}

// Encode returns a msg encoded bytes or error
func (msgs *PartialSignatureMessages) Encode() ([]byte, error) {
	return msgs.MarshalSSZ()
}

// Decode returns error if decoding failed
func (msgs *PartialSignatureMessages) Decode(data []byte) error {
	return msgs.UnmarshalSSZ(data)
}

// GetRoot returns the root used for signing and verification
func (msgs PartialSignatureMessages) GetRoot() ([32]byte, error) {
	return msgs.HashTreeRoot()
}

func (msgs PartialSignatureMessages) Validate() error {
	if len(msgs.Messages) == 0 {
		return errors.New("no PartialSignatureMessages messages")
	}

	signer := msgs.Messages[0].Signer
	for _, msg := range msgs.Messages {
		if signer != msg.Signer {
			return errors.New("inconsistent signers")
		}
	}

	for _, m := range msgs.Messages {
		if err := m.Validate(); err != nil {
			return errors.Wrap(err, "message invalid")
		}
	}
	return nil
}

func (msgs PartialSignatureMessages) GetSigner() (OperatorID, error) {
	if len(msgs.Messages) == 0 {
		return OperatorID(0), errors.New("Can not get signer due to no PartialSignatureMessages messages")
	}
	return msgs.Messages[0].Signer, nil
}

// PartialSignatureMessage is a msg for partial Beacon chain related signatures (like partial attestation, block, randao sigs)
type PartialSignatureMessage struct {
	PartialSignature Signature `ssz-size:"96"` // The Beacon chain partial Signature for a duty
	SigningRoot      [32]byte  `ssz-size:"32"` // the root signed in PartialSignature
	Signer           OperatorID
}

// Encode returns a msg encoded bytes or error
func (pcsm *PartialSignatureMessage) Encode() ([]byte, error) {
	return pcsm.MarshalSSZ()
}

// Decode returns error if decoding failed
func (pcsm *PartialSignatureMessage) Decode(data []byte) error {
	return pcsm.UnmarshalSSZ(data)
}

func (pcsm *PartialSignatureMessage) GetRoot() ([32]byte, error) {
	return pcsm.HashTreeRoot()
}

func (pcsm *PartialSignatureMessage) Validate() error {
	if pcsm.Signer == 0 {
		return errors.New("signer ID 0 not allowed")
	}
	return nil
}
