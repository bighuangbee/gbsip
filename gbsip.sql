

drop table if exists devices;
create table devices
(
    id            bigint auto_increment primary key,
    device_id      varchar(50)                                   not null comment '国标设备id',
    name          varchar(255)                                  not null comment '国标设备名称',
    manufacturer  varchar(255)                                  not null comment '厂商',
    model         varchar(255)                                  not null comment '型号',
    firmware      varchar(255)                                  not null comment '固件版本',
    transport     varchar(50)                                   not null comment '信令传输 UDP',
    status        tinyint(1) unsigned default '2'               not null comment '在线状态 1=在线 2=离线',
    last_alive_time timestamp                                     not null comment '最近一次心跳时间',
    host_address   varchar(50)                                   not null comment 'IPC地址',
    ip            varchar(50)                                   not null comment 'IPC IP地址',
    port          int                                           not null comment 'IPC端口',
    expires       int                                           not null comment '通道更新周期',
    charset       varchar(20)                                   not null comment '字符集 GB2312/UTF-8',
    created_at     timestamp           default CURRENT_TIMESTAMP not null comment '创建时间,注册时间',
    updated_at     timestamp                                     null comment '修改时间',
    deleted_at     timestamp                                     null comment '删除时间',
    constraint device_deviceId_unique
        unique (device_id, deleted_at)
) collate = utf8mb4_general_ci;

drop table if exists channels;
create table channels
(
    id            bigint auto_increment primary key,
    device_id      varchar(50)                                   not null comment '国标设备id',
    channel_id     varchar(50)                                   not null comment '国标通道id',
    name          varchar(255)                                  not null comment '国标设备名称',
    manufacturer  varchar(255)                                  not null comment '厂商',
    model         varchar(255)                                  not null comment '型号',
    firmware      varchar(255)                                  not null comment '固件版本',
    ptz_type      tinyint(1)          default '0'               not null comment '云台类型, 0-未知 1-球机 2-半球机 3-固定枪机 4-遥控枪机',
    ptz_type_text varchar(20)         default ''                not null comment '云台类型文本',
    addr          varchar(200)        default ''                not null comment '位置',
    status        tinyint(1) unsigned default '2'               not null comment '在线状态 1=在线 2=离线',
    lat           decimal(16, 8)                                null comment '纬度 精确到小数点后8位',
    lon           decimal(16, 8)                                null comment '经度 精确到小数点后8位',
    alt           decimal(16, 2)                                null comment '海拔 单位米 精确到小数点后2位',
    owner         varchar(50)                                   not null comment '',
    parent_id      varchar(50)                                   not null comment '',
    register_way   tinyint(1)                                    not null comment '',
    secrecy       int                                           not null comment '',
    stream_num     int                                           not null comment '',
    host_address     varchar(50)                                   not null comment 'IPC地址',
    ip            varchar(50)                                   not null comment 'IPC IP地址',
    port          int                                           not null comment 'IPC端口',
    expires       int                                           not null comment '通道更新周期',
    charset       varchar(20)                                   not null comment '字符集 GB2312/UTF-8',
    created_at     timestamp           default CURRENT_TIMESTAMP not null comment '创建时间,注册时间',
    updated_at     timestamp                                     null comment '修改时间',
    deleted_at     timestamp                                     null comment '删除时间',
    constraint device_deviceId_channelId_unique
        unique (device_id, channel_id, deleted_at)
) collate = utf8mb4_general_ci;

