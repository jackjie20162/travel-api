package handler

import (
    "context"
    "encoding/json"
    "errors"
    "io"
    "net/http"
    "net/url"
    "os"
    "strconv"
    "strings"

    "gitee.com/meinongyihe/travel-api/internal/middleware"
    "gitee.com/meinongyihe/travel-api/internal/paypal"
    "gitee.com/meinongyihe/travel-api/internal/svc"
    "gitee.com/meinongyihe/travel-rpc/travel"
    "github.com/zeromicro/go-zero/rest/httpx"
    "google.golang.org/grpc/metadata"
)

type PaymentHandler struct { svcCtx *svc.ServiceContext }
func NewPaymentHandler(svcCtx *svc.ServiceContext) *PaymentHandler { return &PaymentHandler{svcCtx:svcCtx} }

func (h *PaymentHandler) Create(w http.ResponseWriter,r *http.Request){
    var req CreatePaymentReq
    if err:=httpx.Parse(r,&req);err!=nil{httpx.Error(w,err);return}
    if req.OrderNo==""||req.Provider==""||req.IdempotencyKey==""{httpx.Error(w,errors.New("orderNo, provider and idempotencyKey are required"));return}
    if !strings.EqualFold(req.Provider,"paypal"){httpx.Error(w,errors.New("only paypal is enabled"));return}
    ctx,err:=rpcContext(r);if err!=nil{httpx.Error(w,err);return}
    p,err:=h.svcCtx.PaymentClient.Create(ctx,&travel.CreatePaymentRequest{OrderNo:req.OrderNo,Provider:req.Provider,IdempotencyKey:req.IdempotencyKey});if err!=nil{httpx.Error(w,err);return}
    tenantID,ok:=middleware.TenantID(r.Context());if !ok{httpx.Error(w,errors.New("tenant context is required"));return};merchantID,ok:=middleware.MerchantID(r.Context());if !ok{httpx.Error(w,errors.New("merchant context is required"));return}
    secret:=h.paypalSecret();if secret==""{httpx.Error(w,errors.New("paypal client secret is not configured"));return}
    state:=paypal.SignState(secret,tenantID,merchantID,p.PaymentNo)
    checkout,err:=h.svcCtx.PayPal.CreateOrder(ctx,paypal.CheckoutInput{PaymentNo:p.PaymentNo,Amount:p.Amount,Currency:p.Currency,State:state});if err!=nil{httpx.Error(w,err);return}
    p,err=h.svcCtx.PaymentClient.SetProviderID(ctx,&travel.SetPaymentProviderIDRequest{PaymentNo:p.PaymentNo,ProviderPaymentId:checkout.OrderID});if err!=nil{httpx.Error(w,err);return}
    out:=toPaymentResp(p);out.CheckoutURL=checkout.CheckoutURL;httpx.OkJson(w,out)
}

func (h *PaymentHandler) Get(w http.ResponseWriter,r *http.Request){paymentNo:=r.PathValue("paymentNo");if paymentNo==""{httpx.Error(w,errors.New("invalid payment number"));return};ctx,err:=rpcContext(r);if err!=nil{httpx.Error(w,err);return};p,err:=h.svcCtx.PaymentClient.Get(ctx,&travel.PaymentNoRequest{PaymentNo:paymentNo});if err!=nil{httpx.Error(w,err);return};httpx.OkJson(w,toPaymentResp(p))}

