package tests

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/ssvlabs/ssv-spec/qbft"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/google/go-cmp/cmp"
	typescomparable "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

type MsgProcessingSpecTest struct {
	Name     string
	Runner   ssv.Runner
	Duty     types.Duty
	Messages []*types.SignedSSVMessage
	// DecidedSlashable makes the decided value slashable. Simulates consensus instances running in parallel.
	DecidedSlashable        bool
	PostDutyRunnerStateRoot string
	PostDutyRunnerState     types.Root `json:"-"` // Field is ignored by encoding/json
	// OutputMessages compares pre/ post signed partial sigs to output. We exclude consensus msgs as it's tested in consensus
	OutputMessages         []*types.SignedPartialSignatureMessage
	BeaconBroadcastedRoots []string
	DontStartDuty          bool // if set to true will not start a duty for the runner
	ExpectedError          string
}

func (test *MsgProcessingSpecTest) TestName() string {
	return test.Name
}

// RunAsPartOfMultiTest runs the test as part of a MultiMsgProcessingSpecTest
func (test *MsgProcessingSpecTest) RunAsPartOfMultiTest(t *testing.T) {
	v, c, lastErr := test.runPreTesting()

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}

	network := &testingutils.TestingNetwork{}
	beaconNetwork := &testingutils.TestingBeaconNode{}
	var committee []*types.Operator
	switch test.Runner.(type) {
	case *ssv.CommitteeRunner:
		var runnerInstance *ssv.CommitteeRunner
		for _, runner := range c.Runners {
			runnerInstance = runner
			break
		}
		network = runnerInstance.GetNetwork().(*testingutils.TestingNetwork)
		beaconNetwork = runnerInstance.GetBeaconNode().(*testingutils.TestingBeaconNode)
		committee = c.CommitteeMember.Committee
	default:
		network = v.Network.(*testingutils.TestingNetwork)
		committee = v.CommitteeMember.Committee
		beaconNetwork = test.Runner.GetBeaconNode().(*testingutils.TestingBeaconNode)
	}

	// test output message
	testingutils.ComparePartialSignatureOutputMessages(t, test.OutputMessages, network.BroadcastedMsgs, committee)

	// test beacon broadcasted msgs
	testingutils.CompareBroadcastedBeaconMsgs(t, test.BeaconBroadcastedRoots, beaconNetwork.BroadcastedRoots)

	// post root
	postRoot, err := test.Runner.GetRoot()
	require.NoError(t, err)

	if test.PostDutyRunnerStateRoot != hex.EncodeToString(postRoot[:]) {
		t.Errorf("post runner state not equal: %v", cmp.Diff(test.Runner, test.PostDutyRunnerState, cmp.Exporter(func(p reflect.Type) bool { return true })))
	}
}

// Run as an individual test
func (test *MsgProcessingSpecTest) Run(t *testing.T) {
	test.overrideStateComparison(t)
	test.RunAsPartOfMultiTest(t)
}

func (test *MsgProcessingSpecTest) runPreTesting() (*ssv.Validator, *ssv.Committee, error) {
	var share *types.Share
	ketSetMap := make(map[phase0.ValidatorIndex]*testingutils.TestKeySet)
	if len(test.Runner.GetBaseRunner().Share) == 0 {
		panic("No share in base runner for tests")
	}
	for _, validatorShare := range test.Runner.GetBaseRunner().Share {
		share = validatorShare
		break
	}
	for valIdx, validatorShare := range test.Runner.GetBaseRunner().Share {
		ketSetMap[valIdx] = testingutils.KeySetForShare(validatorShare)
	}

	var v *ssv.Validator
	var c *ssv.Committee
	var lastErr error
	switch test.Runner.(type) {
	case *ssv.CommitteeRunner:
		c = testingutils.BaseCommitteeWithRunner(ketSetMap, test.Runner.(*ssv.CommitteeRunner))

		if !test.DontStartDuty {
			lastErr = c.StartDuty(test.Duty.(*types.CommitteeDuty))
		} else {
			c.Runners[test.Duty.DutySlot()] = test.Runner.(*ssv.CommitteeRunner)
		}

		for _, msg := range test.Messages {
			err := c.ProcessMessage(msg)
			if err != nil {
				lastErr = err
			}
			if test.DecidedSlashable && IsQBFTProposalMessage(msg) {
				for _, validatorShare := range test.Runner.GetBaseRunner().Share {
					test.Runner.GetSigner().(*testingutils.TestingKeyManager).AddSlashableDataRoot(validatorShare.
						SharePubKey, testingutils.TestingAttestationDataRoot[:])
				}
			}
		}

	default:
		v = testingutils.BaseValidator(testingutils.KeySetForShare(share))
		v.DutyRunners[test.Runner.GetBaseRunner().RunnerRoleType] = test.Runner
		v.Network = test.Runner.GetNetwork()

		if !test.DontStartDuty {
			lastErr = v.StartDuty(test.Duty)
		}
		for _, msg := range test.Messages {
			err := v.ProcessMessage(msg)
			if err != nil {
				lastErr = err
			}
		}
	}

	return v, c, lastErr
}

