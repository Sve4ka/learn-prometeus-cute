package integration_test

import (
	"backend/internal/models"
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"strconv"

	//"github.com/ozontech/cute"

	//"github.com/ozontech/cute/asserts/json"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
	"time"
)

type tmp struct {
	User models.User `json:"user"`
}

type create struct {
	ID int `json:"id"`
}

func TestCreate(t *testing.T) {
	// Генерируем случайный email, чтобы тесты не конфликтовали
	noise := uuid.NewString()[:5]
	testUser := map[string]string{
		"name":     "test",
		"email":    "test" + noise + "@test.com",
		"sur_name": "test",
		"pwd":      noise,
	}

	t.Log("Создание запроса создания")
	createUserURL := "http://localhost:8000/user/"
	requestBody, err := json.Marshal(testUser)
	require.NoError(t, err, "Ошибка при сериализации JSON")

	req, err := http.NewRequest("POST", createUserURL, bytes.NewBuffer(requestBody))
	require.NoError(t, err, "Ошибка при создании HTTP-запроса")

	req.Header.Set("Content-Type", "application/json")

	t.Log("Отправляем запрос")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err, "Ошибка при выполнении HTTP-запроса")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Ожидался HTTP 200 при создании пользователя")

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Ошибка при чтении тела ответа")

	var response create
	err = json.Unmarshal(body, &response)
	require.NoError(t, err, "Ошибка при разборе JSON-ответа")

	t.Log("Создан пользователь с ID:", response.ID)

}

func TestCreateGet(t *testing.T) {
	// Генерируем случайный email, чтобы тесты не конфликтовали
	noise := uuid.NewString()[:5]
	testUser := map[string]string{
		"name":     "test",
		"email":    "test" + noise + "@test.com",
		"sur_name": "test",
		"pwd":      noise,
	}

	t.Log("Создание запроса создания")
	createUserURL := "http://localhost:8000/user/"
	requestBody, err := json.Marshal(testUser)
	require.NoError(t, err, "Ошибка при сериализации JSON")

	req, err := http.NewRequest("POST", createUserURL, bytes.NewBuffer(requestBody))
	require.NoError(t, err, "Ошибка при создании HTTP-запроса")

	req.Header.Set("Content-Type", "application/json")

	t.Log("Отправляем запрос")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err, "Ошибка при выполнении HTTP-запроса")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Ожидался HTTP 200 при создании пользователя")

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Ошибка при чтении тела ответа")

	var response create
	err = json.Unmarshal(body, &response)
	require.NoError(t, err, "Ошибка при разборе JSON-ответа")

	t.Log("Создан пользователь с ID:", response.ID)

	getUserURL := "http://localhost:8000/user/" + strconv.Itoa(response.ID)
	t.Log("Создание запроса чтения из ", getUserURL)
	req, err = http.NewRequest("GET", getUserURL, nil)
	require.NoError(t, err, "Ошибка при создании HTTP-запроса")

	t.Log("Отправляем запрос")
	resp, err = client.Do(req)
	require.NoError(t, err, "Ошибка при выполнении HTTP-запроса")
	defer resp.Body.Close()

	t.Log("Читаем ответ")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Ожидался HTTP 200 при получении пользователя")

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err, "Ошибка при чтении тела ответа")

	var user tmp
	err = json.Unmarshal(body, &user)
	require.NoError(t, err, "Ошибка при разборе JSON-ответа")

	assert.Equal(t, response.ID, user.User.ID, "ID пользователя в ответе не совпадает")
	assert.Equal(t, testUser["name"], user.User.Name, "Имя пользователя не совпадает")
	assert.Equal(t, testUser["sur_name"], user.User.SurName, "Фамилия пользователя не совпадает")
	assert.Equal(t, testUser["email"], user.User.Email, "Email пользователя не совпадает")
	t.Log("Пользователь прочитан")
}
