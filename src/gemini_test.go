package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsModelSupported(t *testing.T) {
	tc := []struct {
		name      string
		modelName string
		expected  bool
	}{
		{
			name:      "Test with gemini-pro, should return true",
			modelName: "gemini-pro",
			expected:  true,
		},
		{
			name:      "Test with dummy-model, should return false",
			modelName: "dummy-model",
			expected:  false,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isModelSupported(tt.modelName))
		})
	}
}

func TestSendMessage(t *testing.T) {
	g, err := NewGemini(os.Getenv("GEMINI_API_KEY"), geminiProModel)
	if err != nil {
		log.Println("failed to init gemini")
		return
	}

	g.StartChat([]Content{})
	resp, err := g.SendMessage("This is test")
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}
