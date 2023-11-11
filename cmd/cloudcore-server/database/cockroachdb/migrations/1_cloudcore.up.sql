-- Status enum for soft deletes
CREATE TYPE "status" AS ENUM ('active', 'deleted');

-- Host represents a host machine that is running a cloudcore agent
CREATE TABLE "hosts" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "status" STATUS NOT NULL DEFAULT 'active',
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- "identifier" is the unique identifier for the host and could be anything
    -- based on how the agent is configured. By default it will be the host ID
    -- provided through the operating system.
    "identifier" STRING NOT NULL,

    -- Host metadata
    "hostname" STRING,
    "host_id" STRING,
    "public_ip_address" STRING,
    "private_ip_address" STRING,
    "os_name" STRING,
    "os_family" STRING,
    "os_version" STRING,
    "kernel_architecture" STRING,
    "kernel_version" STRING,
    "cpu_model" STRING,
    "cpu_cores" INTEGER,

    PRIMARY KEY ("id"),
    -- Unique index on identifier
    UNIQUE INDEX "identifier_idx" ("identifier")
);

-- Agent is the cloudcore agent that runs on a host. An agent only reports
-- on a single host, but over the lifetime of a host, there may be multiple
-- agents that report on it.
CREATE TABLE "agents" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "status" STATUS NOT NULL DEFAULT 'active',
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "host_id" UUID NOT NULL,
    "online" BOOL NOT NULL,
    "last_heartbeat_timestamp" TIMESTAMP NOT NULL,

    PRIMARY KEY ("id"),
    FOREIGN KEY ("host_id") REFERENCES "hosts" ("id") ON DELETE CASCADE
);

