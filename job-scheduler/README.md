# Job Scheduler

Design distributed job scheduling system allows users to run scheduled actions. We support specifying timing in the following ways:

- **A certain period of time** It means waiting for X minutes, hours, or days before doing the job.
- **A certain day of the week** For example, if the user selects Monday and Wednesday, but the current day is Tuesday, the system will delay the job until Wednesday because it is the first available day of the week match.
- **A specific day of year** The system will wait until the date matches the user’s inputted date.

## Functional requirements

- Job Scheduler: Design a Job Scheduler as mentioned above.
- Monitoring: Define critical metrics to monitor the performance and latency of the system
- Scalable: The system can be auto-scaled in and out automatically and of course gracefully.
- Fair-share scheduling: As we maintain the same price for all users, we won’t expect a few `whales` eat up all the capacity. The system should also give priority to other small business owners with small audience list. To maximize our paid conversion rate, we should give the best experience to new signed-up users who’s testing the workflow feature with a short time delay.
- The user is able to change value or cancel timers on the fly. 

## Non-functional requirements

- Scalability: The scheduling system should be horizontally scalable, to be able to handle ~5M new/due timers per hour at peak. 
- Reliability: System must be recovered if the system fails or restarts. It should guarantee at-least-once delivery of every single job.
- Timing accuracy: It is required that jobs run within seconds of their scheduled time. The scheduler should have a p95 scheduling deviation below 10 seconds
