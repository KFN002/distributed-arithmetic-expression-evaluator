# _**distributed-arithmetic-expression-evaluator**_

---
## Задача: 

**_Пользователь хочет считать арифметические выражения. Он вводит строку 2 + 2 * 2 и хочет получить в ответ 6. Но наши операции сложения и умножения (также деления и вычитания) выполняются "очень-очень" долго. Поэтому вариант, при котором пользователь делает http-запрос и получает в качетсве ответа результат, невозможна. Более того: вычисление каждой такой операции в нашей "альтернативной реальности" занимает "гигантские" вычислительные мощности. Соответственно, каждое действие мы должны уметь выполнять отдельно и масштабировать эту систему можем добавлением вычислительных мощностей в нашу систему в виде новых "машин". Поэтому пользователь, присылая выражение, получает в ответ идентификатор выражения и может с какой-то периодичностью уточнять у сервера "не посчиталость ли выражение"? Если выражение наконец будет вычислено - то он получит результат. Помните, что некоторые части арфиметического выражения можно вычислять параллельно._**

---

## "Стартуем, я сказала стартуем!"

**Установка:** клонируйте репозиторий или скачайте .zip архив и распакуйте в любое удобное место

**Быстрый запуск:** Запустите файл `Calculator.exe`

**!Перед запуском убедитесь, что у вас последняя версия и порты 8080 и 8050 свободны!**

---

## Что теперь?

* **Куда идти потом?** [http://localhost:8080/](http://localhost:8080/)
* Навигация на страницах имеет говорящие названия:
    * `Create Expression` - создание выражения: **Пользовательский ввод арифметического выражения** -> _ID и статус задачи_.
    * `Expressions` - таблица со всеми выражениями из БД с колонками: _ID, Status, Expression, Result, Creation Date, Completion Date_.
    * `Expression by ID` - получение данных о задаче по ID: **Пользовательский ввод числа** -> _данные о задаче_. Если ID больше, чем есть задач, то будет ошибка (_failed to fetch an expression_).
    * `Edit Time` - Изменение времени выполнения операций: **Пользовательский ввод числа/чисел** -> изменение времени выполнения операций.
    * `Server Data` - Данные о воркерах (серверах/горутинах): _ID "сервера", статус, задание, которое выполняет, последний ответ на запрос о состоянии_.
    * `Project Scheme`- Схема проекта
    * `Log In` - Вход в аккаунт
    * `Sign Up` - Регистрация аккаунта
    * `Logout` - Выход из аккаунта


* (Запуск Ручками - main.go - все: cmd/app - в backend и cmd/server - в calculationServer) Если будет проблема с `gcc`, надо будет его установить и добавить в переменные окружения:
    * [StackOverflow](https://stackoverflow.com/questions/43580131/exec-gcc-executable-file-not-found-in-path-when-trying-go-build)
    * [Discourse GoHugo](https://discourse.gohugo.io/t/golang-newbie-keen-to-contribute/35087)
    * [GitHub Issue](https://github.com/golang/go/issues/47215)

---

## Как все устроено?

1. `backend` - весь бэк по модулям
2. `static` - весь фронт:
    * `assets` - .html templates и .css файлы.
3. `proto` - для gRPC
4. `calculationServer` - подсчет выражений
   
---

## Поподробнее о `backend`, но не слишком:

1. `dataManager` **_- internal_**
    * `librarian.go` - отвечает за работу с базой данных, то есть там находятся все функции для получения/изменения данных в базе.
    * `userManager.go` - работа с данными пользователя в бд
2. `orchestratorAndAgent`
    * `orchestratorAndAgent.go` - занимается получением ответа на выражение: получает выражение через очередь от оркестратора и распределяет его на мелкие таски, считает, а после собирает весь ответ и закидывает в базу данных.
3. `handlers`
    * `handlers.go` - обработка запросов пользователя: получение всех выражений, получение выражения по id, изменение времени работы операции и т.д.
4. `models` **_- pkg_**
    * `expression` - структура выражения.
    * `operations` - структура числовых операций: время их выполнения.
    * `serversData` - структура и пара функций для серверов (горутин) - их состояния, таски. В этом же файле можно масштабироваться (менять переменную - количество серверов (горутин)).
    * `stack` - стек с методами.
    * `JWT` - парсинг и создание токена
    * `userJWT` - работа с sessionStorage и JWT
    * `templateMessage` - структура и ее методы (для работы с шаблонизатором)
5. `cmd`
    * `main.go` **_(app)_** - запуск сервера, раньше там же был и файл handlers, но разумнее оказалось его закинуть в отдельный модуль.
6. `tests` - тесты, везде говорящие имена.
7. `utils` **_- internal_**
    * `utils.go` - вспомогательные функции, в основном проверка корректности выражения.
8. `cacheMaster` **_- internal_**
    * `cacheMaster.go` - кэш для хранения времени операций, чтобы часто не обращаться к базе данных. Удобно, так как время операций назначает пользователь (само не меняется).
9. `queueMaster` **_- internal_**
    * `queueMaster.go` - реализация очереди.
10. `calculator` **_- internal_**
    * `calculator.go` - расчет простых выражений - не используется после 1.02
    * `changeNotation.go` - изменение нотации выражения.

---

## Что делать и как работать?
* Видеоинструкция-экскурсия


https://github.com/KFN002/distributed-arithmetic-expression-evaluator/assets/119512897/d8048a68-eed3-49e9-87ac-0406d7336c5b


---

## Примеры для ввода в поле выражения:

1. _2+2_
2. _3+4_
3. _3*8_
4. _7/7_
5. _3-2_
6. _2-3_
7. _6/9_
8. _6/0_
9. _7**2_
10. _((2-1)))_
11. _2+2*2_
12. _(3-4/2)_
13. _17/(2-2)_
14. _11-4*(3+5)_
15. _81+12*(11+1)_

---
## Что-то непонятно или не работает - лучше звоните Солу!

### TG: @RR7B7


