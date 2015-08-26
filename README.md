# refbook
Easy access to your database reference tables. Standalone or embedded. 

## Status
Planning & Requirement specification

## Concepts
#### Reference Book
Reference book is a plain table usually used as lookup. Reference table can refer to  another reference table (M:1, 1:M). As instance: Towns->Countries->Currencies. 
I recommend to count table as reference book when:
* table has limited amount of rows
* mostly, new rows are introduced by human
* table does not refer to tables with business objects
* table requires only simple CRUD operations
* table's text columns requires multilanuage support


#### RDBMS supported
PostgreSQL 9.x first. 

#### Deployment 
##### Standalone
Standalone http server. Listens tcp port and processing requests. Recommended cases:
* Page rendering on client side
* Reference book items mostly stored in browser local storage
* Standalone server side HTML page renders

##### Embedded
Statical linking into your Go application
* Can work on separate tcp port (requires reverse proxy in front)
* Can add it's own routes into your existing http router and handle requests on the same port
* Refbook features are directly available (without network overhead)
 
## Features
### JSON API

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
* **search**=
  * text - Filters resultset by text
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
* **result**
 *  n/a - Returns JSON. Default behavoir. Content-Type: application/json. 
 *  options - Set of option items. Content-Type: text/html
 *  html - Simple html page, mostly for debug purpose. Content-Type: text/html

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

### Configuration
* Enable permission control (integration with aaa)
* Run it's own http listener, ip:port
* Template file for HTML page composing
* Feature enabler. As instance disable support "?search=" or "enrich="
* Database connection string if standalone deployment
* Reference table names

### Integration
Access control implemented by https://github.com/regorov/aaa 

## TODO
Everything :)

## License
MIT
