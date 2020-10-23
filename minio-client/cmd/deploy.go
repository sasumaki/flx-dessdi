package cmd

import (
	"aiga/api"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/minio/minio-go"
	"github.com/spf13/cobra"
	"github.com/theherk/viper"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Enter local file for upload",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		Deploy(cmd, args)
	},
}

// Deploy to deploy modle
func Deploy(cmd *cobra.Command, args []string) {
	endpoint := viper.GetString("endpoint")
	accesskeyid := viper.GetString("accesskeyid")
	secretKey := viper.GetString("secretKey")

	minioClient, err := api.GetClient(endpoint, accesskeyid, secretKey)
	if err != nil {
		fmt.Println("Issues connecting to your bucket")
		er(err)
	}

	fileToUpload := viper.GetString("file")
	file, err := os.Open(fileToUpload)
	if err != nil {
		er(err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		er(err)
	}

	var bytes int64
	bytes = stat.Size()

	var kilobytes int64
	kilobytes = (bytes / 1024)

	var megabytes float64
	megabytes = float64(kilobytes) / float64(1024)

	base := filepath.Base(fileToUpload)
	stamp := time.Now().UTC().Unix()
	noSuffix := strings.TrimSuffix(base, filepath.Ext(base))
	fmt.Println("Uploading (" + fmt.Sprintf("%.2f", megabytes) + "MB)...")
	folder := noSuffix + "-" + strconv.FormatInt(stamp, 10) + "/"
	progress := pb.New64(bytes)
	progress.Start()
	slug := folder + "model" + filepath.Ext(base)
	_, err = minioClient.FPutObject("models", slug, fileToUpload, minio.PutObjectOptions{
		ContentType: "application/csv",
		Progress:    progress,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully uploaded object: ", "models/"+folder)
}
