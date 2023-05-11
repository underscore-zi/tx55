package configurations

type RestAPI struct {
	Host string
	Port int
	// ApiPrefix contains the first part of the URL path for all API endpoints like /api/v1
	ApiPrefix string
	// GameWebPrefix contains the start of the url path for all the endpoints the console hits. By defualt this should be /us/mgs3/
	GameWebPrefix string
	// TrustedProxies is the list of IPs to trust headers like X-Forwarded-For from
	TrustedProxies []string
	// Database is the primary database for the game server
	Database DatabaseConfig
	// SessionDatabase is purely for storing the sessions
	SessionDatabase DatabaseConfig
	// SessionsSecret is the secret used to secure the sessions
	SessionSecret string
	// Events is the configuration for the events websocket and reporting endpoints
	Events RestAPIEvents
	// RunCronJobs triggers whether the scheduled jobs like rank updates run
	RunCronJobs bool
}

type RestAPIEvents struct {
	// Enabled marks whether the restapi should support the events websocket and reporting endpoints
	Enabled bool
	// AccessTokens are unique strings that game lobbies use as the :token parameter when posting events
	AccessTokens []string
}
