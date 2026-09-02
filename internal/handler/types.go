package handler

type ProductListReq struct { Keyword string `form:"keyword,optional"`; Destination string `form:"destination,optional"`; Page int `form:"page,optional"`; PageSize int `form:"pageSize,optional"` }
type Product struct { Id int64 `json:"id"`; TenantId int64 `json:"tenantId"`; MerchantId int64 `json:"merchantId"`; Code string `json:"code"`; Title string `json:"title"`; Slug string `json:"slug"`; Destination string `json:"destination"`; Description string `json:"description"`; Currency string `json:"currency"`; MinPrice int64 `json:"minPrice"`; Status string `json:"status"` }
type ProductListResp struct { Items []Product `json:"items"`; Total int64 `json:"total"` }
type CreateProductReq struct { Code string `json:"code"`; Title string `json:"title"`; Slug string `json:"slug,optional"`; Destination string `json:"destination,optional"`; Description string `json:"description,optional"`; Currency string `json:"currency,optional"`; MinPrice int64 `json:"minPrice,optional"` }
type UpdateProductReq struct { Code string `json:"code"`; Title string `json:"title"`; Slug string `json:"slug,optional"`; Destination string `json:"destination,optional"`; Description string `json:"description,optional"`; Currency string `json:"currency,optional"`; MinPrice int64 `json:"minPrice,optional"` }
type PublishProductReq struct { Published bool `json:"published"` }
type PackageReq struct { Code string `json:"code"`; Name string `json:"name"` }
type Package struct { Id int64 `json:"id"`; ProductId int64 `json:"productId"`; Code string `json:"code"`; Name string `json:"name"`; Status string `json:"status"` }
type PackageListResp struct { Items []Package `json:"items"` }
type InventoryUpsertReq struct { Date string `json:"date"`; TimeSlot string `json:"timeSlot,optional"`; Capacity int `json:"capacity"`; UnitPrice int64 `json:"unitPrice"`; Currency string `json:"currency,optional"`; Status string `json:"status,optional"` }
type InventoryItem struct { Id int64 `json:"id"`; PackageId int64 `json:"packageId"`; Date string `json:"date"`; TimeSlot string `json:"timeSlot"`; Capacity int32 `json:"capacity"`; Reserved int32 `json:"reserved"`; UnitPrice int64 `json:"unitPrice"`; Currency string `json:"currency"`; Status string `json:"status"` }
type InventoryListResp struct { Items []InventoryItem `json:"items"` }
type InventoryCheckReq struct { PackageId int64 `json:"packageId"`; Date string `json:"date"`; TimeSlot string `json:"timeSlot,optional"`; Quantity int `json:"quantity"` }
type InventoryCheckResp struct { Available bool `json:"available"`; Remaining int32 `json:"remaining"`; UnitPrice int64 `json:"unitPrice"`; Currency string `json:"currency"` }
type CreateOrderReq struct { ProductId int64 `json:"productId"`; PackageId int64 `json:"packageId"`; Date string `json:"date"`; TimeSlot string `json:"timeSlot,optional"`; Quantity int `json:"quantity"`; CustomerEmail string `json:"customerEmail,optional"`; ReservationKey string `json:"reservationKey"` }
type CreateOrderResp struct { Id int64 `json:"id"`; OrderNo string `json:"orderNo"`; Status string `json:"status"`; TotalAmount int64 `json:"totalAmount"`; Currency string `json:"currency"` }
