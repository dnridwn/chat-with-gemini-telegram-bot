package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Part struct {
	Text string `json:"text,omitempty"`
}

type Content struct {
	Parts []Part `json:"parts,omitempty"`
	Role  string `json:"role,omitempty"`
}

type SafetyRating struct {
	Category    string `json:"category,omitempty"`
	Probability string `json:"probability,omitempty"`
}

type SafetySetting struct {
	Category  string `json:"category,omitempty"`
	Threshold string `json:"threshold,omitempty"`
}

type Candidate struct {
	Content       Content        `json:"content,omitempty"`
	FinishReason  string         `json:"finishReason,omitempty"`
	Index         int            `json:"index,omitempty"`
	SafetyRatings []SafetyRating `json:"safetyRatings,omitempty"`
}

type PromptFeedback struct {
	SafetyRatings []SafetyRating `json:"safetyRatings,omitempty"`
}

type GenerateRequest struct {
	Contents       []Content       `json:"contents,omitempty"`
	SafetySettings []SafetySetting `json:"safetySettings,omitempty"`
}

type GenerateResponse struct {
	Candidates     []Candidate    `json:"candidates,omitempty"`
	PromptFeedback PromptFeedback `json:"promptFeedback,omitempty"`
}

func (gr *GenerateResponse) String() string {
	result := ""
	for _, cand := range gr.Candidates {
		for _, part := range cand.Content.Parts {
			result += part.Text
		}
	}
	return result
}

const (
	baseURL        = "https://generativelanguage.googleapis.com/v1beta"
	geminiProModel = "gemini-pro"
)

var (
	supportedModels = []string{
		geminiProModel,
	}
	explicitMode = []SafetySetting{
		{
			Category:  "HARM_CATEGORY_HARASSMENT",
			Threshold: "BLOCK_NONE",
		},
		{
			Category:  "HARM_CATEGORY_HATE_SPEECH",
			Threshold: "BLOCK_NONE",
		},
		{
			Category:  "HARM_CATEGORY_SEXUALLY_EXPLICIT",
			Threshold: "BLOCK_NONE",
		},
		{
			Category:  "HARM_CATEGORY_DANGEROUS_CONTENT",
			Threshold: "BLOCK_NONE",
		},
	}
)

type Gemini struct {
	key       string
	modelName string
	history   []Content
}

func NewGemini(key string, model string) (*Gemini, error) {
	if !isModelSupported(model) {
		return nil, fmt.Errorf("model %s is not supported", model)
	}

	return &Gemini{
		key:       key,
		modelName: model,
	}, nil
}

func (g *Gemini) StartChat(h []Content) {
	g.history = h
}

func (g *Gemini) SendMessage(m string) (GenerateResponse, error) {
	g.addUserInputToHistory(m)

	body := GenerateRequest{
		Contents:       g.history,
		SafetySettings: explicitMode,
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return GenerateResponse{}, err
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", baseURL, g.modelName, g.key)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyJson))
	if err != nil {
		return GenerateResponse{}, err
	}

	defer resp.Body.Close()

	parsedResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return GenerateResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(string(parsedResp))
		return GenerateResponse{}, fmt.Errorf("error status code %d", resp.StatusCode)
	}

	gr := GenerateResponse{}
	if err := json.Unmarshal(parsedResp, &gr); err != nil {
		return GenerateResponse{}, err
	}

	for _, c := range gr.Candidates {
		g.AddToHistory(c.Content)
	}

	return gr, nil
}

func (g *Gemini) AddToHistory(c Content) {
	g.history = append(g.history, c)
}

func (g *Gemini) GetHistory() []Content {
	return g.history
}

func (g *Gemini) addUserInputToHistory(m string) {
	g.history = append(g.history, Content{
		Role: "user",
		Parts: []Part{
			{
				Text: m,
			},
		},
	})
}

func isModelSupported(model string) bool {
	for _, m := range supportedModels {
		if m == model {
			return true
		}
	}
	return false
}
