package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/coolorvi/web-calculator/calc"
)

type CalcRequest struct {
	Expression string `json:"expression"`
}

type CalcResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

func (a *Application) Run() error {
	for {
		log.Println("input expression")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("failed to read expression from console")
		}

		text = strings.TrimSpace(text)
		if text == "exit" {
			log.Println("aplication was successfully closed")
			return nil
		}

		result, err := calc.Tokenize(text)
		if err != nil {
			log.Println(text, " calculation failed wit error: ", err)
		} else {
			log.Println(text, "=", result)
		}
	}
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	var req CalcRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CalcResponse{Error: "Invalid JSON"})
		return
	}

	re := regexp.MustCompile(`^[0-9+\-*/().\s]+$`)
	if !re.MatchString(req.Expression) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(CalcResponse{Error: "Expression is not valid"})
		return
	}

	tokens, err := calc.Tokenize(req.Expression)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(CalcResponse{Error: "Expression is not valid"})
		return
	}

	result, err := calc.ParseTokens(tokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(CalcResponse{Error: "Internal server error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CalcResponse{Result: fmt.Sprintf("%f", result)})
}

func SetupRoutes() {
	http.HandleFunc("/api/v1/calculate", CalculateHandler)
}

func Start() {
	SetupRoutes()
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func (a *Application) RunServer() error {
	http.HandleFunc("/", CalculateHandler)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
