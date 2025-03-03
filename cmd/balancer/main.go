package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync/atomic"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "my.custom.server/gen/proto"
)

type serverImpl struct {
	pb.UnimplementedServiceServer
	counter uint64
	cdnHost string
}

func (s *serverImpl) Method(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	count := atomic.AddUint64(&s.counter, 1)

	parsedURL, err := url.Parse(req.Video)
	if err != nil {
		log.Printf("Ошибка парсинга URL: %v. Входные данные: %s", err, req.Video)
		return nil, fmt.Errorf("неверный URL видео: %v", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		log.Printf("Неверный URL: %s. Ожидается http или https", parsedURL.Scheme)
		return nil, fmt.Errorf("неподдерживаемый протокол: %s", parsedURL.Scheme)
	}

	if count%10 == 0 {
		return &pb.Response{RedirectUrl: req.Video}, nil
	}

	hostParts := strings.Split(parsedURL.Host, ".")
	if len(hostParts) == 0 {
		log.Printf("Неверный host в URL: %s", parsedURL.Host)
		return nil, fmt.Errorf("неверный host в URL: %s", parsedURL.Host)
	}
	subdomain := hostParts[0]

	match, err := regexp.MatchString(`^s\d+$`, subdomain)
	if err != nil {
		log.Printf("Ошибка при проверке поддомена: %v", err)
		return nil, fmt.Errorf("internal error")
	}
	if !match {
		log.Printf("Поддомен не соответствует ожидаемому формату: %s", subdomain)
		return nil, fmt.Errorf("неверный формат поддомена: %s", subdomain)
	}

	newPath := fmt.Sprintf("/%s%s", subdomain, parsedURL.Path)

	cdnURL := url.URL{
		Scheme: "http",
		Host:   s.cdnHost,
		Path:   newPath,
	}

	return &pb.Response{RedirectUrl: cdnURL.String()}, nil
}

func main() {
	cdnHost := os.Getenv("CDN_HOST")
	if cdnHost == "" {
		cdnHost = "cdn.example.com"
	}

	grpcServer := grpc.NewServer()
	pb.RegisterServiceServer(grpcServer, &serverImpl{cdnHost: cdnHost})
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error creating listener: %v", err)
	}

	go func() {
		log.Println("starting gRPC server, port:50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Error starting gRPC server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Stop signal received, terminating server..")

	grpcServer.GracefulStop()
	log.Println("Server stopped")
}
