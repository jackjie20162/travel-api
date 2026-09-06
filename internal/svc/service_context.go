package svc

import (
	"entgo.io/ent/dialect"
	"gitee.com/meinongyihe/travel-api/internal/config"
	"gitee.com/meinongyihe/travel-api/internal/paypal"
	"gitee.com/meinongyihe/travel-rpc/ent"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config           config.Config
	DB               *ent.Client
	CatalogClient    travel.CatalogServiceClient
	InventoryClient  travel.InventoryServiceClient
	OrderClient      travel.OrderServiceClient
	PaymentClient    travel.PaymentServiceClient
	ManagementClient travel.TravelManagementServiceClient
	UserClient       travel.UserServiceClient
	PayPal           *paypal.Provider
}

func NewServiceContext(c config.Config) *ServiceContext {
	// build ent client from DatabaseConf
	var drv dialect.Driver = c.DatabaseConf.NewNoCacheDriver()
	if c.DatabaseConf.Debug {
		drv = dialect.Debug(drv, logx.Info)
	}
	entClient := ent.NewClient(ent.Log(logx.Info), ent.Driver(drv))

	conn := zrpc.MustNewClient(c.TravelRpc)
	return &ServiceContext{
		Config:           c,
		DB:               entClient,
		CatalogClient:    travel.NewCatalogServiceClient(conn.Conn()),
		InventoryClient:  travel.NewInventoryServiceClient(conn.Conn()),
		OrderClient:      travel.NewOrderServiceClient(conn.Conn()),
		PaymentClient:    travel.NewPaymentServiceClient(conn.Conn()),
		ManagementClient: travel.NewTravelManagementServiceClient(conn.Conn()),
		UserClient:       travel.NewUserServiceClient(conn.Conn()),
		PayPal:           paypal.New(paypal.Config{Enabled: c.Payment.PayPal.Enabled, BaseURL: c.Payment.PayPal.BaseURL, ClientID: c.Payment.PayPal.ClientID, Secret: c.Payment.PayPal.Secret, WebhookID: c.Payment.PayPal.WebhookID, ReturnURL: c.Payment.PayPal.ReturnURL, CancelURL: c.Payment.PayPal.CancelURL, SkipVerify: c.Payment.PayPal.SkipVerify}),
	}
}
