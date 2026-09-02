package handler

type ProductListReq struct { Keyword string `form:"keyword,optional"`; Destination string `form:"destination,optional"`; Page int `form:"page,optional"`; PageSize int `form:"pageSize,optional"` }
type Product struct { Id int64 `json:"id"`; TenantId int64 `json:"tenantId"`; MerchantId int64 `json:"merchantId"`; Code string `json:"code"`; Title string `json:"title"`; Slug string `json:"slug"`; Destination string `json:"destination"`; Description string `json:"description"`; Currency string `json:"currency"`; MinPrice int64 `json:"minPrice"`; Status string `json:"status"` }
type ProductListResp struct { Items []Product `json:"items"`; Total int64 `json:"total"` }
type InventoryCheckReq struct { PackageId int64 `json:"packageId"`; Date string `json:"date"`; TimeSlot string `json:"timeSlot,optional"`; Quantity int `json:"quantity"` }
type InventoryCheckResp struct { Available bool `json:"available"`; Remaining int32 `json:"remaining"`; UnitPrice int64 `json:"unitPrice"`; Currency string `json:"currency"` }
type CreateOrderReq struct { ProductId int64 `json:"productId"`; PackageId int64 `json:"packageId"`; Date string `json:"date"`; TimeSlot string `json:"timeSlot,optional"`; Quantity int `json:"quantity"`; CustomerEmail string `json:"customerEmail,optional"`; ReservationKey string `json:"reservationKey"` }
type CreateOrderResp struct { Id int64 `json:"id"`; OrderNo string `json:"orderNo"`; Status string `json:"status"`; TotalAmount int64 `json:"totalAmount"`; Currency string `json:"currency"` }
