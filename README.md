# Strut
Strut is a TUI for Gemini or Chat GPT

# Install
```shell
brew tap postsa/tap
brew install strut
```

## Usage
export your gemini api key, e.g.

```shell
export GEMINI_API_KEY=$(skate get gemini_api_key)
```
```shell
export OPEN_AI_API_KEY=$(skate get open_ai_api_key)
```

## Demo
Ask questions and press tab to browse history.
![Alt Text](demos/demo.gif)

Use `Ctrl+a` to copy all text out of the response pane.

## Vim
Press `Ctrl+e` to open a response in vim
![Alt Text](demos/vim.gif)

## Choosing a New Model
Press `Escape` to choose a different model and `Escape` again to exit.
![Alt Text](demos/choose.gif)