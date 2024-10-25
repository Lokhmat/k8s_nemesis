# Textual use case description

1. Change system settings
   - Actor: Engineer (DevOps / ML engineer)
   - Description: Use case for user to configure load balancer settings in order for it to work properly
   - Preconditions: Load balancer system is up and running inside a k8s cluster
   - Postconditions: Load balancer system is configured as desired by user: scaling rule, load balancing policy, disabling policy
   - Main flow:
     - User opens view with current system settings
     - User changes system settings
     - User saves and applies new system settings
   - Alternative flows: -
   - Exceptions:
     - Failed to save settings - user recieves an error and advised to try again
   - Includes:
     - Set scaling rule
     - Set load balancing policy
     - Set application disabling rule

2. Set scaling rule
   - Actor: Engineer (DevOps / ML engineer)
   - Description: Use case for user to configure load balancer scaling rule
   - Preconditions: Load balancer system is up and running inside a k8s cluster
   - Postconditions: Scaling rule is set by user
   - Main flow:
     - User opens view with current system settings
     - User changes thresholds of k8s metrics based on which system will add new pods and kill pods
     - User saves and applies new scaling rule
   - Alternative flows: -
   - Exceptions:
     - Failed to save - user recieves an error and advised to try again
   - Includes and Extends: -

3. Set load balancing policy
   - Actor: Engineer (DevOps / ML engineer)
   - Description: Use case for user to configure load balancing policy
   - Preconditions: Load balancer system is up and running inside a k8s cluster
   - Postconditions: Load balancing policy is set by user
   - Main flow:
     - User opens view with current system settings
     - User changes currently selected load balancing policy. He chooses it from set of possible policies.
     - User saves and applies new policy
   - Alternative flows: -
   - Exceptions:
     - Failed to save - user recieves an error and advised to try again
   - Includes and Extends: -

4. Set application disabling rule
   - Actor: Engineer (DevOps / ML engineer)
   - Description: Use case for user to configure when application will be disabled fully(scaled to zero instances)
   - Preconditions: Load balancer system is up and running inside a k8s cluster
   - Postconditions: Disabling rule is set by user
   - Main flow:
     - User opens view with current system settings
     - User changes currently selected disabling rule based on thresholds of k8s cluster metrics.
     - User saves and applies new rule
   - Alternative flows: -
   - Exceptions:
     - Failed to save - user recieves an error and advised to try again
   - Includes and Extends: -

5. Proxy request
   - Actor: Client of the Application
   - Description: Use case for proxying and load balancing of requests of clients of the application that system is balancing.
   - Preconditions: Load balancer system is up and running inside a k8s cluster.
   - Postconditions: Client request is processed by application
   - Main flow:
     - Client sends request to an application through the load balancing system
     - Load balancing policy is applied to the request and pod that will process instance is selected
     - Request is proxied to selected pod and it's response is reterned to a client
   - Alternative flows:
     - If there are no currently running pods request will be buffered(see Buffer request use case) and then system will activate cluster(see Activate cluster use case) and proxy request to newly allocated pod
   - Exceptions:
     - External error - return 500 error code to client
   - Includes and Extends: -

6. Buffer request
   - Actor: Client of the Application
   - Description: Use case for buffering client request while there are no active and running application pods
   - Preconditions: Load balancer system is up and running inside a k8s cluster. There are no active application pods and client has submitted request to the system.
   - Postconditions: Client request is buffered until there are no active pods
   - Main flow:
     - Client sends request to an application through the load balancing system
     - Request is saved by the system, but no response is returned to client.
     - When cluster is activated and there are at least one running pod the request is proxied to it and response is returned to user
   - Alternative flows: -
   - Exceptions:
     - If pod is failed to start in preconfigured amount of time - system is failing request
   - Extends:
     - Proxy request 

