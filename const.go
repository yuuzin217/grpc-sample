package grpcsample

const (
	Host = "localhost:50051"

	RootCA_FilePath   = "../ssl/rootCA.pem"
	CertFile_filePath = "../ssl/localhost.pem"
	CertKey_filePath  = "../ssl/localhost-key.pem"

	Dir_storage        = "../storage/"
	Dir_storage_local  = Dir_storage + "local/"
	Dir_storage_remote = Dir_storage + "remote/"

	// キロバイト の単位
	KB = 1024
)
