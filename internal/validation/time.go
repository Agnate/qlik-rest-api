package validation

import "time"

type RuleTime struct {
	Raw    string
	Format string
	Parsed time.Time
}

func NewRuleTime(format string, time string) (*RuleTime, error) {
	v := &RuleTime{
		Raw:    time,
		Format: format,
	}
	return v, v.validate()
}

func (v *RuleTime) validate() error {
	parsed, err := time.Parse(v.Format, v.Raw)
	if err != nil {
		return err
	}
	v.Parsed = parsed
	return nil
}
