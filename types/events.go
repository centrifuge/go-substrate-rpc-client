// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
)

// EventBalancesEndowed is emitted when an account is created with some free balance
type EventBalancesEndowed struct {
	Phase   Phase
	Who     AccountID
	Balance U128
	Topics  []Hash
}

// EventDustLost is emitted when an account is removed with a balance that is
// non-zero but below ExistentialDeposit, resulting in a loss.
type EventBalancesDustLost struct {
	Phase   Phase
	Who     AccountID
	Balance U128
	Topics  []Hash
}

// EventBalancesTransfer is emitted when a transfer succeeded (from, to, value)
type EventBalancesTransfer struct {
	Phase  Phase
	From   AccountID
	To     AccountID
	Value  U128
	Topics []Hash
}

// EventBalanceSet is emitted when a balance is set by root
type EventBalancesBalanceSet struct {
	Phase    Phase
	Who      AccountID
	Free     U128
	Reserved U128
	Topics   []Hash
}

// EventDeposit is emitted when an account receives some free balance
type EventBalancesDeposit struct {
	Phase   Phase
	Who     AccountID
	Balance U128
	Topics  []Hash
}

// EventGrandpaNewAuthorities is emitted when a new authority set has been applied
type EventGrandpaNewAuthorities struct {
	Phase          Phase
	NewAuthorities []struct {
		AuthorityID     AuthorityID
		AuthorityWeight U64
	}
	Topics []Hash
}

// EventGrandpaPaused is emitted when the current authority set has been paused
type EventGrandpaPaused struct {
	Phase  Phase
	Topics []Hash
}

// EventGrandpaResumed is emitted when the current authority set has been resumed
type EventGrandpaResumed struct {
	Phase  Phase
	Topics []Hash
}

// EventImOnlineHeartbeatReceived is emitted when a new heartbeat was received from AuthorityId
type EventImOnlineHeartbeatReceived struct {
	Phase       Phase
	AuthorityID AuthorityID
	Topics      []Hash
}

// EventImOnlineAllGood is emitted when at the end of the session, no offence was committed
type EventImOnlineAllGood struct {
	Phase  Phase
	Topics []Hash
}

// Exposure lists the own and nominated stake of a validator
type Exposure struct {
	Total  UCompact
	Own    UCompact
	Others []IndividualExposure
}

// IndividualExposure contains the nominated stake by one specific third party
type IndividualExposure struct {
	Who   AccountID
	Value UCompact
}

// EventImOnlineSomeOffline is emitted when the end of the session, at least once validator was found to be offline
type EventImOnlineSomeOffline struct {
	Phase                Phase
	IdentificationTuples []struct {
		ValidatorID        AccountID
		FullIdentification Exposure
	}
	Topics []Hash
}

// EventIndicesIndexAssigned is emitted when an index is assigned to an AccountID.
type EventIndicesIndexAssigned struct {
	Phase        Phase
	AccountID    AccountID
	AccountIndex AccountIndex
	Topics       []Hash
}

// EventIndicesIndexFreed is emitted when an index is unassigned.
type EventIndicesIndexFreed struct {
	Phase        Phase
	AccountIndex AccountIndex
	Topics       []Hash
}

// EventOffencesOffence is emitted when there is an offence reported of the given kind happened at the session_index
// and (kind-specific) time slot. This event is not deposited for duplicate slashes
type EventOffencesOffence struct {
	Phase          Phase
	Kind           Bytes16
	OpaqueTimeSlot Bytes
	Topics         []Hash
}

// EventSessionNewSession is emitted when a new session has happened. Note that the argument is the session index,
// not the block number as the type might suggest
type EventSessionNewSession struct {
	Phase        Phase
	SessionIndex U32
	Topics       []Hash
}

// EventStakingReward is emitted when all validators have been rewarded by the first balance; the second is the
// remainder, from the maximum amount of reward.
type EventStakingReward struct {
	Phase     Phase
	Balance   U128
	Remainder U128
	Topics    []Hash
}

