# refbook [![GoDoc](https://godoc.org/github.com/refbook/axkit?status.svg)](https://godoc.org/github.com/axkit/refbook) [![Build Status](https://travis-ci.org/axkit/refbook.svg?branch=master)](https://travis-ci.org/axkit/refbook) [![Coverage Status](https://coveralls.io/repos/github/axkit/refbook/badge.svg)](https://coveralls.io/github/axkit/refbook) [![Go Report Card](https://goreportcard.com/badge/github.com/axkit/refbook)](https://goreportcard.com/report/github.com/axkit/refbook)

Easy access to database reference tables 

## Concept
Reference book is a plain table usually used as lookup. It has following characteristics:
* table has limited amount of rows
* new rows are introduced quite seldom
* primary key is an integer
* name can be in a single language represented by TEXT or multilanguage represented by JSON


### RDBMS supported
PostgreSQL 9.6+. 

## Examples
### Single Language, Plain Reference Table
```
  // create table party_types (
  //    id    int4 not null,
  //    name  text not null,
  //    constraint party_types_pk primary key(id)
  // ); 
  // insert into party_types(id, name) values(1, 'Individual'), (2, 'Organization');
  
  db, err := sql.Open("postgres", constr)
  pt := refbook.New().LoadFromSQL(db, "party_types")
  if pt.Err() != nil {
    // 
  }
  fmt.Println(pt.Name(1))     // Individual
  fmt.Println(pt.IsExist(2))  // true
  fmt.Println(pt.Name(3))     // ?  (as default response for if key not found). See var NotFoundName  
```
### Multi Language, Plain Reference Table
```
  // create table party_types (
  //    id    int4 not null,
  //    name  jsonb not null,
  //    constraint party_types_pk primary key(id)
  // ); 
  // insert into party_types(id, name) 
  // values(1, '{"en" : "Individual", "ru": "Физ.лицо"}'::jsonb), (2, '{"en" : "Organzation", "ru": "Организация"}'::jsonb);
  //
  // 
  db, err := sql.Open("postgres", constr)
  pt := refbook.NewMLRefBook().LoadFromSQL(db, "party_types")
  if pt.Err() != nil {
    // 
  }
  fmt.Println(pt.Lang("ru").Name(1)) // Физ.лицо"
  fmt.Println(pt.Lang("en").Name(2)) // Organzation"
  fmt.Println(pt.Lang("en").Name(3)) // ? (as default response if key not found)
```
### Single Language, Extended Reference Table
```
  // create table event_types (
  //    id          int4 not null,
  //    code        text not null,
  //    name        text not null,
  //    is_critical bool not null default false,
  //    constraint event_types_pk primary key(id)
  // ); 
  // insert into event_types(id, code, name, is_critical) values(1, 'CONLOST', 'Connection Lost', true), 
  // (2, 'SERVERUP', 'Server up', false);
  
  type EventType struct {
    ID          int
    Code        string
    Name        string
    IsCritical  bool  
  }
  var ets []EventType

  db, err := sql.Open("postgres", constr)
  //
  // read rows somehow to the slice et
  //
  et := refbook.New().LoadFromSlice(et, "ID", "Name")
  
  fmt.Println(et.Name(1)) // Individual
  fmt.Println(et.Name(3)) // ? (as default response for if key not found)
```
### Multi Language, Extended Reference Table
```
  // create table event_types (
  //    id          int4 not null,
  //    code        text not null,
  //    name        jsonb not null,
  //    is_critical bool not null default false,
  //    constraint event_types_pk primary key(id)
  // ); 
  // insert into event_types(id, code, name, is_critical) values
  // (1, 'CONLOST', '{"en":"Connection lost", "ru": "Связь потеряна"}'::jsonb, true), 
  // (2, 'SERVERUP', '{"en":"Server up", "ru": "Сервер поднят"}'::jsonb, false);
  
  type EventType struct {
    ID          int
    Code        string
    Name        string
    IsCritical  bool  
  }
  var ets []EventType

  db, err := sql.Open("postgres", constr)
  //
  // read rows somehow to the slice ets
  //
  et := refbook.NewMLRefBook().LoadFromSlice(ets, "ID", "Name")
  
  fmt.Println(et.Lang("ru").Name(1))  // Связь потеряна
  fmt.Println(et.Lang("en").Name(2))  // Server up
  fmt.Println(et.Lang("en").Name(3))  // ? 

  fmt.Println(et.IsExist(2))          // true
```
