package httpNet

import (
	"bytes"
	"easyVending/internal/domain"
	"easyVending/internal/errs"
	"easyVending/internal/storage"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

/*
	type RequestFns struct {
		Date    string   `json:"date"`
		DataKit []string `json:"dataKit"`
	}

	type ResponseFns struct {
		Message string `json:"Message"`
	}

	func ToFns(data *g.SingleRequest) ResponseFns {
		requer := SingleRequest{
			Date:    data.Date,
			DataKit: data.DataKit,
		}

		body, err := json.Marshal(requer)
		if err != nil {
			//logger.Println(err)
		}
		req, err := http.NewRequest("POST", "http://192.168.1.88:8080/sendToFns", bytes.NewReader(body))
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+data.Jwt)
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return ResponseFns{
				Message: "Истекло время запроса. Возможно, проблемы с сервисом/сетью. Попробуйте позже",
			}
		}
		defer resp.Body.Close()

		//refreshToken := resp.Header.Get("Refresh-Token")
		//token := resp.Header.Get("Token")

		var res ResponseFns
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			//logger.Println(err)
			return res
		} else {
			return ResponseFns{
				Message: res.Message,
				//Buf:     res.Buf,
				//RefreshToken: refreshToken,
				//Token:        token,
			}
		}
	}
*/
const baseUrlFns = "https://lknpd.nalog.ru"

type refreshTokenRequest struct {
	DeviceInfo domain.DeviceInfo `json:"deviceInfo"`
	Username   string            `json:"username"`
	Password   string            `json:"password"`
}

type tokensResponse struct {
	Code         string `json:"code"`
	Message      string `json:"message"`
	RefreshToken string `json:"refreshToken"`
	Token        string `json:"token"`
}

type checkRequest struct {
	OperationTime                   string    `json:"operationTime"`
	RequestTime                     string    `json:"requestTime"`
	Services                        []service `json:"services"`
	TotalAmount                     string    `json:"totalAmount"`
	Client                          clientFns `json:"client"`
	PaymentType                     string    `json:"paymentType"`
	IgnoreMaxTotalIncomeRestriction bool      `json:"ignoreMaxTotalIncomeRestriction"`
}

type clientFns struct {
	ContactPhone *string `json:"contactPhone"`
	DisplayName  *string `json:"displayName"`
	Inn          *string `json:"inn"`
	IncomeType   string  `json:"incomeType"`
}

type service struct {
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	Quantity int     `json:"quantity"`
}

type tokenRequest struct {
	DeviceInfo   domain.DeviceInfo `json:"deviceInfo"`
	RefreshToken string            `json:"refreshToken"`
}

func (c *ClientHTTP) LoginFns(inn, pass string, device *domain.DeviceInfo) (string, string, error) {
	payload := refreshTokenRequest{
		DeviceInfo: *device,
		Username:   inn,
		Password:   pass,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", "", errs.Technic
	}

	req, err := http.NewRequest(
		"POST",
		baseUrlFns+"/api/v1/auth/lkfl",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", "", errs.Technic
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.hc.Do(req)
	if err != nil {
		return "", "", errs.NoConnect
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", errs.Technic
	}

	var res tokensResponse
	if err = json.Unmarshal(body, &res); err != nil {
		return "", "", errs.Technic
	}

	if resp.StatusCode != 200 {
		return "", "", errors.New(res.Message)
	}

	return res.RefreshToken, res.Token, nil
}

func (c *ClientHTTP) SendSalesToFns(state *storage.StateApp, drinks []domain.Sale) error {
	/*duplicates := make(map[string]int)
	for i := range drinks {
		duplicates[drinks[i].GoodsName+":"+strconv.Itoa(int(drinks[i].Sum))]++
	}

	for i := range duplicates {
		sep := strings.Split(i, ":")
		price, err := strconv.Atoi(sep[1])
		if err != nil {
			return errs.Technic
		}
		nf := float64(price)

		payload := checkRequest{
			OperationTime: time.Now().Format(time.RFC3339),
			RequestTime:   time.Now().Format(time.RFC3339),
			Services: []service{
				{Name: sep[0],
					Amount:   nf,
					Quantity: duplicates[i]},
			},
			TotalAmount: strconv.Itoa(price * duplicates[i]),
			Client: clientFns{
				ContactPhone: nil,
				DisplayName:  nil,
				Inn:          nil,
				IncomeType:   "FROM_INDIVIDUAL",
			},
			PaymentType:                     "CASH",
			IgnoreMaxTotalIncomeRestriction: false,
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return errs.Technic
		}

		req, err := http.NewRequest(
			"POST",
			baseUrlFns+"/api/v1/income",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			return errs.Technic
		}

		req.Header.Set("Authorization", "Bearer "+state.TokenFns)
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.hc.Do(req)
		if err != nil {
			return errs.NoConnect
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return errs.FnsTokenErr
		}
	}*/

	return nil
}

func (c *ClientHTTP) GetToken(refToken string, device *domain.DeviceInfo) (string, error) {
	payload := tokenRequest{
		DeviceInfo:   *device,
		RefreshToken: refToken,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", errs.Technic
	}

	//url := "https://lknpd.nalog.ru/api/v1/auth/token"

	req, err := http.NewRequest(
		"POST",
		baseUrlFns+"/api/v1/auth/token",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", errs.Technic
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.hc.Do(req)
	if err != nil {
		return "", errs.NoConnect
	}
	defer resp.Body.Close()

	var res tokensResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", errs.Technic
	}

	if resp.StatusCode != 200 {
		return "", errs.FnsRefreshTokenErr
	}

	return res.Token, nil
}
