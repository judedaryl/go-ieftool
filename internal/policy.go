package internal

import (
	"encoding/xml"
	"log"
	"os"
	"strings"

	"com.go.ieftool/internal/tree"
)

func getPolicyDetails(filePath string) (*Policy, error) {
	b2cPolicy := &PolicyB2C{}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(content, b2cPolicy)
	if err != nil {
		return nil, err
	}
	policy := &Policy{
		PolicyId: strings.ToLower(b2cPolicy.PolicyId),
		Path:     filePath,
	}
	if len(b2cPolicy.BasePolicy) > 0 {
		parentPolicyId := b2cPolicy.BasePolicy[0].PolicyId[0].Value
		log.Printf("%s found parent policy id %s", b2cPolicy.PolicyId, parentPolicyId)
		policy.ParentPolicyId = strings.ToLower(parentPolicyId)
	}
	return policy, nil
}

func createTree(policies []Policy) tree.Branch[Policy] {
	log.Println("Building Policy Tree")
	rootPolicy := findRootPolicy(&policies)
	root := tree.NewBranch(rootPolicy)
	recursiveAddBranch(&root, &policies)
	return root
}

func recursiveAddBranch(parent *tree.Branch[Policy], policies *[]Policy) {
	childPolicies := findChildPolicies(policies, parent.Data().PolicyId)
	if len(childPolicies) == 0 {
		return
	}
	for _, child := range childPolicies {
		branch := tree.NewBranch(child)
		recursiveAddBranch(&branch, policies)
		parent.AddChild(branch)
	}
}

func CreateBatchedArray(policies []Policy) [][]Policy {
	root := createTree(policies)
	batch := &[][]Policy{}
	log.Println("Determining batches")
	internalCreateBatch([]tree.Branch[Policy]{root}, batch)
	log.Printf("Found %d batches", len(*batch))
	return *batch
}

func internalCreateBatch(tree []tree.Branch[Policy], policies *[][]Policy) {
	batch := []Policy{}
	for _, branch := range tree {
		batch = append(batch, branch.Data())
	}
	for _, branch := range tree {
		if len(branch.Children()) > 0 {
			internalCreateBatch(branch.Children(), policies)
		}
	}
	*policies = append([][]Policy{batch}, *policies...)
}

func remove[T any](slice []T, s int) []T {
	if len(slice) == 1 {
		return []T{}
	}
	newArr := []T{}
	for i := range slice {
		if i != s {
			newArr = append(newArr, slice[i])
		}
	}
	return newArr
}

func findRootPolicy(policies *[]Policy) Policy {
	var _policy Policy
	for i, policy := range *policies {
		if policy.ParentPolicyId == "" {
			_policy = policy
			*policies = remove(*policies, i)
		}
	}
	return _policy
}

func findChildPolicies(policies *[]Policy, policyId string) []Policy {
	_policies := []Policy{}
	for i, policy := range *policies {
		if policy.ParentPolicyId == policyId {
			_policies = append(_policies, policy)
			*policies = remove(*policies, i)
		}
	}
	return _policies
}
