package credentials

// These values are set at build time using -ldflags
var (
	AccessKey    string
	SecretKey    string
	Region       string
	BucketName   string
)

// HasEmbeddedCredentials returns true if credentials were embedded at build time
func HasEmbeddedCredentials() bool {
	return AccessKey != "" && SecretKey != "" && Region != "" && BucketName != ""
}