## K8s Nemesis


Load balancer with service turn-off


A vanilla k8s distribution provides only basic application scaling and load balancing functions. A problem remains that a particular application deployment cannot be scaled to zero instances, even if not used for a long time. Moreover, ML applications may have specific requirements to the underlying compute resources and may be slow to start. The goal of the project is to implement a cloud-native service that will solve both problems: scaling to zero instances, complex automata based application scaling responding to monitoring events.


