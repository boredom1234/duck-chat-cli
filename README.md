# 🚀 DuckDuckGo AI Chat CLI

Welcome to the **DuckDuckGo AI Chat CLI**! This is a command-line interface (CLI) tool that allows you to interact with various AI models (like GPT-4, Claude 3, Llama, and Mixtral) using DuckDuckGo's chat API. The tool is designed to be simple, interactive, and fun to use, with colorful output and emojis to enhance your experience. 🌈

---

## 📜 Table of Contents

- [🚀 DuckDuckGo AI Chat CLI](#-duckduckgo-ai-chat-cli)
  - [📜 Table of Contents](#-table-of-contents)
  - [✨ Features](#-features)
  - [🛠️ Installation](#️-installation)
  - [🚀 Usage](#-usage)
    - [Starting the Chat](#starting-the-chat)
    - [Available Commands](#available-commands)
  - [🤖 Supported Models](#-supported-models)
  - [📝 Example Usage](#-example-usage)
  - [📄 License](#-license)
  - [🙏 Contributing](#-contributing)

---

## ✨ Features

- **Interactive Chat Interface**: Chat with AI models directly from your terminal. 💬
- **Multiple AI Models**: Choose from GPT-4, Claude 3, Llama, and Mixtral. 🤖
- **Streaming Responses**: Get real-time responses from the AI. ⚡
- **Colorful Output**: Enjoy a visually appealing chat experience with ANSI colors. 🌈
- **Easy Commands**: Use simple commands to reset, clear, or change models. 🎯
- **Cross-Platform**: Works on Windows, macOS, and Linux. 🖥️

---

## 🛠️ Installation

To use the DuckDuckGo AI Chat CLI, you need to have **Go** installed on your system. If you don't have Go installed, you can download it from [here](https://golang.org/dl/).

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/your-username/duckduckgo-ai-chat-cli.git
   cd duckduckgo-ai-chat-cli
   ```

2. **Build the Project**:
   ```bash
   go build -o duckchat
   ```

For Windows:

```bash
go build -o duckchat.exe
```

3. **Run the CLI**:
   ```bash
   ./duckchat
   ```

---

## 🚀 Usage

### Starting the Chat

1. **Run the CLI**:

   ```bash
   ./duckchat
   ```

2. **Select a Model**:
   You will be prompted to select an AI model from the available options:

   ```
   Available Models:
   ╔══════════════════════════════════════╗
   ║ 1. 🤖 GPT-4 Mini (gpt-4o-mini)        ║
   ║ 2. 🎯 Claude 3 Haiku (claude-3-haiku) ║
   ║ 3. 🦙 Llama (llama)                   ║
   ║ 4. 🔮 Mixtral (mixtral)               ║
   ╚══════════════════════════════════════╝
   Select a model (1-4):
   ```

3. **Start Chatting**:
   Once the model is selected, you can start chatting with the AI. Type your message and press `Enter`.

### Available Commands

- **`exit`**: Exit the chat session. 🚪
- **`/model`**: Change the AI model. 🤖
- **`/reset`**: Reset the chat session. 🔄
- **`/clear`**: Clear the terminal screen. 🧹

---

## 🤖 Supported Models

The following AI models are supported:

| Model Name         | Alias            | Description                              |
| ------------------ | ---------------- | ---------------------------------------- |
| **GPT-4 Mini**     | `gpt-4o-mini`    | A compact version of GPT-4.              |
| **Claude 3 Haiku** | `claude-3-haiku` | A fast and efficient model by Anthropic. |
| **Llama**          | `llama`          | Meta's Llama model.                      |
| **Mixtral**        | `mixtral`        | A high-performance model by Mistral.     |

---

## 📝 Example Usage

```bash
$ ./duckchat

     _____       _           _
    |  __ \     | |         | |
    | |  | |_   | |__   __ _| |_
    | |  | | |  | '_ \ / _' | __|
    | |__| | |__| | | | (_| | |_
    |_____/ \___/|_| |_|\__,_|\__|


═══════════════════════════════════════════════════════
Available Models:
╔══════════════════════════════════════╗
║ 1. 🤖 GPT-4 Mini (gpt-4o-mini)        ║
║ 2. 🎯 Claude 3 Haiku (claude-3-haiku) ║
║ 3. 🦙 Llama (llama)                   ║
║ 4. 🔮 Mixtral (mixtral)               ║
╚══════════════════════════════════════╝
Select a model (1-4): 1

═══════════════════════════════════════════════════════
🚀 Chat session started. Commands:
┌─────────────────────────────┐
│ 🚪 Type 'exit' to quit       │
│ 🤖 Type '/model' to change   │
│ 🔄 Type '/reset' to restart  │
│ 🧹 Type '/clear' to clean    │
└─────────────────────────────┘
═══════════════════════════════════════════════════════

You ➤ Hello, how are you?
AI 🤖 I'm just a bunch of code, so I don't have feelings, but I'm here to help! How can I assist you today? 😊
```

---

## 📄 License

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for more details.

---

## 🙏 Contributing

Contributions are welcome! If you have any suggestions, bug reports, or feature requests, please open an issue or submit a pull request. Let's make this tool even better together! 🚀

---

Happy chatting! 🎉
