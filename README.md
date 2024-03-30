# Finance Tracker Bot

The Finance Tracker Bot is a sophisticated tool designed in Go, aimed at providing individuals with a comprehensive
solution for managing their finances. It automates the tracking of expenses, budgets, and generates insightful financial
reports, all stored securely in a local database.

## Features

- **Expense Tracking**: Effortlessly log every expense, categorize them, and keep track of your spending habits.
- **Budget Management**: Set up customizable budgets for different categories and get real-time updates on your budget
  status.
- **Financial Reporting**: Access detailed reports to analyze your spending patterns, savings, and overall financial
  health over time. (coming soon)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing
purposes.

### Prerequisites

- Go 1.15 or later. You can download it from [here](https://golang.org/dl/).
- SQLite3 for the database. Installation instructions can be found [here](https://www.sqlite.org/download.html).

### Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/qfpeeeer/finance-tracker-bot.git
    cd finance-tracker-bot
    ```

2. Install Go dependencies:
    ```bash
    go mod tidy
    ```

3. Initialize the database. Run the following command to create `data.db` and set up the necessary tables or you can run
   the bot, and it will create the database for you.
    ```bash
    sqlite3 data.db < schema.sql
    ```

### Configuration

- **Environment Variables**: Set up your environment variables (if any) in a `.env` file or your preferred configuration
  method.
    - `DATA_FILE_PATH`: Path to your SQLite database file (e.g., `./data.db`).
    - `TELEGRAM_TOKEN`: Telegram Bot API token. You can get one by creating a new bot on Telegram using the
      [BotFather](https://core.telegram.org/bots#6-botfather).

### Running Locally

To start the bot, run:

```bash
go run app/main.go
```

Replace app/main.go with the correct path to your application's entry point.

## Usage

    Starting the Bot: Detailed instructions on how to interact with the bot after it's running. Include any commands or interfaces provided by the bot.
    Adding Expenses: Steps to log an expense using the bot.
    Adding Expenses category: Steps to add a category to the expenses.
    Viewing Reports: How to generate and view financial reports.

## Deployment

Check the deployments folder for scripts and configurations needed to deploy the Finance Tracker Bot. Include specific
instructions for deploying to popular platforms if available.