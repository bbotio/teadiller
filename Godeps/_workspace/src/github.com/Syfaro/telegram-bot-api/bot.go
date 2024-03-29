// Package tgbotapi has bindings for interacting with the Telegram Bot API.
package tgbotapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"teadiller/Godeps/_workspace/src/github.com/technoweenie/multipartstreamer"
	"time"
)

// BotAPI has methods for interacting with all of Telegram's Bot API endpoints.
type BotAPI struct {
	Token  string       `json:"token"`
	Debug  bool         `json:"debug"`
	Self   User         `json:"-"`
	Client *http.Client `json:"-"`
}

// NewBotAPI creates a new BotAPI instance.
// Requires a token, provided by @BotFather on Telegram
func NewBotAPI(token string) (*BotAPI, error) {
	return NewBotAPIWithClient(token, &http.Client{})
}

// NewBotAPIWithClient creates a new BotAPI instance passing an http.Client.
// Requires a token, provided by @BotFather on Telegram
func NewBotAPIWithClient(token string, client *http.Client) (*BotAPI, error) {
	bot := &BotAPI{
		Token:  token,
		Client: client,
	}

	self, err := bot.GetMe()
	if err != nil {
		return &BotAPI{}, err
	}

	bot.Self = self

	return bot, nil
}

// MakeRequest makes a request to a specific endpoint with our token.
// All requests are POSTs because Telegram doesn't care, and it's easier.
func (bot *BotAPI) MakeRequest(endpoint string, params url.Values) (APIResponse, error) {
	resp, err := bot.Client.PostForm(fmt.Sprintf(APIEndpoint, bot.Token, endpoint), params)
	if err != nil {
		return APIResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return APIResponse{}, errors.New(APIForbidden)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return APIResponse{}, err
	}

	if bot.Debug {
		log.Println(endpoint, string(bytes))
	}

	var apiResp APIResponse
	json.Unmarshal(bytes, &apiResp)

	if !apiResp.Ok {
		return APIResponse{}, errors.New(apiResp.Description)
	}

	return apiResp, nil
}

func (bot *BotAPI) makeMessageRequest(endpoint string, params url.Values) (Message, error) {
	resp, err := bot.MakeRequest(endpoint, params)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	bot.debugLog(endpoint, params, message)

	return message, nil
}

// UploadFile makes a request to the API with a file.
//
// Requires the parameter to hold the file not be in the params.
// File should be a string to a file path, a FileBytes struct, or a FileReader struct.
func (bot *BotAPI) UploadFile(endpoint string, params map[string]string, fieldname string, file interface{}) (APIResponse, error) {
	ms := multipartstreamer.New()
	ms.WriteFields(params)

	switch f := file.(type) {
	case string:
		fileHandle, err := os.Open(f)
		if err != nil {
			return APIResponse{}, err
		}
		defer fileHandle.Close()

		fi, err := os.Stat(f)
		if err != nil {
			return APIResponse{}, err
		}

		ms.WriteReader(fieldname, fileHandle.Name(), fi.Size(), fileHandle)
	case FileBytes:
		buf := bytes.NewBuffer(f.Bytes)
		ms.WriteReader(fieldname, f.Name, int64(len(f.Bytes)), buf)
	case FileReader:
		if f.Size != -1 {
			ms.WriteReader(fieldname, f.Name, f.Size, f.Reader)

			break
		}

		data, err := ioutil.ReadAll(f.Reader)
		if err != nil {
			return APIResponse{}, err
		}

		buf := bytes.NewBuffer(data)

		ms.WriteReader(fieldname, f.Name, int64(len(data)), buf)
	default:
		return APIResponse{}, errors.New("bad file type")
	}

	req, err := http.NewRequest("POST", fmt.Sprintf(APIEndpoint, bot.Token, endpoint), nil)
	if err != nil {
		return APIResponse{}, err
	}

	ms.SetupRequest(req)

	res, err := bot.Client.Do(req)
	if err != nil {
		return APIResponse{}, err
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return APIResponse{}, err
	}

	if bot.Debug {
		log.Println(string(bytes))
	}

	var apiResp APIResponse
	json.Unmarshal(bytes, &apiResp)

	if !apiResp.Ok {
		return APIResponse{}, errors.New(apiResp.Description)
	}

	return apiResp, nil
}

