package interfaces

import (
	"github.com/getsentry/sentry-go"
	"github.com/permitio/permit-golang/pkg/permit"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"sync"
	"time"
)

type Conf struct {
	Database *Database
	Server   *Server
	App      *App
}
type App struct {
	JwtSecret string
	PageSize  int
	PrefixUrl string

	RuntimeRootPath string

	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

	ExportSavePath string
	QrCodeSavePath string
	FontSavePath   string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PDP_SERVER   string
	PDP_TOKEN    string
}
type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

type Sentry struct {
	DNS           string
	EnableTracing bool
	BeforeSend    func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event
	DEBUG         bool
}

type GracefulServer struct {
	httpServer  *http.Server
	stopping    chan bool
	status      bool //server 状态
	App         *Application
	MongoClient *mongo.Client
	Permit      *permit.Client
	sync.Mutex
}

type Kafka struct {
	Host         string
	DefaultTopic string
}
