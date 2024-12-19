package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/coolorvi/web-calculator/calc"
)

type CalcRequest struct {
	Expression string `json:"expression"`
}

type CalcResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
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
	http.HandleFunc("/api/v1/calculate", calculateHandler)
}

func main() {
	SetupRoutes()
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
