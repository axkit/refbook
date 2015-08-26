# refbook
Easy access to your database reference tables. Standalone or embedded. 

## Status
Planning & Requirement specification

## Concepts

#### RDBS supported
PostgreSQL first. 

#### Reference Book
Reference book is a plain table typically stores attributes of business objects as set of pairs {id, name}. Reference tables can be more complicated and have references between each other as (M:1). As instance: Towns->Countries->Currencies.

## Use cases
#### Dealing with 

## Features


### JSON API

Content-Type: application/json


Method|URL| Description
----|--------------|----------------------------------------------------------------------
GET |/refbook/| Returns list of reference books available for user
GET |/refbook/:name|Returns reference book metadata and list of available methods       
GET |/refbook/:name/stat|Returns reference book statistics
GET |/refbook/:name/item/    |Returns reference book items 
GET |/refbook/:name/item/:id |Returns single reference book item
GET |/refbook/:name/item/:id/exist|Returns HTTP OK 200 if row exists, 404 if not
POST|/refbook/:name/item/|Add item
PUT |/refbook/:name/item/:id|Update item
DELETE| /refbook/:name/item/:id|Mark as deleted
POST|/refbook/:name/item/:id/erase|Delete rows from database table

#### URL params for /refbook/:name/item
* **cols**= 
  * main - Default behavior. Returns reasonable columns.
  * all - Enrich each refbook item with system columns (createdDt, modifiedDt, etc)
  * colname1,colname2,colname3 - Returns specified columns.
* **enrich**=
  * n/a - No enrich. Default behavior.
  * nameonly - If refbookA refers to refBookB, column refbookB.name added into resultset. 
  * full     - If refbookA refers to refbookB, resultset item extends by refbookA {} object
* **enrichdepth**=
  * n/a, 0 - feature not used.
  * N - parent references level (1, 2, 3...).
* **lang**=
 * n/a - Returns items in basic language
 * $$ - Returns items in specified language (en, cz, ru, etc)
* **from**=
 * n/a - Returns items from cache. Default behavior.
 * db  - Reload cache first. 
* **deleted**
  * n/a - deleted rows are hidden 
  * y - include into a resultset the rows marked as deleted
* **orderby**
  * n/a - ordered by primary key column. Default behavior.
  * column1,column2 - sort by listed order
* result
 *  n/a - application/json. Default behavoir.
 *  options - set of option items
 *  html - simple html page, mostly for debug purpose.

#### URL params for /refbook/:name/item/:id
* **cols**= 
  * main - Default behavior. Returns reasonable columns.
  * all - Enrich each refbook item with system columns (createdDt, modifiedDt, etc)
  * colname1,colname2,colname3 - Returns specified columns.
* **enrich**=
  * n/a - No enrich. Default behavior.
  * nameonly - If refbookA refers to refBookB, column refbookB.name added into resultset. 
  * full     - If refbookA refers to refbookB, resultset item extends by refbookA {} object
* **enrichdepth**=
  * n/a, 0 - feature not used.
  * N - parent references level (1, 2, 3...).
* **lang**=
 * n/a - Returns refbook items in basic language
 * $$ - Returns refbook items in specified language (en, cz, ru, etc)

### Integration
Access control implemented by https://github.com/regorov/aaa 


## TODO
Everything :)

## License
MIT
