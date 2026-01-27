create table menu
(
    id              integer generated always as identity
        constraint menu_pk unique,
    key             varchar(255)             not null
        constraint menu_uk_key
            unique,
    name            varchar(255)             not null,
    price           decimal(10, 2) default 0 not null,
    assembling_time integer,
    cooking_time    integer                  not null
);

comment on column menu.id is 'Идентификатор продукта';

comment on column menu.key is 'Ключ продукта латиницей';

comment on column menu.name is 'Имя продукта';

comment on column menu.price is 'Цена продукта';

comment on column menu.assembling_time is 'Время наполнения (подготовки) в секундах';

comment on column menu.cooking_time is 'Время приготовления в секундах';