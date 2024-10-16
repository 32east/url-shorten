create table if not exists urls (
    id serial primary key not null,
    small_url char(7) unique not null,
    full_url text not null,
    ip_creator text not null,
    created_at timestamp not null
);

create table if not exists clicks (
    id serial primary key not null,
    small_url char(7) not null,
    date timestamp not null,
    ip text not null,
    country text,
    user_agent text not null,
    foreign key (small_url) references urls(small_url) on delete cascade
);

create index if not exists small_url_idx on urls(small_url);