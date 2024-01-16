# gptprep

`gptprep` is a command-line utility for quickly loading files into your clipboard, making it easy to provide context to conversational AI platforms such as ChatGPT. It automates file selection and preparation, significantly reducing the manual effort previously required to set up context for AI interactions.

The motivation behind creating `gptprep` was that I spent a significant amount of time selecting and preparing the context for ChatGPT instead of describing my feature. 

## One-Line Installation

To install or update `gptprep`, run the following command in your terminal. Note that Windows is not currently supported, and users are encouraged to review the install script for transparency on what the installation does.

```sh
curl -sSL https://raw.githubusercontent.com/2start/gptprep/main/install.sh | sudo sh
```

## Usage Examples

Print the manual.

```sh
gptprep -h
```

To prepare the context to generate documentation for this repository, I used the following command:

```sh
gptprep --exclude ".github" --exclude "go.mod" --exclude "go.sum"
```

To help me with developing this tool I used the following command to only load the code files:

```sh
gptprep --extension .go
```

## Configuration

The following table lists the command line parameters supported by `gptprep`:

| Parameter     | Description                                                  |
| ------------- | ------------------------------------------------------------ |
| `--extension` | Specify file extensions to include in the search. Multiple extensions can be specified by repeating the parameter. |
| `--exclude`   | Define patterns or filenames to exclude from the search. Multiple excludes can be specified by repeating the parameter. |

`gptprep` automatically ignores:
- `.git`
- `.gitignore` 
- globs in your `.gitignore`.
- files that are not identified as text file via their mimetype `text/*`.


## Feature Requests

Don't hesitate to ask for a feature. If we agree on something useful I'll implement it quickly or merge your PR.
Just open an ISSUE or submit a pull request.

## Star History

Leave a star if you like the tool. That helps me stay motivated 🤩

[![Star History Chart](https://api.star-history.com/svg?repos=2start/gptprep&type=Date)](https://star-history.com/#2start/gptprep&Date)
