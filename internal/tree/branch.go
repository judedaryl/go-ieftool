package tree

type Branch[T any] struct {
	data     T
	children []Branch[T]
}

func NewBranch[T any](data T) Branch[T] {
	return Branch[T]{data: data, children: []Branch[T]{}}
}

func (branch *Branch[T]) AddChild(data Branch[T]) {
	branch.children = append(branch.children, data)
}

func (branch *Branch[T]) Data() T {
	return branch.data
}

func (branch *Branch[T]) Children() []Branch[T] {
	return branch.children
}
