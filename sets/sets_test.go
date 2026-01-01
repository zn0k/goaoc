package sets

import "testing"

func TestItems(t *testing.T) {
	s := New[int]()
	items := s.Items()
	if len(items) != 0 {
		t.Errorf("length should be 0, is %d", len(items))
	}
}

func TestLen(t *testing.T) {
	s := New[int]()
	items := s.Items()
	if len(items) != s.Len() {
		t.Errorf("s.Len() should be 0, is %d", s.Len())
	}
}

func TestAdd(t *testing.T) {
	s := New[int]()
	s.Add(1)
	if s.Len() != 1 {
		t.Errorf("length should be 1, is %d", s.Len())
	}
	items := s.Items()
	if items[0] != 1 {
		t.Error("s.Items() did not contain 1")
	}
}

func TestTypes(t *testing.T) {
	s := New[int]()
	s.Add(1)
	if s.Len() != 1 {
		t.Errorf("length should be 1, is %d", s.Len())
	}
	items := s.Items()
	if items[0] != 1 {
		t.Error("s.Items() did not contain 1")
	}

	s2 := New[string]()
	s2.Add("test")
	if s2.Len() != 1 {
		t.Errorf("length should be 1, is %d", s2.Len())
	}
	items2 := s2.Items()
	if items2[0] != "test" {
		t.Error("s2.Items() did not contain 1")
	}
}

func TestMany(t *testing.T) {
	s := New[int]()
	s.AddMany([]int{1, 2, 3, 4, 5})
	if s.Len() != 5 {
		t.Errorf("length should be 5, is %d", s.Len())
	}
}

func TestRemove(t *testing.T) {
	s := New[int]()
	s.Add(1)
	s.Add(2)
	s.Remove(1)

	if s.Len() != 1 {
		t.Errorf("length should be 1, is %d", s.Len())
	}
	items := s.Items()
	if items[0] != 2 {
		t.Error("s.Items() did not contain 1")
	}
}

func TestIn(t *testing.T) {
	s := New[int]()
	s.Add(1)
	s.Add(2)

	t1 := s.In(1)
	t2 := s.In(2)
	t3 := s.In(3)

	if t1 != true {
		t.Errorf("s.In(1) does not return true")
	}
	if t2 != true {
		t.Errorf("s.In(2) does not return true")
	}
	if t3 != false {
		t.Errorf("s.In(3) does not return true")
	}
}

func TestNotIn(t *testing.T) {
	s := New[int]()
	s.Add(1)
	s.Add(2)

	t1 := s.NotIn(1)
	t2 := s.NotIn(2)
	t3 := s.NotIn(3)

	if t1 != false {
		t.Errorf("s.NotIn(1) does not return false")
	}
	if t2 != false {
		t.Errorf("s.NotIn(2) does not return false")
	}
	if t3 != true {
		t.Errorf("s.NotIn(3) does not return true")
	}
}

func TestMerge(t *testing.T) {
	s := New[int]()
	s.Add(1)
	s2 := New[int]()
	s2.Add(2)
	s.Merge(s2)

	if s.Len() != 2 {
		t.Errorf("length of s is not 2, got %d", s.Len())
	}
	if s.In(1) != true {
		t.Errorf("s does not contain 1")
	}
	if s.In(2) != true {
		t.Errorf("s does not contain 2")
	}
}

func TestUnion(t *testing.T) {
	s := New[int]()
	s.Add(1)
	s2 := New[int]()
	s2.Add(2)

	s3 := s.Union(s2)

	if s.Len() != 1 {
		t.Errorf("length of s is not 1, got %d", s.Len())
	}
	if s2.Len() != 1 {
		t.Errorf("length of s2 is not 1, got %d", s2.Len())
	}

	if s3.Len() != 2 {
		t.Errorf("length of s3 is not 2, got %d", s3.Len())
	}
	if s3.In(1) != true {
		t.Errorf("s3 does not contain 1")
	}
	if s3.In(2) != true {
		t.Errorf("s3 does not contain 2")
	}
}

