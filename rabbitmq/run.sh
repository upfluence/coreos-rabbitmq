#! /bin/bash
confd -verbose -interval=30 -backend etcd -node 172.17.42.1:4001 -config-file /etc/confd/conf.d/hosts.toml &
confd -verbose -onetime -backend etcd -node 172.17.42.1:4001 -config-file /etc/confd/conf.d/rabbitmq.toml
rabbitmq-plugins enable rabbitmq_management
rabbitmq-server
