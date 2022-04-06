package main

import (
	opentracing4 "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v4"

	"context"
	"fmt"
	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/wsjcko/shopproduct/common"
	pb "github.com/wsjcko/shopproduct/protobuf/pb"
	"go-micro.dev/v4"
	log "go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
)

var (
	clientServiceName = "go.micro.service.shop.product.client"
	serverServiceName = "go.micro.service.shop.product"
)

func main() {
	//注册中心
	consul := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})
	//链路追踪
	t, io, err := common.NewTracer(clientServiceName, "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	service := micro.NewService(
		micro.Name(clientServiceName),
		micro.Version("latest"),
		micro.Address("127.0.0.1:8085"),
		//添加注册中心
		micro.Registry(consul),
		//绑定链路追踪 服务端绑定handle 客户端绑定client
		micro.WrapClient(opentracing4.NewClientWrapper(opentracing.GlobalTracer())),
	)

	// 访问服务器端服务
	productService := pb.NewProductService(serverServiceName, service.Client())

	productAdd := &pb.ProductInfo{
		ProductName:        "imooc",
		ProductSku:         "cap",
		ProductPrice:       1.1,
		ProductDescription: "imooc-cap",
		ProductCategoryId:  1,
		ProductImage: []*pb.ProductImage{
			{
				ImageName: "cap-image",
				ImageCode: "capimage01",
				ImageUrl:  "capimage01",
			},
			{
				ImageName: "cap-image02",
				ImageCode: "capimage02",
				ImageUrl:  "capimage02",
			},
		},
		ProductSize: []*pb.ProductSize{
			{
				SizeName: "cap-size",
				SizeCode: "cap-size-code",
			},
		},
		ProductSeo: &pb.ProductSeo{
			SeoTitle:       "cap-seo",
			SeoKeywords:    "cap-seo",
			SeoDescription: "cap-seo",
			SeoCode:        "cap-seo",
		},
	}
	response, err := productService.AddProduct(context.TODO(), productAdd)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
}
