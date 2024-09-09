# K8s Nemesis


Team members:


Lopatkin Vasiliy @LDFCZ


Roman Kuzmenko @qTronor


Sergey Lokhmatikov @Lokhmat


Andrey Tamplon @AndreyTamplon




## Load balancer with service turn-off


A vanilla k8s distribution provides only basic application scaling and load balancing functions. A problem remains that a particular application deployment cannot be scaled to zero instances, even if not used for a long time. Moreover, ML applications may have specific requirements to the underlying compute resources and may be slow to start. The goal of the project is to implement a cloud-native service that will solve both problems: scaling to zero instances, complex automata based application scaling responding to monitoring events.




## Stakeholders


Developer Team

DevOps

MlOps engineers

Backend engineers



## Features:


1. **Automatically scale K8s cluster by monitoring events, with ability to scale to 0 instances.**

	Implement logic that allows Kubernetes (K8s) to scale applications down to zero instances when they are idle or unused, addressing the problem of unnecessary resource consumption.
2. **Integrate with various monitoring tools to gather real-time data and metrics, which can trigger scaling events.**

	Integration with popular monitoring systems (such as Prometheus, Grafana, etc.) to collect real-time data. These metrics will be used to trigger scaling decisions.
3. **User-Defined Scaling Rules**

	Provide users with the ability to define their own scaling policies based on specific metrics or events. These policies could specify thresholds for scaling up, scaling down, or scaling to zero.
4. **Authentication and authorization**

	Secure the service using authentication via password.
5. **Load balancing**

	Implement a custom load balancing mechanism that distributes traffic or workload across available application instances. The load balancer will dynamically adjust based on the current state of the scaled application, considering scenarios where instances may be scaled to zero.
6. **Reports and alerting**

	Generate detailed reports on scaling events, resource utilization, and application performance over time. Provide real-time alerts to users when scaling events are triggered, or when thresholds are exceeded.
7. **UI with analytics**

	A user-friendly web interface that displays analytics, including historical scaling events, current cluster status, and usage statistics.




## Constraints:
1. **Usability**  
   The system interface should be intuitive for users with varying technical backgrounds.


2. **k8s running environment**  
   Balancer must be able to interact with kubernetes entities
3. **Disaster Recovery**  
   Balancer service must recover automatically after a failure.




