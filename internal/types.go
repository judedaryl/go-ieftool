package internal

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

type Policy struct {
	PolicyId       string `json:"policyId"`
	ParentPolicyId string `json:"parentPolicyId"`
	Path           string `json:"path"`
}

func (p Policy) HasParent() bool {
	return p.ParentPolicyId != ""
}