// EventStakingSlash is emitted when one validator (and its nominators) has been slashed by the given amount
type EventStakingSlash struct {
	Phase     Phase
	AccountID AccountID
	Balance   U128
	Topics    []Hash
}

// EventStakingOldSlashingReportDiscarded is emitted when an old slashing report from a prior era was discarded because
// it could not be processed
type EventStakingOldSlashingReportDiscarded struct {
	Phase        Phase
	SessionIndex U32
	Topics       []Hash
}

// EventSystemExtrinsicSuccessV8 is emitted when an extrinsic completed successfully
//
// Deprecated: EventSystemExtrinsicSuccessV8 exists to allow users to simply implement their own EventRecords struct if
// they are on metadata version 8 or below. Use EventSystemExtrinsicSuccess otherwise
type EventSystemExtrinsicSuccessV8 struct {
	Phase  Phase
	Topics []Hash
}

// EventSystemExtrinsicSuccess is emitted when an extrinsic completed successfully
type EventSystemExtrinsicSuccess struct {
	Phase        Phase
	DispatchInfo DispatchInfo
	Topics       []Hash
}

// DispatchInfo contains a bundle of static information collected from the `#[weight = $x]` attributes.
type DispatchInfo struct {
	// Weight of this transaction
	Weight U32
	// Class of this transaction
	Class DispatchClass
	// PaysFee indicates whether this transaction pays fees
	PaysFee bool
}

// DispatchClass is a generalized group of dispatch types. This is only distinguishing normal, user-triggered
// transactions (`Normal`) and anything beyond which serves a higher purpose to the system (`Operational`).
type DispatchClass struct {
	// A normal dispatch
	IsNormal bool
	// An operational dispatch
	IsOperational bool
}

func (d *DispatchClass) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if b == 0 {
		d.IsNormal = true
	} else if b == 1 {
		d.IsOperational = true
	}
	return err
}

func (d DispatchClass) Encode(encoder scale.Encoder) error {
	var err error
	if d.IsNormal {
		err = encoder.PushByte(0)
	} else if d.IsOperational {
		err = encoder.PushByte(1)
	}
	return err
}

// EventSystemExtrinsicFailedV8 is emitted when an extrinsic failed
//
// Deprecated: EventSystemExtrinsicFailedV8 exists to allow users to simply implement their own EventRecords struct if
// they are on metadata version 8 or below. Use EventSystemExtrinsicFailed otherwise
type EventSystemExtrinsicFailedV8 struct {
	Phase         Phase
	DispatchError DispatchError
	Topics        []Hash
}

// EventSystemExtrinsicFailed is emitted when an extrinsic failed
type EventSystemExtrinsicFailed struct {
	Phase         Phase
	DispatchError DispatchError
	DispatchInfo  DispatchInfo
	Topics        []Hash
}

// EventSystemCodeUpdated is emitted when the runtime code (`:code`) is updated
type EventSystemCodeUpdated struct {
	Phase  Phase
	Topics []Hash
}

// EventSystemNewAccount is emitted when a new account was created
type EventSystemNewAccount struct {
	Phase  Phase
	Who    AccountID
	Topics []Hash
}

// EventSystemKilledAccount is emitted when an account is reaped
type EventSystemKilledAccount struct {
	Phase  Phase
	Who    AccountID
	Topics []Hash
}

// EventAssetIssued is emitted when an asset is issued.
type EventAssetIssued struct {
	Phase   Phase
	AssetID U32
	Who     AccountID
	Balance U128
	Topics  []Hash
}

// EventAssetTransferred is emitted when an asset is transferred.
type EventAssetTransferred struct {
	Phase   Phase
	AssetID U32
	To      AccountID
	From    AccountID
	Balance U128
	Topics  []Hash
}

// EventAssetDestroyed is emitted when an asset is destroyed.
type EventAssetDestroyed struct {
	Phase   Phase
	AssetID U32
	Who     AccountID
	Balance U128
	Topics  []Hash
}

