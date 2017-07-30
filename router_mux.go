package tbot

import (
	"regexp"
	"strings"
)

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
	var messageData map[string]string = nil
	fields := strings.Fields(msg.Data)
	if len(fields) >= 1 {
		if alias, ok := rm.aliases[fields[0]]; ok {
			fields[0] = alias
		}
	}
	path := strings.Join(fields, " ")
	state := rm.storage.Get(msg.ChatID)
	if state == "" {
		state = RouteRoot
	}
	node = rm.findNode(state)
	if node == nil {
		return nil, nil
	}
switch_route:
	switch path {
	case RouteRoot:
		node = rm.root
	case RouteBack:
		if node.parent == nil {
			node = rm.root
		} else {
			node = node.parent
		}
	case RouteRefresh:
	default:
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		path = state + path
		for _, child := range node.children {
			re := regexp.MustCompile(child.handler.pattern)
			matches := re.FindStringSubmatch(path)
			if len(matches) > 0 {
				node = child
				matches = matches[1:]
				if len(matches) > 0 {
					messageData = make(map[string]string)
					for i, match := range matches {
						messageData[child.handler.variables[i]] = match
					}
				}
				break switch_route
			}
		}
		return rm.defaultHandler, nil
	}

	rm.storage.Set(msg.ChatID, nodeToState(node))
	return node.handler, messageData
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