// IsQBFTProposalMessage checks if the message is a QBFT proposal message
func IsQBFTProposalMessage(msg *types.SignedSSVMessage) bool {
	if msg.SSVMessage.MsgType == types.SSVConsensusMsgType {
		qbftMsg := qbft.Message{}
		err := qbftMsg.Decode(msg.SSVMessage.Data)
		if err != nil {
			panic("could not decode message")
		}
		return qbftMsg.MsgType == qbft.ProposalMsgType
	}

	return v, lastErr
}

func (test *MsgProcessingSpecTest) compareBroadcastedBeaconMsgs(t *testing.T) {
	broadcastedRoots := test.Runner.GetBeaconNode().(*testingutils.TestingBeaconNode).BroadcastedRoots
	require.Len(t, broadcastedRoots, len(test.BeaconBroadcastedRoots))
	for _, r1 := range test.BeaconBroadcastedRoots {
		found := false
		for _, r2 := range broadcastedRoots {
			if r1 == hex.EncodeToString(r2[:]) {
				found = true
				break
			}
		}
		require.Truef(t, found, "broadcasted beacon root not found")
	}
}

func (test *MsgProcessingSpecTest) compareOutputMsgs(t *testing.T, v *ssv.Validator) {
	filterPartialSigs := func(messages []*types.SSVMessage) []*types.SSVMessage {
		ret := make([]*types.SSVMessage, 0)
		for _, msg := range messages {
			if msg.MsgType != types.SSVPartialSignatureMsgType {
				continue
			}
			ret = append(ret, msg)
		}
		return ret
	}
	broadcastedSignedMsgs := v.Network.(*testingutils.TestingNetwork).BroadcastedMsgs
	require.NoError(t, testingutils.VerifyListOfSignedSSVMessages(broadcastedSignedMsgs, v.Share.Committee))
	broadcastedMsgs := testingutils.ConvertBroadcastedMessagesToSSVMessages(broadcastedSignedMsgs)
	broadcastedMsgs = filterPartialSigs(broadcastedMsgs)
	require.Len(t, broadcastedMsgs, len(test.OutputMessages))
	index := 0
	for _, msg := range broadcastedMsgs {
		if msg.MsgType != types.SSVPartialSignatureMsgType {
			continue
		}

		msg1 := &types.SignedPartialSignatureMessage{}
		require.NoError(t, msg1.Decode(msg.Data))
		msg2 := test.OutputMessages[index]
		require.Len(t, msg1.Message.Messages, len(msg2.Message.Messages))

		// messages are not guaranteed to be in order so we map them and then test all roots to be equal
		roots := make(map[string]string)
		for i, partialSigMsg2 := range msg2.Message.Messages {
			r2, err := partialSigMsg2.GetRoot()
			require.NoError(t, err)
			if _, found := roots[hex.EncodeToString(r2[:])]; !found {
				roots[hex.EncodeToString(r2[:])] = ""
			} else {
				roots[hex.EncodeToString(r2[:])] = hex.EncodeToString(r2[:])
			}

			partialSigMsg1 := msg1.Message.Messages[i]
			r1, err := partialSigMsg1.GetRoot()
			require.NoError(t, err)

			if _, found := roots[hex.EncodeToString(r1[:])]; !found {
				roots[hex.EncodeToString(r1[:])] = ""
			} else {
				roots[hex.EncodeToString(r1[:])] = hex.EncodeToString(r1[:])
			}
		}
		for k, v := range roots {
			require.EqualValues(t, k, v, "missing output msg")
		}

		// test that slot is correct in broadcasted msg
		require.EqualValues(t, msg1.Message.Slot, msg2.Message.Slot, "incorrect broadcasted slot")

		index++
	}
}

func (test *MsgProcessingSpecTest) overrideStateComparison(t *testing.T) {
	overrideStateComparison(t, test, test.Name, reflect.TypeOf(test).String())
}

