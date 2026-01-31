alter table orders add column status varchar(255) not null default 'accepted';

alter table order_item add column price decimal(10, 2) not null default 0.0