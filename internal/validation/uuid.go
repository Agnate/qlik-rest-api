package validation

import (
	"github.com/google/uuid"
)

type RuleUUID struct {
	Raw    string
	Parsed uuid.UUID
}

func NewRuleUUID(uuid string) (RuleUUID, error) {
	v := RuleUUID{
		Raw: uuid,
	}
	return v, v.validate()
}

func (v *RuleUUID) validate() error {
	parsed, err := uuid.Parse(v.Raw)
	if err != nil {
		return err
	}
	v.Parsed = parsed
	return nil
}
