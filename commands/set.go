package commands

var Registry map[string]interface{}

func Set(parameters []string) (string, error) {
	key := parameters[0]
	value := parameters[1]

	if Registry == nil {
		Registry = make(map[string]interface{})
	}

	Registry[key] = value

	return "+OK", nil
}
