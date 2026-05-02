# To successfully initialize the Redis storage, you'll need to set up custom user settings for your storage.

### You need to implement configuration for Redis

File **redis.conf**:

    bind        ~host~    // (default: localhost or 0.0.0.0)
    port        ~port~    // (default: 6380)
    tcp-backlog ~backlog~ // (default: 511)

File **users.acl**:

    user default off
    user hairdress-root on >~your_password~ ~* +@all // ~your_password~ - paset your password

### The implementation runs on Docker.

### You'll need a Redis Docker image and a makefile from this repository, which will allow you to create a container on an open port with a custom user in a single command.

#### After steps above:

    make redis