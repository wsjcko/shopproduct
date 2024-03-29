package handler

import (
	"context"
	"fmt"
	"github.com/wsjcko/shopproduct/common"
	"github.com/wsjcko/shopproduct/domain/model"
	"github.com/wsjcko/shopproduct/domain/service"
	pb "github.com/wsjcko/shopproduct/protobuf/pb"
)

type ShopProduct struct {
	ProductService service.IProductService
}

// AddProduct 添加商品
func (h *ShopProduct) AddProduct(ctx context.Context, request *pb.ProductInfo, response *pb.ResponseProduct) error {
	productAdd := &model.Product{}
	fmt.Println(request)
	if err := common.SwapTo(request, productAdd); err != nil {
		return err
	}
	fmt.Println(productAdd)
	productID, err := h.ProductService.AddProduct(productAdd)
	if err != nil {
		return err
	}
	response.ProductId = productID
	return nil
}

// FindProductByID 根据ID查找商品
func (h *ShopProduct) FindProductByID(ctx context.Context, request *pb.RequestID, response *pb.ProductInfo) error {
	productData, err := h.ProductService.FindProductByID(request.ProductId)
	if err != nil {
		return err
	}
	if err := common.SwapTo(productData, response); err != nil {
		return err
	}
	return nil
}

// UpdateProduct 商品更新
func (h *ShopProduct) UpdateProduct(ctx context.Context, request *pb.ProductInfo, response *pb.Response) error {
	productAdd := &model.Product{}
	if err := common.SwapTo(request, productAdd); err != nil {
		return err
	}
	err := h.ProductService.UpdateProduct(productAdd)
	if err != nil {
		return err
	}
	response.Msg = "更新成功"
	return nil
}

// DeleteProductByID 根据ID删除对应商品
func (h *ShopProduct) DeleteProductByID(ctx context.Context, request *pb.RequestID, response *pb.Response) error {
	if err := h.ProductService.DeleteProduct(request.ProductId); err != nil {
		return err
	}
	response.Msg = "删除成功"
	return nil
}

// FindAllProduct 查找所有商品
func (h *ShopProduct) FindAllProduct(ctx context.Context, request *pb.RequestAll, response *pb.AllProduct) error {
	productAll, err := h.ProductService.FindAllProduct()
	if err != nil {
		return err
	}

	for _, v := range productAll {
		productInfo := &pb.ProductInfo{}
		err := common.SwapTo(v, productInfo)
		if err != nil {
			return err
		}
		response.ProductInfo = append(response.ProductInfo, productInfo)
	}
	return nil
}
