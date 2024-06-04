package interfaces

type EvaluationConsumer struct {
	UserId          string `json:"userId"`
	EvaluationPoint int32  `json:"evaluationPoint"`
	Date            string `json:"date"` // Format 2006-01-02
}
