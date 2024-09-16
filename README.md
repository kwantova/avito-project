## Описание выполнения тестового задания

# 1. Структура проекта

Все исполняющие программы (.go) находятся в директории package. Внутри нее есть три подкаталога:

__/api__ — все API-эндпоинты. Разнесены по отдельным файлам по смыслу:

offers.go и tenders.go: действия над предложениями и тендерами соотвтетственно.

status-def-t.go и status-def-o.go: логика редактирования статуса тендеров и предложений.

ping.go — пинг сервера.

__/db__ — описание подключения к базе данных. Там описана функция Connect() с проверкой ошибки, и она используется в эндпоинтах.

__/entitydescripts__ — описание структур данных Tender и Offer в файле struct-def.go и структур User и Organization в файле user-org-def.go

__main.go__ — точка входа в проект. Эндпоинтам задаются маршруты с помощью HandleFunc(), запускается прослушивание по порту 8080.

__go.mod__ — файл зависимостей. Указывает пакеты, которые используются в проекте. Каждый раз, когда добавляется/обновляется зависимость, там появляются криптографические хэши, соответствующие зависимости.

__go.sum__ —  хранит контрольные суммы (хэши) всех модулей и их версий, которые используются в проекте. Эти хэши проверяются Go при скачивании и установке зависимостей, чтобы убедиться, что они не были изменены или подделаны. 

# 2. Обоснование выбора библиотек

Для работы с HTTP-серверами выбрана стандартная (встроенная) библиотека __net/http__, т.к. ее функционала было достаточно для реализации задания, да и проект становится более стабильным, если отдавать предпочтение встроенным пакетам.

Для работы с базой данных был выбран __pgx__ - это драйвер стандартной библиотеки database/sql, наточенный на работу конкретно с PostgreSQL. Мне кажется, что он самый понятный, и работа с JSON'ами с ним удобна.

Еще добавлена библиотека __UUID__ для генерации уникальных идентификаторов. При создании таблиц tender и offer в базе данных нужно было учесть, что есть поля, которые являются внешними ключами по отношению к другим полям уже созданных заранее таблиц, которые имеют тип UUID. А чтобы сделать внешний ключ, данные должны иметь одинаковый тип.

# 3. Описание yaml-файлов

YAML — формат файлов для сериализации данных и для конфигурационных файлов. Сериализация - это процесс преобразования данных в такой формат, который удобно передавать по сети. Конфигурационные файлы - файлы, которые содержат параметры и настройки для ПО или системы.

__docker-compose.yml__ — файл с переменными окружения, который я использовала для тестирования запросов на localhost:8080.

__configmap.yaml__ — содержит конфигурационные данные для контейнера в Kubernetes (по большей части, те же переменные окружения)

__secret.yaml__ — содержит чувствительные данные (пароль, хост БД), зашифрованные в base64, чтобы не хранить в открытом виде то, что должно оставаться конфиденциальным.

__deployment.yaml__ — файл, управляющий развертыванием в Kubernetes. Описывает, где и как будет развернут проект.

__ci-cd.yaml__ — добавлен после пуша на GitHub и выполняет CI/CD сборку.

___

Dockerfile — для создания образа в Docker. Там указано, какая версия go используется, для какой архитектуры процессора создается образ и прочее.

# 4. Описание изменений, введенных в структуру БД

В базу данных добавлены две сущности - tender для тендеров и offer для предложений. Вот DDL-запросы, использованные при создании (использовался терминал psql):

CREATE TABLE Tender (
    TenderID UUID PRIMARY KEY,               
    Name VARCHAR(255) NOT NULL, \
    ServiceType VARCHAR(255) NOT NULL, \
    Description TEXT, \
    Status VARCHAR(50) NOT NULL, \
    OrganizationID UUID NOT NULL, \
    AuthorID UUID NOT NULL, \
    Version INTEGER DEFAULT 1, \
    FOREIGN KEY (OrganizationID) REFERENCES organization(id), \
    FOREIGN KEY (AuthorID) REFERENCES employee(id) \
);

CREATE TABLE Offer (
    OfferID UUID PRIMARY KEY, \
    Name VARCHAR(255) NOT NULL, \
    ServiceType VARCHAR(255) NOT NULL, \
    Description TEXT, \
    Status VARCHAR(50) NOT NULL, \
    OrganizationID UUID NOT NULL, \
    TenderID UUID NOT NULL, \
    AuthorID UUID NOT NULL, \
    Version INTEGER DEFAULT 1,\
    FOREIGN KEY (OrganizationID) REFERENCES Organization\
    FOREIGN KEY (TenderID) REFERENCES Tender(TenderID), \
    FOREIGN KEY (AuthorID) REFERENCES Employee(EmployeeID) \
);

_Примечание: Термины offer и bid используются в проекте без какого-либо смыслового различия._

# 5. Как пользоваться проектом

Проект, развернутый с помощью gitlab codenrock, доступен по URL:
https://cnrprod1725741033-team-78604-33124.avito2024.codenrock.com


На момент пуша коммита он все еще доступен и пингуется.

Текущий проект успешно билдится, деплоится и доступен на Kubernetes c небольшой оговоркой: поскольку у меня нет прав доступа на чтение api-secret, я временно убрала secret.yaml и поместила чувствительные данные в config.yaml без шифрования в base64. Также пришлось не подтягивать переменные, связанные с Docker Hub из GitHub'a, а указывать напрямую в deployment.yaml и ci-cd.yaml. В таком случае ошибки сборки CI/CD и ошибок в логах пода не было. Но это небезопасно, поэтому в текущей версии проекта secret.yaml все еще присутствует, и из-за отсутствия прав деплой выдает следующую ошибку:

