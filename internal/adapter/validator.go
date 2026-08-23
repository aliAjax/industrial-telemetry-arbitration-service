package adapter

import "errors"

var ErrRejectedFrame = errors.New("adapter rejected telemetry frame")

type Rule interface{ Check([]byte) bool }

type Validator struct{ rule Rule }

func NewValidator(rule Rule) (*Validator, error) {
	if isNilInterface(rule) {
		return nil, ErrNilAdapter
	}
	return &Validator{rule: rule}, nil
}

func (v *Validator) Validate(payload []byte) error {
	if v == nil || isNilInterface(v.rule) {
		return ErrNilAdapter
	}
	if !v.rule.Check(payload) {
		return ErrRejectedFrame
	}
	return nil
}
