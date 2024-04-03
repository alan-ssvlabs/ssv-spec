package share

import (
	"crypto/rsa"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PartialQuorumWithDuplicate tests msg with unique f+1 signers (but also including duplicates)
func PartialQuorumWithDuplicate() *ShareTest {
	ks := testingutils.Testing4SharesSet()
	share := testingutils.TestingShare(ks)

	msg := testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.NetworkKeys[1], ks.NetworkKeys[1], ks.NetworkKeys[2]}, []types.OperatorID{1, 1, 2})

	return &ShareTest{
		Name:                     "partial quorum with duplicate",
		Share:                    *share,
		Message:                  *msg,
		ExpectedHasPartialQuorum: true,
		ExpectedHasQuorum:        false,
		ExpectedFullCommittee:    false,
		ExpectedError:            "non unique signer",
	}
}
