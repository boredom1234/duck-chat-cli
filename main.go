package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Constants for API endpoints and headers
const (
	statusURL         = "https://duckduckgo.com/duckchat/v1/status"
	chatURL          = "https://duckduckgo.com/duckchat/v1/chat"
	statusHeaders     = "1"
	termsOfServiceURL = "https://duckduckgo.com/aichat/privacy-terms" // Not used in the API version
)

// ANSI color codes
const (
	colorRed    = "\033[31m"
	colorBlue   = "\033[34m"
	colorGreen  = "\033[32m"
	colorReset  = "\033[0m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

// ASCII Art constants
const (
	logo = `
     _____       _           _   
    |  __ \     | |         | |  
    | |  | |_   | |__   __ _| |_ 
    | |  | | |  | '_ \ / _' | __|
    | |__| | |__| | | | (_| | |_ 
    |_____/ \___/|_| |_|\__,_|\__|
                                  
`
	separator = "═══════════════════════════════════════════════════════"
)

// Model represents the AI model used for chat
type Model string

// ModelAlias represents a user-friendly alias for the AI model
type ModelAlias string

// Define available models and their aliases
const (
	GPT4Mini Model = "gpt-4o-mini"
	Claude3  Model = "claude-3-haiku-20240307"
	Llama    Model = "meta-llama/Meta-Llama-3.1-70B-Instruct-Turbo"
	Mixtral  Model = "mistralai/Mixtral-8x7B-Instruct-v0.1"

	GPT4MiniAlias ModelAlias = "gpt-4o-mini"
	Claude3Alias  ModelAlias = "claude-3-haiku"
	LlamaAlias    ModelAlias = "llama"
	MixtralAlias  ModelAlias = "mixtral"
)

// Map model aliases to their corresponding Model values
var modelMap = map[ModelAlias]Model{
	GPT4MiniAlias: GPT4Mini,
	Claude3Alias:  Claude3,
	LlamaAlias:    Llama,
	MixtralAlias:  Mixtral,
}

// Message represents a chat message
type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

// ChatPayload represents the payload sent to the chat API
type ChatPayload struct {
	Model    Model     `json:"model"`
	Messages []Message `json:"messages"`
}

// Chat represents a chat session
type Chat struct {
	OldVqd   string
	NewVqd   string
	Model    Model
	Messages []Message
	Client   *http.Client
}

// NewChat creates a new Chat instance
func NewChat(vqd string, model Model) *Chat {
	return &Chat{
		OldVqd:   vqd,
		NewVqd:   vqd,
		Model:    model,
		Messages: []Message{},
		Client:   &http.Client{},
	}
}

// Fetch sends a chat message and returns the response
func (c *Chat) Fetch(content string) (*http.Response, error) {
	c.Messages = append(c.Messages, Message{Content: content, Role: "user"})
	payload := ChatPayload{
		Model:    c.Model,
		Messages: c.Messages,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %v", err)
	}

	req, err := http.NewRequest("POST", chatURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("x-vqd-4", c.NewVqd)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("%d: Failed to send message. %s. Body: %s", resp.StatusCode, resp.Status, string(body))
	}

	return resp, nil
}

// FetchStream sends a chat message and returns a channel for streaming the response
func (c *Chat) FetchStream(content string) (<-chan string, error) {
	resp, err := c.Fetch(content)
	if err != nil {
		return nil, err
	}

	stream := make(chan string)
	go func() {
		defer resp.Body.Close()
		defer close(stream)

		var text strings.Builder
		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			line := scanner.Text()

			if line == "data: [DONE]" {
				break
			}

			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				var messageData struct {
					Message string `json:"message"`
				}
				if err := json.Unmarshal([]byte(data), &messageData); err != nil {
					log.Printf("Error unmarshaling data: %v\n", err)
					continue
				}

				if messageData.Message != "" {
					text.WriteString(messageData.Message)
					stream <- messageData.Message
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Error reading response body: %v\n", err)
		}

		c.OldVqd = c.NewVqd
		c.NewVqd = resp.Header.Get("x-vqd-4")
		c.Messages = append(c.Messages, Message{Content: text.String(), Role: "assistant"})
	}()

	return stream, nil
}

// Redo resets the chat to the previous state
func (c *Chat) Redo() {
	c.NewVqd = c.OldVqd
	if len(c.Messages) >= 2 {
		c.Messages = c.Messages[:len(c.Messages)-2]
	}
}

// InitChat initializes a new chat session
func InitChat(model ModelAlias) (*Chat, error) {
	req, err := http.NewRequest("GET", statusURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("x-vqd-accept", statusHeaders)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d: Failed to initialize chat. %s", resp.StatusCode, resp.Status)
	}

	vqd := resp.Header.Get("x-vqd-4")
	if vqd == "" {
		return nil, fmt.Errorf("failed to get VQD from response headers")
	}

	return NewChat(vqd, modelMap[model]), nil
}

// clearScreen clears the terminal screen based on the operating system
func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// displayModelOptions displays the available models and prompts the user to select one
func displayModelOptions() ModelAlias {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(colorCyan + separator + colorReset)
	fmt.Println(colorYellow + "Available Models:" + colorReset)
	fmt.Println(colorCyan + "╔══════════════════════════════════════╗")
	fmt.Println("║ 1. 🤖 GPT-4 Mini (gpt-4o-mini)        ║")
	fmt.Println("║ 2. 🎯 Claude 3 Haiku (claude-3-haiku) ║")
	fmt.Println("║ 3. 🦙 Llama (llama)                   ║")
	fmt.Println("║ 4. 🔮 Mixtral (mixtral)               ║")
	fmt.Println("╚══════════════════════════════════════╝" + colorReset)

	for {
		fmt.Print(colorYellow + "Select a model (1-4): " + colorReset)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			return GPT4MiniAlias
		case "2":
			return Claude3Alias
		case "3":
			return LlamaAlias
		case "4":
			return MixtralAlias
		default:
			fmt.Println(colorRed + "Invalid selection. Please choose a number between 1 and 4." + colorReset)
		}
	}
}

