package http

type Handler func(w Headers, r *Request)

type Route struct {
	Handler Handler
	name    string
}

type Router struct {
	routes         map[string]Route
	defaultHandler Handler
}

// NewRouter init and return the new router structure
func NewRouter() *Router {
	return &Router{
		routes:         map[string]Route{},
		defaultHandler: func(w Headers, r *Request) {},
	}
}

// AddRoute creates a new route repsonding to a url and a f function
func (r *Router) AddRoute(url string, f Handler) {
	r.routes[url] = Route{
		Handler: f,
		name:    url,
	}
}

// SetDefaultRoute set the handler function if the route doesn't exists
func (r *Router) SetDefaultRoute(f Handler) {
	r.defaultHandler = f
}
