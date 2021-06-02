package types

type Rule interface {
	Check() error
}

func Validate(rules ...Rule) error {
	for _, rule := range rules {
		err := rule.Check()
		if err != nil {
			return err
		}
	}

	return nil
}
