package config

// Config holds the....config
type Config struct {
	OutputDir       string
	InputFile       string
	KubeConfig      string
	GetK8sPods      bool
	GetK8sSvcs      bool
	ChromeExe       string
	NumberOfWorkers int
	TimeoutSeconds  int
}
