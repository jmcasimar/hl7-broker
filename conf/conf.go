package conf

type AppConfig struct {
	Host             string `required:"true"`
	Port             int    `required:"true"`
	Secret           string `required:"true"`
	UdpDiscoveryHost string `required:"true"`
}

type DBConfig struct {
	Username string `required:"true"`
	Password string `required:"true"`
	Name     string `required:"true"`
	Host     string `required:"true"`
	Port     int    `required:"true"`
}