// EventDemocracyProposed is emitted when a motion has been proposed by a public account.
type EventDemocracyProposed struct {
	Phase         Phase
	ProposalIndex U32
	Balance       U128
	Topics        []Hash
}

// EventDemocracyTabled is emitted when a public proposal has been tabled for referendum vote.
type EventDemocracyTabled struct {
	Phase         Phase
	ProposalIndex U32
	Balance       U128
	Accounts      []AccountID
	Topics        []Hash
}

// EventDemocracyExternalTabled is emitted when an external proposal has been tabled.
type EventDemocracyExternalTabled struct {
	Phase  Phase
	Topics []Hash
}

// VoteThreshold is a means of determining if a vote is past pass threshold.
type VoteThreshold byte

const (
	// SuperMajorityApprove require super majority of approvals is needed to pass this vote.
	SuperMajorityApprove VoteThreshold = 0
	// SuperMajorityAgainst require super majority of rejects is needed to fail this vote.
	SuperMajorityAgainst VoteThreshold = 1
	// SimpleMajority require simple majority of approvals is needed to pass this vote.
	SimpleMajority VoteThreshold = 2
)

func (v *VoteThreshold) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	vb := VoteThreshold(b)
	switch vb {
	case SuperMajorityApprove, SuperMajorityAgainst, SimpleMajority:
		*v = vb
	default:
		return fmt.Errorf("unknown VoteThreshold enum: %v", vb)
	}
	return err
}

func (v VoteThreshold) Encode(encoder scale.Encoder) error {
	return encoder.PushByte(byte(v))
}

// EventDemocracyStarted is emitted when a referendum has begun.
type EventDemocracyStarted struct {
	Phase           Phase
	ReferendumIndex U32
	VoteThreshold   VoteThreshold
	Topics          []Hash
}

// EventDemocracyPassed is emitted when a proposal has been approved by referendum.
type EventDemocracyPassed struct {
	Phase           Phase
	ReferendumIndex U32
	Topics          []Hash
}

// EventDemocracyNotPassed is emitted when a proposal has been rejected by referendum.
type EventDemocracyNotPassed struct {
	Phase           Phase
	ReferendumIndex U32
	Topics          []Hash
}

// EventDemocracyCancelled is emitted when a referendum has been cancelled.
type EventDemocracyCancelled struct {
	Phase           Phase
	ReferendumIndex U32
	Topics          []Hash
}

// EventDemocracyExecuted is emitted when a proposal has been enacted.
type EventDemocracyExecuted struct {
	Phase           Phase
	ReferendumIndex U32
	Result          bool
	Topics          []Hash
}

// EventDemocracyDelegated is emitted when an account has delegated their vote to another account.
type EventDemocracyDelegated struct {
	Phase  Phase
	Who    AccountID
	Target AccountID
	Topics []Hash
}

// EventDemocracyUndelegated is emitted when an account has cancelled a previous delegation operation.
type EventDemocracyUndelegated struct {
	Phase  Phase
	Target AccountID
	Topics []Hash
}

// EventDemocracyVetoed is emitted when an external proposal has been vetoed.
type EventDemocracyVetoed struct {
	Phase       Phase
	Who         AccountID
	Hash        Hash
	BlockNumber BlockNumber
	Topics      []Hash
}

// EventDemocracyPreimageNoted is emitted when a proposal's preimage was noted, and the deposit taken.
type EventDemocracyPreimageNoted struct {
	Phase     Phase
	Hash      Hash
	AccountID AccountID
	Balance   U128
	Topics    []Hash
}

// EventDemocracyPreimageUsed is emitted when a proposal preimage was removed and used (the deposit was returned).
type EventDemocracyPreimageUsed struct {
	Phase     Phase
	Hash      Hash
	AccountID AccountID
	Balance   U128
	Topics    []Hash
}

