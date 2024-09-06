
# Chat With Gemini (Telegram Bot)

**Chat With Gemini** is a telegram bot as an AI Chatbot like ChatGPT, but it is completely free because it uses google's AI model "Gemini".


## Prerequisite

- Telegram bot token
- Gemini API Key
- Docker installed on your computer
- MySQL database (*optional, if you don't use the mysql container*)

## Installation

- create a new database and create the *chat_session* table by running the SQL syntax from [chat_session.sql](/sql/chat_session.sql) (*optional, if you don't use the mysql container*)
- Create new file called *.env* copy from *.env.example*
- Set up *.env* file with your credentials

    ```
    # TELEGRAM BOT
    BOT_TOKEN="<YourBotToken>"

    # GEMINI
    GEMINI_API_KEY="<YourGeminiAPIToken>"

    # DATABASE
    DB_HOST="<YourDBHost>"
    DB_PORT="<YourDBPort>"
    DB_USER="<YourDBUser>"
    DB_PASS="<YourDBPass>"
    DB_NAME="<YourDBName>"
    ```
- Finally, run

    ```
    docker compose up --build
    ```


## Screenshots

![App Screenshot](/assets/chat-with-gemini-telegram-screenshot.png)

