package ui

import (
	"easyVending/internal/errs"

	"gioui.org/layout"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ui *UI) InitLoadIcon() {
	ui.isLoad = true
	go ui.InsertStatus()
}

func (ui *UI) InsertStatus() {
	data := <-ui.a.DataChan

	ui.mu.Lock()
	ui.hasStatusData = true
	ui.statusData = data
	ui.mu.Unlock()
}

func (ui *UI) SendStatusToUI(gtx layout.Context) bool {

	defer func() {
		ui.mu.Lock()
		ui.hasStatusData = false
		ui.mu.Unlock()
	}()

	ui.mu.Lock()
	has := ui.hasStatusData
	res := ui.statusData
	ui.mu.Unlock()

	if has {
		ui.isLoad = false

		ui.getSales()

		if res.Err != nil {
			if st, ok := status.FromError(res.Err); ok {
				switch st.Code() {
				case codes.Unavailable, codes.DeadlineExceeded:
					ui.modalResponse(gtx, nil, errs.NoConnect.Error())
				default:
					ui.modalResponse(gtx, nil, st.Message())
				}
			}

			return has
		}

		if res.Message != "" {
			ui.modalResponse(gtx, nil, "автоматическая отправка чеков была отключена из-за возникшей ошибки: "+res.Message)

			ui.AutoSwt.Value = false
		}

		return has
	}

	return has
}
