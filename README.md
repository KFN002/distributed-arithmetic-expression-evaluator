# distributed-arithmetic-expression-evaluator

---
## Задача: 

Пользователь хочет считать арифметические выражения. Он вводит строку 2 + 2 * 2 и хочет получить в ответ 6. Но наши операции сложения и умножения (также деления и вычитания) выполняются "очень-очень" долго. Поэтому вариант, при котором пользователь делает http-запрос и получает в качетсве ответа результат, невозможна. Более того: вычисление каждой такой операции в нашей "альтернативной реальности" занимает "гигантские" вычислительные мощности. Соответственно, каждое действие мы должны уметь выполнять отдельно и масштабировать эту систему можем добавлением вычислительных мощностей в нашу систему в виде новых "машин". Поэтому пользователь, присылая выражение, получает в ответ идентификатор выражения и может с какой-то периодичностью уточнять у сервера "не посчиталость ли выражение"? Если выражение наконец будет вычислено - то он получит результат. Помните, что некоторые части арфиметического выражения можно вычислять параллельно.

---

## "Стартуем, я сказала стартуем!"

**Быстрый запуск:** Запустите файл `backend/cmd/app/main.go` - все должно работать

**!Перед запуском убедитесь, что у вас последняя версия!**

---

## Что теперь?

* **Куда идти потом?** [http://localhost:8080/](http://localhost:8080/)
* В целом можно не тестить, все работает, да и я постарался)
* Если будет проблема с `gcc`, надо будет его установить и добавить его в переменные окружения:
    * [StackOverflow](https://stackoverflow.com/questions/43580131/exec-gcc-executable-file-not-found-in-path-when-trying-go-build)
    * [Discourse GoHugo](https://discourse.gohugo.io/t/golang-newbie-keen-to-contribute/35087)
    * [GitHub Issue](https://github.com/golang/go/issues/47215)
* Ссылки сверху точно должны помочь в фиксе проблемы.

---

## Как все устроено?

1. `backend` - весь бэк по пакетам.
2. `static` - весь фронт:
    * `assets` - .html templates и .css файлы.
    * `Scripts` - пустой .js файл.
   
---

## Поподробнее, но не слишком:

1. `dataManager`
    * `librarian.go` - отвечает за работу с базой данных, то есть там находятся все функции для получения/изменения данных в базе.
2. `orchestratorAndAgent`
    * `orchestratorAndAgent.go` - занимается получением ответа на выражение: получает выражение через очередь от оркестратора и распределяет его на мелкие таски, считает, а после собирает весь ответ и закидывает в базу данных.
3. `handlers`
    * `handlers.go` - обработка запросов пользователя: получение всех выражений, получение выражения по id, изменение времени работы операции и т.д.
4. `models`
    * `expression` - структура выражения.
    * `operations` - структура числовых операций: время их выполнения.
    * `serversData` - структура и пара функций для серверов (горутин) - их состояния, тасок.
    * `stack` - стек с методами.
5. `cmd`
    * `main.go` - запуск сервера, раньше там же был и файл handlers, но разумнее оказалось его закинуть в отдельный модуль.
6. `tests`
    * `test.go` - тесты, созданные для первых версий сервера, скорее всего файл будет пустым.
7. `utils`
    * `utils.go` - вспомогательные функции, в основном проверка корректности выражения.
8. `cacheMaster`
    * `cacheMaster.go` - кэш для хранения времени операций, чтобы часто не обращаться к базе данных. Удобно, так как время операций назначает пользователь (само не меняется).
9. `queueMaster`
    * `queueMaster.go` - реализация очереди.

---

## Схема работы - на сайте

---

## TG: @RR7B7
