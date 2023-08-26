DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

create table users (
	id integer primary key
);

create table segments (
    id integer primary key generated always as identity,
    name varchar,
    unique (name)
);

create table user_segment (
	id integer primary key generated always as identity,
	user_id bigint,
	segment_id bigint,
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
