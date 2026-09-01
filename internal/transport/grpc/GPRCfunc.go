package grpc

import (
	"context"
	"easyVending/internal/domain"
	kitpb "test/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type GRPCClient struct {
	conn   *grpc.ClientConn
	client kitpb.KitNalogServiceClient
}

func NewGRPCClient(addr string) (*GRPCClient, error) {

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(timeOut, statusInterceptor),
	)
	if err != nil {
		return nil, err
	}
	cli := kitpb.NewKitNalogServiceClient(conn)

	return &GRPCClient{
		conn,
		cli,
	}, nil
}

func (c *GRPCClient) GetConnStatus() bool {
	return IsConnected()
}

func (c *GRPCClient) Close() error {
	return c.conn.Close()
}

func (c *GRPCClient) LoginGrpc(ctx context.Context, ld *domain.AuthData, dev *domain.DeviceInfo) (bool, string) {

	req := &kitpb.Request{
		Data: &kitpb.LoginData{
			CompanyId:   ld.CompanyId,
			UserLogin:   ld.UserLogin,
			PasswordKit: ld.PasswordKit,
			Inn:         ld.INN,
			PasswordFns: ld.PasswordFns,
			TimeOffset:  ld.TimeOffset,
		},
		Device: &kitpb.DeviceInfo{
			SourceDeviceId: dev.SourceDeviceID,
			SourceType:     dev.SourceType,
			AppVersion:     dev.AppVersion,
			MetaDetails: &kitpb.MetaDetails{
				UserAgent: dev.MetaDetails.UserAgent,
			},
		},
	}

	resp, err := c.client.Auth(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.DeadlineExceeded {
			return false, "error: превышено время запроса. Попробуйте позже"
		}
		return false, st.Message()
	}
	return resp.IsValid, resp.Message
}

func (c *GRPCClient) SingleKit(ctx context.Context, sr *domain.SingleRequest) ([]string, string) {
	req := &kitpb.SingleRequest{
		Date:     sr.Date,
		DeviceId: sr.DeviceID,
	}

	resp, err := c.client.SendToKitVend(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			return nil, st.Message()
		}
	}
	return resp.Data, resp.Message
}

func (c *GRPCClient) SingleFns(ctx context.Context, sr *domain.SingleRequest) string {
	req := &kitpb.SingleRequest{
		Date:     sr.Date,
		DataKit:  sr.DataKit,
		DeviceId: sr.DeviceID,
	}

	resp, err := c.client.SendToFns(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			return st.Message()
		}
	}
	return resp.Message
}

func (c *GRPCClient) Start(ctx context.Context, data *domain.AutoModeData) (string, error) {
	req := &kitpb.AutoModeData{
		CompanyID:   data.CompanyID,
		UserLogin:   data.UserLogin,
		PasswordKit: data.PasswordKit,
		INN:         data.INN,
		PasswordFns: data.PasswordFns,
		Device: &kitpb.DeviceInfo{
			SourceDeviceId: data.DeviceInfo.SourceDeviceID,
			SourceType:     data.DeviceInfo.SourceType,
			AppVersion:     data.DeviceInfo.AppVersion,
			MetaDetails: &kitpb.MetaDetails{
				UserAgent: data.DeviceInfo.MetaDetails.UserAgent,
			},
		},
		RefreshToken: data.RefreshToken,
		Token:        data.Token,
		Time: &kitpb.Time{
			Zone:   data.Time.Zone,
			Offset: data.Time.Offset,
		},
	}

	resp, err := c.client.Start(ctx, req)

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			return st.Message(), err
		}
		return err.Error(), err
	}

	return resp.Message, nil
}

func (c *GRPCClient) Stop(ctx context.Context, inn, login, company string) (string, error) {
	req := &kitpb.AutoModeStop{
		INN:       inn,
		CompanyId: company,
		UserLogin: login,
	}

	resp, err := c.client.Stop(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Message, nil
}

func (c *GRPCClient) Sales(ctx context.Context, inn, user, company, date, time string) ([]domain.Sale, error) {
	req := &kitpb.SalesReq{
		INN:       inn,
		UserLogin: user,
		CompanyId: company,
		Date:      date,
		Time:      time,
	}

	resp, err := c.client.Sales(ctx, req)
	if err != nil {
		return nil, err
	}

	var sales []domain.Sale
	for _, v := range resp.Sales {
		sales = append(sales, domain.Sale{
			DateTime:  v.DateTime,
			GoodsName: v.Goods,
		})
	}

	return sales, nil
}

func (c *GRPCClient) WorkingAutoMode(ctx context.Context, sr *domain.SingleRequest) (map[string]int32, string) {
	req := &kitpb.SingleRequest{}

	resp, err := c.client.WorkingAutoMode(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			return nil, st.Message()
		}
	}
	return resp.Values, ""
}

func (c *GRPCClient) GetStatus(ctx context.Context, inn, company, login string) (string, error) {

	req := &kitpb.InnLogin{
		INN:       inn,
		CompanyId: company,
		UserLogin: login,
	}

	resp, err := c.client.CustomerStatus(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Message, nil
}
