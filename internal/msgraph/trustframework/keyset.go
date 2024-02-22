package trustframework

type KeySet struct {
	IDs []string
}

func NewKeySet(ids []string) *KeySet {
	return &KeySet{IDs: ids}
}

func (ks *KeySet) Add(s string) {
	ks.IDs = append(ks.IDs, s)
}

func (ks *KeySet) Remove(s string) {
	for i, e := range ks.IDs {
		if e == s {
			ids := append(ks.IDs[:i], ks.IDs[i+1:]...)
			ks.IDs = ids
			break
		}
	}
}

func (ks *KeySet) Has(s string) bool {
	for _, e := range ks.IDs {
		if e == s {
			return true
		}
	}

	return false
}
