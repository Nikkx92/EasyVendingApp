//go:build !js || !wasm

package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gioui.org/app"
	bolt "go.etcd.io/bbolt"
)

type StateApp struct {
	DeviceId        string
	AutoMode        bool   `json:"AutoMode"`
	IsPaid          bool   `json:"IsPaid"`
	IsAuth          bool   `json:"IsAuth"`
	CompanyId       string `json:"CompanyId"`
	LoginKit        string `json:"LoginKit"`
	PassKit         string `json:"PassKit"`
	Inn             string `json:"Inn"`
	PassFns         string `json:"PassFns"`
	RefreshTokenFns string `json:"RefreshTokenFns"`
	TokenFns        string `json:"TokenFns"`
}

type DrinkSession struct {
	Date    string  `json:"Date"`
	Details Details `json:"Details"`
}

type Details struct {
	TimeEnd string   `json:"TimeEnd"`
	Drinks  []string `json:"Drinks"`
	Merged  bool     `json:"Merged"`
}

func SaveState(s *StateApp) {
	dir, err := Path()
	if err != nil {
		//ui.logger.Println(err)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		//ui.logger.Println(err)
	}

	data, err := json.Marshal(s)
	if err != nil {
		//ui.logger.Println(err)
	}

	filePath := filepath.Join(dir, "state.json")
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		//ui.logger.Println(err)
	}
}

func LoadState() *StateApp {
	dir, err := Path()
	if err != nil {
		//ui.logger.Println(err)
	}
	filePath := filepath.Join(dir, "state.json")

	var st StateApp

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.WriteFile(filePath, []byte("{}"), 0644); err != nil {
				fmt.Println(err)
			}
			return &st
		}
		return &st
	}

	if err = json.Unmarshal(data, &st); err != nil {
		fmt.Println(err)
	}

	return &st

}

func Path() (string, error) {
	dir, err := app.DataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "app_data"), nil
}

func GetSales(db *bolt.DB, limit int) ([]DrinkSession, error) {
	var sessions []DrinkSession

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("sales"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		c := b.Cursor()
		count := 0

		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			var s DrinkSession
			s.Date = string(k)
			if err := json.Unmarshal(v, &s.Details); err != nil {
				fmt.Println(err)
				return err
			}
			sessions = append(sessions, s)
			count++

			if count >= limit {
				break
			}
		}
		return nil
	})
	return sessions, err
}

func SaveSession(db *bolt.DB, session *DrinkSession) {
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("sales"))
		if err != nil {
			return err
		}

		data, err := json.Marshal(session.Details)
		if err != nil {
			return err
		}

		// Используем Date как ключ для O(1) доступа
		return b.Put([]byte(session.Date), data)
	})
	if err != nil {
		fmt.Println("err 108", err)
	}
}

func GetSessionByDate(db *bolt.DB, date string) *Details {
	var session Details
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("sales"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		data := b.Get([]byte(date))
		if data == nil {
			return fmt.Errorf("session not found")
		}

		return json.Unmarshal(data, &session)
	})
	if err != nil {
		return nil
	}
	return &session
}
