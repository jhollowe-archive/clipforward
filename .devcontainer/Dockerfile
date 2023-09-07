ARG VARIANT=1
FROM mcr.microsoft.com/devcontainers/go:${VARIANT}

RUN apt-get update && apt-get install -y --no-install-recommends \
  libx11-dev \
  xvfb \
  && rm -rf /var/lib/apt/lists/*

# run virtual display init on start of the container
COPY .devcontainer/docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh
ENTRYPOINT [ "/docker-entrypoint.sh" ]

# keep the container running since the devcontainer will start its own processes
CMD [ "sleep", "infinity" ]