// GetFileDirectURL returns direct URL to file
//
// Requires fileID
func (bot *BotAPI) GetFileDirectURL(fileID string) (string, error) {
	file, err := bot.GetFile(FileConfig{fileID})

	if err != nil {
		return "", err
	}

	return file.Link(bot.Token), nil
}

// GetMe fetches the currently authenticated bot.
//
// There are no parameters for this method.
func (bot *BotAPI) GetMe() (User, error) {
	resp, err := bot.MakeRequest("getMe", nil)
	if err != nil {
		return User{}, err
	}

	var user User
	json.Unmarshal(resp.Result, &user)

	bot.debugLog("getMe", nil, user)

	return user, nil
}

// IsMessageToMe returns true if message directed to this bot
//
// Requires message
func (bot *BotAPI) IsMessageToMe(message Message) bool {
	return strings.Contains(message.Text, "@"+bot.Self.UserName)
}

// Send will send event(Message, Photo, Audio, ChatAction, anything) to Telegram
//
// Requires Chattable
func (bot *BotAPI) Send(c Chattable) (Message, error) {
	switch c.(type) {
	case Fileable:
		return bot.sendFile(c.(Fileable))
	default:
		return bot.sendChattable(c)
	}
}

func (bot *BotAPI) debugLog(context string, v url.Values, message interface{}) {
	if bot.Debug {
		log.Printf("%s req : %+v\n", context, v)
		log.Printf("%s resp: %+v\n", context, message)
	}
}

func (bot *BotAPI) sendExisting(method string, config Fileable) (Message, error) {
	v, err := config.Values()

	if err != nil {
		return Message{}, err
	}

	message, err := bot.makeMessageRequest(method, v)
	if err != nil {
		return Message{}, err
	}

	return message, nil
}

func (bot *BotAPI) uploadAndSend(method string, config Fileable) (Message, error) {
	params, err := config.Params()
	if err != nil {
		return Message{}, err
	}

	file := config.GetFile()

	resp, err := bot.UploadFile(method, params, config.Name(), file)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	if bot.Debug {
		log.Printf("%s resp: %+v\n", method, message)
	}

	return message, nil
}

func (bot *BotAPI) sendFile(config Fileable) (Message, error) {
	if config.UseExistingFile() {
		return bot.sendExisting(config.Method(), config)
	}

	return bot.uploadAndSend(config.Method(), config)
}

func (bot *BotAPI) sendChattable(config Chattable) (Message, error) {
	v, err := config.Values()
	if err != nil {
		return Message{}, err
	}

	message, err := bot.makeMessageRequest(config.Method(), v)

	if err != nil {
		return Message{}, err
	}

	return message, nil
}

// GetUserProfilePhotos gets a user's profile photos.
//
// Requires UserID.
// Offset and Limit are optional.
func (bot *BotAPI) GetUserProfilePhotos(config UserProfilePhotosConfig) (UserProfilePhotos, error) {
	v := url.Values{}
	v.Add("user_id", strconv.Itoa(config.UserID))
	if config.Offset != 0 {
		v.Add("offset", strconv.Itoa(config.Offset))
	}
	if config.Limit != 0 {
		v.Add("limit", strconv.Itoa(config.Limit))
	}

	resp, err := bot.MakeRequest("getUserProfilePhotos", v)
	if err != nil {
		return UserProfilePhotos{}, err
	}

	var profilePhotos UserProfilePhotos
	json.Unmarshal(resp.Result, &profilePhotos)

	bot.debugLog("GetUserProfilePhoto", v, profilePhotos)

	return profilePhotos, nil
}

