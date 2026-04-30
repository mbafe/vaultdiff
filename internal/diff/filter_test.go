package diff

import (
	"testing"
)

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	changes := []Change{
		{Key: "foo", Type: ChangeAdded},
		{Key: "bar", Type: ChangeRemoved},
		{Key: "baz", Type: ChangeModified},
	}

	got := Filter(changes, FilterOptions{})
	if len(got) != len(changes) {
		t.Fatalf("expected %d changes, got %d", len(changes), len(got))
	}
}

func TestFilter_ByKeyPrefix(t *testing.T) {
	changes := []Change{
		{Key: "db_password", Type: ChangeModified},
		{Key: "db_user", Type: ChangeAdded},
		{Key: "api_key", Type: ChangeRemoved},
	}

	got := Filter(changes, FilterOptions{OnlyKeys: []string{"db_"}})
	if len(got) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(got))
	}
	for _, c := range got {
		if c.Key == "api_key" {
			t.Errorf("api_key should have been filtered out")
		}
	}
}

func TestFilter_ByType_Added(t *testing.T) {
	changes := []Change{
		{Key: "a", Type: ChangeAdded},
		{Key: "b", Type: ChangeRemoved},
		{Key: "c", Type: ChangeModified},
		{Key: "d", Type: ChangeAdded},
	}

	got := Filter(changes, FilterOptions{Types: []ChangeType{ChangeAdded}})
	if len(got) != 2 {
		t.Fatalf("expected 2 added changes, got %d", len(got))
	}
	for _, c := range got {
		if c.Type != ChangeAdded {
			t.Errorf("unexpected change type %v", c.Type)
		}
	}
}

func TestFilter_ByKeyAndType(t *testing.T) {
	changes := []Change{
		{Key: "db_password", Type: ChangeModified},
		{Key: "db_user", Type: ChangeAdded},
		{Key: "api_key", Type: ChangeModified},
	}

	got := Filter(changes, FilterOptions{
		OnlyKeys: []string{"db_"},
		Types:    []ChangeType{ChangeModified},
	})
	if len(got) != 1 {
		t.Fatalf("expected 1 change, got %d", len(got))
	}
	if got[0].Key != "db_password" {
		t.Errorf("expected db_password, got %s", got[0].Key)
	}
}

func TestFilter_EmptyChanges(t *testing.T) {
	got := Filter(nil, FilterOptions{OnlyKeys: []string{"foo"}})
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d", len(got))
	}
}
