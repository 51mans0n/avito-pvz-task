package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestFullFlow(t *testing.T) {
	baseURL := "http://localhost:8080"
	client := &http.Client{Timeout: 5 * time.Second}

	// 1. dummyLogin -> moderator выбран
	modToken, err := getDummyLoginToken(client, baseURL, "moderator")
	if err != nil {
		t.Fatalf("cannot get mod token: %v", err)
	}

	// 2. Создаём pvz
	pvzID, err := createPVZ(client, baseURL, modToken, "Москва")
	if err != nil {
		t.Fatalf("cannot create PVZ: %v", err)
	}
	t.Logf("Created PVZ: %s", pvzID)

	// 3. dummyLogin -> employee выбран
	empToken, err := getDummyLoginToken(client, baseURL, "employee")
	if err != nil {
		t.Fatalf("cannot get employee token: %v", err)
	}

	// 4. создаём ресепшн
	receptionID, err := createReception(client, baseURL, empToken, pvzID)
	if err != nil {
		t.Fatalf("cannot create reception: %v", err)
	}
	t.Logf("Created Reception: %s", receptionID)

	// 5. добавляем 50 продуктов
	for i := 1; i <= 50; i++ {
		if err := createProduct(client, baseURL, empToken, pvzID, "электроника"); err != nil {
			t.Fatalf("cannot create product #%d: %v", i, err)
		}
	}
	t.Logf("Created 50 products")

	// 6. удаление ласт продукта
	if err := deleteLastProduct(client, baseURL, empToken, pvzID); err != nil {
		t.Fatalf("delete last product failed: %v", err)
	}
	t.Logf("Deleted last product")

	// 7. закрываем ресепшн
	if err := closeReception(client, baseURL, empToken, pvzID); err != nil {
		t.Fatalf("close reception failed: %v", err)
	}
	t.Logf("Reception closed successfully")
}

func getDummyLoginToken(c *http.Client, baseURL, role string) (string, error) {
	body := fmt.Sprintf(`{"role":"%s"}`, role)
	req, err := http.NewRequest(http.MethodPost, baseURL+"/dummyLogin", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("dummyLogin status=%d", resp.StatusCode)
	}
	var data struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	return data.Token, nil
}

func createPVZ(c *http.Client, baseURL, token, city string) (string, error) {
	body := fmt.Sprintf(`{"city":"%s"}`, city)
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/pvz", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("createPVZ status=%d", resp.StatusCode)
	}
	var data struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	return data.ID, nil
}

func createReception(c *http.Client, baseURL, token, pvzID string) (string, error) {
	body := fmt.Sprintf(`{"pvzId":"%s"}`, pvzID)
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/receptions", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("createReception status=%d", resp.StatusCode)
	}
	var rec struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rec); err != nil {
		return "", err
	}
	return rec.ID, nil
}

func createProduct(c *http.Client, baseURL, token, pvzID, prodType string) error {
	body := fmt.Sprintf(`{"type":"%s","pvzId":"%s"}`, prodType, pvzID)
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/products", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("createProduct status=%d", resp.StatusCode)
	}
	return nil
}

func deleteLastProduct(c *http.Client, baseURL, token, pvzID string) error {
	url := fmt.Sprintf("%s/pvz/%s/delete_last_product", baseURL, pvzID)
	req, _ := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("deleteLastProduct status=%d", resp.StatusCode)
	}
	return nil
}

func closeReception(c *http.Client, baseURL, token, pvzID string) error {
	url := fmt.Sprintf("%s/pvz/%s/close_last_reception", baseURL, pvzID)
	req, _ := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("closeReception status=%d", resp.StatusCode)
	}
	return nil
}
