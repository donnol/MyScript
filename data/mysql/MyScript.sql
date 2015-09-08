#创建数据库
use myscript;

#创建配置表
drop table if exists t_user;

create table t_user(
	userId integer not null auto_increment,
	name varchar(128) not null,
	phone varchar(128) not null,
	createTime timestamp not null default CURRENT_TIMESTAMP,
	#modifyTime timestamp not null default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP,
	primary key( userId )
)engine=innodb default charset=utf8mb4 auto_increment = 10001;

alter table t_user add index nameIndex(name);
