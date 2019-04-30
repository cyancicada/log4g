package log4g

type Config struct {
	NameSpace           string `json:",optional"`
	Stdout              bool   `json:"stdout,default=true"`
	LogMode             string `json:",options=regular|volume,default=regular"`
	Path                string `json:",default=logs"`
	Compress            bool   `json:",optional"`
	KeepDays            int    `json:",optional"`
	StackCoolDownMillis int    `json:",default=100"`
}