// EventDemocracyPreimageInvalid is emitted when a proposal could not be executed because its preimage was invalid.
type EventDemocracyPreimageInvalid struct {
	Phase           Phase
	Hash            Hash
	ReferendumIndex U32
	Topics          []Hash
}

// EventDemocracyPreimageMissing is emitted when a proposal could not be executed because its preimage was missing.
type EventDemocracyPreimageMissing struct {
	Phase           Phase
	Hash            Hash
	ReferendumIndex U32
	Topics          []Hash
}

// EventDemocracyPreimageReaped is emitted when a registered preimage was removed
// and the deposit collected by the reaper (last item).
type EventDemocracyPreimageReaped struct {
	Phase    Phase
	Hash     Hash
	Provider AccountID
	Balance  U128
	Who      AccountID
	Topics   []Hash
}

// EventDemocracyUnlocked is emitted when an account has been unlocked successfully.
type EventDemocracyUnlocked struct {
	Phase     Phase
	AccountID AccountID
	Topics    []Hash
}

// EventCollectiveProposed is emitted when a motion (given hash) has been proposed (by given account)
// with a threshold (given `MemberCount`).
type EventCollectiveProposed struct {
	Phase         Phase
	Who           AccountID
	ProposalIndex U32
	Proposal      Hash
	MemberCount   U32
	Topics        []Hash
}

// EventCollectiveVote is emitted when a motion (given hash) has been voted on by given account, leaving
// a tally (yes votes and no votes given respectively as `MemberCount`).
type EventCollectiveVoted struct {
	Phase    Phase
	Who      AccountID
	Proposal Hash
	Approve  bool
	YesCount U32
	NoCount  U32
	Topics   []Hash
}

// EventCollectiveApproved is emitted when a motion was approved by the required threshold.
type EventCollectiveApproved struct {
	Phase    Phase
	Proposal Hash
	Topics   []Hash
}

// EventCollectiveDisapproved is emitted when a motion was not approved by the required threshold.
type EventCollectiveDisapproved struct {
	Phase    Phase
	Proposal Hash
	Topics   []Hash
}

// EventCollectiveExecuted is emitted when a motion was executed; `bool` is true if returned without error.
type EventCollectiveExecuted struct {
	Phase    Phase
	Proposal Hash
	Ok       bool
	Topics   []Hash
}

// EventCollectiveMemberExecuted is emitted when a single member did some action;
// `bool` is true if returned without error.
type EventCollectiveMemberExecuted struct {
	Phase    Phase
	Proposal Hash
	Ok       bool
	Topics   []Hash
}

// EventCollectiveClosed is emitted when a proposal was closed after its duration was up.
type EventCollectiveClosed struct {
	Phase    Phase
	Proposal Hash
	YesCount U32
	NoCount  U32
	Topics   []Hash
}

// EventElectionsNewTerm is emitted when a new term with new members.
// This indicates that enough candidates existed, not that enough have has been elected.
// The inner value must be examined for this purpose.
type EventElectionsNewTerm struct {
	Phase      Phase
	NewMembers []struct {
		Member  AccountID
		Balance U128
	}
	Topics []Hash
}

// EventElectionsEmpty is emitted when No (or not enough) candidates existed for this round.
type EventElectionsEmptyTerm struct {
	Phase  Phase
	Topics []Hash
}

// EventElectionsMemberKicked is emitted when a member has been removed.
// This should always be followed by either `NewTerm` or `EmptyTerm`.
type EventElectionsMemberKicked struct {
	Phase  Phase
	Member AccountID
	Topics []Hash
}

// EventElectionsMemberRenounced is emitted when a member has renounced their candidacy.
type EventElectionsMemberRenounced struct {
	Phase  Phase
	Member AccountID
	Topics []Hash
}

// EventElectionsVoterReported is emitted when a voter (first element) was reported (by the second element)
// with the the report being successful or not (third element).
type EventElectionsVoterReported struct {
	Phase            Phase
	Target, Reporter AccountID
	Valid            bool
	Topics           []Hash
}

