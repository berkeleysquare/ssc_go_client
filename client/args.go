package client

import (
	"errors"
	"flag"
	"os"
)

const TEST_USER = "Administrator"
const TEST_PASS = "spectra"
const TEST_DOMAIN = ""
const TEST_URL = "https://localhost/openapi"


type Arguments struct {
	Url string
	Command string
	Bucket string
	Job string
	Name string
	FileName string
	Domain string
	Password string
	Directory string
	IncludeDirectory string
	ExcludeDirectory string
	Share string
	Target string
	Clone string
	ProjectName string
	Prefix string
	Start, Count int
	IgnoreCert bool
	DontWaitForTape bool
	MakeLinks string
	Endpoint, Proxy string
	AccessKey, SecretKey string
}

func ParseArgs() (*Arguments, error) {
	url := flag.String("url", TEST_URL, "Base REST endpoint path for StorCycle server")
	name := flag.String("name", TEST_USER, "user name")
	fileName := flag.String("file_name", "", "user name")
	domain := flag.String("domain", TEST_DOMAIN, "domain")
	password := flag.String("password", TEST_PASS, "user password")
	command := flag.String("command", "", "command to execute; list_commands to list all")
	bucket := flag.String("bucket", "", "bucket name")
	job := flag.String("job", "", "job name")
	projectName := flag.String("project_name", "", "project name")
	prefix := flag.String("prefix", "", "string to start share name")
	start := flag.Int("start", 0, "number/auffix of first share")
	count := flag.Int("count", 1000, "number of shares")
	directory := flag.String("directory", "\\", "directory on share")
	includeDirectory := flag.String("include_directory", "", "explicitly include directory")
	excludeDirectory := flag.String("exclude_directory", "", "exclude directory")
	share := flag.String("share", "", "source location name")
	target := flag.String("target", "", "target location name")
	clone := flag.String("clone", "", "target to clone")
	makeLinks := flag.String("make_links", "none", "none|single|individual|symbolic replace files with HTML or sym links")
	ignoreCert := flag.Bool("ignore_cert", false, "use https with self-signed certificate")
	dontWaitForTape := flag.Bool("dont_wait_for_tape", false, "skip wait for tape placement")
	endpointParam := flag.String("endpoint", "", "Specifies the url to the DS3 server.")
	accessKeyParam := flag.String("access_key", "", "Specifies the access_key for the DS3 user.")
	secretKeyParam := flag.String("secret_key", "", "Specifies the secret_key for the DS3 user.")
	proxyParam := flag.String("proxy", "", "Specifies the HTTP proxy to route through.")
	flag.Parse()

	// Build the arguments object.
	args := Arguments{
		Url: *url,
		Command: *command,
		Bucket: *bucket,
		Job: *job,
		Name: *name,
		FileName: *fileName,
		Password: *password,
		Domain: *domain,
		Prefix: *prefix,
		Directory: *directory,
		IncludeDirectory: *includeDirectory,
		ExcludeDirectory: *excludeDirectory,
		Share: *share,
		Target: *target,
		Clone: *clone,
		ProjectName: *projectName,
		Start: *start,
		Count: *count,
		MakeLinks: *makeLinks,
		IgnoreCert: *ignoreCert,
		DontWaitForTape: *dontWaitForTape,
		Endpoint: paramOrEnv(*endpointParam, "DS3_ENDPOINT"),
		AccessKey: paramOrEnv(*accessKeyParam, "DS3_ACCESS_KEY"),
		SecretKey: paramOrEnv(*secretKeyParam, "DS3_SECRET_KEY"),
		Proxy: paramOrEnv(*proxyParam, "DS3_PROXY"),}

	// Validate required arguments.
	switch {
		case args.Command == "": return nil, errors.New("Must specify a command.")
		default: return &args, nil
	}

}

func ValueOrDefault(argValue string, defaultValue string) string {
	switch {
	case argValue != "":
		return argValue
	case defaultValue != "":
		return defaultValue
	default:
		return ""
	}
}

func paramOrEnv(param, envName string) string {
	env := os.Getenv(envName)
	switch {
	case param != "": return param
	case env != "": return env
	default: return ""
	}
}

