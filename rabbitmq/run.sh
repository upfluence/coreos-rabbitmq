#! /bin/bash

/usr/bin/cluster-bootstrap
/usr/bin/envtmpl -i /opt/rabbitmq.conf.tmpl -o /etc/rabbitmq/rabbitmq.conf

exec rabbitmq-server
