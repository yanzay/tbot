package tbot

import "strings"

const (
	RouteBack    = "<..>"
	RouteRoot    = ""
	RouteRefresh = "<.>"
)

type Node struct {
	route    string
	handler  *Handler
	parent   *Node
	children map[string]*Node
}

func NewNode(parent *Node, route string, handler *Handler) *Node {
	return &Node{
		route:    route,
		handler:  handler,
		children: make(map[string]*Node),
		parent:   parent,
	}
}

// RouterMux is a tree-route multiplexer
type RouterMux struct {
	fileHandler    *Handler
	defaultHandler *Handler
	storage        SessionStorage
	aliases        map[string]string
	root           *Node
}

// NewRouterMux creates new RouterMux
// Takes SessionStorage to store users' sessions state
func NewRouterMux(storage SessionStorage) Mux {
	return &RouterMux{
		storage: storage,
		aliases: make(map[string]string),
		root:    NewNode(nil, RouteRoot, nil),
	}
}

// Handlers returns list of handlers currently presented in mux
func (rm *RouterMux) Handlers() Handlers {
	return Handlers{}
}

// DefaultHandler returns default handler, nil if it's not set
func (rm *RouterMux) DefaultHandler() *Handler {
	return rm.defaultHandler
}

func (rm *RouterMux) FileHandler() *Handler {
	return rm.fileHandler
}

// Mux takes message content and returns corresponding handler
func (rm *RouterMux) Mux(msg *Message) (*Handler, MessageVars) {
	var node *Node
	route := msg.Data
	if _, ok := rm.aliases[route]; ok {
		route = rm.aliases[route]
	}
	if route == RouteRoot {
		node = rm.root
	} else {
		state := rm.storage.Get(msg.ChatID)
		if state == "" {
			state = RouteRoot
		}
		node = rm.findNode(state)
		if node == nil {
			return nil, nil
		}
		switch route {
		case RouteBack:
			node = node.parent
		case RouteRefresh:
		default:
			if child, ok := node.children[route]; ok {
				node = child
			} else {
				return nil, nil
			}
		}
	}
	rm.storage.Set(msg.ChatID, nodeToState(node))
	return node.handler, MessageVars{}
}

func (rm *RouterMux) findNode(path string) *Node {
	node := rm.root
	if path == RouteRoot {
		return node
	}
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return nil
	}
	parts = parts[1:]
	for _, part := range parts {
		node = node.children[part]
		if node == nil {
			return nil
		}
	}
	return node
}

func nodeToState(node *Node) string {
	routes := make([]string, 0)
	if len(node.children) > 0 {
		routes = []string{node.route}
	}
	for node.parent != nil {
		node = node.parent
		routes = append([]string{node.route}, routes...)
	}
	return strings.Join(routes, "/")
}

func back(route string) string {
	routes := strings.Split(route, "/")
	return strings.Join(routes[:len(routes)-1], "/")
}

// HandleFunc adds new handler function to mux.
func (rm *RouterMux) HandleFunc(path string, handler HandlerFunction, description ...string) {
	if path == RouteRoot {
		rm.root.handler = NewHandler(handler, path, description...)
		return
	}
	parts := strings.Split(path, "/")
	node := rm.root
	for _, part := range parts {
		if part == "" {
			continue
		}
		if node.children[part] == nil {
			node.children[part] = NewNode(node, part, NewHandler(handler, path, description...))
		}
		node = node.children[part]
	}
}

// SetAlias sets aliases for specified route.
func (rm *RouterMux) SetAlias(route string, aliases ...string) {
	for _, alias := range aliases {
		rm.aliases[alias] = route
	}
}

// HandleDefault adds new default handler, when nothing matches with message,
func (rm *RouterMux) HandleDefault(handler HandlerFunction, description ...string) {
	rm.defaultHandler = NewHandler(handler, "", description...)
}

func (rm *RouterMux) HandleFile(handler HandlerFunction, description ...string) {
}

func (rm *RouterMux) Reset(chatID int64) {
	rm.storage.Reset(chatID)
}
