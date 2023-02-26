create extension if not exists "uuid-ossp";

create table if not exists users (
    uuid uuid primary key default uuid_generate_v4() not null,
    login character varying(100) not null,
    password character varying(100) not null,
    balance decimal(10, 2) default 0,
    created_at timestamp default now() not null,
    constraint users_unique_login unique (login)
);

create type status as enum ('NEW','PROCESSING', 'INVALID', 'PROCESSED');

create table if not exists orders (
    uuid uuid primary key default uuid_generate_v4() not null,
    number varchar(255) not null,
    status status not null,
    accrual decimal(10, 2) default 0,
    user_uuid uuid not null,
    uploaded_at timestamp default now() not null,
    constraint orders_unique_number unique (number),
    constraint orders_fk_user foreign key (user_uuid) references users (uuid)
);

create table if not exists withdraws (
    uuid uuid primary key default uuid_generate_v4() not null,
    order_number varchar(255) not null,
    amount decimal(10, 2) not null ,
    user_uuid uuid not null,
    processed_at timestamp default now() not null,
    constraint withdraws_unique_order_number unique (order_number),
    constraint withdraws_fk_user foreign key (user_uuid) references users (uuid)
)
