package svc

import (
    "gitee.com/meinongyihe/travel-api/internal/config"
    "gitee.com/meinongyihe/travel-rpc/travel"
    "google.golang.org/grpc"
)

type ServiceContext struct {
    Config          config.Config
    CatalogClient   travel.CatalogServiceClient
    InventoryClient travel.InventoryServiceClient
    OrderClient     travel.OrderServiceClient
    PaymentClient   travel.PaymentServiceClient
}

func NewServiceContext(c config.Config, conn grpc.ClientConnInterface) *ServiceContext {
    return &ServiceContext{
        Config:          c,
        CatalogClient:   travel.NewCatalogServiceClient(conn),
        InventoryClient: travel.NewInventoryServiceClient(conn),
        OrderClient:     travel.NewOrderServiceClient(conn),
        PaymentClient:   travel.NewPaymentServiceClient(conn),
    }
}
