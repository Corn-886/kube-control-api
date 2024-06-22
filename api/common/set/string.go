package set

type Empty struct{}
type String map[string]Empty

func (s String) Delete(items ...string) String {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

func (s String) Has(item string) bool {
	_, contained := s[item]
	return contained
}

// Insert adds items to the set.
func (s String) Insert(items ...string) String {
	for _, item := range items {
		s[item] = Empty{}
	}
	return s
}
