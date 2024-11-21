package yamlconfig

type RabbitMQ_t struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Vhost    string `yaml:"vhost"`
}

type Mysql_t struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Redis_t struct {
	Host     string `yaml:"host"`
	Database int    `yaml:"db"`
}

type myConfig_t struct {
	Name     string                `yaml:"Name"`
	City     string                `yaml:"City"`
	State    string                `yaml:"State"`
	Id       int                   `yaml:"Id"`
	Options  []string              `yaml:"Options"`
	Rabbitmq map[string]RabbitMQ_t `yaml:"Rabbitmq"`
	Mysql    Mysql_t               `yaml:"Mysql"`
	Redis    []Redis_t             `yaml:"Redis"`
}
