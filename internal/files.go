package internal

import (
	"log"
	"os"
	"path/filepath"
)

func GetPolicies(dir string, policies []Policy) []Policy {
	entries, err := os.ReadDir(dir)
	Check(err)
	for _, entry := range entries {
		info, err := entry.Info()
		name := info.Name()
		path := filepath.Join(dir, name)
		Check(err)
		if !entry.IsDir() {
			if filepath.Ext(info.Name()) == ".xml" {
				policy, err := getPolicyDetails(path)
				checkDuplicate(*policy, policies)
				Check(err)
				policies = append(policies, *policy)
			}
		} else {
			policies = GetPolicies(path, policies)
		}
	}
	return policies
}

func checkDuplicate(policy Policy, policies []Policy) {
	hasDuplicate := false
	dupIndex := 0
	for i, p := range policies {
		if p.PolicyId == policy.PolicyId {
			hasDuplicate = true
			dupIndex = i
		}
	}

	if hasDuplicate {
		log.Fatalf("Found duplicate policies\n%s: %s\n%s: %s", policy.PolicyId, policy.Path, policies[dupIndex].PolicyId, policies[dupIndex].Path)
	}
}
