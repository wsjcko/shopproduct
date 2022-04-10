package main

import (
	ratelimit4 "github.com/asim/go-micro/plugins/wrapper/ratelimiter/uber/v4"
	opentracing4 "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v4"

	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/opentracing/opentracing-go"
	cli2 "github.com/urfave/cli/v2"
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
	MICRO_SERVICE_NAME   = "go.micro.service.shop.product"
	MICRO_VERSION        = "latest"
	MICRO_ADDRESS        = "127.0.0.1:8085"
	DOCKER_HOST          = "127.0.0.1"
	MICRO_CONSUL_ADDRESS = "127.0.0.1:8500"
	MICRO_JAEGER_ADDRESS = "127.0.0.1:6831"
	MICRO_QPS            = 100
)

func SetDockerHost(host string) {
	DOCKER_HOST = host
	MICRO_CONSUL_ADDRESS = host + ":8500"
	MICRO_JAEGER_ADDRESS = host + ":6831"
}

func main() {

	function := micro.NewFunction(
		micro.Flags(
			&cli2.StringFlag{ //micro 多个选项 --ip
				Name:  "ip",
				Usage: "docker Host IP(ubuntu)",
				Value: "0.0.0.0",
			},
		),
	)

	function.Init(
		micro.Action(func(c *cli2.Context) error {
			ipstr := c.Value("ip").(string)
			if len(ipstr) > 0 { //后续校验IP
				log.Info("docker Host IP(ubuntu)1111", ipstr)
			}
			SetDockerHost(ipstr)
			return nil
		}),
	)

	log.Info("DOCKER_HOST ", DOCKER_HOST)

	//配置中心
	consulConfig, err := common.GetConsulConfig(MICRO_CONSUL_ADDRESS, "/micro/config")
	if err != nil {
		log.Fatal(err)
	}
	//注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			MICRO_CONSUL_ADDRESS,
		}
	})

	//链路追踪
	t, io, err := common.NewTracer(MICRO_SERVICE_NAME, MICRO_JAEGER_ADDRESS)
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	//数据库设置
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")
	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+"@tcp("+mysqlInfo.Host+":"+mysqlInfo.Port+")/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//禁止副表 gorm默认使用复数映射，go代码的单数、复数struct形式都匹配到复数表中,开启后，将严格匹配，遵守单数形式
	db.SingularTable(true)

	// 设置服务
	srv := micro.NewService(
		micro.Name(MICRO_SERVICE_NAME),
		micro.Version(MICRO_VERSION),
		//这里设置地址和需要暴露的端口
		micro.Address(MICRO_ADDRESS),
		//添加注册中心
		micro.Registry(consulRegistry),
		//绑定链路追踪 服务端绑定handle 客户端绑定client
		micro.WrapHandler(opentracing4.NewHandlerWrapper(opentracing.GlobalTracer())),
		//添加限流
		micro.WrapHandler(ratelimit4.NewHandlerWrapper(MICRO_QPS)),
	)

	//初始化建表，多个表
	// repository.NewProductRepository(db).InitTable()

	productService := service.NewProductService(repository.NewProductRepository(db))

	// Initialise service
	srv.Init()

	// Register Handler
	pb.RegisterShopProductHandler(srv.Server(), &handler.ShopProduct{ProductService: productService})

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
