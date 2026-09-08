package logic

import (
	"context"
	"errors"

	"gitee.com/meinongyihe/travel-api/internal/middleware"
	"gitee.com/meinongyihe/travel-api/internal/svc"
	"gitee.com/meinongyihe/travel-api/internal/types"
	"gitee.com/meinongyihe/travel-rpc/ent/merchantconfig"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveMerchantConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveMerchantConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveMerchantConfigLogic {
	return &SaveMerchantConfigLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SaveMerchantConfigLogic) SaveMerchantConfig(req *types.ConfigSaveReq) (*types.OkResp, error) {
	if req.ConfigType == "" {
		return nil, errors.New("configType is required")
	}
	if len(req.Items) == 0 {
		return nil, errors.New("items is required")
	}

	tenantID, _ := middleware.TenantID(l.ctx)
	merchantID, _ := middleware.MerchantID(l.ctx)

	for _, item := range req.Items {
		if item.ConfigKey == "" {
			return nil, errors.New("configKey is required for all items")
		}

		// Upsert: query first, then update or insert
		exists, err := l.svcCtx.DB.MerchantConfig.Query().
			Where(
				merchantconfig.ConfigTypeEQ(req.ConfigType),
				merchantconfig.ConfigKeyEQ(item.ConfigKey),
				merchantconfig.MerchantIDEQ(merchantID),
			).
			Exist(l.ctx)
		if err != nil {
			return nil, err
		}

		if exists {
			_, err = l.svcCtx.DB.MerchantConfig.Update().
				Where(
					merchantconfig.ConfigTypeEQ(req.ConfigType),
					merchantconfig.ConfigKeyEQ(item.ConfigKey),
					merchantconfig.MerchantIDEQ(merchantID),
				).
				SetConfigValue(item.ConfigValue).
				SetDescription(item.Description).
				Save(l.ctx)
		} else {
			_, err = l.svcCtx.DB.MerchantConfig.Create().
				SetTenantID(tenantID).
				SetMerchantID(merchantID).
				SetConfigType(req.ConfigType).
				SetConfigKey(item.ConfigKey).
				SetConfigValue(item.ConfigValue).
				SetDescription(item.Description).
				Save(l.ctx)
		}
		if err != nil {
			return nil, err
		}
	}

	return &types.OkResp{Ok: true, Msg: "配置已保存"}, nil
}
