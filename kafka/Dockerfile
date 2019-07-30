FROM bitnami/zookeeper:latest
USER root
RUN useradd -rm -d /home/zoo -s /bin/bash -g root -G sudo -u 1001 zoo
RUN mkdir -p /bitnami/zookeeper/data && chown -R 1001:0 /bitnami/* && chmod 755 /bitnami/zookeeper
ENV PATH=/opt/bitnami/java/bin:/opt/bitnami/zookeeper/bin:/opt/bitnami/nami/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin \
    IMAGE_OS=debian-9 \
    NAMI_VERSION=1.0.0-1 \
    GPG_KEY_SERVERS_LIST="ha.pool.sks-keyservers.net hkp://p80.pool.sks-keyservers.net:80 keyserver.ubuntu.com hkp://keyserver.ubuntu.com:80 pgp.mit.edu" \
    TINI_VERSION=v0.13.2 \
    TINI_GPG_KEY=595E85A6B1B4779EA4DAAEC70B588DFF0527A9B7 \
    GOSU_VERSION=1.10 \
    GOSU_GPG_KEY=B42F6819007F00F88E364FD4036A9C25BF357DD4 \
    BITNAMI_IMAGE_VERSION=3.5.5-debian-9-r55 \
    BITNAMI_PKG_CHMOD="-R g+rwX" \
    HOME=/ \
    ALLOW_ANONYMOUS_LOGIN=no \
    BITNAMI_APP_NAME=zookeeper \
    JVMFLAGS="" \
    NAMI_PREFIX=/.nami \
    ZOO_CLIENT_PASSWORD= \
    ZOO_CLIENT_USER= \
    ZOO_ENABLE_AUTH=no \
    ZOO_HEAP_SIZE=1024 \
    ZOO_INIT_LIMIT=10 \
    ZOO_LOG_LEVEL=INFO \
    ZOO_MAX_CLIENT_CNXNS=60 \
    ZOO_PORT_NUMBER=2181 \
    ZOO_RECONFIG_ENABLED=no \
    ZOO_SERVERS="" \
    ZOO_SERVER_ID=1 \
    ZOO_SERVER_PASSWORDS="" \
    ZOO_SERVER_USERS="" \
    ZOO_SYNC_LIMIT=5 \
    ZOO_TICK_TIME=2000

USER zoo
RUN touch /bitnami/zookeeper/.initialized && touch /bitnami/zookeeper/.restored
VOLUME ["/bitnami/zookeeper/data"]
EXPOSE 2181 2888 3888
CMD ["/run.sh"]
ENTRYPOINT ["/app-entrypoint.sh"]

