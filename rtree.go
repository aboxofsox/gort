package gort

import (
	"strings"
)

type rtree struct {
	root *rnode
}

type rnode struct {
	children     map[string]*rnode // children is a map that stores the child nodes of the current node.
	route        *Route            // route is a pointer to the Route associated with the current node.
	isLast       bool              // isLast indicates whether the current node is the last node in a route.
	dynamicChild *rnode            // dynamicChild is a pointer to the dynamic child node of the current node.
	isDynamic    bool              // isDynamic indicates whether the current node is a dynamic node.
}

func newRTree() *rtree {
	return &rtree{
		root: &rnode{
			children: make(map[string]*rnode),
		},
	}
}

// add adds a new route to the rtree.
// It takes a pointer to a Route as input and adds it to the rtree.
// The route's pattern is split into parts using "/" as the delimiter.
// Each part is then used to traverse the rtree and create or update the corresponding nodes.
// If a part is empty, it is skipped.
// If a node for a part does not exist, a new node is created and added to the current node's children.
// If the part is a dynamic part (starts with ":"), the current node's dynamicChild is updated.
// Finally, the last node in the traversal is marked as the last node and its route is set to the input route.
func (t *rtree) add(r *Route) {
	pattern := r.Pattern
	current := t.root
	parts := strings.Split(pattern, "/")[1:]

	for _, part := range parts {
		if part == "" {
			continue
		}

		if _, ok := current.children[part]; !ok {
			newNode := &rnode{
				children:  make(map[string]*rnode),
				isDynamic: strings.HasPrefix(part, ":"),
			}
			current.children[part] = newNode
			if newNode.isDynamic {
				current.dynamicChild = newNode
			}
		}

		current = current.children[part]
	}

	current.isLast = true
	current.route = r
}

// find searches for a route in the rtree based on the given path.
// It returns the corresponding Route if found, otherwise it returns nil.
func (t *rtree) find(path string) *Route {
	current := t.root
	parts := strings.Split(path, "/")[1:]

	for _, part := range parts {
		if part == "" {
			continue
		}

		if next, ok := current.children[part]; ok {
			current = next
		} else if current.dynamicChild != nil {
			current = current.dynamicChild
		} else {
			return nil
		}
	}

	if !current.isLast {
		return nil
	}

	return current.route
}
