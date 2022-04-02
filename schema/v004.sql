alter table t_app_center
    add extension text null;

alter table t_app_center
    add description text null;



ALTER TABLE t_app_center ADD COLUMN server INT COMMENT 'initialized modules of app' AFTER use_status;
