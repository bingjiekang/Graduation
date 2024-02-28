package config

// 缓存数据库
type Redis struct {
	RedisHost        string `mapstructure:"redis-host" json:"redisHost" yaml:"redis-host"` // redis地址
	RedisPort        string `mapstructure:"redis-port" json:"redisPort" yaml:"redis-port"`
	RedisPassword    string `mapstructure:"redis-password" json:"redisPassword" yaml:"redis-password"`
	RedisDb          int    `mapstructure:"redis-db" json:"redisDb" yaml:"redis-db"`
	RedisPoolSize    int    `mapstructure:"redis-pool-size" json:"redisPoolSize" yaml:"redis-pool-size"`
	RedisMaxRetries  int    `mapstructure:"redis-max-retries" json:"redisMaxRetries" yaml:"redis-max-retries"`
	RedisIdleTimeout int    `mapstructure:"redis-idle-timeout" json:"redisIdleTimeout" yaml:"redis-idle-timeout"`
}
