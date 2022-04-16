package wordfilter

import "log"

// Default Wordfilter
var Default *Filter

// Init 初始化 Default Wordfilter
func Init(path string) {
	if Default == nil {
		Default = New()
		if err := Default.LoadWordDict(path); err != nil {
			log.Fatalf("load word filter err %v", err)
		}
	}
}
