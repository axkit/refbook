# refbook
Easy access to your database reference tables. Standalone or embedded.

## Status
Under development


## JSON API

Content-Type: application/json


Method|URL| Description
----|--------------|----------------------------------------------------------------------
GET |/refbook/| Returns list of reference books available for user
GET |/refbook/:name|Returns reference book metadata and list of available methods       
GET |/refbook/:name/stat|Returns reference book statistics
GET |/refbook/:name/row/    |Returns reference book items 
GET |/refbook/:name/row/:id |Returns single reference book item
GET |/refbook/:name/row/:id/exist|Returns HTTP OK 200 if row exists, 404 if not
POST|/refbook/:name/row/|Add item
PUT |/refbook/:name/row/:id|Update item
DELETE| /refbook/:name/row/:id|Mark as deleted
POST|/refbook/:name/row/:id/erase|Delete rows from database table


## Integration
Access control implemented by https://github.com/regorov/aaa 

