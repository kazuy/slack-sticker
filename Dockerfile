# builder
FROM node:13.8.0-buster as builder

# update
RUN apt-get update

# app
FROM golang:1.13.8-buster

# update
RUN apt-get update

# copy from builder
ENV YARN_VERSION=1.21.1
COPY --from=builder /opt/yarn-v$YARN_VERSION /opt/yarn/
COPY --from=builder /usr/local/bin/node /usr/local/bin/
COPY --from=builder /usr/local/lib/node_modules /usr/local/lib/node_modules/
RUN ln -s /opt/yarn/bin/yarn /usr/local/bin/yarn && \
    ln -s /opt/yarn/bin/yarnpkg /usr/local/bin/yarnpkg && \
    ln -s /usr/local/bin/node /usr/local/bin/nodejs && \
    ln -s /usr/local/lib/node_modules/npm/bin/npm-cli.js /usr/local/bin/npm && \
    ln -s /usr/local/lib/node_modules/npm/bin/npx-cli.js /usr/local/bin/npx

# user
RUN useradd -m slack-sticker && \
    gpasswd -a slack-sticker sudo && \
    echo "slack-sticker:slack-sticker" | chpasswd

# workdir
WORKDIR /home/usr/slack-sticker/app
RUN chown slack-sticker:slack-sticker /home/usr/slack-sticker

# user
USER slack-sticker

# install serverless
ENV NPM_CONFIG_PREFIX=/home/usr/slack-sticker/.npm-global \
    PATH=$PATH:/home/usr/slack-sticker/.npm-global/bin
RUN npm install -g serverless