func TestIntersection(t *testing.T) {
	s := New[int]()
	s.AddMany([]int{1, 2, 3})
	s2 := New[int]()
	s2.AddMany([]int{3, 4, 5})

	s3 := s.Intersection(s2)

	if s.Len() != 3 {
		t.Errorf("length of s is not 3, got %d", s.Len())
	}
	if s2.Len() != 3 {
		t.Errorf("length of s2 is not 3, got %d", s2.Len())
	}

	if s3.Len() != 1 {
		t.Errorf("length of s3 is not 1, got %d", s3.Len())
	}
	if s3.In(3) != true {
		t.Errorf("s3 does not contain 3")
	}
}

func TestSymmetric(t *testing.T) {
	s := New[int]()
	s.AddMany([]int{1, 2, 3})
	s2 := New[int]()
	s2.AddMany([]int{2, 3, 4})

	s3 := s.SymmetricDifference(s2)

	if s.Len() != 3 {
		t.Errorf("length of s is not 3, got %d", s.Len())
	}
	if s2.Len() != 3 {
		t.Errorf("length of s2 is not 3, got %d", s2.Len())
	}

	if s3.Len() != 2 {
		t.Errorf("length of s3 is not 1, got %d", s3.Len())
	}
	if s3.In(1) != true {
		t.Errorf("s3 does not contain 1")
	}
	if s3.In(4) != true {
		t.Errorf("s3 does not contain 4")
	}
}

func TestDisjoint(t *testing.T) {
	s := New[int]()
	s.AddMany([]int{1, 2, 3})
	s2 := New[int]()
	s2.AddMany([]int{3, 4, 5})
	s3 := New[int]()
	s3.AddMany([]int{4, 5, 6})

	if s.IsDisjoint(s2) == true {
		t.Errorf("s is claimed to be disjoint to s2, but is not")
	}
	if s.IsDisjoint(s3) != true {
		t.Errorf("s is claimed to not be disjoint to s3, but is")
	}
}

func TestSubset(t *testing.T) {
	s := New[int]()
	s.AddMany([]int{1, 2})
	s2 := New[int]()
	s2.AddMany([]int{1, 2, 3, 4, 5, 6})
	s3 := New[int]()
	s3.AddMany([]int{1, 3, 4, 5})

	if s.IsSubset(s2) != true {
		t.Errorf("s is claimed to not be a subset of s2, but is")
	}
	if s.IsDisjoint(s3) == true {
		t.Errorf("s is claimed to be a subset of s3, but is not")
	}
}

func TestSuperset(t *testing.T) {
	s := New[int]()
	s.AddMany([]int{1, 2, 3, 4, 5, 6})
	s2 := New[int]()
	s2.AddMany([]int{1, 2, 3})
	s3 := New[int]()
	s3.AddMany([]int{1, 3, 4, 7})

	if s.IsSuperset(s2) != true {
		t.Errorf("s is claimed to not be a superset of s2, but is")
	}
	if s.IsSuperset(s3) == true {
		t.Errorf("s is claimed to be a superset of s3, but is not")
	}
}

func TestCopy(t *testing.T) {
	s := New[int]()
	s.AddMany([]int{1, 2})

	c := s.Copy()
	c.Add(3)

	if s.Len() != 2 {
		t.Errorf("length of s is not 2, got %d", s.Len())
	}
	if c.Len() != 3 {
		t.Errorf("length of c is not 3, got %d", s.Len())
	}

	if c.In(1) != true {
		t.Errorf("c does not contain 1")
	}
	if c.In(2) != true {
		t.Errorf("c does not contain 2")
	}
	if c.In(3) != true {
		t.Errorf("c does not contain 3")
	}
}

func TestFromIterable(t *testing.T) {
	s := FromIterable([]int{1, 2, 3, 4, 5})
	if s.Len() != 5 {
		t.Errorf("length should be 5, is %d", s.Len())
	}
}
