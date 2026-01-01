package sets

type Set[T comparable] struct {
	m map[T]struct{}
}

func New[T comparable]() Set[T] {
	s := Set[T]{}
	s.m = make(map[T]struct{})
	return s
}

func FromIterable[T comparable](items []T) Set[T] {
	s := New[T]()
	s.AddMany(items)
	return s
}

func (s *Set[T]) Add(item T) {
	s.m[item] = struct{}{}
}

func (s *Set[T]) AddMany(items []T) {
	for _, item := range items {
		s.m[item] = struct{}{}
	}
}

func (s *Set[T]) Remove(item T) {
	delete(s.m, item)
}

func (s *Set[T]) Items() []T {
	var items []T
	for key := range s.m {
		items = append(items, key)
	}
	return items
}

func (s *Set[T]) Len() int {
	return len(s.m)
}

func (s *Set[T]) In(item T) bool {
	_, ok := s.m[item]
	return ok
}

func (s *Set[T]) NotIn(item T) bool {
	_, ok := s.m[item]
	return !ok
}

func (s *Set[T]) IsDisjoint(other Set[T]) bool {
	intersection := s.Intersection(other)
	return intersection.Len() == 0
}

func (s *Set[T]) IsSubset(other Set[T]) bool {
	result := true
	for key := range s.m {
		_, ok := other.m[key]
		if !ok {
			result = false
			break
		}
	}
	return result
}

func (s *Set[T]) IsSuperset(other Set[T]) bool {
	result := true
	for key := range other.m {
		_, ok := s.m[key]
		if !ok {
			result = false
			break
		}
	}
	return result
}

func (s *Set[T]) Merge(other Set[T]) {
	for key := range other.m {
		s.m[key] = struct{}{}
	}
}

func (s *Set[T]) Union(others ...Set[T]) Set[T] {
	union := New[T]()
	for key := range s.m {
		union.m[key] = struct{}{}
	}
	for _, other := range others {
		for key := range other.m {
			union.m[key] = struct{}{}
		}
	}
	return union
}

func (s *Set[T]) Intersection(others ...Set[T]) Set[T] {
	intersection := New[T]()
	for key := range s.m {
		match := true
		for _, other := range others {
			_, ok := other.m[key]
			if !ok {
				match = false
				break
			}
		}
		if match {
			intersection.m[key] = struct{}{}
		}
	}
	return intersection
}

func (s *Set[T]) Difference(others ...Set[T]) Set[T] {
	difference := New[T]()
	for key := range s.m {
		match := true
		for _, other := range others {
			_, ok := other.m[key]
			if ok {
				match = false
				break
			}
		}
		if match {
			difference.m[key] = struct{}{}
		}
	}
	return difference
}

func (s *Set[T]) SymmetricDifference(others ...Set[T]) Set[T] {
	symmetric := s.Difference(others...)
	right := New[T]()
	right = right.Union(others...)
	symmetric.Merge(right.Difference(*s))
	return symmetric
}

func (s *Set[T]) Copy() Set[T] {
	copy := New[T]()
	copy.Merge(*s)
	return copy
}
