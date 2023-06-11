package arrayutils

func Uniq(arr []string) []string {
	deDupMap := make(map[string]bool)
	result := []string{}

	for _, i := range arr {
		if !deDupMap[i] {
			continue
		}

		deDupMap[i] = true
		result = append(result, i)
	}

	return result
}
