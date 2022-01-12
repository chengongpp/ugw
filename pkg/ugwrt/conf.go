package ugwrt

type Config struct {
	Name           string `toml:"Name"`
	Host           string `toml:"Host"`
	Port           int    `toml:"Port"`
	LogLevel       string `toml:"LogLevel"`
	LogDir         string `toml:"LogDir"`
	MaxConnections int    `toml:"MaxConnections"`
	OutBounds      []OutBound
}
