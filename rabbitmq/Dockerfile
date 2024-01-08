FROM alpine:latest

RUN addgroup -S rabbitmq && \
  adduser -S -h /var/lib/rabbitmq -G rabbitmq rabbitmq

RUN apk upgrade && apk add --no-cache \
		bash \
    curl \
    su-exec \
    erlang

ARG RABBITMQ_VERSION=3.12.11

ENV RABBITMQ_LOGS=- \
    RABBITMQ_SASL_LOGS=- \
    RABBITMQ_HOME=/opt/rabbitmq \
    GPG_KEY=0A9AF2115F4687BD29803A206B73A36E6026DFCA \
    RABBITMQ_VERSION=${RABBITMQ_VERSION} \
    HOME=/var/lib/rabbitmq \
    ENVTMPL_VERSION=0.0.1

ENV ENVTPL_REPO=https://github.com/upfluence/envtmpl/releases/download/v${ENVTMPL_VERSION}/envtmpl-linux-amd64-${ENVTMPL_VERSION} \
    PATH=$RABBITMQ_HOME/sbin:$PATH

RUN curl -sL ${ENVTPL_REPO} > /usr/bin/envtmpl && \
    chmod +x /usr/bin/envtmpl

RUN set -ex; \
	\
	apk add --no-cache --virtual .build-deps \
		ca-certificates \
		libressl \
		tar \
		xz \
	; \
	wget -O rabbitmq-server.tar.xz "https://github.com/rabbitmq/rabbitmq-server/releases/download/v${RABBITMQ_VERSION}/rabbitmq-server-generic-unix-${RABBITMQ_VERSION}.tar.xz"; \
	mkdir -p "$RABBITMQ_HOME"; \
	mkdir -p "/etc/rabbitmq"; \
	tar \
		--extract \
		--verbose \
		--file rabbitmq-server.tar.xz \
		--directory "$RABBITMQ_HOME" \
		--strip-components 1 \
	; \
	rm rabbitmq-server.tar.xz; \
	grep -qE '^SYS_PREFIX=\$\{RABBITMQ_HOME\}$' "$RABBITMQ_HOME/sbin/rabbitmq-defaults"; \
	sed -ri 's!^(SYS_PREFIX=).*$!\1!g' "$RABBITMQ_HOME/sbin/rabbitmq-defaults"; \
	grep -qE '^SYS_PREFIX=$' "$RABBITMQ_HOME/sbin/rabbitmq-defaults"; \
	\
	apk del .build-deps;

ADD erlang.cookie /var/lib/rabbitmq/.erlang.cookie
RUN chown rabbitmq:rabbitmq /var/lib/rabbitmq/.erlang.cookie && \
    chmod 400 /var/lib/rabbitmq/.erlang.cookie && mkdir /etc/rabbitmq/conf.d

RUN rabbitmq-plugins enable --offline rabbitmq_management rabbitmq_prometheus
ADD rabbitmq.conf.tmpl /opt/rabbitmq.conf.tmpl
ADD run.sh /opt
ADD dist/cluster-bootstrap /usr/bin

CMD ["/opt/run.sh"]
