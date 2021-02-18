package deployment

import "fmt"

type ByTimestamp []string

func (s ByTimestamp) Len() int {
	return len(s)
}

func (s ByTimestamp) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByTimestamp) Less(i, j int) bool {
	timestampI, err := getTimestampFromDeployment(s[i])
	if err != nil {
		fmt.Println(err)
		return false
	}
	timestampJ, err := getTimestampFromDeployment(s[j])
	if err != nil {
		fmt.Println(err)
		return false
	}
	// Reversed that latest is element 0
	return timestampI > timestampJ
}
