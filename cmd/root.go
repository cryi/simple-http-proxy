package cmd

import (
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/cryi/simple-http-proxy/proxy"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "shp",
		Short: "Simple http proxy",
		Long:  "Simple http proxy server with support for header injection",
		Run: func(cmd *cobra.Command, args []string) {
			listen, _ := cmd.Flags().GetString("listen")
			if listen == "" {
				log.Error("Invalid listen address")
				os.Exit(1)
			}
			forward, _ := cmd.Flags().GetString("forward")
			if forward == "" {
				log.Error("Invalid forward address")
				os.Exit(2)
			}
			headers := make([]proxy.HttpHeader, 0)
			headersAsStrings, _ := cmd.Flags().GetStringArray("inject-header")
			r := regexp.MustCompile("(.*?):(.*)")

			for _, header := range headersAsStrings {
				match := r.FindStringSubmatch(header)
				if len(match) != 3 {
					continue
				}
				headers = append(headers, proxy.HttpHeader{
					Id:    match[1],
					Value: match[2],
				})
			}
			destinationUrl, err := url.Parse(forward)
			if err != nil {
				log.Error(err, "Invalid forward url!")
			}

			proxy := &proxy.Proxy{
				Destination: destinationUrl,
				Headers:     headers,
			}

			log.Println("Starting proxy server on", listen)
			if err := http.ListenAndServe(listen, proxy); err != nil {
				log.Fatal("ListenAndServe:", err)
			}

		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("listen", "l", "", "Address to listen on.")
	rootCmd.PersistentFlags().StringP("forward", "f", "", "Address to forward requests to.")
	rootCmd.PersistentFlags().StringArrayP("inject-header", "i", make([]string, 0), "List of headers to inject in format '<header>:<value>'")
}
