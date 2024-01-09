package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/upfluence/cfg"
	"github.com/upfluence/pkg/log"
	client "go.etcd.io/etcd/client/v2"
)

type config struct {
	Etcd struct {
		URL string
		Key string
	}

	NodePrefix string
	ConfigFile string
	EtcHosts   string
}

func (c *config) etcdClient() (client.KeysAPI, error) {
	cl, err := client.New(
		client.Config{
			Endpoints: []string{c.Etcd.URL},
			Transport: client.DefaultTransport,
		},
	)

	if err != nil {
		return nil, err
	}

	return client.NewKeysAPI(cl), nil
}

func main() {
	var (
		c = config{
			Etcd: struct {
				URL string
				Key string
			}{URL: "http://localhost:2379", Key: "/discovery/rabbitmq"},
			NodePrefix: "rabbit@",
			ConfigFile: "/etc/rabbitmq/conf.d/cluster.conf",
			EtcHosts:   "/etc/hosts",
		}
		ctx = context.Background()
	)

	if err := cfg.NewDefaultConfigurator().Populate(ctx, &c); err != nil {
		log.WithError(err).Fatal("cant build the config")
	}

	cl, err := c.etcdClient()

	if err != nil {
		log.WithError(err).Fatal("cant build etcd client")
	}

	resp, err := cl.Get(ctx, c.Etcd.Key, &client.GetOptions{Sort: true})

	if err != nil {
		log.WithError(err).Fatal("cant fetch etcd list")
	}

	configFile, err := os.OpenFile(c.ConfigFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {
		log.WithError(err).Fatal("cant create config file")
	}

	etcHostsFile, err := os.OpenFile(c.EtcHosts, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.WithError(err).Fatal("cant create etc/hosts file")
	}

	for i, node := range resp.Node.Nodes {
		key := node.Key

		for _, prefix := range []string{c.Etcd.Key, "/"} {
			key = strings.TrimPrefix(key, prefix)
		}

		fmt.Fprintf(configFile, "cluster_formation.classic_config.nodes.%d = %s%s\n", i+1, c.NodePrefix, key)
		fmt.Fprintf(etcHostsFile, "%s\t%s\n", node.Value, key)
	}

	for _, f := range []*os.File{configFile, etcHostsFile} {
		if err := f.Sync(); err != nil {
			log.WithField(log.Field("file_name", f.Name())).WithError(err).Warning("cant sync file")
		}

		if err := f.Close(); err != nil {
			log.WithField(log.Field("file_name", f.Name())).WithError(err).Warning("cant close file")
		}
	}
}
