package app

import (
	"atomisu.com/ocg-statics/cardApiCall/internal/config"
	appHttp "atomisu.com/ocg-statics/cardApiCall/internal/http"
	"atomisu.com/ocg-statics/cardApiCall/internal/storage"
	"atomisu.com/ocg-statics/cardApiCall/internal/usecase/syncimage"
	_ "github.com/lib/pq"
	"github.com/samber/do"
)

// SetupDIContainer は syncStorageImageWithDB Lambda 用の DI コンテナをセットアップします。
func SetupDIContainer() *do.Injector {
	injector := do.New()
	do.Provide(injector, config.NewDbConnection)
	do.Provide(injector, storage.NewS3StorageClientForDI)
	do.Provide(injector, appHttp.NewStreamClientForDI)
	do.Provide(injector, syncimage.NewSyncStorageImageUseCaseForDI)
	return injector
}
