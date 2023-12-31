package services

import (
	"errors"

	"github.com/sean-david-welch/farmec-v2/server/models"
	"github.com/sean-david-welch/farmec-v2/server/repository"
	"github.com/sean-david-welch/farmec-v2/server/types"
	"github.com/sean-david-welch/farmec-v2/server/utils"
)

type SupplierService struct {
	folder string
	s3Client *utils.S3Client
    repository *repository.SupplierRepository
}

func NewSupplierService(repository *repository.SupplierRepository, s3Client *utils.S3Client, folder string) *SupplierService {
    return &SupplierService{
		repository: repository, 
		s3Client: s3Client, 
		folder: folder,
	}
}

func (service *SupplierService) GetSuppliers() ([]models.Supplier, error) {
    return service.repository.GetSuppliers()
}

func (service *SupplierService) CreateSupplier(supplier *models.Supplier) (*types.SupplierResult, error) {	
	logoImage := supplier.LogoImage
	marketingImage := supplier.MarketingImage

	if logoImage == "" || marketingImage == "" {
		return nil, errors.New("logoImage or marketingImage is empty")
	}

	presignedLogo, logoUrl, err := service.s3Client.GeneratePresignedUrl(service.folder, logoImage)
	if err != nil {
        return nil, err
    }

	presignedMarketing, marketingUrl, err := service.s3Client.GeneratePresignedUrl(service.folder, marketingImage)
	if err != nil {
		return nil, err
	}

	supplier.LogoImage = logoUrl
    supplier.MarketingImage = marketingUrl

    service.repository.CreateSupplier(supplier); if err != nil {
		return nil, err
	}
	
	result := &types.SupplierResult{
		PresignedLogoUrl: presignedLogo,
		LogoUrl: logoUrl,
		PresginedMarketingUrl: presignedMarketing,
		MarketingUrl: marketingUrl,
	}

	return result, nil
}

func (service *SupplierService) GetSupplierById(id string) (*models.Supplier, error) {
    return service.repository.GetSupplierById(id)
}

func (service *SupplierService) UpdateSupplier(id string, supplier *models.Supplier) (*types.SupplierResult, error) {
	logoImage := supplier.LogoImage
	marketingImage := supplier.MarketingImage

	var presignedLogo, logoUrl, presignedMarketing, marketingUrl string
    var err error

    if logoImage != "" {
        presignedLogo, logoUrl, err = service.s3Client.GeneratePresignedUrl(service.folder, logoImage)
        if err != nil {
            return nil, err
        }
        supplier.LogoImage = logoUrl
    }

    if marketingImage != "" {
        presignedMarketing, marketingUrl, err = service.s3Client.GeneratePresignedUrl(service.folder, marketingImage)
        if err != nil {
            return nil, err
        }
        supplier.MarketingImage = marketingUrl
    }


    service.repository.UpdateSupplier(id, supplier); if err != nil {
		return nil, err
	}

	result := &types.SupplierResult{
		PresignedLogoUrl: presignedLogo,
		LogoUrl: logoUrl,
		PresginedMarketingUrl: presignedMarketing,
		MarketingUrl: marketingUrl,
	}
	
	return result, nil
}

func (service *SupplierService) DeleteSupplier(id string) (error) {
	supplier, err := service.repository.GetSupplierById(id); if err != nil {
        return err
    }

	if err := service.s3Client.DeleteImageFromS3(supplier.LogoImage); err != nil {
		return err
	}

	if err := service.s3Client.DeleteImageFromS3(supplier.MarketingImage); err != nil {
		return err
	}

    if err := service.repository.DeleteSupplier(id); err != nil {
		return err
	}

	return nil
}