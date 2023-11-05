# High level design

## Database schema

| jobs                                          |
| --------------------------------------------- |
| id BIGINT                                     |
| user_id BIGINT                                |
| execute_at TIMESTAMP                          |
| message TEXT                                  |
| weight_factor FLOAT DEFAULT 1                 |
| status ENUM(scheduling, failed, cancel, done) |
| created_at TIMESTAMP                          |

## API

**API to create a job**

`POST /api/v1/jobs`

Payload:

```json
{
    [TBD]
}
```

**API to cancel a job**

`POST /api/v1/jobs/{job_id}/cancel`

## Estimation

The scheduling system can handle ~5M new/due jobs per hour at peak = ~2.5M new/due jobs per hour normally = ~2.5M / 60 / 60 = 694 tasks per second.
Each day we will handle around: 2.5M _ 24 = 60M new/due jobs.
Each month we will handle around: 2.5M _ 24 _ 30 = 1.8B new/due jobs.
Each year we will handle around: 1.8B _ 365 = 460.8B new/due jobs.

Each entry will consume 8 bytes (id) + 8 bytes (user_id) + 8 bytes (execute_at) + 255 bytes (I assume each message have 255 bytes) + 10 bytes (status) + 8 bytes (created_at) = 297 bytes (the actual data can be less or more)

- Data consumed each day: 297 bytes \* (60M / 2) = 8.91 gigabytes
- Data consumed each month: 8.91 gigabytes \* 30 = 267.3 gigabytes
- Data consumed each year: 267.3 gigabytes \* 12 = 3.2076 terabytes
- Data consumed each year: 3.2076 terabytes \* 12 = 38.4912 terabytes

It's huge data. A single database instance can not handle such amount of data -> we need to apply **SHARDING**
