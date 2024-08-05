//go:build wireinject
// +build wireinject

package internal

import (
	"manga/pkg/httpserver"
	"manga/pkg/logger"
	"manga/pkg/mongodb"

	"github.com/google/wire"
)

var deps = []interface{}{}

var providerSet wire.ProviderSet = wire.NewSet(
	mongodb.NewMongoDatabase,
	logger.New,
	httpserver.New,
)
