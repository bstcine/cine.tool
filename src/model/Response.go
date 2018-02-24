package model

type ResList struct {
	Code             string     `json:"code"`
	Code_desc        string     `json:"code_desc"`
	Except_case      string     `json:"except_case"`
	Except_case_desc string     `json:"except_case_desc"`
	Result           ResultList `json:"result"`
}

type Res struct {
	Code             string                 `json:"code"`
	Code_desc        string                 `json:"code_desc"`
	Except_case      string                 `json:"except_case"`
	Except_case_desc string                 `json:"except_case_desc"`
	Result           map[string]interface{} `json:"result"`
}

type ResultList struct {
	Cur_page   string      `json:"cur_page"`
	Max_page   string      `json:"max_page"`
	Total_rows string      `json:"total_rows"`
	Rows       interface{} `json:"rows"`
}
