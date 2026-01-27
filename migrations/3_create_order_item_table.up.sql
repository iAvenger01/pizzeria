create table order_item
(
    order_id     uuid not null
        constraint order_item_orders_id_fk
            references orders (id),
    menu_item_id integer not null
        constraint order_item_menu_id_fk
            references menu (id),
    quantity     smallint not null,
    constraint order_item_pk
        unique (order_id, menu_item_id)
);

comment on column order_item.order_id is 'Идентификатор заказа';
comment on column order_item.menu_item_id is 'Идентификатор продукта из меню';
comment on column order_item.quantity is 'Количество продуктов в заказе';