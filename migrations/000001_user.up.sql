create table if not exists "users"
(
    id            serial       not null primary key,
    name          varchar(50)  not null,
    surname       varchar(50)  not null,
    last_name     varchar(50)  not null,
    login         varchar(100) not null unique,
    email         varchar(100) not null unique,
    phone         varchar(11)  not null unique,
    hash_pass     varchar(100) not null,
    is_admin      boolean default false,
    is_blocked    boolean default false,
    confirm_email boolean default false
);
