ALTER TABLE t_app_center ADD COLUMN server INT COMMENT 'initialized modules of app' AFTER use_status;

update t_app_center set use_status = -5;