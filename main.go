package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
)

var (
	listenAddr string
	socks5Addr string
)

var rootCmd = &cobra.Command{
	Use:   "s2h",
	Short: "A simple tool to convert socks5 proxy protocol to http proxy protocol",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Starting Socks5 Proxy Convert Server...")
		logrus.Infof("HTTP Listen Address: %s", listenAddr)
		logrus.Infof("Socks5 Server Address: %s", socks5Addr)

		err := http.ListenAndServe(listenAddr, http.HandlerFunc(serveHTTP))
		if err != nil {
			logrus.Error(err)
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-05-04 15:02:01",
	})

	rootCmd.PersistentFlags().StringVarP(&listenAddr, "listen", "l", "0.0.0.0:8081", "http listen address")
	rootCmd.PersistentFlags().StringVarP(&socks5Addr, "socks5", "s", "127.0.0.1:1080", "remote socks5 listen address")
}
