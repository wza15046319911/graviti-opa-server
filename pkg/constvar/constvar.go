package constvar

const (
	DefaultLimit = 50

	// common
	ApolloAppID      = "100003054"
	ApolloADDR       = "47.103.140.27:8080"
	ApolloCluster    = "default"
	ApolloNamespace  = "online.yaml"
	ApolloConfigType = "yaml"

	EnvDEV = "DEV"
	EnvFAT = "FAT"
	EnvUAT = "UAT"
	EnvPRO = "PRO"

	// config
	ConfigLocal     = "local"
	ConfigLocalPath = "/conf/config.yaml"
	ConfigOnline    = "online"
)

const KeyCtxReqId string = "requestID"
