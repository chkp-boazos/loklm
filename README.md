# LokLM 

![alt text](https://private-user-images.githubusercontent.com/211166084/444964666-9e2b09cc-3a03-4e47-b20d-18fb1c2663cf.png?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnaXRodWIuY29tIiwiYXVkIjoicmF3LmdpdGh1YnVzZXJjb250ZW50LmNvbSIsImtleSI6ImtleTUiLCJleHAiOjE3NDc2MjcyNjYsIm5iZiI6MTc0NzYyNjk2NiwicGF0aCI6Ii8yMTExNjYwODQvNDQ0OTY0NjY2LTllMmIwOWNjLTNhMDMtNGU0Ny1iMjBkLTE4ZmIxYzI2NjNjZi5wbmc_WC1BbXotQWxnb3JpdGhtPUFXUzQtSE1BQy1TSEEyNTYmWC1BbXotQ3JlZGVudGlhbD1BS0lBVkNPRFlMU0E1M1BRSzRaQSUyRjIwMjUwNTE5JTJGdXMtZWFzdC0xJTJGczMlMkZhd3M0X3JlcXVlc3QmWC1BbXotRGF0ZT0yMDI1MDUxOVQwMzU2MDZaJlgtQW16LUV4cGlyZXM9MzAwJlgtQW16LVNpZ25hdHVyZT02ZDk1N2UzYmIzOTkxMWIwM2JjZWZmZmU1NmJlNDljYzExM2FkYmNjYjY0N2M0ZTZjNDE2MDRiOWYyZmUxYzI5JlgtQW16LVNpZ25lZEhlYWRlcnM9aG9zdCJ9.EKm6UaqaZf4eYz6Bg6XsWnWIybsAtrFDObhKueeLbl0 "LoklM LOGO")

---

### LokLM CLI

LokLM (compact, brandable: Local + LLM) CLI is a developer-focused command-line interface designed to streamline the setup, execution, and management of a fully local AI development environment.

---

### Why?

#### PROs:

1. Choose the model you want to use.
2. No reliance on the internet.
3. A free LLM server to run and test your code.

#### CONs:

1. Relies on your computer's resources, especially if they're limited.
2. ChatGPT model is not available.

---

### Let's Get Started!

If you like the idea, what should you do next?
Install Docker and download the binary file.

---

### Setting Up the Environment

Run the following command:

```shell
loklm setup
```

This will pull the Ollama and Jupyter container images, set up the LokLM network on Docker, and establish a directory to store your notebooks and LLM models.

---

### Running the Environment

The next step is to run the containers: Ollama for running LLM models and Jupyter for writing Python code and testing agents.

```shell
loklm start
```

---

### Pulling a Model

Now that we have our environment set up and running, it's time to pull a model and start writing our agent. I usually use `ollama4`, but if you have less powerful hardware, you can use `smollm`, which should run on most computers.

```shell
loklm pull smollm
```

---

### Time to Code

Visit [localhost:8888](http://localhost:8888). We’ll write our first snippet to interact with the LLM.
If it's your first time running Jupyter, you’ll be asked to provide a token.
You can find the token by running:

```shell
loklm jupyterToken
```

Then copy and paste the token.

#### Steps:

1. Create a new Python3 notebook.
2. Install LangChain and LangChain Community:

   ```python
   %pip install langchain langchain_community
   ```
3. Import Ollama:

   ```python
   from langchain.llms import Ollama
   ```
4. Initialize the Ollama LLM:

   ```python
   llm = Ollama(model="smollm", base_url="http://llm:11434")
   ```
5. Run a simple query and print the result:

   ```python
   print(llm("What is LangChain?"))
   ```

