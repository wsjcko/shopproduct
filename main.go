package main

import (
	opentracing4 "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v4"

	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/opentracing/opentracing-go"
	"github.com/wsjcko/shopproduct/common"
	"github.com/wsjcko/shopproduct/domain/repository"
	"github.com/wsjcko/shopproduct/domain/service"
	"github.com/wsjcko/shopproduct/handler"
	pb "github.com/wsjcko/shopproduct/protobuf/pb"
	"go-micro.dev/v4"
	log "go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
)

var (
	serviceName = "go.micro.service.shop.product"
	version     = "latest"
	address     = "127.0.0.1:8085"
)

func main() {
	//配置中心
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		log.Fatal(err)
	}
	//注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	//链路追踪
	t, io, err := common.NewTracer(serviceName, "127.0.0.1:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	//数据库设置
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")
	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+"@/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//禁止副表 gorm默认使用复数映射，go代码的单数、复数struct形式都匹配到复数表中,开启后，将严格匹配，遵守单数形式
	db.SingularTable(true)

	// 设置服务
	srv := micro.NewService(
		micro.Name(serviceName),
		micro.Version(version),
		//这里设置地址和需要暴露的端口
		micro.Address(address),
		//添加注册中心
		micro.Registry(consulRegistry),
		//绑定链路追踪 服务端绑定handle 客户端绑定client
		micro.WrapHandler(opentracing4.NewHandlerWrapper(opentracing.GlobalTracer())),
	)

	//初始化建表，多个表
	// repository.NewProductRepository(db).InitTable()

	productDataService := service.NewProductDataService(repository.NewProductRepository(db))

	// Initialise service
	srv.Init()

	// Register Handler
	pb.RegisterProductHandler(srv.Server(), &handler.Product{ProductDataService: productDataService})

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
