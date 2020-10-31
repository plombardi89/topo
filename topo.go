package topo

import (
	"fmt"
	"sort"
	"strings"
)

type node struct {
	id          string
	connections []string
}

func (n *node) Connect(node string) {
	if n.connections == nil {
		n.connections = make([]string, 0)
	}

	n.connections = append(n.connections, node)
}

func (n *node) Connections() []string {
	return n.connections
}

func (n *node) Disconnect(id string) {
	n.connections = removeIfContains(id, n.connections)
}

func (n *node) Connected(id string) bool {
	for _, c := range n.Connections() {
		if id == c {
			return true
		}
	}

	return false
}

type Graph interface {
	fmt.Stringer
	GetNode(id string) ([]string, bool)
	Nodes() []string
	PutNode(id string)
	PutNodes(id ...string)
	RemoveNode(id string) bool
	Connect(src, dst string)
	Contains(id string) bool
	Sort() ([]string, error)
}

func NewGraph() Graph {
	return &graph{nodes: make(map[string]*node)}
}

func GraphFromMap(m map[string][]string) Graph {
	nodes := make(map[string]*node)

	for k, v := range m {
		if v == nil {
			v = make([]string, 0)
		}

		nodes[k] = &node{id: k, connections: v}
	}

	return &graph{nodes: nodes}
}

type graph struct {
	nodes map[string]*node
}

func (g *graph) GetNode(id string) ([]string, bool) {
	v, ok := g.nodes[id]
	if !ok {
		return nil, false
	}

	return v.connections, true
}

func (g *graph) Nodes() []string {
	var res []string

	for k := range g.nodes {
		res = append(res, k)
	}

	sort.Slice(res, func(i, j int) bool {
		a := res[i]
		b := res[j]

		return a < b
	})

	return res
}

func (g *graph) PutNode(id string) {
	g.nodes[id] = &node{id: id, connections: make([]string, 0)}
}

func (g *graph) PutNodes(ids ...string) {
	for _, id := range ids {
		g.PutNode(id)
	}
}

func (g *graph) RemoveNode(id string) bool {
	if _, ok := g.GetNode(id); !ok {
		return false
	}

	for _, v := range g.Nodes() {
		node, _ := g.nodes[v]
		node.Disconnect(id)
	}

	delete(g.nodes, id)

	return true
}

func (g *graph) Connect(src, dst string) {
	if src == dst {
		return
	}

	if !g.Contains(src) {
		g.PutNode(src)
	}

	if !g.Contains(dst) {
		g.PutNode(dst)
	}

	v := g.nodes[src]
	v.Connect(dst)
}

func (g *graph) Contains(id string) bool {
	_, ok := g.nodes[id]
	return ok
}

func (g *graph) Sort() ([]string, error) {
	var sorted []string

	inboundConnections := make(map[string]int)

	for _, nodeID := range g.Nodes() {
		node := g.nodes[nodeID]
		if _, ok := inboundConnections[node.id]; !ok {
			inboundConnections[node.id] = 0
		}

		for _, c := range node.Connections() {
			if _, ok := inboundConnections[node.id]; !ok {
				inboundConnections[c] = 1
			} else {
				inboundConnections[c]++
			}
		}
	}

	var stack []string
	for k, v := range inboundConnections {
		if v == 0 {
			stack = append(stack, k)
			inboundConnections[k] = -1
		}
	}

	for len(stack) > 0 {
		//var vtxID string
		nodeID := stack[len(stack) - 1]
		stack = stack[:len(stack) - 1]

		node := g.nodes[nodeID]

		for _, c := range node.Connections() {
			inboundConnections[c]--
			if inboundConnections[c] == 0 {
				stack = append(stack, c)
				inboundConnections[c] = -1
			}
		}

		sorted = append(sorted, nodeID)
	}

	if len(g.nodes) != len(sorted) {
		var cycle []string
		for k, v := range inboundConnections {
			if v > 0 {
				cycle = append(cycle, k)
			}
		}

		sort.Slice(cycle, func(i, j int) bool {
			return cycle[i] < cycle[j]
		})

		return nil, fmt.Errorf("graph cycle involving nodes: [%s]", strings.Join(cycle, ", "))
	}

	return sorted, nil
}

func (g *graph) String() string {
	var res string

	for _, nodeID := range g.Nodes() {
		node := g.nodes[nodeID]
		res += fmt.Sprintf("%s -> [%s]\n", node.id, strings.Join(node.Connections(), ", "))
	}

	return res
}

func removeIfContains(needle string, haystack []string) []string {
	var res []string

	for _, v := range haystack {
		if v == needle {
			continue
		}

		res = append(res, v)
	}

	return res
}