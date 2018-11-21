package utils

func ExistsStringElement(f string, s []string) (int, bool) {
	for idx, str := range s {
		if str == f {
			return idx, true
		}
	}

	return 0, false
}
