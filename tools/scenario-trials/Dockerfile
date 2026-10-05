FROM node:24-slim
RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates ripgrep && rm -rf /var/lib/apt/lists/*
ARG CLAUDE_CODE_VERSION=2.1.289
RUN npm install -g @anthropic-ai/claude-code@${CLAUDE_CODE_VERSION}
COPY clue /usr/local/bin/clue
RUN chmod +x /usr/local/bin/clue
USER node
WORKDIR /home/node
RUN git config --global user.name trial && git config --global user.email trial@example.invalid && git config --global init.defaultBranch main
