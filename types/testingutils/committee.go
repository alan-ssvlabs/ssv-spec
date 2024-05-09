package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
)

var BaseCommittee = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet) *ssv.Committee {

	var keySetSample *TestKeySet
	for _, ks := range keySetMap {
		keySetSample = ks
		break
	}

	shareMap := make(map[phase0.ValidatorIndex]*types.Share)
	for valIdx, ks := range keySetMap {
		shareMap[valIdx] = TestingShare(ks, valIdx)
	}

	createRunnerF := func(shareMap map[phase0.ValidatorIndex]*types.Share) *ssv.CommitteeRunner {
		return CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)
	}
	return ssv.NewCommittee(
		*TestingOperator(keySetSample),
		NewTestingVerifier(),
		shareMap,
		createRunnerF,
	)
}

var BaseCommitteeWithRunnerSample = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, runnerSample *ssv.CommitteeRunner) *ssv.Committee {

	var keySetSample *TestKeySet
	for _, ks := range keySetMap {
		keySetSample = ks
		break
	}

	shareMap := make(map[phase0.ValidatorIndex]*types.Share)
	for valIdx, ks := range keySetMap {
		shareMap[valIdx] = TestingShare(ks, valIdx)
	}

	createRunnerF := func(shareMap map[phase0.ValidatorIndex]*types.Share) *ssv.CommitteeRunner {
		return ssv.NewCommitteeRunner(runnerSample.BaseRunner.BeaconNetwork,
			shareMap,
			qbft.NewController(runnerSample.BaseRunner.QBFTController.Identifier,
				runnerSample.BaseRunner.QBFTController.Share,
				runnerSample.BaseRunner.QBFTController.GetConfig()),
			runnerSample.GetBeaconNode(),
			runnerSample.GetNetwork(),
			runnerSample.GetSigner(),
			runnerSample.GetOperatorSigner(),
			runnerSample.GetValCheckF(),
		).(*ssv.CommitteeRunner)
	}

	return ssv.NewCommittee(
		*TestingOperator(keySetSample),
		NewTestingVerifier(),
		shareMap,
		createRunnerF,
	)
}