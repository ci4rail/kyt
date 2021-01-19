FROM bitnami/kubectl:1.20.2

USER root

# install necessary packages
RUN apt-get update && \
    apt-get install -y gosu

# add entrypoint script to switch from root-user
# to user (with current userid)
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh
ENTRYPOINT [ "/usr/local/bin/entrypoint.sh" ]

WORKDIR /home/user/