// GetFile returns a file_id required to download a file.
//
// Requires FileID.
func (bot *BotAPI) GetFile(config FileConfig) (File, error) {
	v := url.Values{}
	v.Add("file_id", config.FileID)

	resp, err := bot.MakeRequest("getFile", v)
	if err != nil {
		return File{}, err
	}

	var file File
	json.Unmarshal(resp.Result, &file)

	bot.debugLog("GetFile", v, file)

	return file, nil
}

// GetUpdates fetches updates.
// If a WebHook is set, this will not return any data!
//
// Offset, Limit, and Timeout are optional.
// To not get old items, set Offset to one higher than the previous item.
// Set Timeout to a large number to reduce requests and get responses instantly.
func (bot *BotAPI) GetUpdates(config UpdateConfig) ([]Update, error) {
	v := url.Values{}
	if config.Offset > 0 {
		v.Add("offset", strconv.Itoa(config.Offset))
	}
	if config.Limit > 0 {
		v.Add("limit", strconv.Itoa(config.Limit))
	}
	if config.Timeout > 0 {
		v.Add("timeout", strconv.Itoa(config.Timeout))
	}

	resp, err := bot.MakeRequest("getUpdates", v)
	if err != nil {
		return []Update{}, err
	}

	var updates []Update
	json.Unmarshal(resp.Result, &updates)

	if bot.Debug {
		log.Printf("getUpdates: %+v\n", updates)
	}

	return updates, nil
}

// RemoveWebhook removes webhook
//
// There are no parameters for this method.
func (bot *BotAPI) RemoveWebhook() (APIResponse, error) {
	return bot.MakeRequest("setWebhook", url.Values{})
}

// SetWebhook sets a webhook.
// If this is set, GetUpdates will not get any data!
//
// Requires URL OR to set Clear to true.
func (bot *BotAPI) SetWebhook(config WebhookConfig) (APIResponse, error) {
	if config.Certificate == nil {
		v := url.Values{}
		v.Add("url", config.URL.String())

		return bot.MakeRequest("setWebhook", v)
	}

	params := make(map[string]string)
	params["url"] = config.URL.String()

	resp, err := bot.UploadFile("setWebhook", params, "certificate", config.Certificate)
	if err != nil {
		return APIResponse{}, err
	}

	var apiResp APIResponse
	json.Unmarshal(resp.Result, &apiResp)

	if bot.Debug {
		log.Printf("setWebhook resp: %+v\n", apiResp)
	}

	return apiResp, nil
}

// GetUpdatesChan starts and returns a channel for getting updates.
//
// Requires UpdateConfig
func (bot *BotAPI) GetUpdatesChan(config UpdateConfig) (<-chan Update, error) {
	updatesChan := make(chan Update, 100)

	go func() {
		for {
			updates, err := bot.GetUpdates(config)
			if err != nil {
				log.Println(err)
				log.Println("Failed to get updates, retrying in 3 seconds...")
				time.Sleep(time.Second * 3)

				continue
			}

			for _, update := range updates {
				if update.UpdateID >= config.Offset {
					config.Offset = update.UpdateID + 1
					updatesChan <- update
				}
			}
		}
	}()

	return updatesChan, nil
}

// ListenForWebhook registers a http handler for a webhook.
func (bot *BotAPI) ListenForWebhook(pattern string) (<-chan Update, http.Handler) {
	updatesChan := make(chan Update, 100)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)

		var update Update
		json.Unmarshal(bytes, &update)

		updatesChan <- update
	})

	http.HandleFunc(pattern, handler)

	return updatesChan, handler
}

// AnswerInlineQuery sends a response to an inline query.
func (bot *BotAPI) AnswerInlineQuery(config InlineConfig) (APIResponse, error) {
	v := url.Values{}

	v.Add("inline_query_id", config.InlineQueryID)
	v.Add("cache_time", strconv.Itoa(config.CacheTime))
	v.Add("is_personal", strconv.FormatBool(config.IsPersonal))
	v.Add("next_offset", config.NextOffset)
	data, err := json.Marshal(config.Results)
	if err != nil {
		return APIResponse{}, err
	}
	v.Add("results", string(data))

	return bot.MakeRequest("answerInlineQuery", v)
}
