package topo_test

import (
	"testing"

	"github.com/plombardi89/topo"
	"github.com/stretchr/testify/assert"
)

func TestGraph(t *testing.T) {
	g := topo.NewGraph()
	g.PutNode("a")
	g.PutNode("b")

	g.PutNodes("c", "d", "e", "f")

	if assert.Equal(t, []string{"a", "b", "c", "d", "e", "f"}, g.Nodes()) {
		assert.True(t, g.Contains("c"))
		assert.False(t, g.Contains("z"))
	}

	assert.True(t, g.RemoveNode("c"))
	assert.False(t, g.RemoveNode("c"))

	if assert.False(t, g.Contains("c")) {
		assert.Equal(t, []string{"a", "b", "d", "e", "f"}, g.Nodes())
	}
}

func TestGraph_Sort(t *testing.T) {
	g := topo.GraphFromMap(map[string][]string{
		"a": {"b", "c"},
		"b": {},
		"c": {"b"},
	})

	sorted, err := g.Sort()
	if assert.NoError(t, err) {
		assert.Equal(t, []string{"a", "c", "b"}, sorted)
	}

	t.Run("detect cycle", func(t *testing.T) {
		g = topo.GraphFromMap(map[string][]string{
			"a": {"b", "c"},
			"b": {},
			"c": {"b", "d"},
			"d": {"a"},
		})

		_, err = g.Sort()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "graph cycle involving nodes: [a, b, c, d]")
		}
	})
}
