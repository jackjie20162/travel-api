package handler

import (
    "net/http"
    "gitee.com/meinongyihe/travel-api/internal/svc"
    "github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, ctx *svc.ServiceContext) {
    products:=NewProductHandler(ctx);inventory:=NewInventoryHandler(ctx);orders:=NewOrderHandler(ctx);payments:=NewPaymentHandler(ctx)
    server.AddRoutes([]rest.Route{
        {Method:http.MethodGet,Path:"/health",Handler:healthHandler},
        {Method:http.MethodGet,Path:"/api/travel/products",Handler:products.List},
        {Method:http.MethodGet,Path:"/api/travel/products/:id",Handler:products.Detail},
        {Method:http.MethodPost,Path:"/api/travel/inventory/check",Handler:inventory.Check},
        {Method:http.MethodPost,Path:"/api/travel/orders",Handler:orders.Create},
        {Method:http.MethodGet,Path:"/api/travel/orders/:orderNo",Handler:orders.Get},
        {Method:http.MethodPost,Path:"/api/travel/payments",Handler:payments.Create},
        {Method:http.MethodGet,Path:"/api/travel/payments/paypal/return",Handler:payments.Return},
        {Method:http.MethodGet,Path:"/api/travel/payments/paypal/cancel",Handler:payments.Cancel},
        {Method:http.MethodPost,Path:"/api/travel/payments/paypal/webhook",Handler:payments.PayPalWebhook},
        {Method:http.MethodGet,Path:"/api/travel/payments/:paymentNo",Handler:payments.Get},
    })
}

func healthHandler(w http.ResponseWriter,r *http.Request){w.Header().Set("Content-Type","application/json");_,_=w.Write([]byte(`{"code":0,"msg":"ok","data":{"service":"travel-api"}}`))}
