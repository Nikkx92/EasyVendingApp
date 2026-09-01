package httpNet

import ()

/*type Request struct {
	Data   g.LoginData       `json:"data"`
	Device domain.DeviceInfo `json:"Device"`
}
type AuthResponse struct {
	IsValid bool   `json:"isValid"`
	Message string `json:"message"`
}

func Login(data *g.LoginData, device *domain.DeviceInfo) (bool, string) {
	requer := Request{
		Data:   *data,
		Device: *device,
	}
	body, err := json.Marshal(requer)
	if err != nil {
		//logger.Println(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://192.168.1.88:8080/auth", bytes.NewReader(body))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, "Истекло время запроса. Возможно, проблемы с сервисом/сетью. Попробуйте позже"
	}
	defer resp.Body.Close()

	var res AuthResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		fmt.Println(err)
		return false, "can't decode response"
	} else {
		return res.IsValid, res.Message
	}
}*/
