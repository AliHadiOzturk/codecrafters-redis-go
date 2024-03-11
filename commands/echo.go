package commands

func Echo(parameters []string) (string, error) {
	var response string = "+"
	for _, v := range parameters {
		response += v
	}

	return response, nil
}
