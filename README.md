# Subscription API

API allows:
- find out current exchange rate of USD in UAH
- subscribe email address to currency rate change dispatch

## Used technologies

- __[Go](https://go.dev/)__ as main language.
- __[gin](https://gin-gonic.com/docs/)__ as web framework.
- __[mockery](https://vektra.github.io/mockery/latest/)__ for generation of interface mocks.
- __[PostgreSQL](https://www.postgresql.org/)__ as main storage.
- __[gRPC](https://grpc.io/)__ for synchronous inter-service communication.
- __[NATS](https://docs.nats.io/)__ for asynchronous inter-service communication.
- __[ExchangeRate-API](https://www.exchangerate-api.com/)__, __[Currency Beacon API](https://currencybeacon.com/)__ and __[Free Currency API](https://github.com/fawazahmed0/exchange-api)__ as third-party APIs for currency rate info.

## Quick start

1. Copy example env file with command below:
```
cp ./.env.example ./.env
```
2. Update `SMTP_EMAIL` and `SMTP_USERNAME` values in `.env`.
3. Get app password from Google ([instruction](https://support.google.com/mail/answer/185833?hl=en)) and set to `SMTP_PASSWORD` value in `.env`.
4. Start project with command (uses docker) below:
```
make start
```


## System design
![system design](docs/system-design.jpg)

### Processes

1. ___Gateway___ is an entry point for external users of `SubscriptionAPI`. It is web server and\
makes requests via gRPC to services for required functionality. For now it allows \
to get exchange rate USD/UAH and subscribe user for daily dispatch of USD/UAH exchange rate.

2. ___Currency Service___ is responsible for fetching currency rates from third-party API. \
In perspective as business requirements grow it will cache rates in some store or \
calculate/aggregate currency-related data.

3. ___Dispatch Service___ is responsible for subscribing users to dispatches, \
sending of dispatches thorough SMTP server to subscribers, and getting info about dispatch.\
In perspective it would be able to create dispatches, change dispatches, \
customize subscriptions etc.

4. ___Customer Service___ is a pretty dumb service and it is responsible only for creation of customer.

5. ___Notification Service___ is responsible for monitoring of subscription-related events and dispatch scheduling. Except of dispatches, when user subscibes for a dispatch, it sends details of the subscription.

6. ___Transaction Manager___ is responsible for SAGA transactions.

7. _[not implemented yet]_ ___Rate Daemon___ is an automatic process that invokes \
updating of exchange rates.

#### More detailed documentation
1. [Gateway](./gateway)
2. [Currency Service](./service/currency)
3. [Dispatch Service](./service/dispatch)
4. [Notification Service](./service/notification)

### Possible alerts to create

| Microservice         | Trigger                                                        | Reason                                    | Severity |
| :------------------- | :------------------------------------------------------------- | :---------------------------------------- | :------- |
| every                | service is not alive for longer than 10min in a row            | availabitlity check                       | minor    |
| every                | percent of used CPU is higher than 70%                         | efficiency check                          | major    |
| gateway              | average response time is longer than 1500ms                    | performance check                         | major    |
| currency-service     | more than 35% of Convert method calls per 30min are failed     | third-parties are misbehaving             | minor    |
| notification-service | more than 10% of consumed messages weren't handled successfuly | check correctness of notification sending | major    |


## TODO
1. Implement rate daemon
2. Cache currency rates to decrease amount of calls to third-APIs
3. to be continued...
