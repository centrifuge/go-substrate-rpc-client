package types

import "github.com/centrifuge/go-substrate-rpc-client/v4/scale"

type OptionExecutionResult struct {
	option
	value ExecutionResult
}

func NewOptionExecutionResult(value ExecutionResult) OptionExecutionResult {
	return OptionExecutionResult{option{hasValue: true}, value}
}

func NewOptionExecutionResultEmpty() OptionExecutionResult {
	return OptionExecutionResult{option: option{hasValue: false}}
}

func (o *OptionExecutionResult) Decode(decoder scale.Decoder) error {
	return decoder.DecodeOption(&o.hasValue, &o.value)
}

func (o OptionExecutionResult) Encode(encoder scale.Encoder) error {
	return encoder.EncodeOption(o.hasValue, o.value)
}

// SetSome sets a value
func (o *OptionExecutionResult) SetSome(value ExecutionResult) {
	o.hasValue = true
	o.value = value
}

// SetNone removes a value and marks it as missing
func (o *OptionExecutionResult) SetNone() {
	o.hasValue = false
	o.value = ExecutionResult{}
}

// Unwrap returns a flag that indicates whether a value is present and the stored value
func (o *OptionExecutionResult) Unwrap() (ok bool, value ExecutionResult) {
	return o.hasValue, o.value
}

type ExecutionResult struct {
	Outcome U32
	Error   XCMError
}

func (e *ExecutionResult) Decode(decoder scale.Decoder) error {
	if err := decoder.Decode(&e.Outcome); err != nil {
		return err
	}

	return decoder.Decode(&e.Error)
}

func (e ExecutionResult) Encode(encoder scale.Encoder) error {
	if err := encoder.Encode(e.Outcome); err != nil {
		return err
	}

	return encoder.Encode(e.Error)
}
