package rainbow

import (
	"sort"
	"testing"
)

func TestReduce(t *testing.T) {
	hash := [20]byte{1, 0, 0, 0}
	res := reduce(hash, 1, 5)
	if res != 10006 {
		t.Errorf("expected 10006, got %d", res)
	}
}

func TestByEndpointSort(t *testing.T) {
	chains := []Chain{
		{Endpoint: 5, Startpoint: 1},
		{Endpoint: 2, Startpoint: 1},
		{Endpoint: 8, Startpoint: 1},
	}
	sort.Sort(ByEndpoint(chains))
	if chains[0].Endpoint != 2 {
		t.Errorf("expected 2, got %d", chains[0].Endpoint)
	}
	if chains[1].Endpoint != 5 {
		t.Errorf("expected 5, got %d", chains[1].Endpoint)
	}
	if chains[2].Endpoint != 8 {
		t.Errorf("expected 8, got %d", chains[2].Endpoint)
	}
}
