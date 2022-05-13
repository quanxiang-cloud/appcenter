/*
Copyright 2022 QuanxiangCloud Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
    app_sign    varchar(30)  null,
    constraint t_app_center_app_name_uindex
        unique (app_name)
);

create table t_app_user_relation
(
    user_id varchar(64) null,
    app_id  varchar(64) null
);

create table t_app_scope
(
    app_id   varchar(64) null,
    scope_id varchar(64) null
)
