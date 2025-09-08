package internal

import (
	"encoding/xml"
	"log"
	"strings"

	"com.go.ieftool/internal/tree"
)

func isPolicy(content []byte) bool {
	b2cPolicy := &PolicyB2C{}
	err := xml.Unmarshal(content, b2cPolicy)
	if err != nil {
		return false
	}
	return b2cPolicy.PolicyId != ""
}

func getPolicyDetails(content []byte, filePath string) (*Policy, error) {
	b2cPolicy := &PolicyB2C{}
	err := xml.Unmarshal(content, b2cPolicy)
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

func createTrees(policies []Policy) []tree.Branch[Policy] {
	log.Println("Building Policy Tree")
	rootPolicies := findRootPolicies(&policies)
	trees := []tree.Branch[Policy]{}
	for _, rootPolicy := range rootPolicies {
		root := tree.NewBranch(rootPolicy)
		recursiveAddBranch(&root, &policies)
		trees = append(trees, root)
	}
	return trees
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
	roots := createTrees(policies)
	batch := &[][]Policy{}
	log.Println("Determining batches")
	internalCreateBatch(roots, batch)
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

func findRootPolicies(policies *[]Policy) []Policy {
	result := []Policy{}
	for _, policy := range *policies {
		if policy.ParentPolicyId == "" {
			result = append(result, policy)
		}
	}
	return result
}

func findChildPolicies(policies *[]Policy, policyId string) []Policy {
	_policies := []Policy{}
	for _, policy := range *policies {
		if policy.ParentPolicyId == policyId {
			_policies = append(_policies, policy)
		}
	}
	return _policies
}
