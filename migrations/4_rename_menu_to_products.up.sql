alter table menu
    rename constraint menu_pk to products_pk;

alter table menu rename constraint menu_uk_key to products_uk_key;

alter table order_item
    rename column menu_item_id to product_id;

alter table menu
    rename to products;

alter table order_item rename constraint order_item_menu_id_fk to order_item_product_id_fk;

alter table order_item add column status varchar(255) not null default 'accepted'