func overrideStateComparison(t *testing.T, test *MsgProcessingSpecTest, name string, testType string) {
	var runner ssv.Runner
	switch test.Runner.(type) {
	case *ssv.CommitteeRunner:
		runner = &ssv.CommitteeRunner{}
	case *ssv.AggregatorRunner:
		runner = &ssv.AggregatorRunner{}
	case *ssv.ProposerRunner:
		runner = &ssv.ProposerRunner{}
	case *ssv.SyncCommitteeAggregatorRunner:
		runner = &ssv.SyncCommitteeAggregatorRunner{}
	case *ssv.ValidatorRegistrationRunner:
		runner = &ssv.ValidatorRegistrationRunner{}
	case *ssv.VoluntaryExitRunner:
		runner = &ssv.VoluntaryExitRunner{}
	default:
		t.Fatalf("unknown runner type")
	}
	basedir, err := os.Getwd()
	require.NoError(t, err)
	runner, err = typescomparable.UnmarshalStateComparison(basedir, name, testType, runner)
	require.NoError(t, err)

	// override
	test.PostDutyRunnerState = runner

	root, err := runner.GetRoot()
	require.NoError(t, err)

	test.PostDutyRunnerStateRoot = hex.EncodeToString(root[:])
}

func (test *MsgProcessingSpecTest) GetPostState() (interface{}, error) {
	_, _, lastErr := test.runPreTesting()
	if lastErr != nil && len(test.ExpectedError) == 0 {
		return nil, lastErr
	}

	return test.Runner, nil
}

// Create alias without duty
type MsgProcessingSpecTestAlias struct {
	Name   string
	Runner ssv.Runner
	// No duty from type types.Duty
	Messages                []*types.SignedSSVMessage
	DecidedSlashable        bool
	PostDutyRunnerStateRoot string
	PostDutyRunnerState     types.Root `json:"-"`
	OutputMessages          []*types.PartialSignatureMessages
	BeaconBroadcastedRoots  []string
	DontStartDuty           bool
	ExpectedError           string
	BeaconDuty              *types.BeaconDuty    `json:"BeaconDuty,omitempty"`
	CommitteeDuty           *types.CommitteeDuty `json:"CommitteeDuty,omitempty"`
}

func (t *MsgProcessingSpecTest) MarshalJSON() ([]byte, error) {
	alias := &MsgProcessingSpecTestAlias{
		Name:                    t.Name,
		Runner:                  t.Runner,
		Messages:                t.Messages,
		DecidedSlashable:        t.DecidedSlashable,
		PostDutyRunnerStateRoot: t.PostDutyRunnerStateRoot,
		PostDutyRunnerState:     t.PostDutyRunnerState,
		OutputMessages:          t.OutputMessages,
		BeaconBroadcastedRoots:  t.BeaconBroadcastedRoots,
		DontStartDuty:           t.DontStartDuty,
		ExpectedError:           t.ExpectedError,
	}

	if t.Duty != nil {
		if beaconDuty, ok := t.Duty.(*types.BeaconDuty); ok {
			alias.BeaconDuty = beaconDuty
		} else if committeeDuty, ok := t.Duty.(*types.CommitteeDuty); ok {
			alias.CommitteeDuty = committeeDuty
		} else {
			return nil, errors.New("can't marshal StartNewRunnerDutySpecTest because t.Duty isn't BeaconDuty or CommitteeDuty")
		}
	}
	byts, err := json.Marshal(alias)

	return byts, err
}

func (t *MsgProcessingSpecTest) UnmarshalJSON(data []byte) error {
	aux := &MsgProcessingSpecTestAlias{}

	// Unmarshal the JSON data into the auxiliary struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.Name = aux.Name
	t.Runner = aux.Runner
	t.DecidedSlashable = aux.DecidedSlashable
	t.Messages = aux.Messages
	t.PostDutyRunnerStateRoot = aux.PostDutyRunnerStateRoot
	t.PostDutyRunnerState = aux.PostDutyRunnerState
	t.OutputMessages = aux.OutputMessages
	t.BeaconBroadcastedRoots = aux.BeaconBroadcastedRoots
	t.DontStartDuty = aux.DontStartDuty
	t.ExpectedError = aux.ExpectedError

	// Determine which type of duty was marshaled
	if aux.BeaconDuty != nil {
		t.Duty = aux.BeaconDuty
	} else if aux.CommitteeDuty != nil {
		t.Duty = aux.CommitteeDuty
	}

	return nil
}
