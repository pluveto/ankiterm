// AnkiConnect Client
package ankicc

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Resp struct {
	Error  *string     `json:"error,omitempty"`
	Result interface{} `json:"result"`
}

type RespDeckGetDeckConfig struct {
	Result struct {
		ID                   int         `json:"id"`
		Mod                  int         `json:"mod"`
		Name                 string      `json:"name"`
		USN                  int         `json:"usn"`
		MaxTaken             int         `json:"max_taken"`
		Autoplay             bool        `json:"autoplay"`
		Timer                int         `json:"timer"`
		Replayq              bool        `json:"replayq"`
		New                  interface{} `json:"new"`
		Rev                  interface{} `json:"rev"`
		Lapse                interface{} `json:"lapse"`
		Dyn                  bool        `json:"dyn"`
		NewMix               int         `json:"new_mix"`
		NewPerDayMinimum     int         `json:"new_per_day_minimum"`
		InterdayLearningMix  int         `json:"interday_learning_mix"`
		ReviewOrder          int         `json:"review_order"`
		NewSortOrder         int         `json:"new_sort_order"`
		NewGatherPriority    int         `json:"new_gather_priority"`
		BuryInterdayLearning bool        `json:"bury_interday_learning"`
		Reminder             interface{} `json:"reminder"`
	} `json:"result"`
}

type DeckStat struct {
	DeckID      int    `json:"deck_id"`
	Name        string `json:"name"`
	NewCount    int    `json:"new_count"`
	LearnCount  int    `json:"learn_count"`
	ReviewCount int    `json:"review_count"`
	TotalInDeck int    `json:"total_in_deck"`
}

type RespDeckGetDeckStatsResult map[string]DeckStat

type RespDeckGetDeckStats struct {
	Result RespDeckGetDeckStatsResult `json:"result"`
}

type RespMediaRetrieveMediaFile struct {
	Result string `json:"result"`
}

type RespMiscellaneousVersion struct {
	Result string `json:"result"`
}

type RespGraphicalDeckNames struct {
	Result []string `json:"result"`
}

type CurrentCard struct {
	Answer     string `json:"answer"`
	Buttons    []int  `json:"buttons"`
	CardID     int64  `json:"cardId"`
	CSS        string `json:"css"`
	DeckName   string `json:"deckName"`
	FieldOrder int    `json:"fieldOrder"`
	Fields     struct {
		Back struct {
			Order int    `json:"order"`
			Value string `json:"value"`
		} `json:"Back"`
		Front struct {
			Order int    `json:"order"`
			Value string `json:"value"`
		} `json:"Front"`
	} `json:"fields"`
	ModelName   string   `json:"modelName"`
	NextReviews []string `json:"nextReviews"`
	Question    string   `json:"question"`
	Template    string   `json:"template"`
}

type RespGraphicalGuiCurrentCard struct {
	Result CurrentCard `json:"result"`
}

type RpcError struct {
	Message string
	Code    int
}

type Client struct {
	BaseURL string
}

func (e *RpcError) Error() string {
	return e.Message
}

// C.Request to AnkiConnect
func (c Client) request(body map[string]interface{}) (interface{}, error) {
	body["version"] = 6
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(c.BaseURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respData Resp
	err = json.Unmarshal(data, &respData)
	if err != nil {
		return nil, err
	}

	if respData.Error != nil && *respData.Error != "" {
		return nil, &RpcError{Message: *respData.Error}
	}

	// AnkiConnect has a rate limit
	time.Sleep(time.Duration(100) * time.Millisecond)
	return respData.Result, nil
}

func (c Client) DeckNames() ([]string, error) {
	body := map[string]interface{}{
		"action": "deckNames",
	}
	result, err := c.request(body)
	if err != nil {
		return nil, err
	}

	var respData RespGraphicalDeckNames
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonData, &respData)
	if err != nil {
		return nil, err
	}

	return respData.Result, nil
}

func (c Client) GetDeckConfig(deck string) (RespDeckGetDeckConfig, error) {
	body := map[string]interface{}{
		"action": "getDeckConfig",
		"params": map[string]string{
			"deck": deck,
		},
	}
	result, err := c.request(body)
	if err != nil {
		return RespDeckGetDeckConfig{}, err
	}

	var respData RespDeckGetDeckConfig
	jsonData, err := json.Marshal(result)
	if err != nil {
		return RespDeckGetDeckConfig{}, err
	}
	err = json.Unmarshal(jsonData, &respData)
	if err != nil {
		return RespDeckGetDeckConfig{}, err
	}

	return respData, nil
}

func (c Client) GetDeckStats(decks []string) (RespDeckGetDeckStats, error) {
	body := map[string]interface{}{
		"action": "getDeckStats",
		"params": map[string][]string{
			"decks": decks,
		},
	}
	result, err := c.request(body)
	if err != nil {
		return RespDeckGetDeckStats{}, err
	}

	var respData RespDeckGetDeckStats
	jsonData, err := json.Marshal(result)
	if err != nil {
		return RespDeckGetDeckStats{}, err
	}
	err = json.Unmarshal(jsonData, &respData)
	if err != nil {
		return RespDeckGetDeckStats{}, err
	}

	return respData, nil
}

func (c Client) GetDeckStat(deck string) (DeckStat, error) {
	result, err := c.GetDeckStats([]string{deck})
	if err != nil {
		return DeckStat{}, err
	}

	return result.Result[deck], nil
}

func (c Client) RetrieveMediaFile(fileName string) (string, error) {
	body := map[string]interface{}{
		"action": "retrieveMediaFile",
		"params": map[string]string{
			"filename": fileName,
		},
	}
	result, err := c.request(body)
	if err != nil {
		return "", err
	}

	var respData RespMediaRetrieveMediaFile
	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(jsonData, &respData)
	if err != nil {
		return "", err
	}

	return respData.Result, nil
}

func (c Client) Version() (string, error) {
	body := map[string]interface{}{
		"action": "version",
	}
	result, err := c.request(body)
	if err != nil {
		return "", err
	}

	var respData RespMiscellaneousVersion
	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(jsonData, &respData)
	if err != nil {
		return "", err
	}

	return respData.Result, nil
}

func (c Client) Sync() error {
	body := map[string]interface{}{
		"action": "sync",
	}
	_, err := c.request(body)
	return err
}

func (c Client) GuiDeckReview(name string) error {
	body := map[string]interface{}{
		"action": "guiDeckReview",
		"params": map[string]string{
			"name": name,
		},
	}
	_, err := c.request(body)
	return err
}

func (c Client) GuiCurrentCard() (*CurrentCard, error) {
	body := map[string]interface{}{
		"action": "guiCurrentCard",
	}
	result, err := c.request(body)
	if err != nil {
		return nil, err
	}

	var respData CurrentCard
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonData, &respData)
	if err != nil {
		return nil, err
	}

	return &respData, nil
}

func (c Client) GuiShowAnswer() (err error) {
	body := map[string]interface{}{
		"action": "guiShowAnswer",
	}
	_, err = c.request(body)
	return
}

func (c Client) GuiAnswerCard(ease int) (err error) {
	body := map[string]interface{}{
		"action": "guiAnswerCard",
		"params": map[string]int{
			"ease": ease,
		},
	}
	_, err = c.request(body)
	return
}
