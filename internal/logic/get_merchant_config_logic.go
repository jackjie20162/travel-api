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

type GetMerchantConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMerchantConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMerchantConfigLogic {
	return &GetMerchantConfigLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetMerchantConfigLogic) GetMerchantConfig() (*types.ConfigGetResp, error) {
	configType := pathVal(l.ctx, "configType")
	if configType == "" {
		return nil, errors.New("missing configType")
	}

	merchantID, _ := middleware.MerchantID(l.ctx)
	l.Infof("GetMerchantConfig: configType=%s, merchantID=%d", configType, merchantID)

	configs, err := l.svcCtx.DB.MerchantConfig.Query().
		Where(
			merchantconfig.ConfigTypeEQ(configType),
			merchantconfig.MerchantIDEQ(merchantID),
		).
		All(l.ctx)
	if err != nil {
		l.Infof("GetMerchantConfig query error: %v", err)
		return nil, err
	}

	l.Infof("GetMerchantConfig: found %d rows", len(configs))

	items := make([]types.ConfigItem, 0, len(configs))
	for _, c := range configs {
		items = append(items, types.ConfigItem{
			ConfigType:  c.ConfigType,
			ConfigKey:   c.ConfigKey,
			ConfigValue: c.ConfigValue,
			Description: c.Description,
		})
	}

	return &types.ConfigGetResp{
		ConfigType: configType,
		Items:      items,
	}, nil
}