// A name was set or reset (which will remove all judgements).
type EventIdentitySet struct {
	Phase    Phase
	Identity AccountID
	Topics   []Hash
}

// A name was cleared, and the given balance returned.
type EventIdentityCleared struct {
	Phase    Phase
	Identity AccountID
	Balance  U128
	Topics   []Hash
}

// A name was removed and the given balance slashed.
type EventIdentityKilled struct {
	Phase    Phase
	Identity AccountID
	Balance  U128
	Topics   []Hash
}

// A judgement was asked from a registrar.
type EventIdentityJudgementRequested struct {
	Phase          Phase
	Sender         AccountID
	RegistrarIndex U32
	Topics         []Hash
}

// A judgement request was retracted.
type EventIdentityJudgementUnrequested struct {
	Phase          Phase
	Sender         AccountID
	RegistrarIndex U32
	Topics         []Hash
}

// A judgement was given by a registrar.
type EventIdentityJudgementGiven struct {
	Phase          Phase
	Target         AccountID
	RegistrarIndex U32
	Topics         []Hash
}

// A registrar was added.
type EventIdentityRegistrarAdded struct {
	Phase          Phase
	RegistrarIndex U32
	Topics         []Hash
}

// EventRecoveryCreated is emitted when a recovery process has been set up for an account
type EventRecoveryCreated struct {
	Phase  Phase
	Who    AccountID
	Topics []Hash
}

// EventRecoveryInitiated is emitted when a recovery process has been initiated for account_1 by account_2
type EventRecoveryInitiated struct {
	Phase   Phase
	Account AccountID
	Who     AccountID
	Topics  []Hash
}

// EventRecoveryVouched is emitted when a recovery process for account_1 by account_2 has been vouched for by account_3
type EventRecoveryVouched struct {
	Phase   Phase
	Lost    AccountID
	Rescuer AccountID
	Who     AccountID
	Topics  []Hash
}

// EventRecoveryClosed is emitted when a recovery process for account_1 by account_2 has been closed
type EventRecoveryClosed struct {
	Phase   Phase
	Who     AccountID
	Rescuer AccountID
	Topics  []Hash
}

// EventRecoveryAccountRecovered is emitted when account_1 has been successfully recovered by account_2
type EventRecoveryAccountRecovered struct {
	Phase   Phase
	Who     AccountID
	Rescuer AccountID
	Topics  []Hash
}

// EventRecoveryRemoved is emitted when a recovery process has been removed for an account
type EventRecoveryRemoved struct {
	Phase  Phase
	Who    AccountID
	Topics []Hash
}

// EventSudoSudid is emitted when a sudo just took place.
type EventSudoSudid struct {
	Phase  Phase
	Result bool
	Topics []Hash
}

// EventSudoKeyChanged is emitted when the sudoer just switched identity; the old key is supplied.
type EventSudoKeyChanged struct {
	Phase     Phase
	AccountID AccountID
	Topics    []Hash
}

// A sudo just took place.
type EventSudoAsDone struct {
	Phase  Phase
	Done   bool
	Topics []Hash
}

// EventTreasuryProposed is emitted when New proposal.
type EventTreasuryProposed struct {
	Phase         Phase
	ProposalIndex U32
	Topics        []Hash
}

// EventTreasurySpending is emitted when we have ended a spend period and will now allocate funds.
type EventTreasurySpending struct {
	Phase           Phase
	BudgetRemaining U128
	Topics          []Hash
}

// EventTreasuryAwarded is emitted when some funds have been allocated.
type EventTreasuryAwarded struct {
	Phase         Phase
	ProposalIndex U32
	Amount        U128
	Beneficiary   AccountID
	Topics        []Hash
}

// EventTreasuryRejected is emitted when s proposal was rejected; funds were slashed.
type EventTreasuryRejected struct {
	Phase         Phase
	ProposalIndex U32
	Amount        U128
	Topics        []Hash
}

