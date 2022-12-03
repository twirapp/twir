FROM archlinux:base-20221127.0.105785

ARG NODE_VERSION="18.12.1"
ARG GO_VERSION="2:1.19.3"
ARG PNPM_VERSION="7.17.1"

ARG USERNAME=tsuwari
ARG USER_UID=1000
ARG USER_GID=$USER_UID
RUN groupadd --gid $USER_GID $USERNAME && \
    useradd --uid $USER_UID --gid $USER_GID -m $USERNAME

RUN pacman -Sy --noconfirm gcc git sudo nodejs=$NODE_VERSION go=$GO_VERSION npm protobuf=21.10-1 && \
    npm i -g pnpm@7.17.1 @go-task/cli

RUN echo "${USERNAME} ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers
    

# switch user on build time
RUN echo "${USERNAME}:tsuwari" | chpasswd
USER ${USERNAME}
ENV PATH="$PATH:/home/$USERNAME/go/bin"
RUN SHELL=/bin/bash pnpm setup
# install neede go tools
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 && \
    go install github.com/mitranim/gow@latest && \
    go install mvdan.cc/gofumpt@latest && \
    go install github.com/segmentio/golines@latest && \
    go install github.com/twitchdev/twitch-cli@latest && \
    go install -v golang.org/x/tools/gopls@latest && \
    go install github.com/swaggo/swag/cmd/swag@3fe9ca22de310099640d9c96cafb2b787a5820f8