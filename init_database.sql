DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

create table users (
	id integer primary key
);

create table segments (
    id integer primary key generated always as identity,
    name varchar not null,
    unique (name)
);

create table user_segment (
	id integer primary key generated always as identity,
	user_id bigint not null,
	segment_id bigint not null,
    ttl timestamp,
    constraint fk_user
        foreign key (user_id)
            references users (id)
                on delete restrict
                on update restrict,
    constraint fk_segment
        foreign key (segment_id)
            references segments (id)
                on delete restrict
                on update restrict,
    unique (user_id, segment_id)
);

create table user_segment_log (
    id integer primary key generated always as identity,
    user_id integer,
    segment_name varchar,
    operation varchar,
    operation_timestamp timestamp
)
