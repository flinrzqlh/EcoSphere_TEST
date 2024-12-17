package handlers

const (
	HUGGING_URL = "https://api-inference.huggingface.co/models/"
	TOKEN       = "hf_xxxxxxxx"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Payload struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	MaxToken int       `json:"max_tokens"`
	Stream   bool      `json:"stream"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

type Response struct {
	Choices []Choice `json:"choices"`
}

type Table struct {
	Date              []string `json:"date"`
	Time              []string `json:"time"`
	Appliance         []string `json:"appliance"`
	EnergyConsumption []string `json:"energy_consumption"`
	Room              []string `json:"room"`
	Status            []string `json:"status"`
}

type EnergyAnalysisPayload struct {
	Query string `json:"query"`
	Table Table  `json:"table"`
}

type EnergyAnalysisResponse struct {
	LeastConsumingAppliance string  `json:"least_consuming"`
	MostConsumingAppliance  string  `json:"most_consuming"`
	LeastConsumptionValue   float64 `json:"least_consumption_value"`
	MostConsumptionValue    float64 `json:"most_consumption_value"`
}
