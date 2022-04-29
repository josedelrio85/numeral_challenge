# Numeral code challenge

Proposed solved exercise

## How to run it

* Export the needed env vars

```sh
export REQUEST_SCHEMA=resources/request_schema.json export BANK_FOLDER=bank export SQLITE_DB_FILE_LOCATION=database.db
```

* Run the program

```sh
go run main.go
```

* Send a request to ingest a payment and follow the stated logic

```sh
curl --request POST 'http://localhost:4567/payment/receive' \
-u "usertest:userpass" \
-d '{
  "debtor_iban": "FR1112739000504482744411A64",
  "debtor_name": "company1",
  "creditor_iban": "DE65500105179799248552",
  "creditor_name": "beneficiary",
  "ammount": 42.99,
  "idempotency_unique_key": "JXJ984XXXZ"
}'
```

* The response will be something like that:

```json
{
  "status":200,
  "success":true,
  "data":{
    "id":"JXJ984XXXZ"
  }
}
```

* You can check the database and see that a new payment was added, and after 5 seconds it will be updated with the response from `resources/bank_response.csv` file.

## How to test it

```sh
go test ./... -cover
```


## Considerations

* It took me a little more than 3 hours because I had never worked with XSD before.

* Because I did not know how to correctly use the XSD file and generate the XML to be stored in the `bank` folder, I decided to create a template and further process it with the corresponding values.

* The structure of the `payments` table is super simple and does not reflect a real case. I think a unique primary key must be created in the `idempotency_unique_key` field, at the very least.

* I have not had enough time to generate all the necessary tests. I include the ones I have been able to create in the time available.

* I have used dummy values for the basic authorisation, obviously these values would have to be obtained through environment variables.

* One of the libraries (https://github.com/xeipuuv/gojsonschema) I have used to validate the received json needs the absolute reference of the file with the schema to be added.