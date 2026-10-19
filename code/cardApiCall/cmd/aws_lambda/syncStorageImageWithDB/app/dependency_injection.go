package app

import (
	"atomisu.com/ocg-statics/cardApiCall/internal/config"
	"atomisu.com/ocg-statics/cardApiCall/internal/usecase/syncimage"
	_ "github.com/lib/pq"
	"github.com/samber/do"
)

func SetupDIContainer() *do.Injector {
	injector := do.New()
	do.Provide(injector, config.NewDbConnection)
	do.Provide(injector, syncimage.NewSyncStorageImageUseCase)
	return injector
}
