create table if not exists urls (
    id serial primary key not null,
    small_url char(7) not null,
    full_url text not null,
    ip_creator text not null,
    created_at timestamp not null,
    clicks int not null default 0
);

create index if not exists small_url_idx on urls(small_url);