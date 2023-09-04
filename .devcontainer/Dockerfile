ARG VARIANT=1
FROM mcr.microsoft.com/devcontainers/go:${VARIANT}

RUN apt-get update && apt-get install -y --no-install-recommends \
  libx11-dev \
  && rm -rf /var/lib/apt/lists/*
