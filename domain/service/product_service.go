package service

import (
	"github.com/wsjcko/shopproduct/domain/model"
	"github.com/wsjcko/shopproduct/domain/repository"
)

type IProductService interface {
	AddProduct(*model.Product) (int64, error)
	DeleteProduct(int64) error
	UpdateProduct(*model.Product) error
	FindProductByID(int64) (*model.Product, error)
	FindAllProduct() ([]model.Product, error)
}

// NewProductService 创建
func NewProductService(productRepository repository.IProductRepository) IProductService {
	return &ProductService{productRepository}
}

type ProductService struct {
	ProductRepository repository.IProductRepository
}

// AddProduct 插入
func (u *ProductService) AddProduct(product *model.Product) (int64, error) {
	return u.ProductRepository.CreateProduct(product)
}

// DeleteProduct 删除
func (u *ProductService) DeleteProduct(productID int64) error {
	return u.ProductRepository.DeleteProductByID(productID)
}

// UpdateProduct 更新
func (u *ProductService) UpdateProduct(product *model.Product) error {
	return u.ProductRepository.UpdateProduct(product)
}

// FindProductByID 查找
func (u *ProductService) FindProductByID(productID int64) (*model.Product, error) {
	return u.ProductRepository.FindProductByID(productID)
}

// FindAllProduct 查找
func (u *ProductService) FindAllProduct() ([]model.Product, error) {
	return u.ProductRepository.FindAll()
}
