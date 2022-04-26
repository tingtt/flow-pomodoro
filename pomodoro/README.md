# flow-pomodoro

## Usage

### With `docker-compose`

#### Variables `.env`

| Name                    | Description                                                              | Default       | Required           |
| ----------------------- | ------------------------------------------------------------------------ | ------------- | ------------------ |
| `PORT`                  | Published port                                                           | 1323          |                    |
| `MYSQL_DATABASE`        | MySQL database name                                                      | flow-pomodoro |                    |
| `MYSQL_USER`            | MySQL user name                                                          | flow-pomodoro |                    |
| `MYSQL_PASSWORD`        | MySQL password                                                           |               | :heavy_check_mark: |
| `MYSQL_ROOT_PASSWORD`   | MySQL root user password                                                 |               |                    |
| `LOG_LEVEL`             | API log level                                                            | 2             |                    |
| `GZIP_LEVEL`            | API Gzip level                                                           | 6             |                    |
| `MYSQL_HOST`            | MySQL host                                                               | db            |                    |
| `MYSQL_PORT`            | MySQL port                                                               | 3306          |                    |
| `JWT_ISSUER`            | JWT issuer                                                               | flow-pomodoro |                    |
| `JWT_SECRET`            | JWT secret                                                               |               | :heavy_check_mark: |
| `SERVICE_URL_PROJECTS`  | The url to [flow-projects](https://gitlab.tingtt.jp/flow/flow-projects). |               | :heavy_check_mark: |
| `SERVICE_URL_TODOS`     | The url to [flow-todos](https://gitlab.tingtt.jp/flow/flow-todos).       |               | :heavy_check_mark: |

```bash
$ docker-compose up
```