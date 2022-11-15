ALTER TABLE `app_center`.`t_app_scope` ADD COLUMN `type` VARCHAR(64);
update t_app_scope set type="user" where type is NULL;