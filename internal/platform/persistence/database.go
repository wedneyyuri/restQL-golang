package persistence

import (
	"context"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/pkg/errors"
)

// Database defines the operations needed to fetch mappings
// and query from an external store.
type Database interface {
	FindMappingsForTenant(ctx context.Context, tenantID string) ([]restql.Mapping, error)
	FindQuery(ctx context.Context, namespace string, name string, revision int) (restql.SavedQuery, error)
}

// NewDatabase constructs a Database compliant value
// from the database plugin registered.
// In case of no plugin, a noop implementation is returned.
func NewDatabase(log restql.Logger) (Database, error) {
	pluginInfo, found := restql.GetDatabasePlugin()
	if !found {
		log.Info("no database plugin provided")
		return noOpDatabase{}, nil
	}

	dbPlugin, err := pluginInfo.New(log)
	if err != nil {
		return noOpDatabase{}, err
	}

	if dbPlugin == nil {
		log.Info("empty database instance returned by plugin", "plugin", pluginInfo.Name)
		return noOpDatabase{}, nil
	}

	database, ok := dbPlugin.(restql.DatabasePlugin)
	if !ok {
		return noOpDatabase{}, errors.Errorf("failed to cast database plugin, unknown type: %T", dbPlugin)
	}

	return database, nil
}

var errNoDatabase = errors.New("no op database")

type noOpDatabase struct{}

func (n noOpDatabase) FindMappingsForTenant(ctx context.Context, tenantID string) ([]restql.Mapping, error) {
	return []restql.Mapping{}, errNoDatabase
}

func (n noOpDatabase) FindQuery(ctx context.Context, namespace string, name string, revision int) (restql.SavedQuery, error) {
	return restql.SavedQuery{}, errNoDatabase
}
