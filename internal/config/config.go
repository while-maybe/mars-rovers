package config

type Config struct {
	FilePath    string
	MinPlateauX int
	MinPlateauY int
}

func Default() Config {
	return Config{
		FilePath:    "",
		MinPlateauX: 2,
		MinPlateauY: 2,
	}
}