// EventTreasuryBurnt is emitted when some of our funds have been burnt.
type EventTreasuryBurnt struct {
	Phase  Phase
	Burn   U128
	Topics []Hash
}

// EventTreasuryRollover is emitted when spending has finished; this is the amount that rolls over until next spend.
type EventTreasuryRollover struct {
	Phase           Phase
	BudgetRemaining U128
	Topics          []Hash
}

// EventTreasuryDeposit is emitted when some funds have been deposited.
type EventTreasuryDeposit struct {
	Phase     Phase
	Deposited U128
	Topics    []Hash
}

// EventTreasuryNewTip is emitted when a new tip suggestion has been opened.
type EventTreasuryNewTip struct {
	Phase  Phase
	Hash   Hash
	Topics []Hash
}

// EventTreasuryTipClosing is emitted when a tip suggestion has reached threshold and is closing.
type EventTreasuryTipClosing struct {
	Phase  Phase
	Hash   Hash
	Topics []Hash
}

// EventTreasuryTipClosed is emitted when a tip suggestion has been closed.
type EventTreasuryTipClosed struct {
	Phase     Phase
	Hash      Hash
	AccountID AccountID
	Balance   U128
	Topics    []Hash
}

// EventTreasuryTipRetracted is emitted when a tip suggestion has been retracted.
type EventTreasuryTipRetracted struct {
	Phase  Phase
	Hash   Hash
	Topics []Hash
}

// EventUtilityBatchInterrupted is emitted when a batch of dispatches did not complete fully.
//Index of first failing dispatch given, as well as the error.
type EventUtilityBatchInterrupted struct {
	Phase         Phase
	Index         U32
	DispatchError DispatchError
	Topics        []Hash
}

// EventUtilityBatchCompleted is emitted when a batch of dispatches completed fully with no error.
type EventUtilityBatchCompleted struct {
	Phase  Phase
	Topics []Hash
}

// EventUtilityNewMultisig is emitted when a new multisig operation has begun.
// First param is the account that is approving, second is the multisig account, third is hash of the call.
type EventUtilityNewMultisig struct {
	Phase    Phase
	Who, ID  AccountID
	CallHash Hash
	Topics   []Hash
}

// TimePoint is a global extrinsic index, formed as the extrinsic index within a block,
// together with that block's height.
type TimePoint struct {
	Height BlockNumber
	Index  U32
}

// EventUtility is emitted when a multisig operation has been approved by someone. First param is the account that is
// approving, third is the multisig account, fourth is hash of the call.
type EventUtilityMultisigApproval struct {
	Phase     Phase
	Who       AccountID
	TimePoint TimePoint
	ID        AccountID
	CallHash  Hash
	Topics    []Hash
}

// DispatchResult can be returned from dispatchable functions
type DispatchResult struct {
	Ok    bool
	Error DispatchError
}

func (d *DispatchResult) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		d.Ok = true
		return nil
	default:
		derr := DispatchError{}
		err = decoder.Decode(&derr)
		if err != nil {
			return err
		}
		d.Error = derr
		return nil
	}
}

func (d DispatchResult) Encode(encoder scale.Encoder) error {
	if d.Ok {
		return encoder.PushByte(0)
	}
	return d.Error.Encode(encoder)
}

// EventUtility is emitted when a multisig operation has been executed. First param is the account that is
// approving, third is the multisig account, fourth is hash of the call to be executed.
type EventUtilityMultisigExecuted struct {
	Phase     Phase
	Who       AccountID
	TimePoint TimePoint
	ID        AccountID
	CallHash  Hash
	Result    DispatchResult
	Topics    []Hash
}

// EventUtility is emitted when a multisig operation has been cancelled. First param is the account that is
// cancelling, third is the multisig account, fourth is hash of the call.
type EventUtilityMultisigCancelled struct {
	Phase     Phase
	Who       AccountID
	TimePoint TimePoint
	ID        AccountID
	CallHash  Hash
	Topics    []Hash
}
