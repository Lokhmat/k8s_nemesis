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


## Features


- Automatically scale K8s cluster by monitoring events, with ability to scale to 0 instances 
- Integrate with various monitoring tools to gather real-time data and metrics, which can trigger scaling events
- User-Defined Scaling Rules
- Authentication and authorization
- Load balancing 
- Reports and alerting
- UI with analytics




## Constraints:
k8s running environment


