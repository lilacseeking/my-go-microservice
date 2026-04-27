// cmd/api/main.go
package main

import (
	"context"
	"fmt"
	"log"
	pb "my-go-microservice/proto"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"my-go-microservice/api/v1/handler"
	"my-go-microservice/internal/cache"
	"my-go-microservice/internal/config"
	"my-go-microservice/internal/database"
	"my-go-microservice/internal/messaging"
	"my-go-microservice/internal/repository"
	"my-go-microservice/internal/service"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("configs")
	if err != nil {
		log.Fatalf("无法加载配置: %v", err)
	}

	// 初始化日志
	initLogger(cfg.Logger)
	logrus.Info("日志系统已启动")

	// 初始化 PostgreSQL 连接
	db, err := database.InitPostgreSQL(&cfg.Database.Postgres)
	if err != nil {
		logrus.WithError(err).Fatal("数据库连接失败")
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()
	logrus.Info("PostgreSQL 数据库连接成功")

	// 初始化 Redis 连接池
	redisPool := cache.NewRedisPool(&cfg.Redis)
	conn := redisPool.Get()
	defer conn.Close()
	if _, err := conn.Do("PING"); err != nil {
		logrus.WithError(err).Fatal("Redis 连接失败")
	}
	logrus.Info("Redis 客户端连接成功")

	// 初始化 Kafka 生产者
	kafkaProducer, err := messaging.NewKafkaProducer(&cfg.Kafka)
	if err != nil {
		logrus.WithError(err).Fatal("Kafka 生产者初始化失败")
	}
	defer kafkaProducer.Close()
	logrus.Info("Kafka 生产者已启动")

	// 初始化数据访问层
	userRepo := repository.NewUserRepository(db)

	// 初始化业务逻辑层
	userService := service.NewUserService(userRepo, redisPool, kafkaProducer)
	logrus.Info("业务逻辑层初始化完成")

	// 启动 HTTP 服务器
	httpServer := &http.Server{
		Addr:    cfg.HTTP.Host + ":" + fmt.Sprint(cfg.HTTP.Port),
		Handler: setupRouter(userService),
	}
	go func() {
		logrus.Infof("HTTP 服务器正在监听 %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Fatal("HTTP 服务器启动失败")
		}
	}()

	// 启动 gRPC 服务器
	lis, err := net.Listen("tcp", cfg.GRPC.Host+":"+fmt.Sprint(cfg.GRPC.Port))
	if err != nil {
		logrus.WithError(err).Fatal("gRPC 服务器监听失败")
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, handler.NewGRPCUserServer(userService))
	go func() {
		logrus.Infof("gRPC 服务器正在监听 %s", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			logrus.WithError(err).Fatal("gRPC 服务器启动失败")
		}
	}()

	// 启动健康检查服务器
	healthCheckServer := startHealthCheckServer(cfg.HTTP.Host, cfg.HTTP.Port+1)

	// 设置信号监听
	_, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	logrus.Info("接收到中断信号，开始优雅关闭...")

	// 执行优雅关闭 - 先取消 context，用于通知其他 goroutine
	cancel()

	// 关闭 HTTP 服务器（使用独立的 context 设置超时）
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logrus.WithError(err).Error("HTTP 服务器关闭失败")
	}
	logrus.Info("HTTP 服务器已关闭")

	// 关闭 gRPC 服务器
	grpcServer.GracefulStop()
	logrus.Info("gRPC 服务器已关闭")

	// 关闭健康检查服务器
	if err := healthCheckServer.Shutdown(shutdownCtx); err != nil {
		logrus.WithError(err).Error("健康检查服务器关闭失败")
	}
	logrus.Info("健康检查服务器已关闭")

	// 关闭 Redis 连接池
	redisPool.Close()
	logrus.Info("Redis 连接池已关闭")

	// 关闭 Kafka 生产者
	kafkaProducer.Close()
	logrus.Info("Kafka 生产者已关闭")

	logrus.Info("所有服务已成功关闭，程序退出")
}

// setupRouter 初始化 Gin 路由
func setupRouter(userService service.UserService) http.Handler {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	v1Group := r.Group("/v1")
	handler.RegisterUserRoutes(v1Group, &handler.UserHandler{UserService: userService})

	return r
}

// initLogger 根据配置初始化日志
func initLogger(loggerCfg config.Logger) {
	level, err := logrus.ParseLevel(loggerCfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	if loggerCfg.Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	if loggerCfg.Output != "stdout" {
		file, err := os.OpenFile(loggerCfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logrus.SetOutput(file)
		}
	}
}

// startHealthCheckServer 启动健康检查HTTP服务
func startHealthCheckServer(host string, port int) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:    host + ":" + fmt.Sprint(port),
		Handler: mux,
	}

	logrus.Infof("健康检查服务器正在监听 %s", server.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Error("健康检查服务器启动失败")
		}
	}()

	return server
}
