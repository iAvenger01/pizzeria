create table employees (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name varchar(255) not null,
    last_name varchar(255) not null,
    status varchar(63) not null,
    work_time  integer not null,
    employment_date date not null default NOW(),
    dismissal_date date
);

alter table employees
    owner to pizzeria;

comment on column employees.id is 'Идентификатор сотрудника';
comment on column employees.first_name is 'Имя сотрудника';
comment on column employees.last_name is 'Фамилия сотрудника';
comment on column employees.status is 'Статус сотрудника';
comment on column employees.work_time is 'Продолжительность рабочего дня (в секундах)';
comment on column employees.employment_date is 'Дата первого трудоустройства сотрудника';
comment on column employees.dismissal_date is 'Дата увольнения сотрудника';

create table cooks
(
    employee_id uuid not null constraint cooks_employee_id_fk
        references employees (id),
    number_processed_order integer not null default 0
);

alter table cooks
    owner to pizzeria;

comment on column cooks.employee_id is 'Идентификатор сотрудника';
comment on column cooks.number_processed_order is 'Количество обработанных заказов';

create table couriers
(
    employee_id uuid not null constraint courier_employee_id_fk
        references employees (id),
    bag_size integer not null default 0,
    number_processed_order integer not null default 0
);

alter table couriers
    owner to pizzeria;

comment on column couriers.employee_id is 'Идентификатор сотрудника';
comment on column couriers.number_processed_order is 'Количество обработанных заказов';

alter table employees rename column shift_duration to work_time;