Error from server (Forbidden): error when retrieving current configuration of:
Resource: "/v1, Resource=secrets", GroupVersionKind: "/v1, Kind=Secret"
Name: "app-secrets", Namespace: "cnrprod1725741033-team-78604"
from server for: "secret.yaml": secrets "app-secrets" is forbidden: User "system:serviceaccount:cnrprod1725741033-team-78604:cnrprod1725741033-team-78604" cannot get resource "secrets" in API group "" in the namespace "cnrprod1725741033-team-78604"
Error: Process completed with exit code 1.

Я пришла к решению, что лучше оставлю сборку на этом моменте, чтобы не хранить конфиденциальные данные в открытом доступе. По логам сборки можете посмотреть, что в предыдущих коммитах и билд, и деплой собирался. 

__Если Вы будете тестировать сборку и ко мне возникнут вопросы, напишите мне на почту или в tg: @kwanto__

К тому же, пришлось бы поменять тип сервиса с ClusterIP на другой, позволяющий публичный доступ, а не только из других подов внутри кластера. Я не знаю, можно ли мне так делать, поэтому решила не рисковать. 

Дальнейшая информация будет описана для URL, полученного при сборке на GitLab Codenrock. Тестирование запросов я проводила с помощью Postman. Во имя краткости изложения будут приведены примеры только для некоторых эндпоинтов. Это не означает, что другие не работают.

## 5.1. Пример отправки пинг-запроса

Отправлено:
GET https://cnrprod1725741033-team-78604-33124.avito2024.codenrock.com/api/ping

Результат:
ok\
200OK 21 ms 206 B

## 5.2 Пример отправки запроса на создание тендера

Отправлено:
GET https://cnrprod1725741033-team-78604-33124.avito2024.codenrock.com/api/tenders/new

C JSON'ом: \
{ \
    "name": "i like kittens and avito", \
    "description": "This is an another one example tender", \
    "service_type": "Consulting", \
    "status": "CREATED", \
    "organization_id": \"550e8400-e29b-41d4-a716-446655440022",\
    "author_id": "550e8400-e29b-41d4-a716-446655440007",\
    "version": 1\
}

Результат:

{ \
    "tender_id": "f7da1376-1dd5-48ac-9be5-6f1925289bd6",\
    "name": "i like kittens and avito",\
    "service_type": "Consulting",\
    "description": "This is an another one example tender",\
    "status": "CREATED",\
    "organization_id": \"550e8400-e29b-41d4-a716-446655440022",\
    "author_id": "550e8400-e29b-41d4-a716-446655440007",\
    "version": 1\
}

200OK 92 ms 505 B

## 5.3. Пример отправки запроса на вывод списка тендеров

Отправлено:
GET https://cnrprod1725741033-team-78604-33124.avito2024.codenrock.com/api/tenders

Результат:

[\
    {\
        "tender_id": "f7da1376-1dd5-48ac-9be5-6f1925289bd6",\
        "name": "i like kittens and avito",\
        "service_type": "This is an another one example tender",\
        "description": "Consulting",\
        "status": "CREATED",\
        "organization_id": \"550e8400-e29b-41d4-a716-446655440022",\
        "author_id": "550e8400-e29b-41d4-a716-446655440007",\
        "version": 1\
    }\
]\

200OK 44 ms 507 B

## 5.4 Пример создания предложения

Отправлено: POST https://cnrprod1725741033-team-78604-33124.avito2024.codenrock.com/api/bids/new

C JSON'ом:

{\
  "name": "New Bid Name",\
  "description": "Description of the new bid",\
  "status": "CREATED",\
  "tender_id": "f7da1376-1dd5-48ac-9be5-6f1925289bd6",\
  "organization_id": "550e8400-e29b-41d4-a716-446655440020",\
  "author_id": "550e8400-e29b-41d4-a716-446655440001"\
}

Результат:

{
    "offer_id": "b465dc0e-cc45-4dd7-85a9-900d60c3152c",\
    "name": "New Bid Name",\
    "service_type": "",\
    "description": "Description of the new bid",\
    "status": "CREATED",\
    "organization_id": "550e8400-e29b-41d4-a716-446655440020",\
    "tender_id": "f7da1376-1dd5-48ac-9be5-6f1925289bd6",\
    "author_id": "550e8400-e29b-41d4-a716-446655440001",\
    "version": 0\
}

200OK 114 ms 522 B

# 6. Про локальную сборку проекта

Проект можно запустить на localhost:8080. Для этого я использовала Docker Desktop и командную строку.

В директории с проектом открываем cmd и вводим:

docker build -t avitoproject:latest .

Это создает образ (image). Точка в конце указывает, что Dockerfile находится в текущем каталоге.

Дальше нужно поднять docker-compose, т.е. сделать видимыми для Docker переменные окружения, откуда он берет информацию о том, на какой порт пробрасываем, как подключиться к БД и прочее.

docker-compose up

Эта команда поднимет все контейнеры, указанные в docker-compose.yaml, и обеспечит их связь между собой.

Далее можно обращаться к локальному серверу с помощью curl в командной строке или через Postman. Например, так:

GET http://127.0.0.1:8080/api/ping


# 7. Послесловие
Огромное спасибо Авито за предоставленную возможность погрузиться в реальный (пусть и очень упрощенный) проект и применить теоретические знания на практике. Надеюсь, у меня получилось не только освежить собственную память в теме бэкенда, HTTP API, Docker'a, но и продемонстрировать свою способность адаптироваться к реальным условиям. Это ценный урок :)

Если потребуется связаться по поводу сборки:\
kwantova@gmail.com
tg: @kwanto




