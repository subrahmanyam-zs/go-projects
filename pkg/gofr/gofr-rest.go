package gofr

type RestReader interface {
	Read(c *Context) (interface{}, error)
}
type RestIndexer interface {
	Index(c *Context) (interface{}, error)
}
type RestCreator interface {
	Create(c *Context) (interface{}, error)
}
type RestUpdater interface {
	Update(c *Context) (interface{}, error)
}
type RestDeleter interface {
	Delete(c *Context) (interface{}, error)
}
type RestPatcher interface {
	Patch(c *Context) (interface{}, error)
}

// REST method adds REST-like routes if the interfaces are satisfied
func (k *Gofr) REST(entity string, handler interface{}) {
	if c, ok := handler.(RestIndexer); ok {
		k.GET("/"+entity, c.Index)
	}

	if c, ok := handler.(RestReader); ok {
		k.GET("/"+entity+"/{id}", c.Read)
	}

	if c, ok := handler.(RestCreator); ok {
		k.POST("/"+entity, c.Create)
	}

	if c, ok := handler.(RestDeleter); ok {
		k.DELETE("/"+entity+"/{id}", c.Delete)
	}

	if c, ok := handler.(RestUpdater); ok {
		k.PUT("/"+entity+"/{id}", c.Update)
	}

	if c, ok := handler.(RestPatcher); ok {
		k.PATCH("/"+entity+"/{id}", c.Patch)
	}
}
