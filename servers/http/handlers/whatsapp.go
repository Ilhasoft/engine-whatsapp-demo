package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/weni/whatsapp-router/config"
	"github.com/weni/whatsapp-router/logger"
	"github.com/weni/whatsapp-router/models"
	"github.com/weni/whatsapp-router/services"
)

type WhatsappHandler struct {
	ContactService services.ContactService
	ChannelService services.ChannelService
	CourierService services.CourierService
}

func (h *WhatsappHandler) HandleIncomingRequests(w http.ResponseWriter, r *http.Request) {
	incomingMsg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error("unexpected server error - " + err.Error())
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, err.Error())
		return
	}

	incomingContact := incomingMsgToContact(string(incomingMsg))
	if incomingContact == nil {
		err := errors.New("request without being from a contact")
		logger.Debug(fmt.Sprintf("%v: %v", err.Error(), string(incomingMsg)))
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, err.Error())
		return
	}

	contact, err := h.ContactService.FindContact(incomingContact)
	if err != nil {
		logger.Debug(err.Error())
	}

	if contact != nil {
		channelId := contact.Channel.Hex()
		channel, err := h.ChannelService.FindChannelById(channelId)
		if err != nil {
			logger.Debug(err.Error())
		}
		if channel != nil {
			channelUUID := channel.UUID
			// RedirectRequest(channelUUID, string(jsonMsg))
			status, err := h.CourierService.RedirectMessage(channelUUID, string(incomingMsg))
			if err != nil {
				logger.Debug(err.Error())
				w.WriteHeader(status)
				fmt.Fprint(w, err)
			}
		}

	} else {
		if possibleToken := extractTextMessage(string(incomingMsg)); possibleToken != "" {
			ch, err := h.ChannelService.FindChannelByToken(possibleToken)
			if err != nil {
				logger.Error(err.Error())
			}
			if ch != nil {
				incomingContact.Channel = ch.ID
				h.ContactService.CreateContact(incomingContact)
				//TODO refactor this to use wpp service
				sendConfirmationMessage(incomingContact)
				w.WriteHeader(http.StatusCreated)
				fmt.Fprint(w, "")
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, errors.New("contact not found and token not valid"))
		}
	}
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	wconfig := config.GetConfig().Whatsapp
	httpClient := &http.Client{}
	reqPath := "/v1/users/login"

	reqURL, _ := url.Parse(wconfig.BaseURL + reqPath)

	req := &http.Request{
		Method: "POST",
		URL:    reqURL,
		Header: map[string][]string{},
		Body:   r.Body,
	}

	req.SetBasicAuth(config.AppConf.Whatsapp.Username, config.AppConf.Whatsapp.Password)

	res, err := httpClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		logger.Error(err.Error())
		return
	}

	var login LoginPayload

	if err := json.NewDecoder(res.Body).Decode(&login); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		logger.Error(err.Error())
		return
	}
	//TODO refactor this to reduce code and simplify
	newToken := login.Users[0].Token

	config.UpdateToken(newToken)
	logger.Info("Whatsapp token update")
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(login)
	for k, v := range res.Header {
		w.Header().Set(k, strings.Join(v, ""))
	}
	fmt.Fprint(w, string(b))
}

func sendConfirmationMessage(contact *models.Contact) {
	message := "Token válido, Whatsapp demo está pronto para sua utilização"
	urn := contact.URN
	payload := fmt.Sprintf(
		`{"to":"%s","type":"text","text":{"body":"%s"}}`,
		urn,
		message,
	)
	payloadBytes := []byte(payload)

	wconfig := config.GetConfig().Whatsapp

	httpClient := &http.Client{}
	reqPath := "/v1/messages"

	reqURL, _ := url.Parse(wconfig.BaseURL + reqPath)
	req := &http.Request{
		Method: "POST",
		URL:    reqURL,
		Header: map[string][]string{
			"Content-Type":  {"application/json; charset=UTF-8"},
			"Authorization": {"Bearer " + wconfig.AuthToken},
		},
		Body: ioutil.NopCloser(bytes.NewReader(payloadBytes)),
	}

	res, err := httpClient.Do(req)
	if err != nil {
		logger.Error(err.Error())
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		logger.Info(string(body))
	}
}

type LoginPayload struct {
	Users []struct {
		Token        string
		ExpiresAfter string
	}
	Meta struct {
		Version   string
		ApiStatus string
	}
}

func incomingMsgToContact(m string) *models.Contact {
	name := extractName(m)
	number := extractNumber(m)
	if name != "" && number != "" {
		return &models.Contact{
			URN:  number,
			Name: name,
		}
	}
	return nil
}

func extractName(m string) string {
	var result map[string][]map[string]map[string]interface{}
	json.Unmarshal([]byte(m), &result)
	if result["contacts"] != nil {
		return result["contacts"][0]["profile"]["name"].(string)
	}
	return ""
}

func extractNumber(m string) string {
	var result map[string][]map[string]interface{}
	json.Unmarshal([]byte(m), &result)
	if result["messages"] != nil {
		return result["messages"][0]["from"].(string)
	}
	return ""
}

func extractTextMessage(m string) string {
	var result map[string][]map[string]map[string]interface{}
	json.Unmarshal([]byte(m), &result)
	if result["messages"] != nil {
		return result["messages"][0]["text"]["body"].(string)
	}
	return ""
}
