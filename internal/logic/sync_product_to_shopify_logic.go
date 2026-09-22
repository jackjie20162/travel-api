package logic

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	shopifyprovider "gitee.com/meinongyihe/travel-api/internal/payment/shopify"
	"gitee.com/meinongyihe/travel-rpc/travel"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncProductToShopifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncProductToShopifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncProductToShopifyLogic {
	return &SyncProductToShopifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

// SyncProductToShopify syncs a travel product to the Shopify store as a Shopify Product.
func (l *SyncProductToShopifyLogic) SyncProductToShopify() (*types.OkResp, error) {
	if l.svcCtx.ShopifyProvider == nil {
		return nil, errors.New("shopify is not configured")
	}

	// Extract product ID from path param.
	idStr := pathVal(l.ctx, "id")
	if idStr == "" {
		return nil, errors.New("missing product id")
	}
	productID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || productID <= 0 {
		return nil, errors.New("invalid product id")
	}

	// Fetch product details from RPC.
	ctx, err := rpcContext(l.ctx)
	if err != nil {
		return nil, err
	}
	product, err := l.svcCtx.CatalogClient.GetProduct(ctx, &travel.ProductIdRequest{Id: productID})
	if err != nil {
		return nil, fmt.Errorf("get product %d: %w", productID, err)
	}

	// Fetch packages as variants.
	packages, err := l.svcCtx.CatalogClient.ListPackages(ctx, &travel.PackageListRequest{ProductId: productID})
	if err != nil {
		return nil, fmt.Errorf("list packages for product %d: %w", productID, err)
	}

	variants := make([]shopifyprovider.VariantInput, 0, len(packages.GetItems()))
	for _, pkg := range packages.GetItems() {
		variants = append(variants, shopifyprovider.VariantInput{
			Title: pkg.Name,
			Price: fmt.Sprintf("%.2f", float64(product.MinPrice)/100),
			SKU:   pkg.Code,
		})
	}

	// Sync to Shopify.
	shopifyID, err := l.svcCtx.ShopifyProvider.SyncProduct(l.ctx, shopifyprovider.ProductInput{
		Title:       product.Title,
		Description: product.Description,
		Vendor:      fmt.Sprintf("Merchant-%d", product.MerchantId),
		ProductType: "Travel",
		Tags:        []string{product.Destination, "travel", "dubai"},
		Images:      splitImages(product.Images),
		Variants:    variants,
	})
	if err != nil {
		return nil, fmt.Errorf("shopify sync failed: %w", err)
	}

	l.Infof("[Shopify] product %d synced as shopify_product=%s", productID, shopifyID)
	return &types.OkResp{Ok: true, Msg: "synced as shopify product " + shopifyID}, nil
}

// splitImages splits a comma-separated image string into a slice.
func splitImages(images string) []string {
	if images == "" {
		return nil
	}
	var result []string
	for _, img := range splitByComma(images) {
		if img != "" {
			result = append(result, img)
		}
	}
	return result
}

func splitByComma(s string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}
