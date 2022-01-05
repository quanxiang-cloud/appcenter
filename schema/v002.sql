/*修改sql*/
ALTER table t_app_center ADD del_flag TINYINT;

UPDATE t_app_center set del_flag = 0;

ALTER table t_app_center ADD delete_time BIGINT;

UPDATE t_app_center set delete_time = 0;

DROP INDEX t_app_center_app_name_uindex ON t_app_center;