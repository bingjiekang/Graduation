package config

type Server struct {
	Zap Zap `mapstructure:"zap" json:"zap" yaml:"zap"`
	// gorm 连接mysql
	Mysql  Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	System System `mapstructure:"system" json:"system" yaml:"system"`
	// oss
	// Local Local `mapstructure:"local" json:"local" yaml:"local"`
}
