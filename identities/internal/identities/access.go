package identities

type accessControl struct {
	g *graph
}

func newAccessControl(g *graph) *accessControl {
	return &accessControl{
		g: g,
	}
}

// HasResourceAccess
func (ac accessControl) HasResourceAccess(userHash, resourceHash string) bool {
	v, err := ac.g.getVertex(resourceHash)
	if err != nil {
		return false
	}

	for _, e := range v.edges {

		if e.relationship == "DENY" {
			return false
		}

		_, err := ac.g.getVertex(e.to)
		if err != nil {
			return false
		}

	}

	return true
}
