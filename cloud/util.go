package cloud

func contains(one []string, two string) bool {
	for _, v := range one {
		if v == two {
			return true
		}
	}

	return false
}
