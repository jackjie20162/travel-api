package logic

import (
	"context"
	"os"
	"testing"

	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/ent"
	"gitee.com/meinongyihe/travel-rpc/ent/product"

	_ "github.com/go-sql-driver/mysql"
)

// 临时验证用：连接本地 travel 库，确认「线路无主图时回退取景点图」。
// 只读，不回写任何数据。运行：
//
//	$env:TRAVEL_DB_DSN='root:123456@tcp(127.0.0.1:13307)/travel?charset=utf8mb4&parseTime=True&loc=Local'
//	go test ./internal/logic -run TestFillMissingProductCovers -v
func TestFillMissingProductCovers(t *testing.T) {
	dsn := os.Getenv("TRAVEL_DB_DSN")
	if dsn == "" {
		t.Skip("TRAVEL_DB_DSN not set")
	}
	client, err := ent.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("open ent: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	rows, err := client.Product.Query().Where(product.MerchantIDEQ(1)).Order(product.ByID()).All(ctx)
	if err != nil {
		t.Fatalf("query products: %v", err)
	}
	items := make([]types.Product, 0, len(rows))
	for _, p := range rows {
		pv := entToProduct(p)
		items = append(items, pv)
	}

	// 复制一份并清掉所有主图，模拟「商户没设置主图」
	cleared := make([]types.Product, len(items))
	copy(cleared, items)
	for i := range cleared {
		cleared[i].CoverImage = ""
	}

	fillMissingProductCovers(ctx, client, cleared)

	filled := 0
	for i := range cleared {
		t.Logf("product %d %-30q cover=%q", cleared[i].Id, cleared[i].Title, cleared[i].CoverImage)
		if cleared[i].CoverImage != "" {
			filled++
		}
	}
	if filled == 0 {
		t.Fatalf("no cover filled at all")
	}
	t.Logf("filled %d/%d products", filled, len(cleared))
}
