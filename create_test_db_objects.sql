create table party_types (
   id    int4 not null,
   name  text not null,
   constraint party_types_pk primary key(id)
);

insert into party_types(id, name) values(1, 'Individual'), 
(2, 'Organization');

create table party_mltypes (
   id    int4 not null,
   name  jsonb not null,
   constraint party_mltypes_pk primary key(id)
); 
insert into party_mltypes(id, name)  values
(1, '{"en" : "Individual",  "ru": "Физ.лицо"}'::jsonb), 
(2, '{"en" : "Organzation", "ru": "Организация"}'::jsonb);
  

create table event_types (
   id          int4 not null,
   code        text not null,
   name        text not null,
   is_critical bool not null default false,
   constraint event_types_pk primary key(id)
); 
insert into event_types(id, code, name, is_critical) values
(1, 'CONLOST', 'Connection Lost', true), 
(2, 'SERVERUP', 'Server up', false);

create table event_mltypes (
   id          int4 not null,
   code        text not null,
   name        jsonb not null,
   is_critical bool not null default false,
   constraint event_mltypes_pk primary key(id)
); 
insert into event_mltypes(id, code, name, is_critical) values
(1, 'CONLOST', '{"en":"Connection lost", "ru": "Связь потеряна"}'::jsonb, true), 
(2, 'SERVERUP', '{"en":"Server up", "ru": "Сервер поднят"}'::jsonb, false);
 