7. Activate cluster
   - Actor: Client of the Application
   - Description: Use case for allocating new pod in a system when there are no currently running pods.
   - Preconditions: Load balancer system is up and running inside a k8s cluster. There are no active application pods and client has submitted request to the system.
   - Postconditions: New up and running pod is created
   - Main flow:
     - Client sends request to an application through the load balancing system
     - System allocates new pod and waits until it will start
     - When pod is active system starts to serve requests to it
   - Alternative flows: -
   - Exceptions:
     - If pod is failed to start in preconfigured amount of time - system will try again
   - Extends:
     - Proxy request 

8. Manage cluster
   - Actor: Scaling trigger
   - Description: Use case for managing amount of currently running pods in a cluster.
   - Preconditions: Load balancer system is up and running inside a k8s cluster. Appliction is up and running on at least one pod.
   - Postconditions: Changes are made in amount of running pods.
   - Main flow:
     - System gets currently running pods of a k8s cluster
     - System is fetching k8s cluster metrics such as CPU usage, RAM usage, etc.
     - If one of preconfigured thresholds for application scaling or disabling is met than system changes amount of cluster pods accordingly
   - Alternative flows: -
   - Exceptions:
     - If system fases an error on any of the steps - it tries again
   - Includes:
     - Get current pods liveliness state
     - Get k8s cluster metrics

9. Get current pods liveliness state
   - Actor: Scaling trigger
   - Description: Use case for getting current cluster up and running pods
   - Preconditions: Load balancer system is up and running inside a k8s cluster. Appliction is up and running on at least one pod.
   - Postconditions: Information of current up and running pods is returned to the system
   - Main flow:
     - System submits request to the k8s cluster in order to get currently running pods
   - Alternative flows: -
   - Exceptions:
     - If system gets an error it returns an error back to the system
   - Includes and Extends: -

10. Get k8s metrics
   - Actor: Scaling trigger
   - Description: Use case for getting current k8s metrics from the system
   - Preconditions: Load balancer system is up and running inside a k8s cluster. Appliction is up and running on at least one pod.
   - Postconditions: Metrics of currently up and running pods is returned to the system
   - Main flow:
     - System submits request to the k8s cluster in order to get metrics of currently running pods
   - Alternative flows: -
   - Exceptions:
     - If system gets an error it returns an error back to the system
   - Includes and Extends: -

11. Decrease cluster pods amount
   - Actor: Scaling trigger
   - Description: Use case for reducing amount of pods running inside the cluster
   - Preconditions: Load balancer system is up and running inside a k8s cluster. Appliction is up and running on at least two pods.
   - Postconditions: Amount of up and running pods inside of the cluster is reduced by one
   - Main flow:
     - System submits request to the k8s cluster in order to reduce amount of application pods by 1
     - System waits until pod is killed and removed
   - Alternative flows: -
   - Exceptions:
     - If system gets an error it tries again
   - Extends:
     - Manage cluster

12. Increase cluster pods amount
   - Actor: Scaling trigger
   - Description: Use case for increasing amount of pods running inside the cluster
   - Preconditions: Load balancer system is up and running inside a k8s cluster. Appliction is up and running on at least one pod.
   - Postconditions: Amount of up and running pods inside of the cluster is increased by one
   - Main flow:
     - System submits request to the k8s cluster in order to allocate and start new application pod
     - System waits until it is created and alive
   - Alternative flows: -
   - Exceptions:
     - If system gets an error it tries again
   - Extends:
     - Manage cluster

13. Disable cluster
   - Actor: Scaling trigger
   - Description: Use case for scaling k8s cluster to 0 instances
   - Preconditions: Load balancer system is up and running inside a k8s cluster. Appliction is up and running on exactly one pod.
   - Postconditions: Application is completely stopped and there are now up and running pods.
   - Main flow:
     - System submits request to the k8s cluster in order to allocate to stop last running pod
     - System waits until k8s will kill and remove pod
     - System validates that there are no currently running application pods
   - Alternative flows: -
   - Exceptions:
     - If system gets an error it tries again
   - Extends:
     - Manage cluster
