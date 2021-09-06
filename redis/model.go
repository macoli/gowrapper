package redis

type AppStandAlone struct {
	AppID     int64
	AppName   string
	Password  string
	Masters   []string
	Slaves    []string
	Instances []string
}

type AppSentinel struct {
	AppID       int64
	AppName     string
	Password    string
	Sentinels   []string
	MasterNames []string
}

type AppCluster struct {
	AppID     int64
	AppName   string
	Password  string
	Instances []string
}

type App struct {
	Type    string
	AppInfo interface{}
}
