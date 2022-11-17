# Задание
Необходимо реализовать микросервис для работы с балансом пользователей (зачисление средств, списание средств, резервирование денег при заказе услуги, перевод резерва на счет компании, а также метод получения баланса пользователя). Сервис должен предоставлять HTTP API и принимать/отдавать запросы/ответы в формате JSON.

## Запуск проекта
Клонируйте репозиторий с помощью git
```git clone https://github.com/arrrrtur/Task.git```

Создайте и запустите докер контейнер коммандой ```docker-compose up --build balance```
Запустите файл `data.sql```, который создаст все необходимые таблицы.


## Структура бд


balance(Содержит id, доступные и зарезервированные средства)  
report(информация об уже предоставленных услугах(для бугалтерии))  
service(id, название и цена услуги)  
order(хранит информацию о заказах)  
operation_type(тип операции операции)  
transaction(хранит информацию обо все транзакциях)  


## Запросы
Я использовал postman

### Метод начисления средств на баланс. Принимает id пользователя и сколько средств зачислить (в случае отсутствия пользователя - он создается). Метод - `POST`
``` http://localhost:8080/balance/top-up-balance ```

Добавить в тело запроса(JSON формат):

```
{
    "balance_id": "5",
    "amount": "5000"
}
```
 
 Добавить в тело запроса(JSON формат):
 Ответ запроса(JSON формат):
 
 ```
 [
    {
        "transaction_id": 0,
        "sender_id": 5,
        "receiver_id": 5,
        "transaction_time": "",
        "transaction_price": 5000,
        "operation_id": 1,
        "status": 3
    },
    {
        "balance_id": 5,
        "amount_on_balance": 20000,
        "amount_on_reserve": 0
    }
]
 ```
 
 
 
 
### Метод резервирования средств с основного баланса на отдельном счете. Принимает id пользователя, ИД услуги, ИД заказа, стоимость. Метод - `PATCH`

```http://localhost:8080/balance/reserve-from-balance```

 
Добавить в тело запроса(JSON формат):

```
{
"balance_id": "2",
"service_id": "5",
"order_id": "5",
"amount": "1100"
}
```

Ответ запроса(JSON формат):
 
 ```
 [
    {
        "transaction_id": 0,
        "sender_id": 2,
        "receiver_id": 2,
        "transaction_time": "",
        "transaction_price": 1100,
        "operation_id": 3,
        "status": 2
    },
    {
        "order_id": 5,
        "balance_Id": 2,
        "service_Id": 5,
        "reserve_time": ""
    },
    {
        "service_id": 5,
        "name": "Какая-то услуга",
        "price": 1100
    },
    {
        "balance_id": 2,
        "amount_on_balance": 4900,
        "amount_on_reserve": 1100
    }
]
 ```


 

### Метод признания выручки – списывает из резерва деньги, добавляет данные в отчет для бухгалтерии. Принимает id пользователя, ИД услуги, ИД заказа, сумму. Метод - `PATCH`

``` http://localhost:8080/balance/has-passed ```

Добавить в тело запроса(JSON формат):

```
{
"balance_id": "2",
"service_id": "5",
"order_id": "5",
"amount": "1100"
}
```

Ответ запроса(JSON формат):

```
[
    {
        "transaction_id": 0,
        "sender_id": 2,
        "receiver_id": 2,
        "transaction_time": "",
        "transaction_price": 1100,
        "operation_id": 3,
        "status": 3
    },
    {
        "balance_id": 2,
        "amount_on_balance": 4900,
        "amount_on_reserve": 0
    }
]
```
 
 
 
 
### Метод получения баланса пользователя. Принимает id пользователя. Метод - `GET`

``` http://localhost:8080/balance/get-balance ```

Добавить в тело запроса(JSON формат):

```
{
    "balance_id":"2"
}
```

Ответ запроса(JSON формат):

```
[
    {
        "balance_id": 2,
        "amount_on_balance": 4900,
        "amount_on_reserve": 0
    }
]
```
 
### Метод разрезервирования денег, если услугу применить не удалось(то есть была отменена). Принимает id пользователя, id заказа, сумму. Метод - `PATCH`

``` http://localhost:8080/balance/cancel-order ```

Добавить в тело запроса(JSON формат):

```
{
    "balance_id": "2",
    "order_id": "6",
    "amount": "500"
}
```

Ответ запроса(JSON формат):

```
[
    {
        "transaction_id": 0,
        "sender_id": 2,
        "receiver_id": 2,
        "transaction_time": "",
        "transaction_price": 500,
        "operation_id": 3,
        "status": 1
    },
    {
        "balance_id": 2,
        "amount_on_balance": 4900,
        "amount_on_reserve": 0
    }
]
```



## Дополнительные задания
### Метод для получения месячного отчета. Принимает год-месяц. Возвращает ссылку на файл `.csv` Метод - `GET`. 

``` http://localhost:8080/report/get-link-report ```

Добавить в тело запроса(JSON формат):

```
{
    "year-month":"2022-11"
}
```

Ответ запроса(JSON формат):

```
[
    "{\"report_link\": \"http://localhost:8080/file/report2022-11.csv\"}"
]
```



### Метод получения списка транзакций. Полуает id, пользователя, порядок сортировки по дате, порядок сортировки по сумме, количество записей на странице, номер страницы. Предусмотрены параметры по умолчанию. Метод - `GET`.

``` http://localhost:8080/balance/get-history ```

Добавить в тело запроса(JSON формат):

```
{
    "balance_id": "100",
    "sort_by_date_order": "DESC",
    "sort_by_amount_order": "DESC",
    "per_page": "5",
    "page": "3"
}
```

Ответ запроса(JSON формат):

```
[
    [
        {
            "transaction_id": 34,
            "sender_id": 100,
            "receiver_id": 100,
            "transaction_time": "2022-11-13 05:21:04.994162 +0300 MSK",
            "transaction_price": 3000,
            "operation_type": "Оплата услуги",
            "status": 3
        },
        {
            "transaction_id": 33,
            "sender_id": 100,
            "receiver_id": 100,
            "transaction_time": "2022-11-13 05:21:01.133444 +0300 MSK",
            "transaction_price": 3000,
            "operation_type": "Оплата услуги",
            "status": 3
        },
        {
            "transaction_id": 32,
            "sender_id": 100,
            "receiver_id": 100,
            "transaction_time": "2022-11-13 05:20:54.751258 +0300 MSK",
            "transaction_price": 3000,
            "operation_type": "Оплата услуги",
            "status": 3
        },
        {
            "transaction_id": 31,
            "sender_id": 100,
            "receiver_id": 100,
            "transaction_time": "2022-11-13 05:20:50.202971 +0300 MSK",
            "transaction_price": 3000,
            "operation_type": "Оплата услуги",
            "status": 3
        },
        {
            "transaction_id": 30,
            "sender_id": 100,
            "receiver_id": 100,
            "transaction_time": "2022-11-13 05:20:43.977516 +0300 MSK",
            "transaction_price": 3000,
            "operation_type": "Оплата услуги",
            "status": 3
        }
    ]
]
```

