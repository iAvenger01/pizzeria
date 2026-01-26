create table orders
(
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamp default NOW() not null,
    updated_at timestamp default NOW() not null
);

alter table orders
    owner to pizzeria;