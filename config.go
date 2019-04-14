package log4g

type Config struct {
	NameSpace           string `json:",optional"`
	LogMode             string `json:",options=regular|console|volume,default=regular"`
	Path                string `json:",default=logs"`
	Compress            bool   `json:",optional"`
	KeepDays            int    `json:",optional"`
	StackCoolDownMillis int    `json:",default=100"`
}
