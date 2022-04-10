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
	clientServiceName    = "go.micro.service.shop.product.client"
	serverServiceName    = "go.micro.service.shop.product"
	MICRO_CONSUL_ADDRESS = "127.0.0.1:8500"
	MICRO_JAEGER_ADDRESS = "127.0.0.1:6831"
	MICRO_ADDRESS        = "127.0.0.1:8005"
)

func main() {
	//注册中心
	consul := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			MICRO_CONSUL_ADDRESS,
		}
	})
	//链路追踪
	t, io, err := common.NewTracer(clientServiceName, MICRO_JAEGER_ADDRESS)
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	service := micro.NewService(
		micro.Name(clientServiceName),
		micro.Version("latest"),
		micro.Address(MICRO_ADDRESS),
		//添加注册中心
		micro.Registry(consul),
		//绑定链路追踪 服务端绑定handle 客户端绑定client
		micro.WrapClient(opentracing4.NewClientWrapper(opentracing.GlobalTracer())),
	)

	// 访问服务器端服务
	productService := pb.NewShopProductService(serverServiceName, service.Client())

	productAdd := &pb.ProductInfo{
		ProductName:        "imooc",
		ProductSku:         "cap121",
		ProductPrice:       1.1,
		ProductDescription: "imooc-cap",
		ProductCategoryId:  1,
		ProductImage: []*pb.ProductImage{
			{
				ImageName: "cap-image2",
				ImageCode: "capimage022",
				ImageUrl:  "capimage02",
			},
			{
				ImageName: "cap-image03",
				ImageCode: "capimage03",
				ImageUrl:  "capimage03",
			},
		},
		ProductSize: []*pb.ProductSize{
			{
				SizeName: "cap-size11",
				SizeCode: "cap-size-code11",
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
