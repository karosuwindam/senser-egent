# この資料について

C105用にトレース処理を取り除いた資料


以下コマンドで、Prometheusが取得した最後のデータを取得可能
curl "localhost:9090/api/v1/query?query=senser_press_value"
curl "localhost:9090/api/v1/query?query=senser_tmp_value"
curl "localhost:9090/api/v1/query?query=senser_hum_value"