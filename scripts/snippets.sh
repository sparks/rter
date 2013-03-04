curl -i -X POST -d '{"Type":"web", "AuthorID":1}' http://localhost:8080/items/
curl -i -X POST -d '{"Type":"somethingelse", "AuthorID":1}' http://localhost:8080/items/1
curl -i -X GET http://localhost:8080/items/
curl -i -X GET http://localhost:8080/items/1
curl -i -X DELETE http://localhost:8080/items/1