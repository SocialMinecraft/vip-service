drop table if exists memberships;
drop table if exists accounts;

create table memberships (
    id bigserial primary key,
    email VARCHAR(100) not null unique,
    user_id VARCHAR(100) unique, /* can be null, if not claimed */
    start date NOT NULL,
    expire date NOT NULL
);

create table accounts (
    id bigserial primary key,
    user_id VARCHAR(100) NOT NULL,
    minecraft_uuid VARCHAR(100) NOT NULL unique
)