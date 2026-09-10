package xfer

// StaticRegistry stores static XFER decoder mappings for Nox objects. See DefaultRegistry.
type StaticRegistry struct {
	ByType map[string]Type
	ByID   map[int]Type
}

func (r *StaticRegistry) XferByObjectType(typ string) Type {
	return r.ByType[typ]
}

func (r *StaticRegistry) XferByObjectTypeID(id int) Type {
	return r.ByID[id]
}
