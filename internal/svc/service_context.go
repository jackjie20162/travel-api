package svc

import (
	"entgo.io/ent/dialect"
	"gitee.com/meinongyihe/travel-api/internal/config"
	"gitee.com/meinongyihe/travel-api/internal/paypal"
	paypaladapter "gitee.com/meinongyihe/travel-api/internal/payment/paypal"
	"gitee.com/meinongyihe/travel-api/internal/payment"
	shopifyprovider "gitee.com/meinongyihe/travel-api/internal/payment/shopify"
	stripeprovider "gitee.com/meinongyihe/travel-api/internal/payment/stripe"
	"gitee.com/meinongyihe/travel-rpc/ent"
	"gitee.com/meinongyihe/travel-rpc/travel"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config           config.Config
	DB               *ent.Client
	Redis            *redis.Redis
	CatalogClient    travel.CatalogServiceClient
	InventoryClient  travel.InventoryServiceClient
	OrderClient      travel.OrderServiceClient
	PaymentClient    travel.PaymentServiceClient
	ManagementClient travel.TravelManagementServiceClient
	UserClient       travel.UserServiceClient
	ReviewClient     travel.ReviewServiceClient
	CurrencyClient   travel.CurrencyServiceClient
	PayPal           *paypal.Provider          // legacy PayPal provider (backward compat)
	StripeProvider   *stripeprovider.Provider  // Stripe provider
	ShopifyProvider  *shopifyprovider.Provider // Shopify provider
}

func NewServiceContext(c config.Config) *ServiceContext {
	// build ent client from DatabaseConf
	var drv dialect.Driver = c.DatabaseConf.NewNoCacheDriver()
	if c.DatabaseConf.Debug {
		drv = dialect.Debug(drv, logx.Info)
	}
	entClient := ent.NewClient(ent.Log(logx.Info), ent.Driver(drv))

	conn := zrpc.MustNewClient(c.TravelRpc)

	ctx := &ServiceContext{
		Config:           c,
		DB:               entClient,
		Redis:            c.Redis.NewRedis(),
		CatalogClient:    travel.NewCatalogServiceClient(conn.Conn()),
		InventoryClient:  travel.NewInventoryServiceClient(conn.Conn()),
		OrderClient:      travel.NewOrderServiceClient(conn.Conn()),
		PaymentClient:    travel.NewPaymentServiceClient(conn.Conn()),
		ManagementClient: travel.NewTravelManagementServiceClient(conn.Conn()),
		UserClient:       travel.NewUserServiceClient(conn.Conn()),
		ReviewClient:     travel.NewReviewServiceClient(conn.Conn()),
		CurrencyClient:   travel.NewCurrencyServiceClient(conn.Conn()),
		PayPal:           paypal.New(paypal.Config{Enabled: c.Payment.PayPal.Enabled, BaseURL: c.Payment.PayPal.BaseURL, ClientID: c.Payment.PayPal.ClientID, Secret: c.Payment.PayPal.Secret, WebhookID: c.Payment.PayPal.WebhookID, ReturnURL: c.Payment.PayPal.ReturnURL, CancelURL: c.Payment.PayPal.CancelURL, SkipVerify: c.Payment.PayPal.SkipVerify}),
	}

	// Register payment providers into the unified registry.
	if c.Payment.PayPal.Enabled {
		adapter := paypaladapter.New(paypal.Config{
			Enabled:    c.Payment.PayPal.Enabled,
			BaseURL:    c.Payment.PayPal.BaseURL,
			ClientID:   c.Payment.PayPal.ClientID,
			Secret:     c.Payment.PayPal.Secret,
			WebhookID:  c.Payment.PayPal.WebhookID,
			ReturnURL:  c.Payment.PayPal.ReturnURL,
			CancelURL:  c.Payment.PayPal.CancelURL,
			SkipVerify: c.Payment.PayPal.SkipVerify,
		})
		payment.Register(adapter)
		logx.Info("[Payment] PayPal provider registered")
	}

	if c.Payment.Stripe.Enabled {
		stripeProvider := stripeprovider.New(stripeprovider.Config{
			Enabled:        c.Payment.Stripe.Enabled,
			SecretKey:      c.Payment.Stripe.SecretKey,
			PublishableKey: c.Payment.Stripe.PublishableKey,
			WebhookSecret:  c.Payment.Stripe.WebhookSecret,
		})
		ctx.StripeProvider = stripeProvider
		payment.Register(stripeProvider)
		logx.Info("[Payment] Stripe provider registered")
	}

	if c.Shopify.Enabled {
		shopifyProvider := shopifyprovider.New(shopifyprovider.Config{
			Enabled:       c.Shopify.Enabled,
			StoreDomain:   c.Shopify.StoreDomain,
			AccessToken:   c.Shopify.AccessToken,
			APIVersion:    c.Shopify.APIVersion,
			WebhookSecret: c.Shopify.WebhookSecret,
			SyncEnabled:   c.Shopify.SyncEnabled,
		})
		ctx.ShopifyProvider = shopifyProvider
		payment.Register(shopifyProvider)
		logx.Info("[Payment] Shopify provider registered")
	}

	return ctx
}
