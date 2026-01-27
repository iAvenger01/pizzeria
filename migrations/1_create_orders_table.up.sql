create table orders
(
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    address    smallint not null,
    status     varchar(63) not null,
    created_at timestamp default NOW() not null,
    updated_at timestamp default NOW() not null
);

alter table orders
    owner to pizzeria;

comment on column orders.id is 'Идентификатор заказа';

comment on column orders.address is 'Адрес доставки заказа';

comment on column orders.status is 'Статус заказа';

comment on column orders.created_at is 'Дата и время создания заказа';

comment on column orders.updated_at is 'Дата и время последнего обновления заказа';
