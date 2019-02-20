# fanssp


モックサーバを4台でスタートさせる

```
$ docker-compose build
$ docker-compose up -d --scale mock_dsp=4
```

curl http://10.100.100.20/req -X POST -H "Content-Type: application/json" -d '{"ssp_name": "hoge", "request_time": "yyyyMMdd-HHMMSS.ssss", "request_id": "sssssss", "app_id":
123}'