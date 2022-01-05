-- 测试自动化
create table t_app_center
(
    id          varchar(64)  not null
        primary key,
    app_name    varchar(80)  null,
    access_url  varchar(200) null,
    app_icon    text         null,
    create_by   varchar(64)  null,
    update_by   varchar(64)  null,
    create_time bigint       null,
    update_time bigint       null,
    use_status  bigint       null,
    constraint t_app_center_app_name_uindex
        unique (app_name)
);

create table t_app_user_relation
(
    user_id varchar(64) null,
    app_id  varchar(64) null
);

create table t_app_scope(
    app_id varchar(64) null ,
    scope_id varchar(64) null
)