func (h *PaymentHandler) Return(w http.ResponseWriter,r *http.Request){
    state:=r.URL.Query().Get("state");token:=r.URL.Query().Get("token");if state==""||token==""{httpx.Error(w,errors.New("paypal return state or token is missing"));return}
    tenantID,merchantID,paymentNo,err:=paypal.ParseState(h.paypalSecret(),state);if err!=nil{httpx.Error(w,err);return}
    capture,err:=h.svcCtx.PayPal.CaptureOrder(r.Context(),token);if err!=nil{httpx.Error(w,err);return};if !strings.EqualFold(capture.Status,"COMPLETED")||capture.CaptureID==""{httpx.Error(w,errors.New("paypal payment is not completed"));return}
    ctx:=paymentContext(r.Context(),tenantID,merchantID);p,err:=h.svcCtx.PaymentClient.MarkPaid(ctx,&travel.MarkPaymentPaidRequest{PaymentNo:paymentNo,ProviderPaymentId:token});if err!=nil{httpx.Error(w,err);return}
    if h.svcCtx.Config.Payment.FrontendReturnURL!=""{u,_:=url.Parse(h.svcCtx.Config.Payment.FrontendReturnURL);q:=u.Query();q.Set("paymentNo",p.PaymentNo);q.Set("status",p.Status);u.RawQuery=q.Encode();http.Redirect(w,r,u.String(),http.StatusFound);return};httpx.OkJson(w,toPaymentResp(p))
}

func (h *PaymentHandler) Cancel(w http.ResponseWriter,r *http.Request){state:=r.URL.Query().Get("state");if state==""{httpx.Error(w,errors.New("paypal cancel state is missing"));return};_,_,paymentNo,err:=paypal.ParseState(h.paypalSecret(),state);if err!=nil{httpx.Error(w,err);return};if h.svcCtx.Config.Payment.FrontendReturnURL!=""{u,_:=url.Parse(h.svcCtx.Config.Payment.FrontendReturnURL);q:=u.Query();q.Set("paymentNo",paymentNo);q.Set("status","CANCELLED");u.RawQuery=q.Encode();http.Redirect(w,r,u.String(),http.StatusFound);return};httpx.OkJson(w,map[string]any{"paymentNo":paymentNo,"status":"CANCELLED"})}

func (h *PaymentHandler) PayPalWebhook(w http.ResponseWriter,r *http.Request){
    body,err:=io.ReadAll(r.Body);if err!=nil{httpx.Error(w,err);return};if err:=h.svcCtx.PayPal.VerifyWebhook(r.Context(),r.Header,json.RawMessage(body));err!=nil{httpx.Error(w,err);return}
    var event struct{ID string `json:"id"`;EventType string `json:"event_type"`;Resource struct{CustomID string `json:"custom_id"`;Status string `json:"status"`;SupplementaryData struct{RelatedIDs struct{OrderID string `json:"order_id"`} `json:"related_ids"`} `json:"supplementary_data"`} `json:"resource"`}
    if err:=json.Unmarshal(body,&event);err!=nil{httpx.Error(w,err);return}
    if event.EventType=="PAYMENT.CAPTURE.COMPLETED"&&strings.EqualFold(event.Resource.Status,"COMPLETED")&&event.Resource.CustomID!=""{tenantID,merchantID,paymentNo,err:=paypal.ParseState(h.paypalSecret(),event.Resource.CustomID);if err!=nil{httpx.Error(w,err);return};ctx:=paymentContext(r.Context(),tenantID,merchantID);if _,err:=h.svcCtx.PaymentClient.MarkPaid(ctx,&travel.MarkPaymentPaidRequest{PaymentNo:paymentNo,ProviderPaymentId:event.Resource.SupplementaryData.RelatedIDs.OrderID});err!=nil{httpx.Error(w,err);return}}
    w.WriteHeader(http.StatusNoContent)
}

func (h *PaymentHandler) paypalSecret()string{if h.svcCtx.Config.Payment.PayPal.Secret!=""{return h.svcCtx.Config.Payment.PayPal.Secret};return os.Getenv("PAYPAL_CLIENT_SECRET")}
func paymentContext(ctx context.Context,tenantID,merchantID int64)context.Context{return metadata.NewOutgoingContext(ctx,metadata.Pairs("x-tenant-id",fmtInt64(tenantID),"x-merchant-id",fmtInt64(merchantID)))}
func fmtInt64(v int64)string{return strconv.FormatInt(v,10)}
func toPaymentResp(p *travel.Payment) PaymentResp{return PaymentResp{Id:p.Id,PaymentNo:p.PaymentNo,OrderNo:p.OrderNo,Provider:p.Provider,ProviderPaymentId:p.ProviderPaymentId,Amount:p.Amount,Currency:p.Currency,Status:p.Status}}
