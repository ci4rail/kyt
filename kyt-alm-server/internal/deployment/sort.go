package deployment

import "fmt"

// ByTimestamp is a custom sort format that can sort customer deployments by their timestamp
type ByTimestamp []string

func (s ByTimestamp) Len() int {
	return len(s)
}

func (s ByTimestamp) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByTimestamp) Less(i, j int) bool {
	timestampI, err := getTimestampFromDeployment(s[i])
	// if timestamp cannot be found for any reason, order way back
	if err != nil {
		fmt.Println(err)
		return false
	}
	// if timestamp cannot be found for any reason, order way back
	timestampJ, err := getTimestampFromDeployment(s[j])
	if err != nil {
		fmt.Println(err)
		return true
	}
	// Reversed that latest is element 0
	return timestampI > timestampJ
}
