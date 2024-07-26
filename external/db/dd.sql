create table recipes
(
    id          serial
        primary key,
    created_at  timestamp default now() not null,
    title       varchar(100)            not null
        constraint uq_title
            unique,
    description text,
    instruction text,
    publish     boolean   default false not null
);