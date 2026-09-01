package httpNet

import (
	"bytes"
	"crypto/md5"
	"easyVending/internal/domain"
	"easyVending/internal/errs"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
	type ResponseKit struct {
		Data     []string `json:"data"`
		Buf      string   `json:"buffer"`
		ErrorKit string
	}

	type SingleRequest struct {
		Jwt      string   `json:"jwt"`
		Date     string   `json:"date"`
		DataKit  []string `json:"dataKit"`
		AutoMode bool     `json:"autoMode"`
		IsPaid   bool     `json:"isPaid"`
	}

func ToKit(data *g.SingleRequest) ([]string, string) {

		requer := SingleRequest{
			Date: data.Date,
		}

		body, err := json.Marshal(requer)
		if err != nil {
			//logger.Println(err)
		}
		req, err := httpNet.NewRequest("POST", "http://192.168.1.88:8080/sendToKitVend", bytes.NewReader(body))
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+data.Jwt)
		client := &httpNet.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return nil, "Истекло время запроса. Возможно, проблемы с сервисом/сетью. Попробуйте позже"

		}
		defer resp.Body.Close()

		//kitError := resp.Header.Get("Error-Kit")

		var res struct {
			Data    []string `json:"data"`
			Message string   `json:"message"`
		}
		if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return nil, fmt.Sprintf("%s", err)
		} else {
			return res.Data, res.Message
		}
	}

	func GetSales(data *g.SingleRequest) (map[string]map[string]int32, string) {
		req, err := httpNet.NewRequest("GET", "http://192.168.1.88:8080/sales", nil)
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+data.Jwt)

		client := &httpNet.Client{
			Timeout: 5 * time.Second,
		}
		resp, err := client.Do(req)
		var res struct {
			Data    map[string]map[string]int32 `json:"data"`
			Message string                      `json:"message"`
		}
		if err != nil {
			return nil, "Истекло время запроса. Возможно, проблемы с сервисом/сетью. Попробуйте позже"
			//fmt.Println(err)
			/*return ResponseFns{
				Message: "Истекло время запроса. Возможно, проблемы с сервисом/сетью. Попробуйте позже",
			}
		}
		defer resp.Body.Close()

		if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
			fmt.Println(err)
			return nil, "can't decode response"
			//logger.Println(err)
		} else {
			return res.Data, ""
		}
	}
*/
const baseUrlKit = "https://api2.kit-invest.ru"

type auth struct {
	CompanyId string `json:"companyid"`
	RequestId string `json:"requestid"`
	UserLogin string `json:"userlogin"`
	Sign      string `json:"sign"`
}

type filter struct {
	UpDate string `json:"update"`
	ToDate string `json:"todate"`
}

type request struct {
	Auth   auth   `json:"Auth"`
	Filter filter `json:"Filter"`
}

type response struct {
	ErrorMessage string        `json:"ErrorMessage"`
	ResultCode   int           `json:"ResultCode"`
	Sales        []domain.Sale `json:"Sales"`
}

func hashing(c, p string) (h, u string) {
	uniqueNumb := strconv.FormatInt(time.Now().UnixNano(), 10)
	data := c + p + uniqueNumb
	hash := md5.Sum([]byte(data))
	sign := hex.EncodeToString(hash[:])
	return sign, uniqueNumb
}

func (c *ClientHTTP) LoginKitVending(companyId, userLogin, pass string) error {
	sign, uniqNum := hashing(companyId, pass)
	r := request{
		Auth: auth{
			CompanyId: companyId,
			RequestId: uniqNum,
			UserLogin: userLogin,
			Sign:      sign,
		},
	}
	jsonData, err := json.Marshal(r)
	if err != nil {
		return errs.Technic
	}

	req, err := http.NewRequest(
		"POST",
		baseUrlKit+"/APIService.svc/GetModems",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return errs.Technic
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return errs.NoConnect
	}
	defer resp.Body.Close()

	var res response
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return errs.Technic
	}
	if res.ErrorMessage != "" {
		return errors.New(res.ErrorMessage)
	}
	return nil
}

func (c *ClientHTTP) GetDataKitVending(companyId, userLogin, password, date string) ([]domain.Sale, error) {
	sign, uniqNum := hashing(companyId, password)
	var upDate string
	var toDate string
	sep := strings.Split(date, "--")
	upDate = sep[0]
	toDate = sep[1]

	requ := request{
		Auth: auth{
			CompanyId: companyId,
			RequestId: uniqNum,
			UserLogin: userLogin,
			Sign:      sign,
		},
		Filter: filter{
			UpDate: upDate,
			ToDate: toDate,
		},
	}

	jsonData, err := json.Marshal(requ)
	if err != nil {
		return nil, errs.Technic
	}

	req, err := http.NewRequest(
		"POST",
		baseUrlKit+"/APIService.svc/GetSales",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, errs.Technic
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, errs.NoConnect
	}
	defer resp.Body.Close()

	var r response

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, errs.Technic
	}

	if r.ResultCode != 0 {
		return nil, errors.New(r.ErrorMessage)
	}

	var sales []domain.Sale

	for i := range r.Sales {
		sepDate := strings.Split(r.Sales[i].DateTime, " ")
		r.Sales[i].DateTime = sepDate[0]
		sales = append(sales, r.Sales[i])
	}

	if len(sales) == 0 {
		return nil, nil
	}

	return sales, nil
}
