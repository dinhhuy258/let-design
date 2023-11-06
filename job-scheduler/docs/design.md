# Design

## Database schema

| jobs                                          |
| --------------------------------------------- |
| id BIGINT                                     |
| user_id BIGINT                                |
| execute_at TIMESTAMP                          |
| message TEXT                                  |
| weight_factor FLOAT DEFAULT 1                 |
| status ENUM(scheduled, running, failed, cancelled, completed) |
| created_at TIMESTAMP                          |

## API

**API to create a job**

POST `/api/v1/jobs`

Payload

```json
{
    "message": "Hello world",
    "execute_at": "2023-11-11T01:05:00+07:00"
}
```

Response

```json
{
    "success": true,
    "data": {
        "job_id": 1,
        "execute_at": "2023-11-11T01:05:00+07:00"
    }
}
```

**API to cancel a job**

POST `/api/v1/jobs/{job_id}/cancel`

## Estimation

The scheduling system can handle ~5M new/due jobs per hour at peak.
The system should only reach this peak a few times a day (around 9:00AM or 6:00PM)
At normal, It's likely to run at 10% of its capacity (500000 new/due jobs per hours)

- Each day we will handle around: 5M * 2 (I assume that we reach peak 2 hour per day) + 500000 * 22 = ~20M new/ due job
- Each month we will handle around: 20M * 30 = 600M new/due jobs.
- Each year we will handle around: 20M * 365 = 7.2B new/due jobs.

Each entry will consume 
8 bytes (id) + 8 bytes (user_id) + 8 bytes (execute_at) + 
100 bytes (I assume each message have 100 bytes) + 
4 bytes (weight_factor) + 4 bytes (status) + 8 bytes (created_at) = 140 bytes (the actual data can be less or more)

- Data consumed each day: 140 bytes * (20M / 2) = 1.4 gigabytes
- Data consumed each month: 1.4 gigabytes * 30 = 42 gigabytes
- Data consumed each year: 42 gigabytes * 12 = ~500 gigabytes
- Data consumed 10 years: 500 gigabytes * 10 = ~5 terabytes 

## How to handle 5M due jobs per hour at peak

5M jobs per hour = 1388 jobs per seconds

