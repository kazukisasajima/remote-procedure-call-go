package protocol

type Request struct {
	Method      string        `json:"method"`
	Params      []interface{} `json:"params"`
	ParamTypes  []string      `json:"param_types"`
	ID          int           `json:"id"`
}

type Response struct {
	Results     interface{} `json:"results"`
	ResultType  string      `json:"result_type"`
	ID          int         `json:"id"`
	Error       string      `json:"error,omitempty"`
}
