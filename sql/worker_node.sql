create table WORKER_NODE
(
    ID          bigint auto_increment comment 'auto increment id'
        primary key,
    HOST_NAME   varchar(64) not null comment 'host name',
    PORT        varchar(64) not null comment 'port',
    TYPE        int         not null comment 'node type: ACTUAL or CONTAINER',
    LAUNCH_DATE date        not null comment 'launch date',
    MODIFIED    datetime    not null,
    CREATED     datetime    not null
)
    comment 'DB WorkerID Assigner for UID Generator';