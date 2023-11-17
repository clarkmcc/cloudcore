-- Status enum for soft deletes
CREATE TYPE IF NOT EXISTS "status" AS ENUM ('active', 'deleted');

-- If cloud-hosted, a tenant represents a user and allows for better optimized
-- or geo-located queries using data-domiciling techniques.
CREATE TABLE IF NOT EXISTS "tenant" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "status" STATUS NOT NULL DEFAULT 'active',
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "name" STRING NOT NULL,
    "description" STRING,

    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "user" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "status" STATUS NOT NULL DEFAULT 'active',
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "tenant_id" UUID NOT NULL,
    "subject" STRING NOT NULL,

    PRIMARY KEY ("id"),
    FOREIGN KEY ("tenant_id") REFERENCES "tenant" ("id") ON DELETE CASCADE,
    UNIQUE ("subject") WHERE "status" = 'active'
);

CREATE TABLE IF NOT EXISTS "project" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "status" STATUS NOT NULL DEFAULT 'active',
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "tenant_id" UUID NOT NULL,
    "name" STRING NOT NULL,
    "description" STRING,

    PRIMARY KEY ("id"),
    FOREIGN KEY ("tenant_id") REFERENCES "tenant" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "user_project" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "status" STATUS NOT NULL DEFAULT 'active',
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "project_id" UUID NOT NULL,
    "user_id" UUID NOT NULL,

    PRIMARY KEY ("id"),
    FOREIGN KEY ("project_id") REFERENCES "project" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE,
    UNIQUE ("project_id", "user_id")
);

-- Host represents a host machine that is running a cloudcore agent
CREATE TABLE IF NOT EXISTS "host" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "status" STATUS NOT NULL DEFAULT 'active',
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "project_id" UUID NOT NULL,

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
    FOREIGN KEY ("project_id") REFERENCES "project" ("id") ON DELETE CASCADE,
    UNIQUE INDEX "identifier_idx" ("identifier")
);

-- Agent is the cloudcore agent that runs on a host. An agent only reports
-- on a single host, but over the lifetime of a host, there may be multiple
-- agents that report on it.
CREATE TABLE IF NOT EXISTS "agent" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "status" STATUS NOT NULL DEFAULT 'active',
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "project_id" UUID NOT NULL,

    "host_id" UUID NOT NULL,
    "online" BOOL NOT NULL,
    "last_heartbeat_timestamp" TIMESTAMP NOT NULL,

    PRIMARY KEY ("id"),
    FOREIGN KEY ("project_id") REFERENCES "project" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("host_id") REFERENCES "host" ("id") ON DELETE CASCADE
);

-- Pre-shared key for agents to authenticate with
CREATE TABLE IF NOT EXISTS "agent_psk" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "status" STATUS NOT NULL DEFAULT 'active',
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "project_id" UUID NOT NULL,

    "name" STRING NOT NULL,
    "description" STRING,
    "key" STRING NOT NULL DEFAULT gen_random_uuid(),
    -- The number of times this PSK can be used before it cannot be used again
    "uses_remaining" INTEGER NOT NULL DEFAULT 1,
    -- The timestamp when this PSK expires and can no longer be used
    "expiration" TIMESTAMP,

    PRIMARY KEY ("id"),
    UNIQUE INDEX "key_idx" ("key")
);

-- A group of agents that can be targeted for reasons
CREATE TABLE IF NOT EXISTS "agent_group" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "status" STATUS NOT NULL DEFAULT 'active',
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "project_id" UUID NOT NULL,

    "name" STRING NOT NULL,
    "description" STRING,

    PRIMARY KEY ("id"),
    FOREIGN KEY ("project_id") REFERENCES "project" ("id") ON DELETE CASCADE
);

-- Associates an agent with an agent group
CREATE TABLE IF NOT EXISTS "agent_group_member" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "status" STATUS NOT NULL DEFAULT 'active',
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "project_id" UUID NOT NULL,

    "agent_id" UUID NOT NULL,
    "agent_group_id" UUID NOT NULL,

    PRIMARY KEY ("id"),
    FOREIGN KEY ("project_id") REFERENCES "project" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("agent_id") REFERENCES "agent" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("agent_group_id") REFERENCES "agent_group" ("id") ON DELETE CASCADE,
    UNIQUE INDEX "agent_id_agent_group_id_idx" ("agent_id", "agent_group_id")
);

-- Associates a pre-shared key with an agent group. When an agent registers
-- using a PSK, then we should automatically add it to the agent group.
CREATE TABLE IF NOT EXISTS "agent_group_psk" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "status" STATUS NOT NULL DEFAULT 'active',
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "project_id" UUID NOT NULL,

    "agent_group_id" UUID NOT NULL,
    "agent_psk_id" UUID NOT NULL,

    PRIMARY KEY ("id"),
    FOREIGN KEY ("project_id") REFERENCES "project" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("agent_group_id") REFERENCES "agent_group" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("agent_psk_id") REFERENCES "agent_psk" ("id") ON DELETE CASCADE,
    UNIQUE INDEX "agent_group_id_agent_psk_id_idx" ("agent_group_id", "agent_psk_id")
);

-- Create the default data
INSERT INTO "tenant" ("name", "description") VALUES ('Default', 'Default tenant');
INSERT INTO "project" ("tenant_id", "name", "description") VALUES ((SELECT "id" FROM "tenant" WHERE "name" = 'Default'), 'Default', 'Default project');

