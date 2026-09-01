package httpNet

/*func SendStartRequest(data g.SingleRequest) string {
	req, err := http.NewRequest("GET", "http://192.168.1.88:8080/start",
		nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+data.Jwt)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		//log.Println("Ошибка запуска:", err)
		return "Истекло время запроса. Возможно, проблемы с сервисом/сетью. Попробуйте позже"
	}
	defer resp.Body.Close()

	var res struct {
		Message string `json:"message"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "ошибка декодирования ответа"
	} else {
		return res.Message
	}

}

func SendStopRequest(data g.SingleRequest) string {
	req, _ := http.NewRequest("GET", "http://192.168.1.88:8080/stop",
		nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+data.Jwt)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		//log.Println("Ошибка остановки:", err)
		return "Истекло время запроса. Возможно, проблемы с сервисом/сетью. Попробуйте позже"
	}
	defer resp.Body.Close()

	var res struct {
		Message string `json:"message"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "ошибка декодирования ответа"
	} else {
		return res.Message
	}

	//log.Println("Горутина остановлена на сервере")
}*/
