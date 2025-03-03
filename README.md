# Сервис-балансировщик видео-трафика

## Описание

Cервис принимает gRPC вызовы с параметром video и согласно алгоритму
балансировки, возвращает в ответ URL для редиректа:

- Каждый 10-й вызов возвращат оригинальный URL.

- Остальные вызовы перенаправляются на CDN, где URL формируется как:

  http://{CDN_HOST}/{subdomain}{original_path} Например, для входного URL:

  http://s1.origin-cluster/video/123/xcg2djHckad.m3u8 и CDN_HOST=cdn.example.com
  будет возвращён:

  http://cdn.example.com/s1/video/123/xcg2djHckad.m3u8

Результаты нагрузочного тестирования:

Summary: Count: 9911 Total: 998.86 ms Slowest: 13.41 ms Fastest: 0.28 ms
Average: 1.15 ms Requests/sec: 9922.32

Response time histogram: 0.285 [1] | 1.598 [8777]
|∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎ 2.911 [492] |∎∎ 4.224 [189] |∎ 5.537
[140] |∎ 6.850 [49] | 8.163 [62] | 9.475 [68] | 10.788 [74] | 12.101 [38] |
13.414 [18] |

Latency distribution: 10 % in 0.44 ms 25 % in 0.52 ms 50 % in 0.67 ms 75 % in
1.02 ms 90 % in 1.75 ms 95 % in 3.97 ms 99 % in 10.05 ms

Status code distribution: [OK] 9908 responses  
 [Canceled] 3 responses

Error distribution: [3] rpc error: code = Canceled desc = grpc: the client
connection is closing

Для продакшн:

- разнесенеие сервиса на логические части в структуре проекта

- нужно добавить логирование

- сбор метрик и мониторинга

- безопасность и аутентификация если нужна, в идеале какая-нибудь защита от ddos

- расширить конфиги

- тесты

- рассмотреть перспективу масшатибрования, для нагрузок

p.s. гитхаб пустой, так как новый для github copilot
