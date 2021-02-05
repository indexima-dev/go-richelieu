package constants

const (
	AppName               string = "Richelieu"
	AppDescription        string = "Data generator that respects cardinality and schema structures#."
	AppVersion            string = "0.1"
	DefaultPlanFile       string = "plan.json"
	DefaultSchemaFile     string = "./input/createDataspace.txt"
	DefaultDictionaryFile string = "./input/showDictionaries.csv"
	DefaultTableCountFile string = "./input/showTableCount.csv"
	DefaultS3Repository   string = "s3://indexima-data/dummy_data/generated1/source"
)
