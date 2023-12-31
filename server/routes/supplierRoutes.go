package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/sean-david-welch/farmec-v2/server/controllers"
	"github.com/sean-david-welch/farmec-v2/server/lib"
	"github.com/sean-david-welch/farmec-v2/server/middleware"
	"github.com/sean-david-welch/farmec-v2/server/repository"
	"github.com/sean-david-welch/farmec-v2/server/services"
	"github.com/sean-david-welch/farmec-v2/server/utils"
)

func InitializeSuppliers(router *gin.Engine, db *sql.DB, s3Client *utils.S3Client, firebase *lib.Firebase) {
    supplierRepository := repository.NewSupplierRepository(db)
    supplierService := services.NewSupplierService(supplierRepository, s3Client, "suppliers")
    supplierController := controllers.NewSupplierContoller(supplierService)

	adminMiddleware := middleware.NewAdminMiddleware(firebase)

    SupplierRoutes(router, supplierController, adminMiddleware)
}

func SupplierRoutes(router *gin.Engine, supplierController *controllers.SupplierController, adminMiddleware *middleware.AdminMiddleware) {
	supplierGroup := router.Group("/api/suppliers")

	supplierGroup.GET("", supplierController.GetSuppliers)
	supplierGroup.GET("/:id", supplierController.GetSupplierByID) 
	
	protected := supplierGroup.Group("/")
	protected.Use(adminMiddleware.Middleware())
	{
		protected.POST("", supplierController.CreateSupplier)
		protected.PUT("/:id", supplierController.UpdateSupplier) 
		protected.DELETE("/:id", supplierController.DeleteSupplier) 	
	}
}