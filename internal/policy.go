package internal

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"com.go.ieftool/internal/tree"
)

type PolicyB2C struct {
	PolicyId   string          `xml:"PolicyId,attr"`
	BasePolicy []BasePolicyB2C `xml:"BasePolicy"`
}

type BasePolicyB2C struct {
	PolicyId []PolicyIdB2C `xml:"PolicyId"`
}

type PolicyIdB2C struct {
	Value string `xml:",chardata"`
}

func (p *PolicyB2C) isPolicy(content []byte) bool {
	err := xml.Unmarshal(content, p)
	if err != nil {
		return false
	}
	return p.PolicyId != ""
}

type Policy struct {
	PolicyId       string `json:"policyId"`
	ParentPolicyId string `json:"parentPolicyId"`
	Path           string `json:"path"`
}

func NewPolicy(content []byte, filePath string) (*Policy, error) {
	b2cPolicy := &PolicyB2C{}
	if !b2cPolicy.isPolicy(content) {
		return nil, errors.New(fmt.Sprintf("File %s is not a policy", filePath))
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

func (p Policy) HasParent() bool {
	return p.ParentPolicyId != ""
}

type Policies []Policy

func NewPoliciesFromDir(d string) Policies {
	ps := Policies{}
	ps.get(d)

	return ps
}

func (ps *Policies) get(d string) {
	entries, err := os.ReadDir(d)
	Check(err)
	for _, entry := range entries {
		info, err := entry.Info()
		name := info.Name()
		path := filepath.Join(d, name)
		Check(err)
		if !entry.IsDir() {
			if filepath.Ext(info.Name()) == ".xml" {
				content, err := os.ReadFile(path)
				if err != nil {
					Check(err)
				}
				p, err := NewPolicy(content, path)
				Check(err)
				ps.checkDuplicate(p)
				*ps = append(*ps, *p)
			}
		} else {
			ps.get(path)
		}
	}
}

func (ps *Policies) checkDuplicate(policy *Policy) {
	hasDuplicate := false
	for _, p := range *ps {
		if p.PolicyId == policy.PolicyId {
			hasDuplicate = true
		}
	}

	if hasDuplicate {
		log.Fatalf("Found duplicate policies %s: %s", policy.PolicyId, policy.Path)
	}
}

func (ps *Policies) GetBatch() [][]Policy {
	log.Println("Building Policy Tree")
	rp := ps.findRoot()
	r := tree.NewBranch(rp)
	ps.recursiveAddBranch(&r)

	batch := &[][]Policy{}
	log.Println("Determining batches")
	ps.getBatch([]tree.Branch[Policy]{r}, batch)
	log.Printf("Found %d batches", len(*batch))

	return *batch
}

func (ps *Policies) findRoot() Policy {
	var _policy Policy
	for i, policy := range *ps {
		if policy.ParentPolicyId == "" {
			_policy = policy
			*ps = remove(*ps, i)
		}
	}

	return _policy
}

func (ps *Policies) recursiveAddBranch(parent *tree.Branch[Policy]) {
	childPolicies := ps.findChildPolicies(parent.Data().PolicyId)
	if len(childPolicies) == 0 {
		return
	}
	for _, child := range childPolicies {
		branch := tree.NewBranch(child)
		ps.recursiveAddBranch(&branch)
		parent.AddChild(branch)
	}
}

func (ps *Policies) findChildPolicies(policyId string) []Policy {
	var _policies []Policy
	for _, policy := range *ps {
		if policy.ParentPolicyId == policyId {
			_policies = append(_policies, policy)
		}
	}
	return _policies
}

func (ps *Policies) getBatch(tree []tree.Branch[Policy], policies *[][]Policy) {
	var batch []Policy
	for _, branch := range tree {
		batch = append(batch, branch.Data())
	}
	for _, branch := range tree {
		if len(branch.Children()) > 0 {
			ps.getBatch(branch.Children(), policies)
		}
	}
	*policies = append([][]Policy{batch}, *policies...)
}
