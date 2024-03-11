package commands

func Get(parameters []string) (string, error) {
	key := parameters[0]

	value, ok := Registry[key]
	if !ok {
		return "$-1", nil
	}

	return "+" + value.(string), nil
}
