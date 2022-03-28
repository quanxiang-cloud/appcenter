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

-- auto-generated definition
create table t_app_template
(
    id           varchar(64)  not null,
    name         varchar(80)  null comment 'template name',
    app_icon     text         null comment 'app icon',
    path         varchar(200) null comment 'file server path',
    source_id    varchar(64)  null comment 'source app id',
    source_name  varchar(80)  null comment 'source app name',
    version      varchar(64)  null comment 'template version',
    group_id     varchar(64)  null comment 'group id',
    created_by   varchar(64)  null,
    created_name varchar(64)  null,
    created_time bigint       null,
    updated_by   varchar(64)  null,
    updated_name varchar(64)  null,
    updated_time bigint       null,
    status       int          null comment 'publish status：0:private 1:public',
    constraint t_app_template_id_uindex
        unique (id)
)
    comment 'app template table';

alter table t_app_template
    add primary key (id);

