package catalog

import (
	"cafestartups/internal/domain"
)

// Loader validates fixture data at the application boundary and returns the
// domain card catalog. Keeping fixture loading here prevents the HTTP layer
// from knowing how catalog data is encoded.
type Loader struct{}

func (Loader) Load(data []byte) ([]domain.Card, error) {
	return domain.LoadCatalog(data)
}

func Load(data []byte) ([]domain.Card, error) {
	return Loader{}.Load(data)
}
