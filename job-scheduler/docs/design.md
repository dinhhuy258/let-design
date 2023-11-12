# Design

## Functional Requirements

- Design one time scheduled task system
- Tasks should be excuted in fair share order.

## Non-Functional Requirements

**Scalability**

The scheduling system should be horizontally scalable, to be able to handle ~5M new/due timers per hour at peak.

**Durability**

A job should not be lost and must persist.

**Availability**

It should be possible all the time to schedule and perform tasks. A task should be executed a minimum number of times.

## High-Level Design

![](https://user-images.githubusercontent.com/17776979/282286336-ba5941a2-c3a3-48cd-9d48-36909eda2af1.png) 

**Job service**

It will handle requests to create jobs. The service will generate a unique ID for each job, which can be accomplished using [Twitter Snowflake](https://developer.twitter.com/ja/docs/basics/twitter-ids) or [ksuid](https://github.com/segmentio/ksuid)... Based on the `job_id`, the service will determine the database shard into which the task shall be placed.

**Job scheduler worker**

Each job scheduler has its own configuration database shard ids, and they will query and list all available tasks on the given shard ids. Tasks are scheduled in fair share mode, and the scheduler will publish tasks to the Kafka queue. We use a `Coordination Service}` (e.g., ZooKeeper or etcd) to keep track of the list of DB shards and assign shards to workers on startup. This ensures that each shard is assigned to only one worker.

**Job executor worker**

Each job executor will consume messages from a Kafka topic partition. We need to set up the Kafka topic with a sufficient number of partitions. When a worker finishes executing a job, it will update the execution status in the database. If it fails to handle the job, we can retry based on the given error. If it still fails, we can mark the job as failed and move it to the dead letter queue for investigating later.

**Database**

We will use `RDBMS` for ACID properties, especially transactions. We will shard the database several shards to distribute the data and load. We will use the master-slave replication for every partition in a semi-synchronous manner to balance between consistency and performance.

**Coordination Service**

We can use ZooKeeper or etcd as coordination service. It stores the shard ids, job scheduler workers, job executor worker information.

**Message Queues**

We can use `Kafka` as message broker for: 

- Independently scaling the nodes and producer/consumer.
- We isolate the producer and consumer from each other.
- If a consumer node crashes, some other node can process the message.
- Performance 
- Ordering of messages is done by Kafka.

## API

## Create user

POST `/api/v1/users`

**Payload**

```json
{
  "username": "username",
  "password": "password",
  "job_weight": 1.0
}
```

## Login

POST `/api/v1/auth/login`

**Payload**

```json
{
  "username": "username",
  "password": "password"
}
```

**Response**

```json
{
  "data": {
    "expire": "2023-11-12T13:29:26.925212+07:00",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTk3NzA1NjYsIm9yaWdfaWF0IjoxNjk5NzY2OTY2LCJ1c2VyX2lkIjoyfQ.aSkY9Os1lP4oOspaF5H433APXFoqFsC3htMN6jzF10I",
    "user": {
      "id": 2,
      "username": "username",
      "job_weight": 1
    }
  },
  "success": true
}
```

## Refresh token

POST `/api/v1/auth/refresh`

**Headers**

- Bearer token

**Response**

```json
{
  "data": {
    "expire": "2023-11-12T13:29:26.925212+07:00",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTk3NzA1NjYsIm9yaWdfaWF0IjoxNjk5NzY2OTY2LCJ1c2VyX2lkIjoyfQ.aSkY9Os1lP4oOspaF5H433APXFoqFsC3htMN6jzF10I"
  },
  "success": true
}
```

## Get jobs

GET `/api/v1/me/jobs`

**Headers**

- Bearer token

**Response**

```json
{
  "data": [
    {
      "id": 1,
      "user_id": 2,
      "message": "Hello world",
      "status": "running",
      "execute_at": "2023-11-11T23:08:37.086851Z",
      "shard_id": 1
    }
  ],
  "success": true
}
```

## Create job

POST `/api/v1/me/jobs`

**Headers**

- Bearer token

**Payload**

```json
{
  "message": "Hello world",
  "execute_at": "2023-11-11T23:08:37.086851+07:00"
}
```

**Response**

```json
{
  "data": {
    "id": 2,
    "user_id": 2,
    "message": "Hello world",
    "status": "created",
    "execute_at": "2023-11-11T23:08:37.086851+07:00",
    "shard_id": 1
  },
  "success": true
}
```

## Cancel job

POST `/api/v1/me/jobs/{job_id}/cancel`

**Headers**

- Bearer token

## GET job

GET `/api/v1/me/jobs/{job_id}`

**Headers**

- Bearer token

**Response**

```json
{
  "data": {
    "id": 1,
    "user_id": 2,
    "message": "Hello world",
    "status": "running",
    "execute_at": "2023-11-11T23:08:37.086851Z",
    "shard_id": 1
  },
  "success": true
}
```

## How to scale job-scheduler-worker

Scaling the `job-scheduler-worker` is not as straightforward as merely adding a new instance. For instance, consider two `job-scheduler-workers` named `worker1` and `worker2`. Both workers handle jobs concurrently, and without proper task assignment definitions, there's a risk of `worker1` and `worker2` processing the same task, negating the benefits of adding a new worker instance.

To solve this challenge, we must assign a `shard_id` to each created job. This `shard_id` can be generated by creating a lexicographically unique `job_id`. With the `job_id`, we can calculate the `shard_id` for the given job. One implementation option is to utilize database sharding. In this demonstration, I will include the `shard_id` column in the jobs table.

Upon starting up, each `job-scheduler-worker` will be assigned a `shard_id` to handle jobs. Similarly, when a `job-scheduler-worker` shuts down, we need to define a mechanism to reassign shard_ids for all running workers.

**Note**: Scaling the `job-executor-worker` is quite easier compared to scaling the `job-scheduler-worker` when using Kafka. Here, what we need to do is define the number of partitions for the Kafka topic, perhaps 10 if we estimate that, in the future, there will be at most 10 `job-executor-workers` added to the system. All `job-executor-workers` should belong to a single consumer group. If a new `job-executor-worker` is added or removed, Kafka's architecture will assist in reassigning partitions to the new or existing workers.

## Database schema

![](https://user-images.githubusercontent.com/17776979/282299391-5980379c-9fe9-4d23-9d02-5a0b2d2a17e6.png)

## Estimation

The requirements of the system that it can handle ~5M new/due jobs per hour at peak.

I assume that each day there are 5M new tasks are created.

Each entry will consume: 
8 bytes (id) +
8 bytes (user_id) + 
8 bytes (execute_at) +
100 bytes (message) +
4 bytes (status) + 
8 bytes (created_at) + 
8 bytes (updated_at) 
= 144 bytes (the actual data can be less or more)

- Data consumed each day: 144 bytes \* 5M = 720 megabytes
- Data consumed each month: 720 megabytes \* 30 = 21.6 gigabytes
- Data consumed each year: 21.6 gigabytes \* 12 = 259.2 gigabytes
- Data consumed 10 years: 259.2 gigabytes \* 10 = 2.59200 terabytes

## Monitors

- Define dashboards to keep track of CPU and RAM usage for the `job-service`, `job-scheduler-worker`, and `job-executor-worker`.
- Define a dashboard to monitor the number of requests per second for the `job-service`.
- Define a dashboard to monitor p99, p95 latency for the `job-service`
- Define a dashboard to track p99, p95 job latency. Job latency is calculated as the time it takes for a job to finish execution minus the time the job is scheduled.
- Set up an alert if the job latency is greater than or equal to 10 seconds.
- Define a dashboard to monitor the number of processed jobs per second.
- Define a dashboard to monitor the number of new tasks added per second.
- Set up an alert for jobs that fail to handle, indicating when a job is moved to the dead letter queue. 
- Establish an alert for a 500 Internal Server Error status code in the `job-service`.
