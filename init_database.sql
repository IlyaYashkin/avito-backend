DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

create table users (
	id serial primary key,
	created_at timestamp with time zone
);

create table segments (
    id serial primary key,
    created_at timestamp with time zone,
    name varchar
);

create table user_segment (
	id serial primary key,
	created_at timestamp with time zone,
	user_id bigint,
	segment_id bigint,
    constraint fk_user
        foreign key (user_id)
            references users (id),
    constraint fk_segment
        foreign key (segment_id)
            references segments (id)
);