func main() {
	var chat *Chat
	var err error

	clearScreen()
	fmt.Print(colorCyan + logo + colorReset)

	initializeChat := func() {
		modelAlias := displayModelOptions()
		chat, err = InitChat(modelAlias)
		if err != nil {
			fmt.Printf(colorRed+"Error initializing chat: %v\n"+colorReset, err)
			return
		}
		fmt.Println(colorCyan + separator + colorReset)
		fmt.Println(colorYellow + "🚀 Chat session started. Commands:" + colorReset)
		fmt.Println(colorCyan + "┌─────────────────────────────┐")
		fmt.Println("│ 🚪 Type 'exit' to quit       │")
		fmt.Println("│ 🤖 Type '/model' to change   │")
		fmt.Println("│ 🔄 Type '/reset' to restart  │")
		fmt.Println("│ 🧹 Type '/clear' to clean    │")
		fmt.Println("└─────────────────────────────┘" + colorReset)
		fmt.Println(colorCyan + separator + colorReset)
	}

	initializeChat()
	if err != nil {
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(colorBlue + "You ➤ " + colorReset)
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)

		switch userInput {
		case "exit":
			fmt.Println(colorYellow + "👋 Goodbye!" + colorReset)
			return
		case "/model":
			modelAlias := displayModelOptions()
			chat.Model = modelMap[modelAlias]
			fmt.Println(colorGreen + "✨ Model changed successfully!" + colorReset)
			continue
		case "/reset":
			clearScreen()
			fmt.Print(colorCyan + logo + colorReset)
			initializeChat()
			if err != nil {
				return
			}
			continue
		case "/clear":
			clearScreen()
			fmt.Print(colorCyan + logo + colorReset)
			continue
		}

		stream, err := chat.FetchStream(userInput)
		if err != nil {
			fmt.Printf(colorRed+"Error fetching response: %v\n"+colorReset, err)
			continue
		}

		fmt.Print(colorGreen + "AI 🤖 " + colorReset)
		for chunk := range stream {
			fmt.Print(colorGreen + chunk + colorReset)
		}
		fmt.Println("\n" + colorCyan + separator + colorReset)
	}
}