
drop table if exists `administrator`;

create table if not exists `administrator` (
  `id` bigint unsigned not null auto_increment,
  `created_at` timestamp null default null,
  `updated_at` timestamp null default null,
  `deleted_at` timestamp null default null,
  `status` tinyint not null default 1,
  `mobile` varchar(11) not null,
  `password` varchar(255) not null,
  `nickname` varchar(32) not null,
  `avatar_url` varchar(128) null default null,
  primary key (`id`),
  unique key `uk_mobile` (`mobile`)
);

drop table if exists `user`;

create table if not exists `user` (
  `id` bigint unsigned not null auto_increment,
  `created_at` timestamp null default null,
  `updated_at` timestamp null default null,
  `deleted_at` timestamp null default null,
  `status` tinyint not null default 1,
  `mobile` varchar(11) null default null,
  `password` varchar(255) null default null,
  `nickname` varchar(32) null default null,
  `avatar_url` varchar(128) null default null,
  primary key (`id`),
  unique key `uk_mobile` (`mobile`)
);

drop table if exists `wechat_mini_program_user`;

create table if not exists `wechat_mini_program_user` (
  `id` bigint unsigned not null auto_increment,
  `created_at` timestamp null default null,
  `updated_at` timestamp null default null,
  `deleted_at` timestamp null default null,
  `open_id` varchar(32) not null,
  `union_id` varchar(32) null default null,
  `session_key` varchar(32) not null,
  `nickname` varchar(32) null default null,
  `avatar_url` varchar(128) null default null,
  `gender` tinyint not null default 0,
  `province` varchar(16) null default null,
  `city` varchar(16) null default null,
  `country` varchar(16) null default null,
  `language` varchar(8) null default null,
  `user_id` bigint unsigned not null,
  primary key (`id`),
  unique key `uk_wechat` (`open_id`),
  unique key `uk_user` (`user_id`)
);
