# RTER System Design - The road ahead

*by Alexander Eichhorn (echa@cim.mcgill.ca, echa@kidtsunami.com)*

The early RTER prototype is about to evolve into a large-scale Internet services that people are supposed to rely on. It is important to understand that RTER lives in a larger ecosystem of end users, emergency operation centers, legal constraints, mobile devices, web standards, server and networking technologies.

More important than the functional fit to its purpose is to design the entire architecture for robustness, scalability and security from ground up. Scalability, robustness and security don't come as libraries or 3rd party services you can simply plug-in. They are a philosphy behind each single line of code.

Developing an Internet-scale system is not done in a single shot. It rather is a continuous process of failing, understanding, learning and improving. Eventually this will lead to a useful and secure application if executed with proper scrunity and management oversight.

Here I have assembled a set of recommendations and hints that can guide your ongoing development. In particular I have selected the currently best available Go library implementations for integrating our backend with tools that help us scale and become more robust.

Try. Test often. Fail often. Try harder.


## Strategic hints

It's important for good scalability to keep any and all __state__ at the client side. When servers must keep state your design is broken. Treat each HTTP request as an independent call. Even if originating from the same client, subsequent requests may be executed at different application servers. Only then you can balance load over multiple backends and provide uninterrupted service when one of the backends fails.

Sometime it is not entirely possible to have stateless servers, when, for example, you have long-running tasks or continuous queries such as large file uploads or WebSocket notifications. In such cases time out connections (see architecture patterns below) and make clients reconnect. This at least catches the case when clients are too slow or became disconnected for any reason.


- don't keep state across calls in the backend servers (the right place for state is in your data stores)
- never expect a remote component is alive when making a call; prepare for failure, timeout and retry calls, reestablish connections, and switch to alternative instances when reaching a retry threshold
- give meaningful error response to clients so they can recover from timeouts and backend failures by retrying an operation
- make all API calls idempotent; if ordering is an issue, use Lamport clocks to timestamp calls with sequence numbers at the origin
- use libraries that support *connection pools* for interfacing between your backend and storage servers or 3rd party services
- offload expensive or insecure computations to asynchronous services using message queues
- continuously (at all times) measure and analyse API queries to understand your typical load and when the API is in trouble
- regularly (from time to time) analyse queries to data store backends to find slowest queries and hot spots


## Recommended Tools

* [Gorilla](http://www.gorillatoolkit.org/) Go Web Toolkit (Context, Session Cookies) used for login session management
* [Redistore](https://github.com/boj/redistore) A session store backend for gorilla/sessions
* [Nginx](http://nginx.org/), a scalable web-server I propose to use in [reverse proxy mode](http://www.cyberciti.biz/faq/howto-linux-unix-setup-nginx-ssl-proxy/) for [terminating SSL/TLS connections](http://wiki.nginx.org/SSL-Offloader)
* [Redis](http://redis.io/) Scalable Key Value Store, used among other things to store counters for API-level quota support
* [Redigo](https://github.com/garyburd/redigo) Go client for Redis
* [Gearman](http://gearman.org/) a distributed job queue which is a good way for scaling out our video server distribution
* [gearman-go](https://github.com/mikespook/gearman-go) Gearman API for Golang
* Amazon S3, which is one option for scalable delivery of video segments
* [Goamz](https://wiki.ubuntu.com/goamz) Go Client for Amazon Webservices
* [qbs](https://github.com/coocood/qbs) a Query by Struct ORM


```
                          +-------------+
              +---------> |     CDN     | <--------------------------------+
              |           +-------------+                                  |
              v                                   ___________________      |
          +------+       +-------------+        +-------------------+|     |
          |Client| <---> |SSL-Nginx:443| <----> |Go API HTTP_mode:80||     |
          +------+       +-------------+        +-------------------+      |
                                                           |               |
             +-------------+--------------+-----------+----+----------+    |
             |             |              |           |               |    |
          _______      ___________     _______        |               |    v
        +-------+|  +------------+|  +-------+|  +---------+  +---------------+
        |  DB   ||  | App Caches ||  | redis ||  | gearman |  | Cloud Storage |
        +-------+   +------------+   +-------+   +---------+  +---------------+
            |                            |            |               ^
            v                            v         _________          |
       +---------+                  +---------+  +---------+|         |
       | Backup  |                  | Backup  |  | Workers ||<--------+
       +---------+                  +---------+  +---------+
```


There is also a large number of useful online services that help making the life of developers and operators less painful:

- [ApiGee](http://apigee.com) API gateway and insights
- [Drone IO](https://drone.io/) continuous integration testing service
- [Travis CI](https://travis-ci.org/) continuous integration testing service


# Architecture Patterns

I came across some quite useful patterns for structuring distributed applications. Here's a non-exhaustive list

- Golang [timeout pattern](http://blog.golang.org/go-concurrency-patterns-timing-out-and)
- Rate limiting with redis [here](http://redis.io/commands/incr) and [here](http://blog.domaintools.com/2013/04/rate-limiting-with-redis/)
- Sharding and Id's at [Instagram](http://instagram-engineering.tumblr.com/post/10853187575/sharding-ids-at-instagram)
- Fast key-value mapping in Redis for [millions of items](http://instagram-engineering.tumblr.com/post/12202313862/storing-hundreds-of-millions-of-simple-key-value-pairs)

## API Monitoring - Why, How and What?

The purpose of monitoring is to gain an understanding about the state and health of a system during normal operation and in extreme situations.

### Questions important to APP Developers
- Is the API error prone?
- Which API errors is my application seeing?
- How does the API usually perform?
- Is the API slow now?
- Which API methods are slow?
- Does the API have a quota?
- Is my app violating the API quota?
- How often does the API go down?
- Is the API down now?
- When will the API be back up?
- Why was the API down?

__Key Indicators__ (include them in an *is alive* API endpoint and in a developer dashboard)
- errors
- performance
- availability
- quota

### Questions important to API Desiners
- Which are our top applications?
- Who are our top application users?
- Who are our best application developers?
- Which API methods are most popular?
- How much API capacity will we need next year?
- Why is the API down?
- Why is the API slow?
- Why is the API throwing errors?
- Why is the API traffic spiking?
- Why did the API traffic disappear?

__Key indicators__  (include them in an *is alive* API endpoint and in an operator dashboard)
- application users
- applications
- developers
- API quality
- internal systems

## How to monitor and output data in a scalable way?
- every API call is an independent HTTP request, so keep them independent
- measure internal system performance in every call
  - e.g. use a context such as [Gorilla Context](http://www.gorillatoolkit.org/pkg/context) to store performance data along the call chain
  - combine performance data at the end of a call before returning a response
- put performance data into private response headers (x-rter-performance)
 - DB query times
 - message bus response times
 - external callout response times
- strip and log these response headers at frontend web proxies
- data-mine the monitoring data with separate real-time query engines
- display real-time data on operators and developer dashboards

###  Examples

Google Appengine statistics measuring the load of an application node

```
type Statistics struct {
    // CPU records the CPU consumed by this instance, in megacycles.
    CPU struct {
        Total   float64
        Rate1M  float64 // consumption rate over one minute
        Rate10M float64 // consumption rate over ten minutes
    }
    // RAM records the memory used by the instance, in megabytes.
    RAM struct {
        Current    float64
        Average1M  float64 // average usage over one minute
        Average10M float64 // average usage over ten minutes
    }
}
```


