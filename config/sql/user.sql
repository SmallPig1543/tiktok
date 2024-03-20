CREATE TABLE tiktok.user
(
    `id`         INT AUTO_INCREMENT PRIMARY KEY,
    `username`   varchar(100),
    `password`   varchar(100),
    `avatar`     varchar(100),
    `otp_secret` varchar(100),
    `created_at` timestamp default CURRENT_TIMESTAMP,
    `updated_at` timestamp default CURRENT_TIMESTAMP,
    `deleted_at` timestamp
) engine = InnoDB
  default charset = utf8mb4;