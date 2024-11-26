package yamlconfig

import "time"

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
	Host string `yaml:"host"`
	Db   int    `yaml:"db"`
}

type Path_t struct {
	Path string `yaml:"path"`
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
	MyArray  []Path_t              `yaml:"MyArray"`
}

// ============== Workflow Types =================
type Workflow_t struct {
	// Id            string                        `yaml:"id"`
	// SchemaVersion string                        `yaml:"schema_version"`
	Name          string                        `yaml:"name"`
	Description   string                        `yaml:"description"`
	Version       string                        `yaml:"version"`
	InitiatePaths []WorkflowInitiate_t          `yaml:"initiate"`
	Behaviors     map[string]WorkflowBehavior_t `yaml:"paths"`
	Services      map[string]TwistyService_t    `yaml:"twisty-services"`
}

type WorkflowInitiate_t struct {
	Path string `yaml:"path"`
}

type WorkflowPath_t struct {
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	PathType    string           `yaml:"type"`
	Steps       []WorkflowStep_t `yaml:"steps"`
}

type WorkflowBehavior_t struct {
	Name        string                     `yaml:"name"`
	Description string                     `yaml:"description"`
	PathType    string                     `yaml:"type"`
	FirstStep   string                     `yaml:"firststep"`
	Steps       map[string]WorkflowStep_t  `yaml:"steps"`
	Services    map[string]TwistyService_t `yaml:"services"`
}

type WorkflowStep_t struct {
	Name string `yaml:"name"`
	// Id            int32                 `yaml:"id"`
	Description   string                `yaml:"description"`
	StepTypeStr   string                `yaml:"type"` // This is a reserved field for future use
	StepType      StepType_t            `yaml:"-"`
	Command       string                `yaml:"command"` // Only Valid in workflow-service --- Does this have any use?
	Parameters    map[string]string     `yaml:"parameters"`
	WaitFor       []string              `yaml:"waitfor"` // Only Valid in workflow-merge
	TimeoutStr    string                `yaml:"timeout"` // eg. 200ms, 10s, 5m  Will be converted to int (miliseconds) in the code
	Timeout       time.Duration         `yaml:"-"`
	Exception     string                `yaml:"exception"`
	Queue         string                `yaml:"queue"` // Only Valid in workflow-service
	Next          []WorkflowNext_t      `yaml:"next"`
	MergeStrategy []MergeStrategyFunc_t `yaml:"mergestrategy"`
}

type WorkflowNext_t struct {
	PathStep string     `yaml:"step"` // Defines the next step. Use "exit" to end the workflow
	StepType StepType_t `yaml:"-"`
	Branch   string     `yaml:"branch"` // Branch is required when there are more than 1 next steps (excluding exits)
}

type MergeStrategyFunc_t struct {
	Action string `yaml:"action"`
	Field  string `yaml:"field"`
}

type TwistyService_t struct {
	Name       string `yaml:"name"`
	Queue      string `yaml:"queue"`
	Exchange   string `yaml:"exchange"`
	RoutingKey string `yaml:"routingkey"`
}

type WorkflowExit_t struct {
	ReplyTo string `yaml:"replyto"`
}

type Merge_t struct {
	Completed bool      `json:"completed"`
	Timestamp time.Time `json:"timestamp"`
	Parent    string    `json:"parent"`
	Payload   []byte    `json:"payload"`
}

type StepType_t int

const (
	STEPTYPE_TWISTYGOSERVICE  StepType_t = 10
	STEPTYPE_INTERNALFUNCTION StepType_t = 20
	STEPTYPE_MERGE            StepType_t = 30
	STEPTYPE_EXCEPTION        StepType_t = 80
	STEPTYPE_EXIT             StepType_t = 90
	STEPTYPE_TERMINATE        StepType_t = 95
)
